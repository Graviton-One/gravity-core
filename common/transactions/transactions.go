package transactions

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/common/account"
	"gravity-hub/common/keys"
	"gravity-hub/common/score"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dgraph-io/badger"
	_ "github.com/tendermint/tendermint/crypto/ed25519"
	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
)

type TxFunc string

const (
	Commit            TxFunc = "commit"
	Reveal            TxFunc = "reveal"
	AddOracle         TxFunc = "addOracle"
	AddOracleInNebula TxFunc = "addOracleInNebula"
	SignResult        TxFunc = "signResult"
	NewRound          TxFunc = "newRound"
	Vote              TxFunc = "vote"
)

type Transaction struct {
	Id           string
	SenderPubKey string
	Signature    string
	Func         TxFunc
	Timestamp    time.Time
	ChainType    account.ChainType
	Args         string
}

func New(pubKey []byte, funcName TxFunc, chainType account.ChainType, privKey tendermintCrypto.PrivKeyEd25519, args []byte) (*Transaction, error) {
	tx := &Transaction{
		SenderPubKey: hexutil.Encode(pubKey),
		Args:         hexutil.Encode(args),
		Func:         funcName,
		ChainType:    chainType,
		Timestamp:    time.Now(),
	}
	tx.Hash()

	err := tx.Sign(privKey)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (tx *Transaction) Hash() {
	tx.Id = hexutil.Encode(crypto.Keccak256(tx.MarshalBytesWithoutSig()))
}

func (tx *Transaction) Sign(privKey tendermintCrypto.PrivKeyEd25519) error {
	txIdeBytes, err := hexutil.Decode(tx.Id)
	if err != nil {
		return err
	}
	sign, err := account.Sign(privKey, txIdeBytes)
	if err != nil {
		return err
	}
	tx.Signature = hexutil.Encode(sign)
	return nil
}

func (tx *Transaction) MarshalBytesWithoutSig() []byte {
	var result []byte
	result = append(result, tx.Id[:]...)
	result = append(result, tx.SenderPubKey[:]...)
	result = append(result, tx.Func...)
	result = append(result, byte(tx.ChainType))
	result = append(result, tx.Args...)

	time := tx.Timestamp.Unix()
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(time))
	result = append(result, b[:]...)

	return result
}

func UnmarshalJson(data []byte) (*Transaction, error) {
	tx := new(Transaction)
	err := json.Unmarshal(data, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tx *Transaction) IsValid(ethClient *ethclient.Client, wavesClient *client.Client, db *badger.DB, ctx context.Context) error {
	if !tx.isValidSigns() {
		return errors.New("invalid signature")
	}

	switch tx.Func {
	case Commit:
		return tx.isValidCommit(db)
	case Reveal:
		return tx.isValidReveal(db)
	case AddOracle:
		return nil
	case AddOracleInNebula:
		return tx.isValidAddOracleInNebula(db)
	case SignResult:
		return tx.isValidSignResult(db)
	case NewRound:
		return tx.isValidNewRound(ethClient, wavesClient, db, ctx)
	default:
		return errors.New(fmt.Sprintf("function '%s' is not found", string(tx.Func)))
	}
}

func (tx *Transaction) isValidSigns() bool {
	pubKeyBytes, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return false
	}

	sigBytes, err := hexutil.Decode(tx.Signature)
	if err != nil {
		return false
	}
	txIdBytes, err := hexutil.Decode(tx.Id)
	if err != nil {
		return false
	}

	return ed25519.Verify(pubKeyBytes, txIdBytes, sigBytes)
}

func (tx *Transaction) isValidAddOracleInNebula(db *badger.DB) error {
	if len(tx.Args) == 64 {
		return errors.New("invalid args size")
	}

	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	pubKey := args[32:]

	key := []byte(keys.FormOraclesByNebulaKey(nebulaAddress))

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			return nil
		}

		oraclesByNebula := make(map[string]string)
		b, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &oraclesByNebula)
		if err != nil {
			return err
		}

		if _, ok := oraclesByNebula[hexutil.Encode(pubKey)]; !ok {
			return nil
		}

		return errors.New("validator is exist")
	})

	pubKeyOwner, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}

	scoreKey := keys.FormScoreKey(pubKeyOwner)
	var scoreValue float32
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(scoreKey))
		if err != nil {
			return err
		}

		err = item.Value(func(val []byte) error {
			scoreValue = score.UInt64ToFloat32Score(binary.BigEndian.Uint64(val))
			return err
		})
		return err
	})

	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	if scoreValue < 0 {
		return errors.New("invalid score. score <= 0")
	}

	return nil
}

func (tx *Transaction) isValidCommit(db *badger.DB) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}
	length := args[0]
	nebula := args[1 : 1+length]
	height := args[1+length : 9+length]
	pubKey := args[41+length:]

	key := keys.FormCommitKey(nebula, binary.BigEndian.Uint64(height), pubKey)
	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}

		return errors.New("commit is exist")
	})
	return nil
}

func (tx *Transaction) isValidReveal(db *badger.DB) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	commit := args[0:32]
	length := int(args[32])
	nebula := args[33 : 33+length]
	height := args[33+length : 41+length]
	lengthReveal := int(args[41+length])
	reveal := args[42+length : lengthReveal+42+length]
	pubKey := args[lengthReveal+42+length:]

	revealKey := keys.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)

	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(revealKey))
		if err == badger.ErrKeyNotFound {
			return nil
		}

		return errors.New("reveal is exist")
	})

	var commitBytes []byte
	keyCommit := keys.FormCommitKey(nebula, binary.BigEndian.Uint64(height), pubKey)
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(keyCommit))
		if err == badger.ErrKeyNotFound {
			return errors.New("commit is not exist")
		}
		if err != nil {
			return err
		}
		return item.Value(func(value []byte) error {
			commitBytes = value
			return nil
		})
	})
	if err != nil {
		return err
	}

	expectedHash := crypto.Keccak256(reveal)
	if !bytes.Equal(commitBytes, expectedHash[:]) {
		return errors.New("invalid reveal")
	}
	return nil
}

func (tx *Transaction) isValidSignResult(db *badger.DB) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	length := args[0]
	nebulaAddress := args[1 : 1+length]
	heightBytes := args[1+length : 9+length]
	resultHash := args[9+length : 41+length]
	signBytes := args[41+length:]

	height := binary.BigEndian.Uint64(heightBytes)
	prefix := strings.Join([]string{string(keys.RevealKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", height)}, "_")

	var reveals []uint64
	db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Seek([]byte(prefix)); it.ValidForPrefix([]byte(prefix)); it.Next() {
			item := it.Item()
			item.Value(func(v []byte) error {
				reveals = append(reveals, binary.BigEndian.Uint64(v))
				return nil
			})
		}
		return nil
	})

	var average uint64
	for _, v := range reveals {
		average += v
	}
	value := uint64(float64(average) / float64(len(reveals)))

	bytesValue := make([]byte, 8)
	binary.BigEndian.PutUint64(bytesValue, value)
	hash := crypto.Keccak256(bytesValue)

	if bytes.Compare(resultHash, hash[:]) != 0 {
		return errors.New("invalid result hash")
	}
	senderPubKeyBytes, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}

	key := []byte(keys.FormOraclesByValidatorKey(senderPubKeyBytes))
	oracles := make(map[account.ChainType][]byte)

	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}

		b, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &oracles)
		if err != nil {
			return err
		}
		return nil
	})

	switch tx.ChainType {
	case account.Ethereum:
		if !crypto.VerifySignature(senderPubKeyBytes, resultHash, oracles[account.Ethereum][0:64]) {
			return errors.New("invalid result hash sign")
		}
	case account.Waves:
		var senderPubKey [32]byte
		copy(senderPubKey[:], oracles[account.Waves])
		sign := wavesCrypto.Signature{}
		copy(sign[:], signBytes)

		if !wavesCrypto.Verify(senderPubKey, sign, resultHash) {
			return errors.New("invalid result hash sign")
		}
	default:
		return errors.New("invalid result hash sign")
	}

	return err
}

func (tx *Transaction) isValidNewRound(ethClient *ethclient.Client, wavesClient *client.Client, db *badger.DB, ctx context.Context) error {
	var key string
	var height uint64
	key = keys.FormBlockKey(tx.ChainType, height)
	switch tx.ChainType {
	case account.Ethereum:
		ethHeight, err := ethClient.BlockByNumber(ctx, nil)
		if err != nil {
			return err
		}
		height = ethHeight.NumberU64()
	case account.Waves:
		wavesHeight, _, err := wavesClient.Blocks.Height(ctx)
		if err != nil {
			return err
		}
		height = wavesHeight.Height
	default:
		return errors.New("invalid chain type")
	}

	err := db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			return nil
		}
		return errors.New("eth height is exist")
	})
	if err != nil {
		return err
	}
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}
	txHeightBytes := args[:8]
	txHeight := binary.BigEndian.Uint64(txHeightBytes)
	if txHeight != height {
		return errors.New("invalid height")
	}
	return nil
}

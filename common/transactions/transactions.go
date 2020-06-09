package transactions

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/common/account"
	"gravity-hub/common/keys"
	"strings"

	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dgraph-io/badger"
	_ "github.com/tendermint/tendermint/crypto/ed25519"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
)

type TxFunc string
type ChainType byte

const (
	Commit       TxFunc = "commit"
	Reveal       TxFunc = "reveal"
	AddValidator TxFunc = "addValidator"
	SignResult   TxFunc = "signResult"
	NewRound     TxFunc = "newRound"
)
const (
	Ethereum ChainType = iota
	Waves
)

type Transaction struct {
	Id           string
	SenderPubKey string
	Signature    string
	Func         TxFunc
	ChainType    account.ChainType
	Args         string
}

func New(pubKey []byte, funcName TxFunc, chainType account.ChainType, privKey []byte, args []byte) (*Transaction, error) {
	tx := &Transaction{
		SenderPubKey: hex.EncodeToString(pubKey),
		Args:         hex.EncodeToString(args),
		Func:         funcName,
		ChainType:    chainType,
	}
	tx.Hash()

	err := tx.Sign(privKey)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (tx *Transaction) Hash() {
	tx.Id = hex.EncodeToString(crypto.Keccak256(tx.MarshalBytesWithoutSig()))
}

func (tx *Transaction) Sign(privKey []byte) error {
	txIdeBytes, err := hex.DecodeString(tx.Id)
	if err != nil {
		return err
	}
	sign := account.Sign(privKey, txIdeBytes)
	tx.Signature = hex.EncodeToString(sign)
	return nil
}

func (tx *Transaction) MarshalBytesWithoutSig() []byte {
	var result []byte
	result = append(result, tx.Id[:]...)
	result = append(result, tx.SenderPubKey[:]...)
	result = append(result, tx.Func...)
	result = append(result, byte(tx.ChainType))
	result = append(result, tx.Args...)
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
	case AddValidator:
		return tx.isValidAddValidator(db)
	case SignResult:
		return tx.isValidSignResult(db)
	case NewRound:
		return tx.isValidNewRound(ethClient, wavesClient, db, ctx)
	default:
		return errors.New(fmt.Sprintf("function '%s' is not found", string(tx.Func)))
	}
}

func (tx *Transaction) isValidSigns() bool {
	pubKeyBytes, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return false
	}

	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}
	txIdBytes, err := hex.DecodeString(tx.Id)
	if err != nil {
		return false
	}
	switch tx.ChainType {
	case account.Ethereum:
		return crypto.VerifySignature(pubKeyBytes, txIdBytes, sigBytes[0:64])
	case account.Waves:
		pubKey := wavesCrypto.PublicKey{}
		copy(pubKey[:], pubKeyBytes)
		sig := wavesCrypto.Signature{}
		copy(sig[:], sigBytes)

		return wavesCrypto.Verify(pubKey, sig, txIdBytes)
	default:
		return false
	}
}

func (tx *Transaction) isValidAddValidator(db *badger.DB) error {
	if len(tx.Args) == 64 {
		return errors.New("invalid args size")
	}

	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	pubKey := args[32:]
	key := keys.FormValidatorKey(nebulaAddress, pubKey)

	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			return nil
		}

		return errors.New("validator is exist")
	})

	return err
}

func (tx *Transaction) isValidCommit(db *badger.DB) error {
	if len(tx.Args) == 72 {
		return errors.New("invalid commit size")
	}
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}
	nebula := args[0:32]
	height := args[32:40]
	sender, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}

	key := keys.FormCommitKey(nebula, binary.BigEndian.Uint64(height), sender)
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
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	commit := args[0:32]
	nebula := args[32:64]
	height := args[64:72]
	reveal := args[72:]
	revealKey := keys.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)

	err = db.View(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(revealKey))
		if err == badger.ErrKeyNotFound {
			return nil
		}

		return errors.New("reveal is exist")
	})

	sender, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}

	var commitBytes []byte
	keyCommit := keys.FormCommitKey(nebula, binary.BigEndian.Uint64(height), sender)
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
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	heightBytes := args[32:40]
	resultHash := args[40:72]
	signBytes := args[72:]

	height := binary.BigEndian.Uint64(heightBytes)
	prefix := strings.Join([]string{string(keys.RevealKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", height)}, "_")

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
	senderPubKeyBytes, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}
	switch tx.ChainType {
	case account.Ethereum:
		if !crypto.VerifySignature(senderPubKeyBytes, resultHash, signBytes[0:64]) {
			return errors.New("invalid result hash sign")
		}
	case account.Waves:
		var senderPubKey [32]byte
		copy(senderPubKey[:], senderPubKeyBytes)
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
	args, err := hex.DecodeString(tx.Args)
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

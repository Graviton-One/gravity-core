package transaction

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/ledger-node/state"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/dgraph-io/badger"
)

type TxFunc string

const (
	Commit       TxFunc = "commit"
	Reveal       TxFunc = "reveal"
	AddValidator TxFunc = "addValidator"
	SignResult   TxFunc = "signResult"
)

type Transaction struct {
	Id           string
	SenderPubKey string
	Signature    string
	Func         TxFunc
	Args         string
}

func (tx *Transaction) MarshalBytesWithoutSig() []byte {
	var result []byte
	result = append(result, tx.Id[:]...)
	result = append(result, tx.SenderPubKey[:]...)
	result = append(result, tx.Func...)
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

func (tx *Transaction) IsValid(db *badger.DB) error {
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

	return crypto.VerifySignature(pubKeyBytes, txIdBytes, sigBytes[0:64])
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
	key := state.FormValidatorKey(nebulaAddress, pubKey)

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

	key := state.FormCommitKey(nebula, binary.BigEndian.Uint64(height), sender)
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
	revealKey := state.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)

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
	keyCommit := state.FormCommitKey(nebula, binary.BigEndian.Uint64(height), sender)
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
	prefix := strings.Join([]string{string(state.RevealKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", height)}, "_")

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
	binary.LittleEndian.PutUint64(bytesValue, value)
	hash := crypto.Keccak256(bytesValue)
	validationMsg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(hash))
	currentResultHash := crypto.Keccak256(append([]byte(validationMsg), hash...))

	if bytes.Compare(resultHash, currentResultHash[:]) != 0 {
		return errors.New("invalid result hash")
	}
	senderPubKeyBytes, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}
	if !crypto.VerifySignature(senderPubKeyBytes, resultHash, signBytes[0:64]) {
		return errors.New("invalid result hash sign")
	}

	return err
}

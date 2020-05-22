package transaction

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/state"
	"strings"

	"github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/dgraph-io/badger"
)

type TxFunc string

const (
	Commit       TxFunc = "commit"
	Reveal       TxFunc = "reveal"
	AddValidator TxFunc = "addValidator"
	SignResult   TxFunc = "signResult"
)

const (
	HashKeySize   = 32
	PublicKeySize = 32
	SignatureSize = 64
)

type Transaction struct {
	Id           string //[HashKeySize]byte
	SenderPubKey string // [PublicKeySize]byte
	Signature    string // [SignatureSize]byte
	Func         TxFunc //[]byte
	Args         string //Args
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

	switch TxFunc(tx.Func) {
	case Commit:
		return tx.isValidCommit()
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
	pubKey := crypto.PublicKey{}
	copy(pubKey[:], pubKeyBytes)

	sigBytes, err := hex.DecodeString(tx.Signature)
	if err != nil {
		return false
	}
	sig := crypto.Signature{}
	copy(sig[:], sigBytes)

	return crypto.Verify(pubKey, sig, tx.MarshalBytesWithoutSig())
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

func (tx *Transaction) isValidCommit() error {
	if len(tx.Args) == 72 {
		return errors.New("invalid commit size")
	}
	return nil
}

func (tx *Transaction) isValidReveal(db *badger.DB) error {
	if len(tx.Args) == 32 {
		return errors.New("invalid reveal size")
	}
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	commit := args[0:32]
	nebula := args[32:64]
	height := args[64:72]
	reveal := args[72:]
	key := state.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)
	sha256.Sum256(reveal)

	var commitTxBytes []byte
	var commitTx Transaction
	err = db.View(func(txn *badger.Txn) error {
		commitTxItem, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		commitTxBytes, err = commitTxItem.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	err = json.Unmarshal(commitTxBytes, &commitTx)
	if err != nil {
		return err
	}

	argsCommitTx, err := hex.DecodeString(commitTx.Args)
	if err != nil {
		return err
	}

	expectedHash := sha256.Sum256(reveal)
	if !bytes.Equal(argsCommitTx[0:32], expectedHash[:]) {
		return errors.New("invalid reveal")
	}
	return nil
}

func (tx *Transaction) isValidSignResult(db *badger.DB) error {
	if len(tx.Args) == 136 {
		return errors.New("invalid args size")
	}

	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	heightBytes := args[32:40]
	signBytes := args[40:72]
	resultHash := args[72:]

	sign := crypto.Signature{}
	copy(sign[:], signBytes)

	height := binary.BigEndian.Uint64(heightBytes)
	prefix := strings.Join([]string{string(state.RevealKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", height)}, "_")

	var realResultHash []byte
	err = db.View(func(txn *badger.Txn) error {
		iterator := txn.NewIterator(badger.IteratorOptions{Prefix: []byte(prefix)})
		hash := sha256.New()
		for iterator.Valid() {

			iterator.Next()
			reveal, err := iterator.Item().ValueCopy(nil)
			if err != nil {
				return err
			}
			hash.Write(reveal)
		}
		realResultHash = hash.Sum(nil)

		return errors.New("validator is exist")
	})

	if bytes.Compare(resultHash, realResultHash) != 0 {
		return errors.New("invalid result hash")
	}
	senderPubKeyBytes, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}
	var senderPubKey [32]byte
	copy(senderPubKey[:], senderPubKeyBytes)
	if !crypto.Verify(senderPubKey, sign, resultHash) {
		return errors.New("invalid result hash sign")
	}

	return err
}

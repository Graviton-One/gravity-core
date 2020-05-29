package transaction

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"gravity-hub/ledger-node/state"

	"github.com/dgraph-io/badger"
)

func (tx *Transaction) SetState(currentBatch *badger.Txn) error {
	switch tx.Func {
	case Commit:
		return tx.SetStateCommit(currentBatch)
	case Reveal:
		return tx.SetStateReveal(currentBatch)
	case AddValidator:
		return tx.SetStateAddValidator(currentBatch)
	case SignResult:
		return tx.SetStateSignResult(currentBatch)
	default:
		return errors.New(fmt.Sprintf("function '%s' is not found", string(tx.Func)))
	}
}

func (tx *Transaction) SetStateCommit(currentBatch *badger.Txn) error {
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}
	nebula := args[0:32]
	height := args[32:40]
	commit := args[40:]
	sender, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}

	key := state.FormCommitKey(nebula, binary.BigEndian.Uint64(height), sender)
	return currentBatch.Set([]byte(key), commit)
}

func (tx *Transaction) SetStateReveal(currentBatch *badger.Txn) error {
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	commit := args[0:32]
	nebula := args[32:64]
	height := args[64:72]
	reveal := args[72:]

	key := state.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)
	return currentBatch.Set([]byte(key), reveal)
}

func (tx *Transaction) SetStateAddValidator(currentBatch *badger.Txn) error {
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	pubKey := args[32:]
	key := state.FormValidatorKey(nebulaAddress, pubKey)
	return currentBatch.Set([]byte(key), []byte{1})
}

func (tx *Transaction) SetStateSignResult(currentBatch *badger.Txn) error {
	args, err := hex.DecodeString(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	height := args[32:40]
	signBytes := args[72:]

	sender, err := hex.DecodeString(tx.SenderPubKey)
	if err != nil {
		return err
	}

	keySign := state.FormSignResultKey(nebulaAddress, binary.BigEndian.Uint64(height), sender)
	return currentBatch.Set([]byte(keySign), signBytes)
}

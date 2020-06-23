package transactions

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/common/keys"
	"gravity-hub/score-calculator/models"

	"github.com/ethereum/go-ethereum/common/hexutil"

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
	case NewRound:
		return tx.SetStatesNewRound(currentBatch)
	case Vote:
		return tx.SetVote(currentBatch)
	default:
		return errors.New(fmt.Sprintf("function '%s' is not found", string(tx.Func)))
	}
}

func (tx *Transaction) SetStateCommit(currentBatch *badger.Txn) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}
	nebula := args[0:32]
	height := args[32:40]
	commit := args[40:]
	sender, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}

	key := keys.FormCommitKey(nebula, binary.BigEndian.Uint64(height), sender)
	return currentBatch.Set([]byte(key), commit)
}

func (tx *Transaction) SetStateReveal(currentBatch *badger.Txn) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	commit := args[0:32]
	nebula := args[32:64]
	height := args[64:72]
	reveal := args[72:]

	key := keys.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)
	return currentBatch.Set([]byte(key), reveal)
}

func (tx *Transaction) SetStateAddValidator(currentBatch *badger.Txn) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	chainType := args[:1]
	nebulaAddress := args[1:33]
	pubKey := args[33:]
	key := []byte(keys.FormValidatorKey(nebulaAddress, pubKey))
	err = currentBatch.Set(key, chainType)
	if err != nil {
		return err
	}

	key = []byte(keys.FormNebulaeByValidatorKey(pubKey))
	item, err := currentBatch.Get(key)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	var nebulae []string
	if err != badger.ErrKeyNotFound {
		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(value, &nebulae)
		if err != nil {
			return err
		}
	}
	nebulae = append(nebulae, hexutil.Encode(nebulaAddress))

	b, err := json.Marshal(nebulae)
	if err != nil {
		return err
	}

	err = currentBatch.Set(key, b)
	if err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) SetStateSignResult(currentBatch *badger.Txn) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	nebulaAddress := args[:32]
	height := args[32:40]
	signBytes := args[72:]

	sender, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}

	keySign := keys.FormSignResultKey(nebulaAddress, binary.BigEndian.Uint64(height), sender)
	return currentBatch.Set([]byte(keySign), signBytes)
}

func (tx *Transaction) SetStatesNewRound(currentBatch *badger.Txn) error {
	var key string
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}
	txTcHeightBytes := args[:8]
	txTcHeight := binary.BigEndian.Uint64(txTcHeightBytes)

	key = keys.FormBlockKey(tx.ChainType, txTcHeight)
	currentBatch.Set([]byte(key), args[8:16])
	return nil
}

func (tx *Transaction) SetVote(currentBatch *badger.Txn) error {
	var key string
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}
	pubKey, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}
	var votes []models.Vote
	err = json.Unmarshal(args, &votes)
	if err != nil {
		return err
	}

	b, err := json.Marshal(votes)
	if err != nil {
		return err
	}

	key = keys.FormVoteKey(pubKey)
	currentBatch.Set([]byte(key), b)
	return nil
}

package transactions

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/account"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/keys"
	"github.com/Gravity-Hub-Org/proof-of-concept/score-calculator/models"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/dgraph-io/badger"
)

func (tx *Transaction) SetState(currentBatch *badger.Txn) error {
	switch tx.Func {
	case Commit:
		return tx.SetStateCommit(currentBatch)
	case Reveal:
		return tx.SetStateReveal(currentBatch)
	case AddOracleInNebula:
		return tx.SetStateAddOracleInNebula(currentBatch)
	case AddOracle:
		return tx.SetStateAddOracle(currentBatch)
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
	length := args[0]
	nebula := args[1 : 1+length]
	height := args[1+length : 9+length]
	commit := args[9+length : 41+length]
	pubKey := args[41+length:]

	key := keys.FormCommitKey(nebula, binary.BigEndian.Uint64(height), pubKey)
	return currentBatch.Set([]byte(key), commit)
}

func (tx *Transaction) SetStateReveal(currentBatch *badger.Txn) error {
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
	//	pubKey := args[lengthReveal+41+length:]

	key := keys.FormRevealKey(nebula, binary.BigEndian.Uint64(height), commit)
	return currentBatch.Set([]byte(key), reveal)
}

func (tx *Transaction) SetStateAddOracleInNebula(currentBatch *badger.Txn) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	length := args[0:1][0]
	nebulaAddress := args[1 : int(length)+1]
	pubKey := args[int(length)+1:]

	key := []byte(keys.FormOraclesByNebulaKey(nebulaAddress))
	item, err := currentBatch.Get(key)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	var b []byte
	oraclesByNebula := make(map[string]string)
	if item != nil {
		b, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &oraclesByNebula)
		if err != nil {
			return err
		}
	}
	oraclesByNebula[hexutil.Encode(pubKey)] = tx.SenderPubKey

	b, err = json.Marshal(&oraclesByNebula)
	if err != nil {
		return err
	}

	err = currentBatch.Set(key, b)
	if err != nil {
		return err
	}

	return nil
}

func (tx *Transaction) SetStateAddOracle(currentBatch *badger.Txn) error {
	args, err := hexutil.Decode(tx.Args)
	if err != nil {
		return err
	}

	chainType := args[:1]
	pubKey := args[1:]

	pubKeyOwner, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}
	key := []byte(keys.FormOraclesByValidatorKey(pubKeyOwner))
	item, err := currentBatch.Get(key)
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	oracles := make(map[account.ChainType][]byte)
	if err != badger.ErrKeyNotFound {
		value, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(value, &oracles)
		if err != nil {
			return err
		}
	}

	oracles[account.ChainType(chainType[0])] = pubKey
	b, err := json.Marshal(oracles)
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

	length := args[0]
	nebulaAddress := args[1 : 1+length]
	height := args[1+length : 9+length]
	signBytes := args[9+32+length:]

	senderPubKeyBytes, err := hexutil.Decode(tx.SenderPubKey)
	if err != nil {
		return err
	}

	key := []byte(keys.FormOraclesByValidatorKey(senderPubKeyBytes))
	oracles := make(map[account.ChainType][]byte)

	item, err := currentBatch.Get(key)
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

	var senderPubKey []byte
	switch tx.ChainType {
	case account.Ethereum:
		senderPubKey = oracles[account.Ethereum]
	case account.Waves:
		senderPubKey = oracles[account.Waves]
	}

	keySign := keys.FormSignResultKey(nebulaAddress, binary.BigEndian.Uint64(height), senderPubKey)
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

	for _, v := range votes {
		pubKey, err := hexutil.Decode(v.Target)
		if err != nil {
			return err
		}
		if _, err := currentBatch.Get([]byte(keys.FormScoreKey(pubKey))); err == badger.ErrKeyNotFound {
			err := currentBatch.Set([]byte(keys.FormScoreKey(pubKey)), make([]byte, 8, 8))
			if err != nil {
				return err
			}
		}
	}

	b, err := json.Marshal(votes)
	if err != nil {
		return err
	}

	key = keys.FormVoteKey(pubKey)
	currentBatch.Set([]byte(key), b)
	return nil
}

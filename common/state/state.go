package state

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"

	"github.com/Gravity-Tech/proof-of-concept/common/transactions"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/Gravity-Tech/proof-of-concept/common/account"
	"github.com/Gravity-Tech/proof-of-concept/common/storage"

	"github.com/ethereum/go-ethereum/crypto"
)

type SubRound int64

const (
	SubRoundDuration = 2
	SubRoundCount    = 3

	CommitSubRound SubRound = iota
	RevealSubRound
	ResultSubRound
)

var (
	ErrInvalidSign       = errors.New("invalid signature")
	ErrFuncNotFound      = errors.New("function is not found")
	ErrRevealIsExist     = errors.New("reveal is exist")
	ErrCommitIsExist     = errors.New("commit is exist")
	ErrCommitIsNotExist  = errors.New("commit is not exist")
	ErrInvalidReveal     = errors.New("invalid reveal")
	ErrNewRound          = errors.New("round is exist")
	ErrAddOracleInNebula = errors.New("oracle was added in nebula")
	ErrInvalidScore      = errors.New("invalid score. score <= 0")
	ErrInvalidChainType  = errors.New("invalid chain type")
	ErrInvalidHeight     = errors.New("invalid height")
	ErrInvalidSubRound   = errors.New("invalid sub round")
)

func calculateSubRound(height int64) SubRound {
	return SubRound(height % (SubRoundCount * SubRoundDuration) / SubRoundDuration)
}

func SetState(tx *transactions.Transaction, storage *storage.Storage, height int64, ethClient *ethclient.Client, wavesClient *client.Client, ctx context.Context) error {
	if isValidSigns(tx) {
		return ErrInvalidSign
	}

	switch tx.Func {
	case transactions.Commit:
		if calculateSubRound(height) != CommitSubRound {
			return ErrInvalidSubRound
		}
		return commit(storage, tx)
	case transactions.Reveal:
		if calculateSubRound(height) != RevealSubRound {
			return ErrInvalidSubRound
		}
		return reveal(storage, tx)
	case transactions.Result:
		if calculateSubRound(height) != ResultSubRound {
			return ErrInvalidSubRound
		}
		return result(storage, tx)
	case transactions.AddOracleInNebula:
		return addOracleInNebula(storage, tx)
	case transactions.AddOracle:
		return addOracle(storage, tx)
	case transactions.NewRound:
		return newRound(storage, tx, height, ethClient, wavesClient, ctx)
	case transactions.Vote:
		return vote(storage, tx)
	default:
		return ErrFuncNotFound
	}
}

func commit(store *storage.Storage, tx *transactions.Transaction) error {
	nebula := tx.Args[0].Value.([]byte)
	height := tx.Args[1].Value.(int64)
	commit := tx.Args[2].Value.([]byte)
	pubKey := tx.Args[3].Value.([]byte)

	_, err := store.CommitHash(nebula, height, pubKey)
	if err == storage.ErrKeyNotFound {
		err := store.SetCommitHash(nebula, height, pubKey, commit)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		return ErrCommitIsExist
	}

	return nil
}

func reveal(store *storage.Storage, tx *transactions.Transaction) error {
	commit := tx.Args[0].Value.([]byte)
	nebula := tx.Args[1].Value.([]byte)
	height := tx.Args[2].Value.(int64)
	reveal := tx.Args[3].Value.([]byte)
	pubKey := tx.Args[4].Value.([]byte)

	_, err := store.Reveal(nebula, height, commit)
	if err == storage.ErrKeyNotFound {
		var commitBytes []byte
		_, err := store.CommitHash(nebula, height, pubKey)

		if err == storage.ErrKeyNotFound {
			return ErrCommitIsNotExist
		} else if err != nil {
			return err
		}

		expectedHash := crypto.Keccak256(reveal)
		if !bytes.Equal(commitBytes, expectedHash[:]) {
			return ErrInvalidReveal
		}

		return store.SetReveal(nebula, height, commit, reveal)
	} else if err != nil {
		return err
	} else {
		return ErrRevealIsExist
	}
}

func addOracleInNebula(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaAddress := tx.Args[0].Value.([]byte)
	pubKey := tx.Args[1].Value.([]byte)

	oraclesByNebula, err := store.OraclesByNebula(nebulaAddress)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}

	var oraclePubKey account.OraclesPubKey
	copy(oraclePubKey[:], pubKey)

	if _, ok := oraclesByNebula[oraclePubKey]; ok {
		return ErrAddOracleInNebula
	}

	score, err := store.Score(tx.SenderPubKey)
	if err != nil {
		return err
	}

	if score < 0 {
		return ErrInvalidScore
	}

	oraclesByNebula[oraclePubKey] = true

	err = store.SetOraclesByNebula(nebulaAddress, oraclesByNebula)
	if err != nil {
		return err
	}

	return nil
}

func addOracle(store *storage.Storage, tx *transactions.Transaction) error {
	chainType := account.ChainType(tx.Args[0].Value.(byte))
	pubKey := tx.Args[1].Value.([]byte)

	oracles, err := store.OraclesByValidator(tx.SenderPubKey)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}
	if err == storage.ErrKeyNotFound {
		oracles = make(storage.OraclesByTypeMap)
	}

	var oraclePubKey account.OraclesPubKey
	copy(oraclePubKey[:], pubKey)

	oracles[chainType] = oraclePubKey

	err = store.SetOraclesByValidator(tx.SenderPubKey, oracles)
	if err != nil {
		return err
	}

	return nil
}

func result(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaAddress := tx.Args[0].Value.([]byte)
	height := tx.Args[1].Value.(int64)
	signBytes := tx.Args[2].Value.([]byte)
	chainType := account.ChainType(tx.Args[3].Value.(byte))

	oracles, err := store.OraclesByValidator(tx.SenderPubKey)
	if err != nil {
		return err
	}

	return store.SetResult(nebulaAddress, height, oracles[chainType], signBytes)
}

func newRound(store *storage.Storage, tx *transactions.Transaction, ledgerHeight int64, ethClient *ethclient.Client, wavesClient *client.Client, ctx context.Context) error {
	chainType := account.ChainType(tx.Args[0].Value.(byte))
	tcHeight := tx.Args[1].Value.(int64)

	_, err := store.RoundHeight(chainType, ledgerHeight)
	if err != storage.ErrKeyNotFound {
		return ErrNewRound
	}

	var height uint64
	switch chainType {
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
		return ErrInvalidChainType
	}

	if uint64(tcHeight) != height {
		return ErrInvalidHeight
	}

	return store.SetNewRound(chainType, ledgerHeight, uint64(tcHeight))
}

func vote(store *storage.Storage, tx *transactions.Transaction) error {
	votesBytes := tx.Args[0].Value.([]byte)

	var votes []storage.Vote
	err := json.Unmarshal(votesBytes, &votes)
	if err != nil {
		return err
	}

	return store.SetVote(tx.SenderPubKey, votes)
}

func isValidSigns(tx *transactions.Transaction) bool {
	return ed25519.Verify(tx.SenderPubKey[:], tx.Id.Bytes(), tx.Signature[:])
}

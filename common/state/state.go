package state

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"

	"github.com/ethereum/go-ethereum/crypto"
)

type SubRound int64

const (
	SubRoundCount = 4

	CommitSubRound SubRound = iota
	RevealSubRound
	ResultSubRound
	SendToTargetChain
)

var (
	ErrInvalidSign        = errors.New("invalid signature")
	ErrFuncNotFound       = errors.New("function is not found")
	ErrRevealIsExist      = errors.New("reveal is exist")
	ErrCommitIsExist      = errors.New("commit is exist")
	ErrCommitIsNotExist   = errors.New("commit is not exist")
	ErrInvalidReveal      = errors.New("invalid reveal")
	ErrNewRound           = errors.New("round is exist")
	ErrAddOracleInNebula  = errors.New("oracle was added in nebula")
	ErrInvalidScore       = errors.New("invalid score. score <= 0")
	ErrInvalidChainType   = errors.New("invalid chain type")
	ErrInvalidHeight      = errors.New("invalid height")
	ErrInvalidSubRound    = errors.New("invalid sub round")
	ErrInvalidNebulaOwner = errors.New("invalid nebula owner")
	ErrNebulaNotFound     = errors.New("nebula not found")
)

func CalculateSubRound(id uint64) SubRound {
	return SubRound(id % SubRoundCount)
}

func SetState(tx *transactions.Transaction, store *storage.Storage, ethClient *ethclient.Client, wavesClient *client.Client, ctx context.Context) error {
	if err := isValidSigns(store, tx); err != nil {
		return err
	}

	height, err := store.LastHeight()
	if err != nil {
		return err
	}

	switch tx.Func {
	case transactions.Commit:
		if CalculateSubRound(height) != CommitSubRound {
			return ErrInvalidSubRound
		}
		return commit(store, tx)
	case transactions.Reveal:
		if CalculateSubRound(height) != RevealSubRound {
			return ErrInvalidSubRound
		}
		return reveal(store, tx)
	case transactions.Result:
		if CalculateSubRound(height) != ResultSubRound {
			return ErrInvalidSubRound
		}
		return result(store, tx)
	case transactions.AddOracleInNebula:
		return addOracleInNebula(store, tx)
	case transactions.AddOracle:
		return addOracle(store, tx)
	case transactions.NewRound:
		return newRound(store, tx, height, ethClient, wavesClient, ctx)
	case transactions.Vote:
		return vote(store, tx)
	case transactions.SetNebula:
		return setNebula(store, tx)
	default:
		return ErrFuncNotFound
	}
}

func commit(store *storage.Storage, tx *transactions.Transaction) error {
	nebula := account.BytesToNebulaId(tx.Args[0].Value.([]byte))
	pulseId := tx.Args[1].Value.(int64)
	commit := tx.Args[2].Value.([]byte)
	pubKey := tx.Args[3].Value.(account.OraclesPubKey)

	_, err := store.CommitHash(nebula, pulseId, pubKey)
	if err == storage.ErrKeyNotFound {
		err := store.SetCommitHash(nebula, pulseId, pubKey, commit)
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
	nebula := account.BytesToNebulaId(tx.Args[1].Value.([]byte))
	pulseId := tx.Args[2].Value.(int64)
	reveal := tx.Args[3].Value.([]byte)
	pubKey := tx.Args[4].Value.(account.OraclesPubKey)

	_, err := store.Reveal(nebula, pulseId, commit)
	if err == storage.ErrKeyNotFound {
		var commitBytes []byte
		_, err := store.CommitHash(nebula, pulseId, pubKey)

		if err == storage.ErrKeyNotFound {
			return ErrCommitIsNotExist
		} else if err != nil {
			return err
		}

		expectedHash := crypto.Keccak256(reveal)
		if !bytes.Equal(commitBytes, expectedHash[:]) {
			return ErrInvalidReveal
		}

		return store.SetReveal(nebula, pulseId, commit, reveal)
	} else if err != nil {
		return err
	} else {
		return ErrRevealIsExist
	}
}

func addOracleInNebula(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaAddress := account.BytesToNebulaId(tx.Args[0].Value.([]byte))
	pubKey := tx.Args[1].Value.([]byte)

	nebulae, err := store.Nebulae()
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}

	nebula, ok := nebulae[nebulaAddress]
	if !ok {
		return ErrNebulaNotFound
	}

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

	if score < nebula.MinScore {
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

	oracles, err := store.OraclesByConsul(tx.SenderPubKey)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}
	if err == storage.ErrKeyNotFound {
		oracles = make(storage.OraclesByTypeMap)
	}

	var oraclePubKey account.OraclesPubKey
	copy(oraclePubKey[:], pubKey)

	oracles[chainType] = oraclePubKey

	err = store.SetOraclesByConsul(tx.SenderPubKey, oracles)
	if err != nil {
		return err
	}

	return nil
}

func result(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaAddress := account.BytesToNebulaId(tx.Args[0].Value.([]byte))
	pulseId := tx.Args[1].Value.(int64)
	signBytes := tx.Args[2].Value.([]byte)
	chainType := account.ChainType(tx.Args[3].Value.(byte))

	oracles, err := store.OraclesByConsul(tx.SenderPubKey)
	if err != nil {
		return err
	}

	return store.SetResult(nebulaAddress, pulseId, oracles[chainType], signBytes)
}

func newRound(store *storage.Storage, tx *transactions.Transaction, ledgerHeight uint64, ethClient *ethclient.Client, wavesClient *client.Client, ctx context.Context) error {
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

func setNebula(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaId := account.BytesToNebulaId(tx.Args[0].Value.([]byte))
	nebulaInfoBytes := tx.Args[1].Value.([]byte)

	nebulae, err := store.Nebulae()
	if err != nil {
		return err
	}

	if v, ok := nebulae[nebulaId]; ok {
		if v.Owner != tx.SenderPubKey {
			return ErrInvalidNebulaOwner
		}
	}

	var nebulaInfo storage.NebulaInfo
	err = json.Unmarshal(nebulaInfoBytes, &nebulaInfo)
	if err != nil {
		return err
	}

	nebulae[nebulaId] = nebulaInfo

	return store.SetNebula(nebulaId, nebulaInfo)
}

func isValidSigns(store *storage.Storage, tx *transactions.Transaction) error {
	score, err := store.Score(tx.SenderPubKey)
	if err != nil || score < 0 {
		return ErrInvalidScore
	}

	if ed25519.Verify(tx.SenderPubKey[:], tx.Id.Bytes(), tx.Signature[:]) {
		return ErrInvalidSign
	}
	return nil
}

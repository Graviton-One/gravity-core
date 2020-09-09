package state

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"

	"github.com/Gravity-Tech/gravity-core/common/adaptors"

	"github.com/Gravity-Tech/gravity-core/ledger/scheduler"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
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
	ErrInvalidSign         = errors.New("invalid signature")
	ErrFuncNotFound        = errors.New("function is not found")
	ErrRevealIsExist       = errors.New("reveal is exist")
	ErrCommitIsExist       = errors.New("commit is exist")
	ErrCommitIsNotExist    = errors.New("commit is not exist")
	ErrInvalidReveal       = errors.New("invalid reveal")
	ErrNewRound            = errors.New("round is exist")
	ErrAddOracleInNebula   = errors.New("oracle was added in nebula")
	ErrInvalidScore        = errors.New("invalid score. score <= 0")
	ErrInvalidChainType    = errors.New("invalid chain type")
	ErrInvalidHeight       = errors.New("invalid height")
	ErrInvalidSubRound     = errors.New("invalid sub round")
	ErrInvalidNebulaOwner  = errors.New("invalid nebula owner")
	ErrNebulaNotFound      = errors.New("nebula not found")
	ErrSignIsExistNotFound = errors.New("sign is exist")
	ErrRoundIsExist        = errors.New("round is exist")
)

func CalculateSubRound(id uint64) SubRound {
	return SubRound(id % SubRoundCount)
}

func SetState(tx *transactions.Transaction, store *storage.Storage, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, ctx context.Context) error {
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
		return newRound(store, tx, height, adaptors, ctx)
	case transactions.Vote:
		return vote(store, tx)
	case transactions.SetNebula:
		return setNebula(store, tx)
	case transactions.SignNewConsuls:
		return signNewConsuls(store, tx)
	case transactions.SignNewOracles:
		return signNewOracles(store, tx)
	case transactions.ApproveLastRound:
		return approveLastRound(store, adaptors, height, ctx)
	default:
		return ErrFuncNotFound
	}
}

func commit(store *storage.Storage, tx *transactions.Transaction) error {
	nebula := account.BytesToNebulaId(tx.Value(0).([]byte))
	pulseId := tx.Value(1).(int64)
	commit := tx.Value(2).([]byte)
	pubKeyBytes := tx.Value(3).([]byte)
	var pubKey account.OraclesPubKey
	copy(pubKey[:], pubKeyBytes)

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
	commit := tx.Value(0).([]byte)
	nebula := account.BytesToNebulaId(tx.Value(1).([]byte))
	pulseId := tx.Value(2).(int64)
	reveal := tx.Value(3).([]byte)
	pubKeyBytes := tx.Value(4).([]byte)
	var pubKey account.OraclesPubKey
	copy(pubKey[:], pubKeyBytes)

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
	nebulaAddress := account.BytesToNebulaId(tx.Value(0).([]byte))
	pubKeyBytes := tx.Value(1).([]byte)
	var pubKey account.OraclesPubKey
	copy(pubKey[:], pubKeyBytes)

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

	if _, ok := oraclesByNebula[pubKey]; ok {
		return ErrAddOracleInNebula
	}

	score, err := store.Score(tx.SenderPubKey)
	if err != nil {
		return err
	}

	if score < nebula.MinScore {
		return ErrInvalidScore
	}

	oraclesByNebula[pubKey] = true

	err = store.SetOraclesByNebula(nebulaAddress, oraclesByNebula)
	if err != nil {
		return err
	}

	return nil
}

func result(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaAddress := account.BytesToNebulaId(tx.Value(0).([]byte))
	pulseId := tx.Value(1).(int64)
	signBytes := tx.Value(2).([]byte)
	chainType := account.ChainType(tx.Value(3).([]byte)[0])

	oracles, err := store.OraclesByConsul(tx.SenderPubKey)
	if err != nil {
		return err
	}

	return store.SetResult(nebulaAddress, pulseId, oracles[chainType], signBytes)
}

func newRound(store *storage.Storage, tx *transactions.Transaction, ledgerHeight uint64, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, ctx context.Context) error {
	chainType := account.ChainType(tx.Value(0).([]byte)[0])
	tcHeight := tx.Value(1).(int64)

	_, err := store.RoundHeight(chainType, ledgerHeight)
	if err != storage.ErrKeyNotFound {
		return ErrNewRound
	}

	height, err := adaptors[chainType].GetHeight(ctx)
	if err != nil {
		return err
	}

	if uint64(tcHeight) != height {
		return ErrInvalidHeight
	}

	return store.SetNewRound(chainType, ledgerHeight, uint64(tcHeight))
}

func vote(store *storage.Storage, tx *transactions.Transaction) error {
	votesBytes := tx.Value(0).([]byte)

	var votes []storage.Vote
	err := json.Unmarshal(votesBytes, &votes)
	if err != nil {
		return err
	}

	return store.SetVote(tx.SenderPubKey, votes)
}

func setNebula(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaId := account.BytesToNebulaId(tx.Value(0).([]byte))
	nebulaInfoBytes := tx.Value(0).([]byte)

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
func addOracle(store *storage.Storage, tx *transactions.Transaction) error {
	chainType := account.ChainType(tx.Value(0).([]byte)[0])
	pubKey := tx.Value(1).([]byte)

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
func signNewConsuls(store *storage.Storage, tx *transactions.Transaction) error {
	chainType := account.ChainType(tx.Value(0).([]byte)[0])
	roundId := tx.Value(1).(int64)
	sign := tx.Value(2).([]byte)

	_, err := store.SignConsulsByConsul(tx.SenderPubKey, chainType, roundId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	} else if err == nil {
		return ErrSignIsExistNotFound
	}

	err = store.SetSignConsuls(tx.SenderPubKey, chainType, roundId, sign)
	if err != nil {
		return err
	}

	return nil
}
func signNewOracles(store *storage.Storage, tx *transactions.Transaction) error {
	roundId := tx.Value(0).(int64)
	sign := tx.Value(1).([]byte)
	nebulaAddress := account.BytesToNebulaId(tx.Value(2).([]byte))

	_, err := store.SignOraclesByConsul(tx.SenderPubKey, nebulaAddress, roundId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	} else if err == nil {
		return ErrSignIsExistNotFound
	}

	err = store.SetSignOracles(tx.SenderPubKey, nebulaAddress, roundId, sign)
	if err != nil {
		return err
	}

	return nil
}
func approveLastRound(store *storage.Storage, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, height uint64, ctx context.Context) error {
	roundId := height / scheduler.CalculateScoreInterval

	lastRound, err := store.LastRoundApproved()
	if err != nil {
		return err
	}

	if lastRound >= roundId {
		return ErrRoundIsExist
	}
	isExist := true
	for _, v := range adaptors {
		isExist, err = v.RoundExist(int64(roundId), ctx)
		if err != nil {
			return err
		}

		if isExist {
			lastRoundContract, err := v.LastRound(ctx)
			if err != nil {
				return err
			}
			isExist = lastRoundContract < roundId
		}
	}

	err = store.SetLastRoundApproved(roundId)
	if err != nil {
		return err
	}
	return nil
}

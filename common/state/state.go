package state

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"errors"

	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/common/hashing"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/ledger/scheduler"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

type SubRound int64

const (
	CommitSubRound SubRound = iota
	RevealSubRound
	ResultSubRound
	SendToTargetChain

	SubRoundCount = 4
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
	ErrSignIsExist        = errors.New("sign is exist")
	ErrRoundIsExist       = errors.New("round is exist")
)

func CalculateSubRound(tcHeight uint64, blocksInterval uint64) SubRound {
	return SubRound((tcHeight / (blocksInterval / SubRoundCount)) % SubRoundCount)
}

func SetState(tx *transactions.Transaction, store *storage.Storage, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, isSync bool, ctx context.Context) error {

	if err := isValidSigns(store, tx); err != nil {
		zap.L().Sugar().Error(err.Error())
		return err
	}

	height, err := store.LastHeight()
	if err != nil {
		zap.L().Sugar().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debugf("SetState func[%s]", tx.Func)
	//scheduler.PublishMessage("example.topic", []byte(fmt.Sprintf("SetState func[%s]", tx.Func)))
	switch tx.Func {
	case transactions.Commit:
		return persistCommit(store, tx)
	case transactions.Reveal:
		return persistReveal(store, tx)
	case transactions.Result:
		return persistResult(store, tx)
	case transactions.AddOracleInNebula:
		return addOracleInNebula(store, tx)
	case transactions.AddOracle:
		return addOracle(store, tx)
	case transactions.NewRound:
		return persistNewRound(store, tx, height, adaptors, ctx)
	case transactions.Vote:
		return vote(store, tx)
	case transactions.AddNebula:
		return setNebula(store, tx)
	case transactions.DropNebula:
		return dropNebula(store, tx)
	case transactions.SignNewConsuls:
		return signNewConsuls(store, tx)
	case transactions.SignNewOracles:
		return signNewOracles(store, tx)
	case transactions.ApproveLastRound:
		return approveLastRound(store, adaptors, height, isSync, ctx)
	case transactions.SetSolanaRecentBlock:
		return setSolanaRecentBlock(store, tx)
	case transactions.SetNebulaCustomParams:
		return setNebulaCustomParams(store, tx)
	case transactions.DropNebulaCustomParams:
		return dropNebulaCustomParams(store, tx)
	default:
		return ErrFuncNotFound
	}
}

func persistCommit(store *storage.Storage, tx *transactions.Transaction) error {
	nebula := account.BytesToNebulaId(tx.Value(0).([]byte))
	pulseId := tx.Value(1).(int64)
	tcHeight := tx.Value(2).(int64)
	commit := tx.Value(3).([]byte)
	pubKeyBytes := tx.Value(4).([]byte)
	var pubKey account.OraclesPubKey
	copy(pubKey[:], pubKeyBytes)

	_, err := store.CommitHash(nebula, tcHeight, pulseId, pubKey)
	if err == storage.ErrKeyNotFound {
		err := store.SetCommitHash(nebula, tcHeight, pulseId, pubKey, commit)
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

func persistReveal(store *storage.Storage, tx *transactions.Transaction) error {
	commit := tx.Value(0).([]byte)
	nebula := account.BytesToNebulaId(tx.Value(1).([]byte))
	pulseId := tx.Value(2).(int64)
	height := tx.Value(3).(int64)
	reveal := tx.Value(4).([]byte)
	pubKeyBytes := tx.Value(5).([]byte)
	chainType := tx.Value(6).(int64)
	var pubKey account.OraclesPubKey
	copy(pubKey[:], pubKeyBytes)
	zap.L().Sugar().Debug("State reveal", commit, nebula, pulseId, height, reveal, pubKeyBytes)

	_, err := store.Reveal(nebula, height, pulseId, commit, pubKey)
	if err == storage.ErrKeyNotFound {
		commitBytes, err := store.CommitHash(nebula, height, pulseId, pubKey)

		if err == storage.ErrKeyNotFound {
			return ErrCommitIsNotExist
		} else if err != nil {
			return err
		}

		expectedHash := hashing.WrappedKeccak256(reveal[:], account.ChainType(chainType))
		if !bytes.Equal(commitBytes, expectedHash[:]) {
			return ErrInvalidReveal
		}

		return store.SetReveal(nebula, height, pulseId, commit, pubKey, reveal)
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

	nebula, err := store.NebulaInfo(nebulaAddress)
	if err != nil {
		return err
	}

	oraclesByNebula, err := store.OraclesByNebula(nebulaAddress)
	if err == storage.ErrKeyNotFound {
		oraclesByNebula = make(storage.OraclesMap)
	} else if err != nil {
		return err
	}

	zap.L().Sugar().Debug("ORACLES BY NEBULA", oraclesByNebula)

	if _, ok := oraclesByNebula[pubKey.ToString(nebula.ChainType)]; ok {
		return ErrAddOracleInNebula
	}

	score, err := store.Score(tx.SenderPubKey)
	if err != nil {
		return err
	}

	if score < nebula.MinScore {
		return ErrInvalidScore
	}

	oraclesByNebula[pubKey.ToString(nebula.ChainType)] = nebula.ChainType

	err = store.SetOraclesByNebula(nebulaAddress, oraclesByNebula)
	if err != nil {
		return err
	}

	return nil
}

func persistResult(store *storage.Storage, tx *transactions.Transaction) error {
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

func persistNewRound(store *storage.Storage, tx *transactions.Transaction, ledgerHeight uint64, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, ctx context.Context) error {
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

func dropNebula(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaId := account.BytesToNebulaId(tx.Value(0).([]byte))
	return store.DropNebula(nebulaId)
}

func setSolanaRecentBlock(store *storage.Storage, tx *transactions.Transaction) error {
	round := tx.Value(0).(int)
	blockHash := tx.Value(1).([]byte)
	return store.SetSolanaRecentBlock(round, blockHash)
}

func setNebula(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaId := account.BytesToNebulaId(tx.Value(0).([]byte))
	nebulaInfoBytes := tx.Value(1).([]byte)

	nebula, err := store.NebulaInfo(nebulaId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}

	if err == nil && nebula.Owner != tx.SenderPubKey {
		return ErrInvalidNebulaOwner
	}

	var nebulaInfo storage.NebulaInfo
	err = json.Unmarshal(nebulaInfoBytes, &nebulaInfo)
	if err != nil {
		return err
	}

	return store.SetNebula(nebulaId, nebulaInfo)
}

func isValidSigns(store *storage.Storage, tx *transactions.Transaction) error {
	score, err := store.Score(tx.SenderPubKey)
	if err != nil || score < 0 {
		if err != nil {
			zap.L().Error(err.Error())
		}
		return ErrInvalidScore
	}

	if ed25519.Verify(tx.SenderPubKey[:], tx.Id.Bytes(), tx.Signature[:]) {
		return ErrInvalidSign
	}
	return nil
}
func addOracle(store *storage.Storage, tx *transactions.Transaction) error {
	zap.L().Debug("adding oracle")
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

	// _, err := store.SignConsulsByConsul(tx.SenderPubKey, chainType, roundId)
	// if err != nil && err != storage.ErrKeyNotFound {
	// 	return err
	// } else if err == nil {
	// 	return ErrSignIsExist
	// }

	err := store.SetSignConsuls(tx.SenderPubKey, chainType, roundId, sign)
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
		//return ErrSignIsExist
	}

	err = store.SetSignOracles(tx.SenderPubKey, nebulaAddress, roundId, sign)
	if err != nil {
		return err
	}

	return nil
}
func approveLastRound(store *storage.Storage, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, height uint64, isSync bool, ctx context.Context) error {
	roundId := uint64(scheduler.CalculateRound(int64(height)))

	lastRound, err := store.LastRoundApproved()
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}

	if lastRound >= roundId {
		return ErrRoundIsExist
	}

	if !isSync {
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
	}

	err = store.SetLastRoundApproved(roundId)
	if err != nil {
		return err
	}
	return nil
}

func setNebulaCustomParams(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaId := account.BytesToNebulaId(tx.Value(0).([]byte))
	nebulaCustomParamsBytes := tx.Value(1).([]byte)

	nebula, err := store.NebulaInfo(nebulaId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}

	if err == nil && nebula.Owner != tx.SenderPubKey {
		return ErrInvalidNebulaOwner
	}

	var nebulaCustomParams storage.NebulaCustomParams
	err = json.Unmarshal(nebulaCustomParamsBytes, &nebulaCustomParams)
	if err != nil {
		return err
	}

	return store.SetNebulaCustomParams(nebulaId, nebulaCustomParams)
}

func dropNebulaCustomParams(store *storage.Storage, tx *transactions.Transaction) error {
	nebulaId := account.BytesToNebulaId(tx.Value(0).([]byte))
	return store.DropNebulaCustomParams(nebulaId)
}

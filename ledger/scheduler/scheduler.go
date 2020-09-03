package scheduler

import (
	"context"
	"fmt"
	"sort"

	"github.com/Gravity-Tech/gravity-core/common/adaptors"

	"github.com/Gravity-Tech/gravity-core/common/account"
	calculator "github.com/Gravity-Tech/gravity-core/common/score"
	"github.com/Gravity-Tech/gravity-core/common/storage"

	"github.com/dgraph-io/badger"
)

const (
	CalculateScoreInterval = 20
	ValidatorCount         = 5
	OracleCount            = 5
)

type Scheduler struct {
	Adaptors map[account.ChainType]*adaptors.IBlockchainAdaptor
	Ledger   *account.LedgerValidator
	ctx      context.Context
}

func New(adaptors map[account.ChainType]*adaptors.IBlockchainAdaptor, ledger *account.LedgerValidator, ctx context.Context) (*Scheduler, error) {
	return &Scheduler{
		Ledger:   ledger,
		Adaptors: adaptors,
		ctx:      ctx,
	}, nil
}

func (scheduler *Scheduler) HandleBlock(height int64, store *storage.Storage) error {
	lastHeight, err := store.LastHeight()
	if err != nil {
		return err
	}

	if height%CalculateScoreInterval == 0 {
		if err := scheduler.calculateScores(store); err != nil {
			return err
		}
	}

	if uint64(height) < lastHeight {
		return nil
	}

	roundId := height / CalculateScoreInterval
	if height%CalculateScoreInterval < CalculateScoreInterval/2 {
		for k, v := range scheduler.Adaptors {
			err := scheduler.signConsulsResult(roundId, k, store)
			if err != nil {
				return err
			}

			for _, v := range v.Nebulae {
				err := scheduler.signOracleResultByNebula(roundId, v, k, store)
				if err != nil {
					continue
				}
			}
		}
	} else if height%CalculateScoreInterval > CalculateScoreInterval/2 {
		for k, v := range scheduler.Adaptors {
			err := scheduler.sendConsulsToGravityContract(roundId, k, store)
			if err != nil {
				return err
			}

			for _, v := range v.Nebulae {
				err := scheduler.sendOraclesToNebula(v, k, roundId, store)
				if err != nil {
					continue
				}
			}
		}
	}

	return nil
}

func (scheduler *Scheduler) signConsulsResult(roundId int64, chainType account.ChainType, store *storage.Storage) error {
	_, err := store.SignConsulsResultByConsul(scheduler.Ledger.PubKey, chainType, roundId)
	if err == badger.ErrKeyNotFound {
		return nil
	} else if err != nil {
		return err
	}

	scores, err := store.Scores()
	if err != nil {
		return err
	}

	var sortedScores []storage.Consul
	for k, v := range scores {
		sortedScores = append(sortedScores, storage.Consul{
			PubKey: k,
			Value:  v,
		})
	}

	sort.SliceStable(sortedScores, func(i, j int) bool {
		return sortedScores[i].Value > sortedScores[j].Value
	})

	var newConsuls []storage.Consul
	for _, v := range sortedScores {
		newConsuls = append(newConsuls, v)
		if len(newConsuls) >= ValidatorCount {
			break
		}
	}

	var consulsAddresses []account.OraclesPubKey
	for _, v := range newConsuls {
		oraclesByConsul, err := store.OraclesByConsul(v.PubKey)
		if err != nil {
			return err
		}

		consulsAddresses = append(consulsAddresses, oraclesByConsul[chainType])
	}

	sign, err := scheduler.Adaptors[chainType].SignConsuls(consulsAddresses)
	if err != nil {
		return err
	}

	err = store.SetConsuls(newConsuls)
	if err != nil {
		return err
	}

	err = store.SetSignConsulsResult(scheduler.Ledger.PubKey, chainType, roundId, sign)
	if err != nil {
		return err
	}

	return nil
}
func (scheduler *Scheduler) signOracleResultByNebula(roundId int64, nebulaId account.NebulaId, chainType account.ChainType, store *storage.Storage) error {
	_, err := store.SignOraclesResultByConsul(scheduler.Ledger.PubKey, nebulaId, roundId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	} else if err == nil {
		return nil
	}

	oraclesByNebula, err := store.OraclesByNebula(nebulaId)
	if err != nil {
		return err
	}

	lastIndex, err := store.NebulaOraclesIndex()
	if err != nil {
		return err
	}

	var newOracles []account.OraclesPubKey
	var oracles []account.OraclesPubKey
	newOraclesMap := make(storage.OraclesMap)

	for k, _ := range oraclesByNebula {
		oracles = append(oracles, k)
	}

	newIndex := lastIndex + 1
	if newIndex >= uint64(len(oracles)) {
		newIndex = 0
	}

	if newIndex+OracleCount > uint64(len(oracles)) {
		newOracles = oracles[newIndex:]
		newOracles = append(newOracles, newOracles[:OracleCount-len(newOracles)]...)
	} else {
		newOracles = oracles[newIndex : newIndex+OracleCount]
	}

	for _, v := range newOracles {
		newOraclesMap[v] = true
	}

	sign, err := scheduler.Adaptors[chainType].SignOracles(nebulaId, newOracles)
	if err != nil {
		return err
	}

	err = store.SetSignOraclesResult(scheduler.Ledger.PubKey, nebulaId, roundId, sign)
	if err != nil {
		return err
	}

	err = store.SetBftOraclesByNebula(nebulaId, newOraclesMap)
	if err != nil {
		return err
	}

	return nil
}

func (scheduler *Scheduler) sendConsulsToGravityContract(round int64, chainType account.ChainType, store *storage.Storage) error {
	prevConsuls, err := store.PrevConsuls()
	if err != nil {
		return err
	}

	var signs [][]byte
	for _, v := range prevConsuls {
		sign, err := store.SignConsulsResultByConsul(v.PubKey, chainType, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			var empty [32]byte
			signs = append(signs, empty[:])
			continue
		}

		signs = append(signs, sign)
	}

	newConsuls, err := store.Consuls()
	if err != nil {
		return err
	}

	var newConsulsAddresses []account.OraclesPubKey
	for _, v := range newConsuls {
		oracles, err := store.OraclesByConsul(v.PubKey)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			newConsulsAddresses = append(newConsulsAddresses, account.OraclesPubKey{})
			continue
		}

		pubKey := oracles[chainType]
		newConsulsAddresses = append(newConsulsAddresses, pubKey)
	}

	id, err := scheduler.Adaptors[chainType].SendConsulsToGravityContract(newConsulsAddresses, signs, round, scheduler.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx consuls update (%d): %s \n", chainType, id)
	return nil
}
func (scheduler *Scheduler) sendOraclesToNebula(nebulaId account.NebulaId, chainType account.ChainType, round int64, store *storage.Storage) error {
	prevConsuls, err := store.PrevConsuls()
	if err != nil {
		return err
	}

	var signs [][]byte
	for _, v := range prevConsuls {
		sign, err := store.SignOraclesResultByConsul(v.PubKey, nebulaId, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			var empty [32]byte
			signs = append(signs, empty[:])
			continue
		}

		signs = append(signs, sign)
	}

	oracles, err := store.OraclesByNebula(nebulaId)
	if err != nil {
		return err
	}

	var oraclesAddresses []account.OraclesPubKey
	for k, _ := range oracles {
		oraclesAddresses = append(oraclesAddresses, k)
	}

	tx, err := scheduler.Adaptors[chainType].SendOraclesToNebula(nebulaId, oraclesAddresses, signs, round, scheduler.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx nebula (%s) oracles update: %s \n", nebulaId, tx)
	return nil

}

func (scheduler *Scheduler) calculateScores(store *storage.Storage) error {
	voteMap, err := store.Votes()
	if err != nil {
		return err
	}

	scores, err := store.Scores()
	if err != nil {
		return err
	}

	newScores, err := calculator.Calculate(scores, voteMap)
	if err != nil {
		return err
	}

	for k, v := range newScores {
		err := store.SetScore(k, v)
		if err != nil {
			return err
		}

		if v <= 0 {
			oracles, err := store.OraclesByConsul(k)
			if err != nil {
				return err
			}

			for _, oracle := range oracles {
				nebulae, err := store.NebulaeByOracle(oracle)
				if err != nil {
					return err
				}

				for _, nebulaId := range nebulae {
					oracles, err := store.OraclesByNebula(nebulaId)
					if err != nil {
						return err
					}

					delete(oracles, oracle)
					err = store.SetOraclesByNebula(nebulaId, oracles)
					if err != nil {
						return err
					}
				}

				err = store.SetNebulaeByOracle(oracle, []account.NebulaId{})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/Gravity-Tech/gravity-core/common/account"
	calculator "github.com/Gravity-Tech/gravity-core/common/score"
	"github.com/Gravity-Tech/gravity-core/common/storage"
)

const (
	CalculateScoreInterval = 20
	OracleCount            = 5
)

type Scheduler struct {
	httpSchedulerHost string
	ctx               context.Context
}

type ConsulInfo struct {
	ConsulIndex int
	TotalCount  int
	IsConsul    bool
}

func New(httpSchedulerHost string, ctx context.Context) (*Scheduler, error) {
	return &Scheduler{
		httpSchedulerHost: "http://" + httpSchedulerHost,
		ctx:               ctx,
	}, nil
}

func (scheduler *Scheduler) HandleBlock(height int64, store *storage.Storage) error {
	err := scheduler.sendRqForProcessing(height)
	if err != nil {
		return err
	}

	roundId := height / CalculateScoreInterval

	if height%CalculateScoreInterval == 0 || height == 1 {
		if err := scheduler.calculateScores(store); err != nil {
			return err
		}

		if err := scheduler.updateConsulsAndCandidate(store, roundId-1); err != nil {
			return err
		}

		nebulae, err := store.Nebulae()
		if err != nil {
			return err
		}

		for k, _ := range nebulae {
			err = scheduler.updateOracles(k, store)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (scheduler *Scheduler) updateConsulsAndCandidate(store *storage.Storage, roundId int64) error {
	lastRound, err := store.LastRoundApproved()
	if err != nil {
		return err
	}

	if lastRound != uint64(roundId) {
		return nil
	}

	validatorCount, err := store.ConsulsCount()
	if err != nil {
		return err
	}

	newConsuls, err := store.ConsulsCandidate()
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}
	if err != storage.ErrKeyNotFound {
		err = store.SetConsuls(newConsuls)
		if err != nil {
			return err
		}
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
	var consulsCandidate []storage.Consul
	for _, v := range sortedScores {
		consulsCandidate = append(consulsCandidate, v)
		if len(consulsCandidate) >= validatorCount {
			break
		}
	}
	err = store.SetConsulsCandidate(consulsCandidate)
	if err != nil {
		return err
	}
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

	nebulaeInfo, err := store.Nebulae()
	if err != nil {
		return err
	}

	for k, v := range newScores {
		err := store.SetScore(k, v)
		if err != nil {
			return err
		}

		oracles, err := store.OraclesByConsul(k)
		if err != nil {
			return err
		}

		for _, oracle := range oracles {
			nebulae, err := store.NebulaeByOracle(oracle)
			if err != nil && err != storage.ErrKeyNotFound {
				return err
			}
			if err == storage.ErrKeyNotFound {
				break
			}

			var newNebulae []account.NebulaId
			for _, nebulaId := range nebulae {
				oracles, err := store.OraclesByNebula(nebulaId)
				if err != nil {
					return err
				}

				if v < nebulaeInfo[nebulaId].MinScore || v <= 0 {
					delete(oracles, oracle)
					err = store.SetOraclesByNebula(nebulaId, oracles)
					if err != nil {
						return err
					}
					continue
				}
				newNebulae = append(newNebulae, nebulaId)
			}

			err = store.SetNebulaeByOracle(oracle, newNebulae)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (scheduler *Scheduler) updateOracles(nebulaId account.NebulaId, store *storage.Storage) error {
	oraclesByNebula, err := store.OraclesByNebula(nebulaId)
	if err != nil {
		return err
	}

	lastIndex, err := store.NebulaOraclesIndex(nebulaId)
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

	err = store.SetBftOraclesByNebula(nebulaId, newOraclesMap)
	if err != nil {
		return err
	}

	return nil
}

func (scheduler *Scheduler) sendRqForProcessing(height int64) error {
	rqUrl := fmt.Sprintf("%v/%v?height=%d", scheduler.httpSchedulerHost, "/process", height)

	req, err := http.NewRequestWithContext(scheduler.ctx, "POST", rqUrl, nil)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}

package scheduler

import (
	"bytes"
	"context"
	"fmt"
	"sort"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/adaptors"

	"github.com/Gravity-Tech/gravity-core/common/account"
	calculator "github.com/Gravity-Tech/gravity-core/common/score"
	"github.com/Gravity-Tech/gravity-core/common/storage"
)

type NebulaToUpdate struct {
	Id        string
	ChainType account.ChainType
}
type ManualUpdateStruct struct {
	Active      bool
	UpdateQueue []NebulaToUpdate
}

var EventBus *gochannel.GoChannel
var GlobalScheduler Scheduler
var GlobalStorage *storage.Storage
var SchedulerEventServer *EventServer
var ManualUpdate ManualUpdateStruct

const (
	HardforkHeight = 95574

	StarValueForNewRound      = 1000
	CalculateScoreInterval    = 100
	NewCalculateScoreInterval = 9600
	OracleCount               = 5
)

type Scheduler struct {
	Adaptors map[account.ChainType]adaptors.IBlockchainAdaptor
	Ledger   *account.LedgerValidator
	ctx      context.Context
	client   *gravity.Client
}

type ConsulInfo struct {
	ConsulIndex int
	TotalCount  int
	IsConsul    bool
}

func (ma ManualUpdateStruct) Disable() {
	ma.Active = false
	ma.UpdateQueue = []NebulaToUpdate{}
}
func New(adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, ledger *account.LedgerValidator, localHost string, ctx context.Context) (*Scheduler, error) {
	ManualUpdate.Active = false
	ManualUpdate.UpdateQueue = []NebulaToUpdate{}
	EventBus = gochannel.NewGoChannel(
		gochannel.Config{},
		watermill.NewStdLogger(false, false),
	)

	client, err := gravity.New(localHost)
	if err != nil {
		return nil, err
	}

	messages, err := EventBus.Subscribe(ctx, "ledger.events")
	if err != nil {
		panic(err)
	}
	SchedulerEventServer = NewEventServer()
	go SchedulerEventServer.Serve(messages)
	GlobalScheduler = Scheduler{
		Ledger:   ledger,
		Adaptors: adaptors,
		ctx:      ctx,
		client:   client,
	}
	return &GlobalScheduler, nil
}

func CalculateRound(height int64) int64 {
	if height >= HardforkHeight {
		return height/NewCalculateScoreInterval + StarValueForNewRound
	}
	// exists only for backward compatibility
	if height >= 77852 {
		return height/21600 + StarValueForNewRound
	}

	return height / CalculateScoreInterval
}
func IsRoundStart(height int64) bool {
	if height >= HardforkHeight {
		return height%NewCalculateScoreInterval == 0
	}
	// exists only for backward compatibility
	if height >= 77852 {
		return height%21600 == 0
	}

	return height%CalculateScoreInterval == 0
}

func (scheduler *Scheduler) HandleBlock(height int64, store *storage.Storage, isSync bool, isConsul bool) error {
	if !isSync && isConsul {
		PublishMessage("ledger.events", SchedulerEvent{
			Name: "handle_block",
			Params: map[string]interface{}{
				"height": height,
			},
		})
		//go scheduler.process(height)
	}

	roundId := CalculateRound(height)

	if height%100 == 0 || height == 1 {
		if err := scheduler.calculateScores(store); err != nil {
			zap.L().Error(err.Error())
			return err
		}

		if err := scheduler.updateConsulsAndCandidate(store, roundId-1); err != nil {
			zap.L().Error(err.Error())
			return err
		}

		nebulae, err := store.Nebulae()
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}

		for k, v := range nebulae {
			zap.L().Sugar().Debugf("Iterate nebule: %s", k)
			nebulaId, err := account.StringToNebulaId(k, v.ChainType)
			if err != nil {
				fmt.Printf("Error:%s\n", err.Error())
				continue
			}
			err = scheduler.UpdateOracles(roundId, nebulaId, store)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (scheduler *Scheduler) updateConsulsAndCandidate(store *storage.Storage, roundId int64) error {
	lastRound, err := store.LastRoundApproved()
	if err != nil && err != storage.ErrKeyNotFound {
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
	if len(newConsuls) <= 0 {
		return nil
	}
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
		if sortedScores[i].Value == sortedScores[j].Value {
			return bytes.Compare(sortedScores[i].PubKey[:], sortedScores[j].PubKey[:]) == 1
		} else {
			return sortedScores[i].Value > sortedScores[j].Value
		}
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

	for k, v := range newScores {
		err := store.SetScore(k, v)
		if err != nil {
			return err
		}
	}

	return nil
}
func (scheduler *Scheduler) UpdateOracles(roundId int64, nebulaId account.NebulaId, store *storage.Storage) error {
	zap.L().Debug("updateOracles called")
	nebulaInfo, err := store.NebulaInfo(nebulaId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	oraclesByNebula, err := store.OraclesByNebula(nebulaId)
	if err == storage.ErrKeyNotFound {
		return nil
	} else if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	var newOracles []account.OraclesPubKey
	var oracles []account.OraclesPubKey
	newOraclesMap := make(storage.OraclesMap)

	for k, v := range oraclesByNebula {
		oracleAddress, err := account.StringToOraclePubKey(k, v)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		oracles = append(oracles, oracleAddress)
	}

	if len(oracles) <= OracleCount {
		newOracles = append(newOracles, oracles...)
	} else {
		newIndex := int(roundId) % (len(oracles) - 1)
		if newIndex+OracleCount > len(oracles) {
			newOracles = oracles[newIndex:]
			count := OracleCount - len(newOracles)
			newOracles = append(newOracles, oracles[:count]...)
		} else {
			newOracles = oracles[newIndex : newIndex+OracleCount]
		}
	}

	for _, v := range newOracles {
		newOraclesMap[v.ToString(nebulaInfo.ChainType)] = nebulaInfo.ChainType
	}

	err = store.SetBftOraclesByNebula(nebulaId, newOraclesMap)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

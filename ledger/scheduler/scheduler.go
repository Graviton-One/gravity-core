package scheduler

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/Gravity-Tech/gravity-core/common/account"
	ghClient "github.com/Gravity-Tech/gravity-core/common/client"
	calculator "github.com/Gravity-Tech/gravity-core/common/score"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/Gravity-Tech/gravity-core/ledger/blockchain"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/dgraph-io/badger"
)

const (
	CalculateScoreInterval = 20
	ValidatorCount         = 5
	OracleCount            = 5
)

type LedgerValidator struct {
	PrivKey ed25519.PrivKeyEd25519
	PubKey  account.ConsulPubKey
}
type Scheduler struct {
	Blockchains map[account.ChainType]blockchain.IBlockchain
	Ledger      *LedgerValidator
	GhNode      string
	ctx         context.Context
	nebulae     map[account.ChainType][][]byte
}

func New(blockchains map[account.ChainType]blockchain.IBlockchain, ghNode string, ctx context.Context, ledger *LedgerValidator, nebulae map[account.ChainType][][]byte) (*Scheduler, error) {

	return &Scheduler{
		Ledger:      ledger,
		Blockchains: blockchains,
		ctx:         ctx,
		GhNode:      strings.Replace(ghNode, "tcp", "http", 1),
		nebulae:     nebulae,
	}, nil
}

func (scheduler *Scheduler) HandleBlock(height int64, store *storage.Storage) error {
	go scheduler.setPrivKeys(store, account.Ethereum) //TODO: refactoring
	go scheduler.setPrivKeys(store, account.Waves)    //TODO: refactoring

	roundId := height / CalculateScoreInterval
	if height%CalculateScoreInterval == 0 {
		if err := scheduler.calculateScores(store); err != nil {
			return err
		}
	} else if height%CalculateScoreInterval < CalculateScoreInterval/2 {
		for k, _ := range scheduler.Blockchains {
			err := scheduler.signConsulsResult(roundId, k, store)
			if err != nil {
				return err
			}

			for _, v := range scheduler.nebulae[k] {
				err := scheduler.signOracleResultByNebula(roundId, v, account.Waves, store)
				if err != nil {
					continue
				}
			}
		}
	} else if height%CalculateScoreInterval > CalculateScoreInterval/2 {
		for k, _ := range scheduler.Blockchains {
			err := scheduler.sendConsulsToGravityContract(roundId, k, store)
			if err != nil {
				return err
			}

			for _, v := range scheduler.nebulae[k] {
				err := scheduler.sendOraclesToNebula(v, k, roundId, store)
				if err != nil {
					continue
				}
			}
		}
	}

	return nil
}

func (scheduler *Scheduler) setPrivKeys(storage *storage.Storage, chainType account.ChainType) error {
	oracles, err := storage.OraclesByConsul(scheduler.Ledger.PubKey)
	if err != nil {
		return err
	}

	if _, ok := oracles[chainType]; ok {
		return nil
	}

	args := []transactions.Args{
		{
			Value: chainType,
		},
		{
			Value: scheduler.Blockchains[chainType].PubKey(),
		},
	}

	tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.AddOracle, scheduler.Ledger.PrivKey, args)
	if err != nil {
		return err
	}

	ghClient, err := ghClient.New(scheduler.GhNode)
	if err != nil {
		return err
	}

	err = ghClient.SendTx(tx)
	if err != nil {
		return err
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

	sign, err := scheduler.Blockchains[chainType].SignConsuls(consulsAddresses)
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
func (scheduler *Scheduler) signOracleResultByNebula(roundId int64, nebulaId []byte, chainType account.ChainType, store *storage.Storage) error {
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

	sign, err := scheduler.Blockchains[chainType].SignOracles(nebulaId, newOracles)
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
		sign, err := store.SignConsulsResultByConsul(v.PubKey, account.Waves, round)
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

	id, err := scheduler.Blockchains[chainType].SendConsulsToGravityContract(newConsulsAddresses, signs, round, scheduler.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx consuls update (%d): %s \n", chainType, id)
	return nil
}
func (scheduler *Scheduler) sendOraclesToNebula(nebulaId []byte, chainType account.ChainType, round int64, store *storage.Storage) error {
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

	tx, err := scheduler.Blockchains[chainType].SendOraclesToNebula(nebulaId, oraclesAddresses, signs, round, scheduler.ctx)
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

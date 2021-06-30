package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

func (scheduler *Scheduler) Process(height int64) {
	zap.L().Debug("Called process func")
	err := scheduler.processByHeight(height)
	if err != nil {
		fmt.Printf("Error:%s\n", err)
	}
}
func (scheduler *Scheduler) processByHeight(height int64) error {

	roundId := CalculateRound(height)

	consulInfo, err := scheduler.consulInfo()
	if err != nil {
		return err
	}

	//Refresh targetchains pubkeys
	if height%20 == 0 {
		scheduler.updateTargetChainsPubKeys()
	}

	senderIndex := int64(CalculateRound(height)) % int64(consulInfo.TotalCount)

	zap.L().Sugar().Debugf("Sender index: %d", senderIndex)
	consuls, err := scheduler.client.Consuls()
	if err != nil {
		return err
	}
	senderConsul := consuls[senderIndex]
	oraclesBySenderConsul, err := scheduler.client.OraclesByValidator(senderConsul.PubKey)
	if err != nil {
		return err
	}
	if height%50 != 0 {
		return nil
	}
	go func() {
		time.Sleep(time.Second * 10)
		ManualUpdate.Disable()
	}()
	isExist := true
	var wg sync.WaitGroup
	for k, v := range scheduler.Adaptors {
		_ck, _cv := k, v
		wg.Add(1)
		go func(wg *sync.WaitGroup, k account.ChainType, v adaptors.IBlockchainAdaptor) {
			defer wg.Done()
			index := roundId % int64(consulInfo.TotalCount)

			lastRound, err := v.LastRound(scheduler.ctx)
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			zap.L().Sugar().Debugf("RoundId %d, lastRound [%s] - %d", roundId, k, lastRound)
			isExist = uint64(roundId) == lastRound

			nebulae, err := scheduler.client.Nebulae()
			if err != nil {
				zap.L().Error(err.Error())
				return
			}
			//var nebulaWG sync.WaitGroup

			for nk, val := range nebulae {
				if val.ChainType != k {
					continue
				}

				if ManualUpdate.Active {
					zap.L().Sugar().Debug("Check for manual update the Nebula ", ManualUpdate)
					found := false
					nindex := int(-1)
					for i, n := range ManualUpdate.UpdateQueue {
						if n.Id == nk && n.ChainType == val.ChainType {
							found = true
							nindex = i
							break
						}
					}
					if found {
						ManualUpdate.UpdateQueue = append(ManualUpdate.UpdateQueue[:nindex], ManualUpdate.UpdateQueue[nindex+1:]...)
						if len(ManualUpdate.UpdateQueue) == 0 {
							ManualUpdate.Active = false
						}
					} else {
						zap.L().Sugar().Debugf("Manual updated Nebula [%s] not found", nk)
						continue
					}
				} else {
					if uint64(roundId) <= lastRound {
						zap.L().Debug("roundId <= lastRound")
						return
					}
				}

				payload := map[string]interface{}{
					"nebula_key": nk,
					"round_id":   roundId,
					"sender":     oraclesBySenderConsul[k],
					"is_sender":  index == int64(consulInfo.ConsulIndex),
					"chain_type": val.ChainType,
				}
				PublishMessage("ledger.events", SchedulerEvent{
					Name:   "update_oracles",
					Params: payload,
				})
			}

			if uint64(roundId) <= lastRound {
				zap.L().Debug("roundId <= lastRound")
				return
			}

			var consulsWG sync.WaitGroup
			consulsWG.Add(1)
			go func() {
				defer consulsWG.Done()
				err = scheduler.signConsulsResult(roundId, _ck, oraclesBySenderConsul[_ck])
				if err != nil {
					zap.L().Error(err.Error())
				}
				time.Sleep(time.Second * 5)
				if index == int64(consulInfo.ConsulIndex) {
					err = scheduler.sendConsulsToGravityContract(roundId, _ck)
					if err != nil {
						zap.L().Error(err.Error())
					}
				}
			}()

			consulsWG.Wait()
		}(&wg, _ck, _cv)
	}
	wg.Wait()

	lastRound, err := scheduler.client.LastRoundApproved()
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	}
	if isExist && uint64(roundId) > lastRound && senderIndex == int64(consulInfo.ConsulIndex) {
		tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.ApproveLastRound, scheduler.Ledger.PrivKey)
		if err != nil {
			return err
		}
		err = scheduler.client.SendTx(tx)
		if err != nil {
			return err
		}
	}

	return nil
}
func (scheduler *Scheduler) consulInfo() (*ConsulInfo, error) {
	consuls, err := scheduler.client.Consuls()
	if err != nil {
		return nil, err
	}

	isConsul := false
	consulIndex := 0
	for i, consul := range consuls {
		if scheduler.Ledger.PubKey == consul.PubKey {
			isConsul = true
			consulIndex = i
			break
		}
	}

	return &ConsulInfo{
		ConsulIndex: consulIndex,
		TotalCount:  len(consuls),
		IsConsul:    isConsul,
	}, nil
}
func (scheduler *Scheduler) signConsulsResult(roundId int64, chainType account.ChainType, sender account.OraclesPubKey) error {
	_, err := scheduler.client.SignNewConsulsByConsul(scheduler.Ledger.PubKey, chainType, roundId)
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	}

	consuls, err := scheduler.client.ConsulsCandidate()
	if err != nil {
		return err
	}

	var consulsAddresses []*account.OraclesPubKey
	for i := 0; i < OracleCount; i++ {
		if i >= len(consuls) {
			consulsAddresses = append(consulsAddresses, nil)
			continue
		}
		v := consuls[i]
		oraclesByConsul, err := scheduler.client.OraclesByValidator(v.PubKey)
		if err == gravity.ErrValueNotFound {
			consulsAddresses = append(consulsAddresses, nil)
			continue
		} else if err != nil {
			return err
		}

		oracle := oraclesByConsul[chainType]
		consulsAddresses = append(consulsAddresses, &oracle)
	}
	sign, err := scheduler.Adaptors[chainType].SignConsuls(consulsAddresses, roundId, sender)
	if err != nil {
		return err
	}
	tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.SignNewConsuls, scheduler.Ledger.PrivKey)
	if err != nil {
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{
			Value: []byte{byte(chainType)},
		},
		transactions.IntValue{
			Value: roundId,
		},
		transactions.BytesValue{
			Value: sign,
		},
	})
	err = scheduler.client.SendTx(tx)
	if err != nil {
		return err
	}
	return nil
}
func (scheduler *Scheduler) signOraclesByNebula(roundId int64, nebulaId account.NebulaId, chainType account.ChainType, sender account.OraclesPubKey) error {
	// _, err := scheduler.client.SignNewOraclesByConsul(scheduler.Ledger.PubKey, chainType, nebulaId, roundId)
	// if err != nil && err != gravity.ErrValueNotFound {
	// 	zap.L().Error(err.Error())
	// 	return err
	// } else if err == nil {
	// 	zap.L().Debug("Returning from func signOraclesByNebula")
	// 	//return nil
	// }
	bftOraclesByNebula, err := scheduler.client.BftOraclesByNebula(chainType, nebulaId)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	var newOracles []*account.OraclesPubKey
	for k, v := range bftOraclesByNebula {
		oracleAddress, err := account.StringToOraclePubKey(k, v)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		newOracles = append(newOracles, &oracleAddress)
	}
	for i := len(newOracles); i < OracleCount; i++ {
		newOracles = append(newOracles, nil)
	}
	zap.L().Sugar().Debugf("[%s] Signing oracles", chainType)
	zap.L().Sugar().Debug("NebulaId: ", nebulaId, "Oracles: ", newOracles, "Round: ", roundId, "Sender: ", sender)
	sign, err := scheduler.Adaptors[chainType].SignOracles(nebulaId, newOracles, roundId, sender)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debugf("[%s] Oracles signed - %s", chainType, sign)
	tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.SignNewOracles, scheduler.Ledger.PrivKey)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.IntValue{
			Value: roundId,
		},
		transactions.BytesValue{
			Value: sign,
		},
		transactions.BytesValue{
			Value: nebulaId[:],
		},
	})
	err = scheduler.client.SendTx(tx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}
func (scheduler *Scheduler) sendConsulsToGravityContract(round int64, chainType account.ChainType) error {
	if scheduler.ctx == nil {
		zap.L().Debug("Context is nil")
	}
	exist, err := scheduler.Adaptors[chainType].RoundExist(round, scheduler.ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if exist {
		return nil
	}

	lastRound, err := scheduler.Adaptors[chainType].LastRound(scheduler.ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	if round <= int64(lastRound) {
		return nil
	}

	consuls, err := scheduler.client.Consuls()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	newConsuls, err := scheduler.client.ConsulsCandidate()
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	realSignCount := 0

	signs := make(map[account.OraclesPubKey][]byte)
	for i := 0; i < OracleCount; i++ {
		if i >= len(consuls) {
			break
		}
		v := consuls[i]

		oracles, err := scheduler.client.OraclesByValidator(v.PubKey)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		}

		oraclePubKey, ok := oracles[chainType]
		if !ok {
			continue
		}

		sign, err := scheduler.client.SignNewConsulsByConsul(v.PubKey, chainType, round)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		}
		if err == gravity.ErrValueNotFound {
			continue
		}

		signs[oraclePubKey] = sign
		realSignCount++
	}

	if realSignCount < len(consuls)*2/3 {
		return nil
	}

	var newConsulsAddresses []*account.OraclesPubKey
	for i := 0; i < OracleCount; i++ {
		if i >= len(newConsuls) {
			newConsulsAddresses = append(newConsulsAddresses, nil)
			continue
		}
		v := newConsuls[i]
		oraclesByConsul, err := scheduler.client.OraclesByValidator(v.PubKey)
		if err == gravity.ErrValueNotFound {
			newConsulsAddresses = append(newConsulsAddresses, nil)
			continue
		} else if err != nil {
			return err
		}

		oracle := oraclesByConsul[chainType]
		newConsulsAddresses = append(newConsulsAddresses, &oracle)
	}

	id, err := scheduler.Adaptors[chainType].SendConsulsToGravityContract(newConsulsAddresses, signs, round, scheduler.ctx)
	if err != nil {
		zap.L().Sugar().Errorf("Chain [%s] err: %s", chainType, err.Error())
		return err
	}
	if id != "" {
		err := scheduler.Adaptors[chainType].WaitTx(id, scheduler.ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Tx consuls update (%s): %s\n", chainType.String(), id)
	}
	return nil
}
func (scheduler *Scheduler) sendOraclesToNebula(nebulaId account.NebulaId, chainType account.ChainType, round int64) error {
	consuls, err := scheduler.client.Consuls()
	if err != nil {
		return err
	}

	realSignCount := 0
	signs := make(map[account.OraclesPubKey][]byte)
	for i := 0; i < OracleCount; i++ {
		if i >= len(consuls) {
			break
		}
		v := consuls[i]

		oracles, err := scheduler.client.OraclesByValidator(v.PubKey)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		}

		oraclePubKey, ok := oracles[chainType]
		if !ok {
			continue
		}

		sign, err := scheduler.client.SignNewOraclesByConsul(v.PubKey, chainType, nebulaId, round)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		}
		if err == gravity.ErrValueNotFound {
			continue
		}

		signs[oraclePubKey] = sign
		realSignCount++
	}

	if realSignCount < len(consuls)*2/3 {
		return nil
	}

	bftOraclesByNebula, err := scheduler.client.BftOraclesByNebula(chainType, nebulaId)
	if err != nil {
		return err
	}
	var newOracles []*account.OraclesPubKey
	for k, v := range bftOraclesByNebula {
		oracleAddress, err := account.StringToOraclePubKey(k, v)
		if err != nil {
			return err
		}
		newOracles = append(newOracles, &oracleAddress)
	}
	for i := len(newOracles); i < OracleCount; i++ {
		newOracles = append(newOracles, nil)
	}

	tx, err := scheduler.Adaptors[chainType].SetOraclesToNebula(nebulaId, newOracles, signs, round, scheduler.ctx)
	if err != nil {
		return err
	}
	if tx != "" {
		err := scheduler.Adaptors[chainType].WaitTx(tx, scheduler.ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Tx nebula (%s) oracles update: %s \n", nebulaId.ToString(chainType), tx)
	}

	return nil
}

func (scheduler *Scheduler) updateTargetChainsPubKeys() {
	for chain, adaptor := range scheduler.Adaptors {
		scheduler.setConsulTargetChainPubKey(adaptor.PubKey(), chain)
	}
}

func (scheduler *Scheduler) setConsulTargetChainPubKey(oracle account.OraclesPubKey, chainType account.ChainType) error {
	zap.L().Debug("Start adding oracles")
	oracles, err := scheduler.client.OraclesByValidator(scheduler.Ledger.PubKey)
	if err != nil && err != gravity.ErrValueNotFound {
		zap.L().Error(err.Error())
		return err
	}
	//zap.L().Sugar().Debug("Oracles", oracles)

	if _, ok := oracles[chainType]; ok {
		zap.L().Sugar().Debugf("pubkey for chain [%s] exists", chainType)
		return nil
	}
	zap.L().Debug("Creating transaction")
	tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.AddOracle, scheduler.Ledger.PrivKey)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Debug("Adding values")
	tx.AddValues([]transactions.Value{
		transactions.BytesValue{
			Value: []byte{byte(chainType)},
		},
		transactions.BytesValue{
			Value: oracle[:],
		},
	})
	zap.L().Debug("Sending transaction")
	err = scheduler.client.SendTx(tx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

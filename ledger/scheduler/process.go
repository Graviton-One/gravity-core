package scheduler

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

func (scheduler *Scheduler) process(height int64) {
	err := scheduler.processByHeight(height)
	if err != nil {
		fmt.Printf("Error:%s\n", err)
	}
}
func (scheduler *Scheduler) processByHeight(height int64) error {
	logger, _ := zap.NewDevelopment()
	roundId := CalculateRound(height)
	logger.Sugar().Debugf("Process by height, round id: %d", roundId)
	consulInfo, err := scheduler.consulInfo()
	if err != nil {
		return err
	}
	zap.L().Sugar().Debug("Check round")
	isExist := true
	if IsRoundStart(height) {
		logger.Sugar().Debugf("Round start")
		roundId := int64(CalculateRound(height) - 1)
		logger.Sugar().Debugf("Round started, id: %d", roundId)
		index := roundId % int64(consulInfo.TotalCount)

		if index == int64(consulInfo.ConsulIndex) {
			for k, v := range scheduler.Adaptors {

				lastRound, err := v.LastRound(scheduler.ctx)
				logger.Sugar().Debugf("Round: chain: %d, Last round: %d", k, lastRound)
				if err != nil {
					return err
				}
				if uint64(roundId) <= lastRound {
					logger.Sugar().Debugf("roundid <= lastround")
					continue
				}

				err = scheduler.sendConsulsToGravityContract(roundId, k)
				if err != nil {
					return err
				}
			}

			nebulae, err := scheduler.client.Nebulae()
			if err != nil {
				return err
			}

			for k, v := range nebulae {
				nebulaId, err := account.StringToNebulaId(k, v.ChainType)
				if err != nil {
					fmt.Printf("Error:%s\n", err.Error())
					continue
				}

				err = scheduler.sendOraclesToNebula(nebulaId, v.ChainType, roundId)
				if err != nil {
					fmt.Printf("SendOraclesToNebula Error: %s\n", err.Error())
					continue
				}
			}
		}
	}
	zap.L().Sugar().Debug("Check adaptors")
	for k, v := range scheduler.Adaptors {
		zap.L().Sugar().Debug("Iterate adaptors: ", k)
		lastRound, err := v.LastRound(scheduler.ctx)
		if err != nil {
			return err
		}
		isExist = uint64(roundId) == lastRound
		if uint64(roundId) <= lastRound {
			continue
		}

		err = scheduler.signConsulsResult(roundId, k)
		if err != nil {
			return err
		}

		nebulae, err := scheduler.client.Nebulae()
		if err != nil {
			return err
		}

		for k, v := range nebulae {
			zap.L().Sugar().Debug("Iterate Nebula: ", k, v.ChainType)
			nebulaId, err := account.StringToNebulaId(k, v.ChainType)
			if err != nil {
				fmt.Printf("Error:%s\n", err.Error())
				continue
			}
			err = scheduler.signOraclesByNebula(roundId, nebulaId, v.ChainType)
			if err != nil {
				continue
			}

		}
	}

	lastRound, err := scheduler.client.LastRoundApproved()
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	}
	senderIndex := height % int64(consulInfo.TotalCount)
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
func (scheduler *Scheduler) signConsulsResult(roundId int64, chainType account.ChainType) error {
	_, err := scheduler.client.SignNewConsulsByConsul(scheduler.Ledger.PubKey, chainType, roundId)
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	} else if err == nil {
		return nil
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

	sign, err := scheduler.Adaptors[chainType].SignConsuls(consulsAddresses, roundId)
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
func (scheduler *Scheduler) signOraclesByNebula(roundId int64, nebulaId account.NebulaId, chainType account.ChainType) error {
	zap.L().Sugar().Debug("SignOraclesByNebula chain: ", chainType)
	_, err := scheduler.client.SignNewOraclesByConsul(scheduler.Ledger.PubKey, chainType, nebulaId, roundId)
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	} else if err == nil {
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
	zap.L().Sugar().Debug("SignOracles chain: ", chainType, scheduler.Adaptors)
	sign, err := scheduler.Adaptors[chainType].SignOracles(nebulaId, newOracles)
	if err != nil {
		return err
	}

	tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.SignNewOracles, scheduler.Ledger.PrivKey)
	if err != nil {
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
		return err
	}

	return nil
}
func (scheduler *Scheduler) sendConsulsToGravityContract(round int64, chainType account.ChainType) error {

	logger, err := zap.NewDevelopment()
	logger.Sugar().Debugf("Send consuls to gravity contract Round: %d, Chain: %d\n", round, chainType)
	exist, err := scheduler.Adaptors[chainType].RoundExist(round, scheduler.ctx)
	if err != nil {
		logger.Sugar().Debugf("Exist error: %s\n", err.Error())
		return err
	}

	if exist {
		logger.Sugar().Debugf("Round exist")
		return nil
	}

	lastRound, err := scheduler.Adaptors[chainType].LastRound(scheduler.ctx)
	if err != nil {
		logger.Sugar().Debugf("last round error: %s", err.Error())
		return err
	}
	logger.Sugar().Debugf("last round: %d , chain type: %d", lastRound, chainType)
	if round <= int64(lastRound) {
		return nil
	}

	consuls, err := scheduler.client.Consuls()
	if err != nil {
		return err
	}

	newConsuls, err := scheduler.client.ConsulsCandidate()
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

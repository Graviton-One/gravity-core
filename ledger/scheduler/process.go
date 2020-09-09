package scheduler

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/gravity"

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
	roundId := height / CalculateScoreInterval

	consulInfo, err := scheduler.consulInfo()
	if err != nil {
		return err
	}

	isExist := true
	senderIndex := height % int64(consulInfo.TotalCount)
	for k, v := range scheduler.Adaptors {
		lastRound, err := v.LastRound(scheduler.ctx)
		if err != nil {
			return err
		}
		isExist = uint64(roundId) == lastRound
		if uint64(roundId) <= lastRound {
			continue
		}

		if height%CalculateScoreInterval < CalculateScoreInterval/2 {
			err := scheduler.signConsulsResult(roundId, k)
			if err != nil {
				return err
			}

			nebulae, err := scheduler.client.Nebulae()
			if err != nil {
				return err
			}

			for nebulaId, _ := range nebulae {
				err := scheduler.signOraclesByNebula(roundId, nebulaId, k)
				if err != nil {
					continue
				}
			}
		} else {
			if senderIndex != int64(consulInfo.ConsulIndex) {
				continue
			}

			err := scheduler.sendConsulsToGravityContract(roundId, k)
			if err != nil {
				return err
			}

			nebulae, err := scheduler.client.Nebulae()
			if err != nil {
				return err
			}

			for nebulaId, _ := range nebulae {
				err := scheduler.sendOraclesToNebula(nebulaId, k, roundId)
				if err != nil {
					continue
				}
			}
		}
	}

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
	var newOracles []account.OraclesPubKey
	for pubKey := range bftOraclesByNebula {
		newOracles = append(newOracles, pubKey)
	}
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
	exist, err := scheduler.Adaptors[chainType].RoundExist(round, scheduler.ctx)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	lastRound, err := scheduler.Adaptors[chainType].LastRound(scheduler.ctx)
	if err != nil {
		return err
	}

	if round <= int64(lastRound) {
		return nil
	}

	consuls, err := scheduler.client.Consuls()
	if err != nil {
		return err
	}

	var signs [][]byte
	var empty [65]byte
	for i := 0; i < OracleCount; i++ {
		if i >= len(consuls) {
			signs = append(signs, empty[:])
			continue
		}

		v := consuls[i]
		sign, err := scheduler.client.SignNewConsulsByConsul(v.PubKey, chainType, round)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		}
		if err == gravity.ErrValueNotFound {
			signs = append(signs, empty[:])
			continue
		}

		signs = append(signs, sign)
	}

	newConsuls, err := scheduler.client.ConsulsCandidate()
	if err != nil {
		return err
	}

	var newConsulsAddresses []*account.OraclesPubKey
	for i := 0; i < OracleCount; i++ {
		if i >= len(consuls) {
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

	var signs [][]byte
	for _, v := range consuls {
		sign, err := scheduler.client.SignNewOraclesByConsul(v.PubKey, chainType, nebulaId, round)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		}

		if err == gravity.ErrValueNotFound {
			var empty [32]byte
			signs = append(signs, empty[:])
			continue
		}

		signs = append(signs, sign)
	}

	oracles, err := scheduler.client.OraclesByNebula(nebulaId, chainType)
	if err != nil {
		return err
	}

	var oraclesAddresses []account.OraclesPubKey
	for k, _ := range oracles {
		oraclesAddresses = append(oraclesAddresses, k)
	}

	tx, err := scheduler.Adaptors[chainType].SetOraclesToNebula(nebulaId, oraclesAddresses, signs, round, scheduler.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx nebula (%s) oracles update: %s \n", nebulaId, tx)
	return nil
}

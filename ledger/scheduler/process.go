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
		fmt.Printf("Process by height\n")
		fmt.Printf("Error:%s\n", err)
	}
}
func (scheduler *Scheduler) processByHeight(height int64) error {
	roundId := CalculateRound(height)

	consulInfo, err := scheduler.consulInfo()
	if err != nil {
		fmt.Printf("consulInfo error\n")
		return err
	}

	isExist := true
	if IsRoundStart(height) {
		roundId := int64(CalculateRound(height) - 1)

		index := roundId % int64(consulInfo.TotalCount)

		if index == int64(consulInfo.ConsulIndex) {
			for k, v := range scheduler.Adaptors {
				lastRound, err := v.LastRound(scheduler.ctx)
				if err != nil {
					fmt.Printf("lastround error\n")
					return err
				}
				if uint64(roundId) <= lastRound {
					continue
				}

				err = scheduler.sendConsulsToGravityContract(roundId, k)
				if err != nil {
					fmt.Printf("sendConsulsToGravityContract:%d -> %d\n", roundId, k)
					return err
				}
			}

			nebulae, err := scheduler.client.Nebulae()
			if err != nil {
				fmt.Printf("client nebulae error\n")
				return err
			}

			for k, v := range nebulae {
				nebulaId, err := account.StringToNebulaId(k, v.ChainType)
				if err != nil {
					fmt.Printf("Nebula Map: %s -> %d\n", k, v.ChainType)
					fmt.Printf("Error:%s\n", err.Error())
					continue
				}

				err = scheduler.sendOraclesToNebula(nebulaId, v.ChainType, v.ChainSelector, roundId)
				if err != nil {
					fmt.Printf("SendOraclesToNebula Error: %s\n", err.Error())
					continue
				}
			}
		}
	}

	for k, v := range scheduler.Adaptors {

		lastRound, err := v.LastRound(scheduler.ctx)
		if err != nil {
			fmt.Printf("Adapter lastround Error: %s\n", err.Error())
			return err
		}
		isExist = uint64(roundId) == lastRound
		if uint64(roundId) <= lastRound {
			continue
		}

		err = scheduler.signConsulsResult(roundId, k)
		if err != nil {
			fmt.Printf("SignConsulsResult: %d -> %d\n", roundId, k)
			fmt.Printf("SignConsulsResult Error: %s\n", err.Error())
			return err
		}

		nebulae, err := scheduler.client.Nebulae()
		if err != nil {
			fmt.Printf("ClientNebula Error: %s\n", err.Error())
			return err
		}

		for k2, v := range nebulae {
			nebulaId, err := account.StringToNebulaId(k2, v.ChainType)
			if err != nil {
				fmt.Printf("String to Nebula ID:%s -> %d\n", k2, v.ChainType)
				fmt.Printf("Error:%s\n", err.Error())
				continue
			}
			err = scheduler.signOraclesByNebula(roundId, nebulaId, v.ChainType, k)
			if err != nil {
				continue
			}

		}
	}

	lastRound, err := scheduler.client.LastRoundApproved()
	if err != nil && err != gravity.ErrValueNotFound {
		fmt.Printf("Approve Error: %s\n", err.Error())
		return err
	}
	senderIndex := height % int64(consulInfo.TotalCount)
	if isExist && uint64(roundId) > lastRound && senderIndex == int64(consulInfo.ConsulIndex) {
		tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.ApproveLastRound, scheduler.Ledger.PrivKey)
		if err != nil {
			fmt.Printf("SendTx Error: %s\n", err.Error())
			fmt.Printf("PubKey: %s\n", string(scheduler.Ledger.PrivKey.PubKey().Bytes()))
			return err
		}
		err = scheduler.client.SendTx(tx)
		if err != nil {
			fmt.Printf("SendTxCall Error: %s\n", err.Error())
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
		fmt.Printf("SignNewConsulsByConsul Error: %s\n", err.Error())
		return err
	} else if err == nil {
		return nil
	}

	consuls, err := scheduler.client.ConsulsCandidate()
	if err != nil {
		fmt.Printf("ConsulsCandidate Error: %s\n", err.Error())
		return err
	}
	fmt.Printf("Consuls candidate length: %d\n", len(consuls))
	var consulsAddresses []*account.OraclesPubKey
	for i := 0; i < OracleCount; i++ {
		if i >= len(consuls) {
			consulsAddresses = append(consulsAddresses, nil)
			continue
		}
		v := consuls[i]
		oraclesByConsul, err := scheduler.client.OraclesByValidator(v.PubKey)
		fmt.Println("Consul: ")
		fmt.Println(v.PubKey)
		fmt.Println("Oracles: ")
		fmt.Println(oraclesByConsul)
		if err == gravity.ErrValueNotFound {
			consulsAddresses = append(consulsAddresses, nil)
			continue
		} else if err != nil {
			fmt.Printf("OracleByValidator Error: %s\n", err.Error())
			return err
		}

		fmt.Printf("CHAIN TYPE FOR %d\n", scheduler.Adaptors[chainType].ChainType())
		oracle := oraclesByConsul[scheduler.Adaptors[chainType].ChainType()]
		consulsAddresses = append(consulsAddresses, &oracle)
		fmt.Println("oracle")
		fmt.Println(oracle)
	}
	fmt.Println("Consuls Addresses: ")
	fmt.Println(consulsAddresses)
	sign, err := scheduler.Adaptors[chainType].SignConsuls(consulsAddresses, roundId)
	if err != nil {
		fmt.Printf("Adaptor SignConsuls Error: %s\n", err.Error())
		return err
	}
	tx, err := transactions.New(scheduler.Ledger.PubKey, transactions.SignNewConsuls, scheduler.Ledger.PrivKey)
	if err != nil {
		fmt.Printf("Transactions Error: %s\n", err.Error())
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
		fmt.Printf("SignNewConsulsByConsul SendTx Error: %s\n", err.Error())
		return err
	}
	return nil
}
func (scheduler *Scheduler) signOraclesByNebula(roundId int64, nebulaId account.NebulaId, chainType account.ChainType, chainSelector account.ChainType) error {
	_, err := scheduler.client.SignNewOraclesByConsul(scheduler.Ledger.PubKey, chainType, chainSelector, nebulaId, roundId)
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	} else if err == nil {
		return nil
	}

	bftOraclesByNebula, err := scheduler.client.BftOraclesByNebula(chainType, chainSelector, nebulaId)
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
func (scheduler *Scheduler) sendOraclesToNebula(nebulaId account.NebulaId, chainType account.ChainType, chainSelector account.ChainType, round int64) error {
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

		sign, err := scheduler.client.SignNewOraclesByConsul(v.PubKey, chainType, chainSelector, nebulaId, round)
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

	bftOraclesByNebula, err := scheduler.client.BftOraclesByNebula(chainType, chainSelector, nebulaId)
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

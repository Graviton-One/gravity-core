package rpc

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Gravity-Tech/gravity-core/ledger/scheduler"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

func (s *Scheduler) ListenRpcServer() {
	http.HandleFunc("/process", s.BlockHandler)
	err := http.ListenAndServe(s.Host, nil)
	if err != nil {
		fmt.Printf("Error Private RPC: %s", err.Error())
	}
}

func (s *Scheduler) BlockHandler(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["height"]

	if !ok || len(keys[0]) < 1 {
		http.Error(w, "invalid query", http.StatusBadRequest)
		return
	}

	height, err := strconv.ParseUint(keys[0], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.HandleBlock(int64(height))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *Scheduler) HandleBlock(height int64) error {
	roundId := height / scheduler.CalculateScoreInterval

	consulInfo, err := s.consulInfo()
	if err != nil {
		return err
	}

	for k, v := range s.Adaptors {
		lastRound, err := v.LastRound(s.ctx)
		if err != nil {
			return err
		}

		if uint64(roundId) <= lastRound {
			continue
		}

		if height%scheduler.CalculateScoreInterval < scheduler.CalculateScoreInterval/2 {
			err := s.signConsulsResult(roundId, k)
			if err != nil {
				return err
			}

			nebulae, err := s.client.Nebulae()
			if err != nil {
				return err
			}

			for nebulaId, _ := range nebulae {
				err := s.signOraclesByNebula(roundId, nebulaId, k)
				if err != nil {
					continue
				}
			}
		} else {
			senderIndex := height % int64(consulInfo.TotalCount)
			if senderIndex != int64(consulInfo.ConsulIndex) {
				return nil
			}

			err := s.sendConsulsToGravityContract(roundId, k)
			if err != nil {
				return err
			}

			nebulae, err := s.client.Nebulae()
			if err != nil {
				return err
			}

			for nebulaId, _ := range nebulae {
				err := s.sendOraclesToNebula(nebulaId, k, roundId)
				if err != nil {
					continue
				}
			}
		}
	}

	return nil
}

func (s *Scheduler) updateTargetChain(round int64, isConsul bool, consulIndex int, consulsCount int) error {
	if round < 1 {
		return nil
	}
	for k, v := range s.Adaptors {
		exist, err := v.RoundExist(round, s.ctx)
		if err != nil {
			return err
		}
		if exist {
			break
		}

		if isConsul {
			err = s.signConsulsResult(round, k)
			if err != nil {
				return err
			}

			height, err := v.GetHeight(s.ctx)
			if err != nil {
				return err
			}
			senderIndex := height % uint64(consulsCount)
			if uint64(consulIndex) == senderIndex {
				err = s.sendConsulsToGravityContract(round, k)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (s *Scheduler) consulInfo() (*ConsulInfo, error) {
	consuls, err := s.client.Consuls()
	if err != nil {
		return nil, err
	}

	isConsul := false
	consulIndex := 0
	for i, consul := range consuls {
		if s.Ledger.PubKey == consul.PubKey {
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
func (s *Scheduler) signConsulsResult(roundId int64, chainType account.ChainType) error {
	_, err := s.client.SignNewConsulsByConsul(s.Ledger.PubKey, chainType, roundId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	} else if err == nil {
		return nil
	}

	consuls, err := s.client.ConsulsCandidate()
	if err != nil {
		return err
	}

	var consulsAddresses []account.OraclesPubKey
	for _, v := range consuls {
		oraclesByConsul, err := s.client.OraclesByValidator(v.PubKey)
		if err != nil {
			return err
		}

		consulsAddresses = append(consulsAddresses, oraclesByConsul[chainType])
	}

	sign, err := s.Adaptors[chainType].SignConsuls(consulsAddresses, roundId)
	if err != nil {
		return err
	}
	tx, err := transactions.New(s.Ledger.PubKey, transactions.SignNewConsuls, s.Ledger.PrivKey)
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
	err = s.client.SendTx(tx)
	if err != nil {
		return err
	}
	return nil
}
func (s *Scheduler) signOraclesByNebula(roundId int64, nebulaId account.NebulaId, chainType account.ChainType) error {
	_, err := s.client.SignNewOraclesByConsul(s.Ledger.PubKey, chainType, nebulaId, roundId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	} else if err == nil {
		return nil
	}

	bftOraclesByNebula, err := s.client.BftOraclesByNebula(chainType, nebulaId)
	if err != nil {
		return err
	}
	var newOracles []account.OraclesPubKey
	for pubKey := range bftOraclesByNebula {
		newOracles = append(newOracles, pubKey)
	}
	sign, err := s.Adaptors[chainType].SignOracles(nebulaId, newOracles)
	if err != nil {
		return err
	}

	tx, err := transactions.New(s.Ledger.PubKey, transactions.SignNewOracles, s.Ledger.PrivKey)
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
	err = s.client.SendTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) sendConsulsToGravityContract(round int64, chainType account.ChainType) error {
	exist, err := s.Adaptors[chainType].RoundExist(round, s.ctx)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	lastRound, err := s.Adaptors[chainType].LastRound(s.ctx)
	if err != nil {
		return err
	}

	if round <= int64(lastRound) {
		return nil
	}

	consuls, err := s.client.Consuls()
	if err != nil {
		return err
	}

	var signs [][]byte
	for _, v := range consuls {
		sign, err := s.client.SignNewConsulsByConsul(v.PubKey, chainType, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			var empty [65]byte
			signs = append(signs, empty[:])
			continue
		}

		signs = append(signs, sign)
	}

	newConsuls, err := s.client.Consuls()
	if err != nil {
		return err
	}

	var newConsulsAddresses []account.OraclesPubKey
	for _, v := range newConsuls {
		oracles, err := s.client.OraclesByValidator(v.PubKey)
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

	id, err := s.Adaptors[chainType].SendConsulsToGravityContract(newConsulsAddresses, signs, round, s.ctx)
	if err != nil {
		return err
	}
	if id != "" {
		err := s.Adaptors[chainType].WaitTx(id, s.ctx)
		if err != nil {
			return err
		}

		fmt.Printf("Tx consuls update (%s): %s\n", chainType.String(), id)
	}
	return nil
}

func (s *Scheduler) sendOraclesToNebula(nebulaId account.NebulaId, chainType account.ChainType, round int64) error {
	consuls, err := s.client.Consuls()
	if err != nil {
		return err
	}

	var signs [][]byte
	for _, v := range consuls {
		sign, err := s.client.SignNewOraclesByConsul(v.PubKey, chainType, nebulaId, round)
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

	oracles, err := s.client.OraclesByNebula(nebulaId, chainType)
	if err != nil {
		return err
	}

	var oraclesAddresses []account.OraclesPubKey
	for k, _ := range oracles {
		oraclesAddresses = append(oraclesAddresses, k)
	}

	tx, err := s.Adaptors[chainType].SetOraclesToNebula(nebulaId, oraclesAddresses, signs, round, s.ctx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx nebula (%s) oracles update: %s \n", nebulaId, tx)
	return nil
}

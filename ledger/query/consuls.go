package query

import (
	"encoding/json"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
)

type SignByConsulRq struct {
	ConsulPubKey string
	ChainType    account.ChainType
	NebulaId     string
	RoundId      int64
}

func consuls(store *storage.Storage) ([]storage.Consul, error) {
	v, err := store.Consuls()
	if err != nil && err != storage.ErrKeyNotFound {
		return nil, err
	}

	return v, nil
}

func consulsCandidate(store *storage.Storage) ([]storage.Consul, error) {
	v, err := store.ConsulsCandidate()
	if err != nil && err != storage.ErrKeyNotFound {
		return nil, err
	}

	return v, nil
}

func signNewConsulsByConsul(store *storage.Storage, value []byte) ([]byte, error) {
	var rq SignByConsulRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	consul, err := account.HexToValidatorPubKey(rq.ConsulPubKey)
	if err != nil {
		return nil, err
	}

	v, err := store.SignConsulsByConsul(consul, rq.ChainType, rq.RoundId)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func signNewOraclesByConsul(store *storage.Storage, value []byte) ([]byte, error) {
	var rq SignByConsulRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	consul, err := account.HexToValidatorPubKey(rq.ConsulPubKey)
	if err != nil {
		return nil, err
	}

	nebula, err := account.StringToNebulaId(rq.NebulaId, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.SignOraclesByConsul(consul, nebula, rq.RoundId)
	if err != nil {
		return nil, err
	}

	return v, nil
}

package query

import (
	"encoding/json"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
)

type ByValidatorRq struct {
	PubKey string
}

type ByNebulaRq struct {
	ChainType     account.ChainType
	NebulaAddress string
}

type ResultsRq struct {
	ChainType     account.ChainType
	Height        uint64
	NebulaAddress string
}

func oraclesByValidator(store *storage.Storage, value []byte) (storage.OraclesByTypeMap, error) {
	var rq ByValidatorRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	pubKey, err := account.HexToValidatorPubKey(rq.PubKey)
	if err != nil {
		return nil, err
	}

	v, err := store.OraclesByConsul(pubKey)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func oraclesByNebula(store *storage.Storage, value []byte) (storage.OraclesMap, error) {
	var rq ByNebulaRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebula, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.OraclesByNebula(nebula)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func bftOraclesByNebula(store *storage.Storage, value []byte) (storage.OraclesMap, error) {
	var rq ByNebulaRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebula, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.BftOraclesByNebula(nebula)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func results(store *storage.Storage, value []byte) ([][]byte, error) {
	var rq ResultsRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebula, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.Results(nebula, rq.Height)
	if err != nil {
		return nil, err
	}

	return v, nil
}

package query

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/mr-tron/base58"

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

	pubKey, err := account.HexToPubKey(rq.PubKey)
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

	var nebula []byte
	switch rq.ChainType {
	case account.Ethereum:
		nebula, err = hexutil.Decode(rq.NebulaAddress)
		if err != nil {
			return nil, err
		}
	case account.Waves:
		nebula, err = base58.Decode(rq.NebulaAddress)
		if err != nil {
			return nil, err
		}
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

	var nebula []byte
	switch rq.ChainType {
	case account.Ethereum:
		nebula, err = hexutil.Decode(rq.NebulaAddress)
		if err != nil {
			return nil, err
		}
	case account.Waves:
		nebula, err = base58.Decode(rq.NebulaAddress)
		if err != nil {
			return nil, err
		}
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

	var nebula []byte
	switch rq.ChainType {
	case account.Ethereum:
		nebula, err = hexutil.Decode(rq.NebulaAddress)
		if err != nil {
			return nil, err
		}
	case account.Waves:
		nebula, err = base58.Decode(rq.NebulaAddress)
		if err != nil {
			return nil, err
		}
	}

	v, err := store.Results(nebula, rq.Height)
	if err != nil {
		return nil, err
	}

	return v, nil
}

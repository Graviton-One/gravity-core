package query

import (
	"encoding/json"

	"github.com/mr-tron/base58"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
)

type OraclesByValidatorRq struct {
	PubKey string
}

type OraclesByNebulaRq struct {
	ChainType     account.ChainType
	NebulaAddress string
}

func oraclesByValidator(store *storage.Storage, value []byte) (storage.OraclesByTypeMap, error) {
	var rq OraclesByValidatorRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	pubKey, err := account.HexToPubKey(rq.PubKey)
	if err != nil {
		return nil, err
	}

	v, err := store.OraclesByValidator(pubKey)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func oraclesByNebula(store *storage.Storage, value []byte) (storage.OraclesMap, error) {
	var rq OraclesByNebulaRq
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

package query

import (
	"encoding/json"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func allValidators(store *storage.Storage, _ []byte) ([]byte, error) {
	scores, err := store.Scores()
	result := make([]string, len(scores))

	if err != nil {
		return make([]byte, 0), err
	}

	for consulPubKey, _ := range scores {
		result = append(result, hexutil.Encode(consulPubKey[:]))
	}

	encoded, err := json.Marshal(result)

	if err != nil {
		return make([]byte, 0), err
	}

	return encoded, nil
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

func results(store *storage.Storage, value []byte) ([]string, error) {
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

func nebulaOraclesIndex(store *storage.Storage, value []byte) (uint64, error) {
	var rq ByNebulaRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return 0, err
	}

	nebula, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return 0, err
	}

	v, err := store.NebulaOraclesIndex(nebula)
	if err != nil {
		return 0, err
	}

	return v, nil
}

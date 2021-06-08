package query

import (
	"encoding/json"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/storage"
)

type RoundHeightRq struct {
	ChainType    account.ChainType
	LedgerHeight uint64
}
type CommitHashRq struct {
	ChainType     account.ChainType
	NebulaAddress string
	Height        int64
	PulseId       int64
	OraclePubKey  string
}
type RevealRq struct {
	ChainType     account.ChainType
	NebulaAddress string
	OraclePubKey  string
	Height        int64
	PulseId       int64
	CommitHash    string
}
type ResultRq struct {
	ChainType     account.ChainType
	NebulaAddress string
	Height        int64
	OraclePubKey  string
}

func roundHeight(store *storage.Storage, value []byte) (uint64, error) {
	var rq RoundHeightRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return 0, err
	}

	return store.RoundHeight(rq.ChainType, rq.LedgerHeight)
}
func commitHash(store *storage.Storage, value []byte) ([]byte, error) {
	var rq CommitHashRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebulaAddress, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	oraclePubKey, err := account.StringToOraclePubKey(rq.OraclePubKey, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.CommitHash(nebulaAddress, rq.Height, rq.PulseId, oraclePubKey)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func reveal(store *storage.Storage, value []byte) ([]byte, error) {
	var rq RevealRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebulaAddress, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	commitHash, err := hexutil.Decode(rq.CommitHash)
	if err != nil {
		return nil, err
	}

	oraclePubKey, err := account.StringToOraclePubKey(rq.OraclePubKey, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.Reveal(nebulaAddress, rq.Height, rq.PulseId, commitHash, oraclePubKey)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func reveals(store *storage.Storage, value []byte) ([]string, error) {
	var rq RevealRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebulaAddress, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.Reveals(nebulaAddress, rq.Height, rq.PulseId)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func result(store *storage.Storage, value []byte) ([]byte, error) {
	var rq ResultRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebulaAddress, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	oraclePubKey, err := account.StringToOraclePubKey(rq.OraclePubKey, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.Result(nebulaAddress, rq.Height, oraclePubKey)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func nebulae(store *storage.Storage) (storage.NebulaMap, error) {
	v, err := store.Nebulae()
	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, storage.ErrKeyNotFound
	}

	return v, nil
}

func nebulaInfo(store *storage.Storage, value []byte) (*storage.NebulaInfo, error) {
	var rq ByNebulaRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebulaAddress, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.NebulaInfo(nebulaAddress)
	if err != nil {
		return nil, err
	}

	return v, nil
}
func nebulaCustomParams(store *storage.Storage, value []byte) (*storage.NebulaCustomParams, error) {
	var rq ByNebulaRq
	err := json.Unmarshal(value, &rq)
	if err != nil {
		return nil, err
	}

	nebulaAddress, err := account.StringToNebulaId(rq.NebulaAddress, rq.ChainType)
	if err != nil {
		return nil, err
	}

	v, err := store.NebulaCustomParams(nebulaAddress)
	if err != nil {
		return nil, err
	}

	return v, nil
}

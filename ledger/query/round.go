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
	OraclePubKey  string
}
type RevealRq struct {
	NebulaAddress string
	Height        int64
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

	nebulaAddress, err := hexutil.Decode(rq.NebulaAddress)
	if err != nil {
		return nil, err
	}

	oraclePubKey, err := hexutil.Decode(rq.OraclePubKey)
	if err != nil {
		return nil, err
	}

	v, err := store.CommitHash(nebulaAddress, rq.Height, account.BytesToOraclePubKey(oraclePubKey, rq.ChainType))
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

	nebulaAddress, err := hexutil.Decode(rq.NebulaAddress)
	if err != nil {
		return nil, err
	}

	commitHash, err := hexutil.Decode(rq.CommitHash)
	if err != nil {
		return nil, err
	}

	v, err := store.Reveal(nebulaAddress, rq.Height, commitHash)
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

	nebulaAddress, err := hexutil.Decode(rq.NebulaAddress)
	if err != nil {
		return nil, err
	}

	oraclePubKey, err := hexutil.Decode(rq.OraclePubKey)
	if err != nil {
		return nil, err
	}

	v, err := store.Result(nebulaAddress, rq.Height, account.BytesToOraclePubKey(oraclePubKey, rq.ChainType))
	if err != nil {
		return nil, err
	}

	return v, nil
}

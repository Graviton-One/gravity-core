package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/dgraph-io/badger"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formResultKey(nebulaId account.NebulaId, height int64, oraclePubKey account.OraclesPubKey) []byte {
	return formKey(string(SignResultKey), hexutil.Encode(nebulaId), fmt.Sprintf("%d", height), hexutil.Encode(oraclePubKey[:]))
}

func (storage *Storage) Result(nebulaId account.NebulaId, height int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	b, err := storage.getValue(formResultKey(nebulaId, height, oraclePubKey))
	if err != nil {
		return nil, err
	}

	return b, err
}
func (storage *Storage) Results(nebulaId account.NebulaId, height uint64) ([][]byte, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := formKey(string(SignResultKey), hexutil.Encode(nebulaId), fmt.Sprintf("%d", height))
	var values [][]byte
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		item.Value(func(v []byte) error {
			values = append(values, v)
			return nil
		})
	}

	return values, nil
}
func (storage *Storage) SetResult(nebulaId account.NebulaId, height int64, oraclePubKey account.OraclesPubKey, reveal []byte) error {
	return storage.setValue(formResultKey(nebulaId, height, oraclePubKey), reveal)
}

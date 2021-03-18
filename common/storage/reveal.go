package storage

import (
	"encoding/base64"
	"fmt"

	"github.com/dgraph-io/badger"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formRevealKey(nebulaId account.NebulaId, height int64, pulseId int64, commitHash []byte, oraclePubKey account.OraclesPubKey) []byte {
	return formKey(string(RevealKey), hexutil.Encode(nebulaId[:]), fmt.Sprintf("%d", height), fmt.Sprintf("%d", pulseId), hexutil.Encode(commitHash), hexutil.Encode(oraclePubKey[:]))
}

func (storage *Storage) Reveal(nebulaId account.NebulaId, height int64, pulseId int64, commitHash []byte, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	zap.L().Sugar().Debugf("Reveal key: %s", formRevealKey(nebulaId, height, pulseId, commitHash, oraclePubKey))
	b, err := storage.getValue(formRevealKey(nebulaId, height, pulseId, commitHash, oraclePubKey))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) Reveals(nebulaId account.NebulaId, height int64, pulseId int64) ([]string, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := formKey(string(RevealKey), hexutil.Encode(nebulaId[:]), fmt.Sprintf("%d", height), fmt.Sprintf("%d", pulseId))
	zap.L().Sugar().Debugf("Reveals key prefix: %s", prefix)
	var values []string
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		item.Value(func(v []byte) error {
			values = append(values, base64.StdEncoding.EncodeToString(v))
			return nil
		})
	}

	return values, nil
}

func (storage *Storage) SetReveal(nebulaId account.NebulaId, height int64, pulseId int64, commitHash []byte, oraclePubKey account.OraclesPubKey, reveal []byte) error {
	zap.L().Sugar().Debugf("SetReveal key: %s", formRevealKey(nebulaId, height, pulseId, commitHash, oraclePubKey))
	return storage.setValue(formRevealKey(nebulaId, height, pulseId, commitHash, oraclePubKey), reveal)
}

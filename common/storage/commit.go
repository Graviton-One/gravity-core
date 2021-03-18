package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formCommitKey(nebulaId account.NebulaId, tcHeight int64, pulseId int64, oraclePubKey account.OraclesPubKey) []byte {
	return formKey(string(CommitKey), hexutil.Encode(nebulaId[:]), fmt.Sprintf("%d", tcHeight), fmt.Sprintf("%d", pulseId), hexutil.Encode(oraclePubKey[:]))
}

func (storage *Storage) CommitHash(nebulaId account.NebulaId, tcHeight int64, pulseId int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	zap.L().Sugar().Debugf("CommitHash key: %s", formCommitKey(nebulaId, tcHeight, pulseId, oraclePubKey))
	b, err := storage.getValue(formCommitKey(nebulaId, tcHeight, pulseId, oraclePubKey))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetCommitHash(nebulaId account.NebulaId, tcHeight int64, pulseId int64, oraclePubKey account.OraclesPubKey, commit []byte) error {
	zap.L().Sugar().Debugf("SetCommitHash key: %s", formCommitKey(nebulaId, tcHeight, pulseId, oraclePubKey))
	return storage.setValue(formCommitKey(nebulaId, tcHeight, pulseId, oraclePubKey), commit)
}

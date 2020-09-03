package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formCommitKey(nebulaId account.NebulaId, pulseId int64, oraclePubKey account.OraclesPubKey) []byte {
	return formKey(string(CommitKey), hexutil.Encode(nebulaId[:]), fmt.Sprintf("%d", pulseId), hexutil.Encode(oraclePubKey[:]))
}

func (storage *Storage) CommitHash(nebulaId account.NebulaId, pulseId int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	b, err := storage.getValue(formCommitKey(nebulaId, pulseId, oraclePubKey))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetCommitHash(nebulaId account.NebulaId, pulseId int64, oraclePubKey account.OraclesPubKey, commit []byte) error {
	return storage.setValue(formCommitKey(nebulaId, pulseId, oraclePubKey), commit)
}

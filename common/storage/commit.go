package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formCommitKey(nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey) []byte {
	return formKey(string(CommitKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", height), hexutil.Encode(oraclePubKey[:]))
}

func (storage *Storage) CommitHash(nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	b, err := storage.getValue(formCommitKey(nebulaAddress, height, oraclePubKey))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetCommitHash(nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey, commit []byte) error {
	return storage.setValue(formCommitKey(nebulaAddress, height, oraclePubKey), commit)
}

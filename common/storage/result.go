package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formResultKey(nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey) []byte {
	return formKey(string(SignResultKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", height), hexutil.Encode(oraclePubKey[:]))
}

func (storage *Storage) Result(nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	b, err := storage.getValue(formResultKey(nebulaAddress, height, oraclePubKey))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetResult(nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey, reveal []byte) error {
	return storage.setValue(formResultKey(nebulaAddress, height, oraclePubKey), reveal)
}

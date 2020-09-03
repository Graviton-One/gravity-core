package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formRevealKey(nebulaId account.NebulaId, height int64, commitHash []byte) []byte {
	return formKey(string(RevealKey), hexutil.Encode(nebulaId), fmt.Sprintf("%d", height), hexutil.Encode(commitHash))
}

func (storage *Storage) Reveal(nebulaId account.NebulaId, height int64, commitHash []byte) ([]byte, error) {
	b, err := storage.getValue(formRevealKey(nebulaId, height, commitHash))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetReveal(nebulaId account.NebulaId, height int64, commitHash []byte, reveal []byte) error {
	return storage.setValue(formRevealKey(nebulaId, height, commitHash), reveal)
}

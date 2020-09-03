package storage

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formRevealKey(nebulaId account.NebulaId, pulseId int64, commitHash []byte) []byte {
	return formKey(string(RevealKey), hexutil.Encode(nebulaId[:]), fmt.Sprintf("%d", pulseId), hexutil.Encode(commitHash))
}

func (storage *Storage) Reveal(nebulaId account.NebulaId, pulseId int64, commitHash []byte) ([]byte, error) {
	b, err := storage.getValue(formRevealKey(nebulaId, pulseId, commitHash))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetReveal(nebulaId account.NebulaId, pulseId int64, commitHash []byte, reveal []byte) error {
	return storage.setValue(formRevealKey(nebulaId, pulseId, commitHash), reveal)
}

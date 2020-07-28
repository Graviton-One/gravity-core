package storage

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formRevealKey(nebulaAddress []byte, height int64, commitHash []byte) []byte {
	return formKey(string(RevealKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", height), hexutil.Encode(commitHash))
}

func (storage *Storage) Reveal(nebulaAddress []byte, height int64, commitHash []byte) ([]byte, error) {
	b, err := storage.getValue(formRevealKey(nebulaAddress, height, commitHash))
	if err != nil {
		return nil, err
	}

	return b, err
}

func (storage *Storage) SetReveal(nebulaAddress []byte, height int64, commitHash []byte, reveal []byte) error {
	return storage.setValue(formRevealKey(nebulaAddress, height, commitHash), reveal)
}

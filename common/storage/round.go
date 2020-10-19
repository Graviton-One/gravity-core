package storage

import (
	"encoding/binary"
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

func formNewRoundKey(chainType account.ChainType, ledgerHeight uint64) []byte {
	return formKey(string(BlockKey), chainType.String(), fmt.Sprintf("%d", ledgerHeight))
}

func (storage *Storage) RoundHeight(chainType account.ChainType, ledgerHeight uint64) (uint64, error) {
	b, err := storage.getValue(formNewRoundKey(chainType, ledgerHeight))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), err
}

func (storage *Storage) SetNewRound(chainType account.ChainType, ledgerHeight uint64, tcHeight uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], tcHeight)
	return storage.setValue(formNewRoundKey(chainType, ledgerHeight), b[:])
}

package storage

import (
	"encoding/binary"

	"github.com/Gravity-Tech/proof-of-concept/common/account"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func formScoreKey(validatorAddress account.PubKey) []byte {
	return formKey(string(ScoreKey), hexutil.Encode(validatorAddress[:]))
}

func (storage *Storage) Score(validatorAddress account.PubKey) (uint64, error) {
	b, err := storage.getValue(formScoreKey(validatorAddress))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), err
}

func (storage *Storage) SetScore(validatorAddress account.PubKey, score uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], score)
	err := storage.setValue(formScoreKey(validatorAddress), b[:])
	if err != nil {
		return err
	}

	return err
}

package storage

import (
	"encoding/binary"
	"strings"

	"github.com/dgraph-io/badger"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/keys"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ScoresByValidatorMap map[account.ValidatorPubKey]uint64

func formScoreKey(validatorAddress account.ValidatorPubKey) []byte {
	return formKey(string(ScoreKey), hexutil.Encode(validatorAddress[:]))
}
func parseScoreKey(value []byte) account.ValidatorPubKey {
	b := []byte(strings.Split(string(value), Separator)[1])
	var pubKey account.ValidatorPubKey
	copy(pubKey[:], b[:])
	return pubKey
}

func (storage *Storage) Score(validatorAddress account.ValidatorPubKey) (uint64, error) {
	b, err := storage.getValue(formScoreKey(validatorAddress))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}

func (storage *Storage) SetScore(validatorAddress account.ValidatorPubKey, score uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], score)
	err := storage.setValue(formScoreKey(validatorAddress), b[:])
	if err != nil {
		return err
	}

	return err
}

func (storage *Storage) Scores() (ScoresByValidatorMap, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(keys.ScoreKey)
	scores := make(ScoresByValidatorMap)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		item.Value(func(v []byte) error {
			scores[parseScoreKey(k)] = binary.BigEndian.Uint64(v)
			return nil
		})
	}

	return scores, nil
}

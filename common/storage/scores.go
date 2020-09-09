package storage

import (
	"encoding/binary"
	"strings"

	"github.com/dgraph-io/badger"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type ScoresByConsulMap map[account.ConsulPubKey]uint64

func formScoreKey(pubKey account.ConsulPubKey) []byte {
	return formKey(string(ScoreKey), hexutil.Encode(pubKey[:]))
}
func parseScoreKey(value []byte) (account.ConsulPubKey, error) {
	hex := []byte(strings.Split(string(value), Separator)[1])
	key, err := hexutil.Decode(string(hex))
	if err != nil {
		return [32]byte{}, err
	}
	var pubKey account.ConsulPubKey
	copy(pubKey[:], key[:])
	return pubKey, nil
}

func (storage *Storage) Score(pubKey account.ConsulPubKey) (uint64, error) {
	b, err := storage.getValue(formScoreKey(pubKey))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}
func (storage *Storage) SetScore(pubKey account.ConsulPubKey, score uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], score)
	err := storage.setValue(formScoreKey(pubKey), b[:])
	if err != nil {
		return err
	}

	return err
}

func (storage *Storage) Scores() (ScoresByConsulMap, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(ScoreKey)
	scores := make(ScoresByConsulMap)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		item.Value(func(v []byte) error {
			pubKey, err := parseScoreKey(k)
			if err != nil {
				return err
			}
			scores[pubKey] = binary.BigEndian.Uint64(v)
			return nil
		})
	}

	return scores, nil
}

package storage

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type Vote struct {
	PubKey account.ConsulPubKey
	Score  uint64
}

type VoteByConsulMap map[account.ConsulPubKey][]Vote

func formVoteKey(pubKey account.ConsulPubKey) []byte {
	return formKey(string(VoteKey), hexutil.Encode(pubKey[:]))
}
func parseVoteKey(value []byte) (account.ConsulPubKey, error) {
	hex := []byte(strings.Split(string(value), Separator)[1])
	key, err := hexutil.Decode(string(hex))
	if err != nil {
		return [32]byte{}, err
	}
	var pubKey account.ConsulPubKey
	copy(pubKey[:], key[:])
	return pubKey, nil
}

func (storage *Storage) Vote(pubKey account.ConsulPubKey) ([]Vote, error) {
	b, err := storage.getValue(formVoteKey(pubKey))
	if err != nil {
		return nil, err
	}

	var votes []Vote
	err = json.Unmarshal(b, &votes)
	if err != nil {
		return nil, err
	}
	return votes, err
}
func (storage *Storage) SetVote(pubKey account.ConsulPubKey, votes []Vote) error {
	return storage.setValue(formVoteKey(pubKey), votes)
}

func (storage *Storage) Votes() (VoteByConsulMap, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(VoteKey)
	votes := make(VoteByConsulMap)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		item.Value(func(v []byte) error {
			var vote []Vote
			err := json.Unmarshal(v, &vote)
			if err != nil {
				return err
			}
			pubKey, err := parseVoteKey(k)
			if err != nil {
				return err
			}
			votes[pubKey] = vote
			return nil
		})
	}

	return votes, nil
}

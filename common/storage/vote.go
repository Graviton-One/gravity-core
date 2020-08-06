package storage

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/keys"
)

type Vote struct {
	Target account.ValidatorPubKey
	Score  uint64
}

type VoteByValidatorMap map[account.ValidatorPubKey][]Vote

func formVoteKey(validator account.ValidatorPubKey) []byte {
	return formKey(string(VoteKey), hexutil.Encode(validator[:]))
}
func parseVoteKey(value []byte) account.ValidatorPubKey {
	b := []byte(strings.Split(string(value), Separator)[1])
	var pubKey account.ValidatorPubKey
	copy(pubKey[:], b[:])
	return pubKey
}

func (storage *Storage) Vote(validator account.ValidatorPubKey) ([]Vote, error) {
	b, err := storage.getValue(formVoteKey(validator))
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

func (storage *Storage) SetVote(validator account.ValidatorPubKey, votes []Vote) error {
	return storage.setValue(formVoteKey(validator), votes)
}

func (storage *Storage) Votes() (VoteByValidatorMap, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(keys.VoteKey)
	votes := make(VoteByValidatorMap)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		item.Value(func(v []byte) error {
			var vote []Vote
			err := json.Unmarshal(v, &vote)
			if err != nil {
				return err
			}
			votes[parseVoteKey(k)] = vote
			return nil
		})
	}

	return votes, nil
}

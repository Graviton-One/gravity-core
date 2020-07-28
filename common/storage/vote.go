package storage

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/proof-of-concept/common/account"
)

type Vote struct {
	Target string
	Score  float32
}

func FormVoteKey(validator account.PubKey) []byte {
	return formKey(string(VoteKey), hexutil.Encode(validator[:]))
}

func (storage *Storage) Vote(validator account.PubKey) ([]Vote, error) {
	b, err := storage.getValue(FormVoteKey(validator))
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

func (storage *Storage) SetVote(validator account.PubKey, votes []Vote) error {
	return storage.setValue(FormVoteKey(validator), votes)
}

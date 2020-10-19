package score

import (
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/score/trustgraph"
	"github.com/Gravity-Tech/gravity-core/common/storage"
)

const (
	Accuracy = 100
)

func UInt64ToFloat32Score(score uint64) float32 {
	return float32(score) / Accuracy
}
func Float32ToUInt64Score(score float32) uint64 {
	return uint64(score * Accuracy)
}

type Actor struct {
	Name      account.ConsulPubKey
	InitScore uint64
}

func Calculate(initScores storage.ScoresByConsulMap, votes storage.VoteByConsulMap) (storage.ScoresByConsulMap, error) {
	group := trustgraph.NewGroup()

	var newValidators []int
	idByValidator := make(map[account.ConsulPubKey]int)
	validatorById := make(map[int]account.ConsulPubKey)

	index := 0
	for k, v := range initScores {
		idByValidator[k] = index
		validatorById[index] = k
		err := group.InitialTrust(idByValidator[k], UInt64ToFloat32Score(v))
		if err != nil {
			return nil, err
		}
		index++
	}

	for voter, _ := range initScores {
		existVote := make(map[account.ConsulPubKey]bool)
		for _, vote := range votes[voter] {
			if voter == vote.PubKey {
				continue
			}
			if _, ok := idByValidator[vote.PubKey]; !ok {
				idByValidator[vote.PubKey] = index
				validatorById[index] = vote.PubKey
				err := group.InitialTrust(idByValidator[vote.PubKey], 0)
				if err != nil {
					return nil, err
				}
				newValidators = append(newValidators, index)
				index++
			}
			err := group.Add(idByValidator[voter], idByValidator[vote.PubKey], UInt64ToFloat32Score(vote.Score))
			if err != nil {
				return nil, err
			}
			existVote[vote.PubKey] = true
		}
		for validator, _ := range initScores {
			if existVote[validator] || voter == validator {
				continue
			}

			err := group.Add(idByValidator[voter], idByValidator[validator], UInt64ToFloat32Score(initScores[validator]))
			if err != nil {
				return nil, err
			}
		}
	}
	for _, v := range newValidators {
		for validator, _ := range initScores {
			err := group.Add(v, idByValidator[validator], UInt64ToFloat32Score(initScores[validator]))
			if err != nil {
				return nil, err
			}
		}
	}

	out := group.Compute()

	score := make(storage.ScoresByConsulMap)
	for i, v := range out {
		score[validatorById[i]] = Float32ToUInt64Score(v)
	}
	return score, nil
}

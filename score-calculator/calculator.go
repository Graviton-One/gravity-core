package score_calculator

import (
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/score-calculator/trustgraph"
)

const (
	Accuracy = 10
)

func UInt64ToFloat32Score(score uint64) float32 {
	return float32(score) / Accuracy
}
func Float32ToUInt64Score(score float32) uint64 {
	return uint64(score * Accuracy)
}

type Actor struct {
	Name      account.ValidatorPubKey
	InitScore uint64
}

func Calculate(initScores storage.ScoresByValidatorMap, votes storage.VoteByValidatorMap) (storage.ScoresByValidatorMap, error) {
	group := trustgraph.NewGroup()
	idByValidator := make(map[account.ValidatorPubKey]int)
	validatorById := make(map[int]account.ValidatorPubKey)
	count := 0
	for k, v := range initScores {
		idByValidator[k] = count
		validatorById[count] = k
		err := group.InitialTrust(idByValidator[k], UInt64ToFloat32Score(v))
		if err != nil {
			return nil, err
		}
		count++
	}

	for k, _ := range initScores {
		existVote := make(map[account.ValidatorPubKey]bool)
		for _, scoreV := range votes[k] {
			if k == scoreV.Target {
				continue
			}
			err := group.Add(idByValidator[k], idByValidator[scoreV.Target], UInt64ToFloat32Score(scoreV.Score))
			if err != nil {
				return nil, err
			}
			existVote[scoreV.Target] = true
		}
		for validator, _ := range initScores {
			if existVote[validator] {
				continue
			}
			if k == validator {
				continue
			}
			err := group.Add(idByValidator[k], idByValidator[validator], 1)
			if err != nil {
				return nil, err
			}
		}
	}

	out := group.Compute()

	score := make(storage.ScoresByValidatorMap)
	for i, v := range out {
		score[validatorById[i]] = Float32ToUInt64Score(v)
	}
	return score, nil
}

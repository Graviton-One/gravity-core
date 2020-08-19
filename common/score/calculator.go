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
	idByValidator := make(map[account.ConsulPubKey]int)
	validatorById := make(map[int]account.ConsulPubKey)
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
		existVote := make(map[account.ConsulPubKey]bool)
		for _, scoreV := range votes[k] {
			if k == scoreV.PubKey {
				continue
			}
			err := group.Add(idByValidator[k], idByValidator[scoreV.PubKey], UInt64ToFloat32Score(scoreV.Score))
			if err != nil {
				return nil, err
			}
			existVote[scoreV.PubKey] = true
		}
		for validator, _ := range initScores {
			if existVote[validator] {
				continue
			}
			if k == validator {
				continue
			}
			err := group.Add(idByValidator[k], idByValidator[validator], 0)
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

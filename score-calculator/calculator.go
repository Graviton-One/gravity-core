package score_calculator

import (
	"gravity-hub/score-calculator/models"
	"gravity-hub/score-calculator/trustgraph"
)

func Calculate(actors []models.Actor, votes map[string][]models.Vote) (map[string]float32, error) {
	group := trustgraph.NewGroup()
	actorsScore := make(map[string]int)
	for i, v := range actors {
		actorsScore[v.Name] = i
		err := group.InitialTrust(i, v.InitScore)
		if err != nil {
			return nil, err
		}
	}

	for _, v := range actors {
		existVote := make(map[string]bool)
		for _, scoreV := range votes[v.Name] {
			if v.Name == scoreV.Target {
				continue
			}
			err := group.Add(actorsScore[v.Name], actorsScore[scoreV.Target], scoreV.Score)
			if err != nil {
				return nil, err
			}
			existVote[scoreV.Target] = true
		}
		for _, actor := range actors {
			if existVote[actor.Name] {
				continue
			}
			if v.Name == actor.Name {
				continue
			}
			err := group.Add(actorsScore[v.Name], actorsScore[actor.Name], 1)
			if err != nil {
				return nil, err
			}
		}
	}

	out := group.Compute()

	score := make(map[string]float32)
	for i, v := range out {
		score[actors[i].Name] = v
	}
	return score, nil
}

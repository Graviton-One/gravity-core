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
	for k, v := range votes {
		for _, scoreV := range v {
			err := group.Add(actorsScore[k], actorsScore[scoreV.Target], scoreV.Score)
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

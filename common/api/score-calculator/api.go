package score_calculator

import (
	"context"
	"encoding/json"
	"fmt"
	api2 "gravity-hub/common/api"
	"gravity-hub/gravity-score-calculator/models"
	"net/http"
	"time"
)

type Client struct {
	host   string
	client *http.Client
}

func NewClient(host string) *Client {
	return &Client{host: host, client: &http.Client{Timeout: time.Second * 10}}
}

func (client *Client) Calculate(actors []models.Actor, votes map[string][]models.Vote, ctx context.Context) (*models.Response, error) {
	rq := models.Request{
		Votes:  votes,
		Actors: actors,
	}
	data, err := api2.Do(client.client, fmt.Sprintf("%s/api/calculate", client.host), api2.POST, rq, ctx)

	var response models.Response
	err = json.Unmarshal(data, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

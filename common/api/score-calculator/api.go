package score_calculator

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/score-calculator/models"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	host   string
	client *http.Client
	ctx    context.Context
}

func NewClient(host string, ctx context.Context) *Client {
	return &Client{
		host:   host,
		client: &http.Client{Timeout: time.Second * 10},
		ctx:    ctx,
	}
}

func (client *Client) Calculate(actors []models.Actor, votes map[string][]models.Vote) (*models.Response, error) {
	rq := models.Request{
		Votes:  votes,
		Actors: actors,
	}
	data, err := do(client.client, fmt.Sprintf("%s/api/calculate", client.host), "POST", rq, client.ctx)

	var response models.Response
	err = json.Unmarshal(data, &response)

	if err != nil {
		return nil, err
	}

	return &response, nil
}

func do(client *http.Client, url string, method string, v interface{}, ctx context.Context) ([]byte, error) {
	var data []byte
	var err error
	if v != nil {
		data, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		if resp != nil {
			return nil, err
		} else {
			return nil, err
		}
	}

	defer resp.Body.Close()
	rsBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return rsBody, errors.New(string(rsBody))
	}
	return rsBody, nil
}

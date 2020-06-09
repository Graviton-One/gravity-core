package score_calculator

import (
	"context"
	"encoding/json"
	"fmt"
	"gravity-hub/gh-node/api"
)

type Client struct {
	url string
}

func (ghClient *Client) Get(ctx context.Context) (ResponseBlock, error) {
	data, err := api.Do(ghClient.client, fmt.Sprintf("%s/block", ghClient.nodeUrl), api.POST, nil, ctx)
	if err != nil {
		return ResponseBlock{}, err
	}

	var response ResponseBlock
	err = json.Unmarshal(data, &response)
	if err != nil {
		return ResponseBlock{}, err
	}

	return response, nil
}

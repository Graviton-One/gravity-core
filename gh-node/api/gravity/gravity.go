package gravity

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gh-node/api"
	"gh-node/transaction"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	nodeUrl string
	client  *http.Client
}

const (
	CodeError = 1
)

var (
	KeyNotFound = errors.New("key not found")
)

func NewClient(nodeUrl string) *Client {
	return &Client{nodeUrl: nodeUrl, client: &http.Client{Timeout: time.Second * 3}}
}

func (ghClient *Client) SendTx(transaction *transaction.Transaction, ctx context.Context) error {
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	rs, err := api.Do(ghClient.client, fmt.Sprintf("%s/broadcast_tx_commit?tx=\"%s\"", ghClient.nodeUrl, strings.ReplaceAll(string(txBytes), "\"", "\\\"")), api.POST, nil, ctx)
	var response ResponseTx
	err = json.Unmarshal(rs, &response)
	if err != nil {
		return err
	}
	if response.Result.CheckTx.Code == CodeError {
		return errors.New(response.Result.CheckTx.Info)
	}
	return err
}

func (ghClient *Client) GetKey(key string, ctx context.Context) ([]byte, error) {
	data, err := api.Do(ghClient.client, fmt.Sprintf(ghClient.nodeUrl+"/abci_query?path=\"%s\"&data=\"%s\"", "key", key), api.GET, nil, ctx)
	if err != nil {
		return nil, err
	}

	var response ResponseGet
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	if response.Result.Response.Code == CodeError {
		return nil, KeyNotFound
	}

	value, err := base64.StdEncoding.DecodeString(response.Result.Response.Value)
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (ghClient *Client) GetByPrefix(prefix string, ctx context.Context) (map[string][]byte, error) {
	data, err := api.Do(ghClient.client, fmt.Sprintf(ghClient.nodeUrl+"/abci_query?path=\"%s\"&data=\"%s\"", "prefix", prefix), api.GET, nil, ctx)
	if err != nil {
		return nil, err
	}

	var response ResponseGet
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	if response.Result.Response.Code == CodeError {
		return nil, errors.New(response.Result.Response.Info)
	}
	value, err := base64.StdEncoding.DecodeString(response.Result.Response.Value)
	if err != nil {
		return nil, err
	}

	values := make(map[string][]byte)
	err = json.Unmarshal(value, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

func (ghClient *Client) GetBlock(ctx context.Context) (ResponseBlock, error) {
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

package gravity

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"gravity-hub/common/transactions"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

type Client struct {
	HttpClient *rpchttp.HTTP
}

const (
	ErrorCode = 500
)

var (
	KeyNotFound = errors.New("key not found")
)

func NewClient(rpcAddress string) (*Client, error) {
	client, err := rpchttp.New(rpcAddress, "/websocket")
	if err != nil {
		return nil, err
	}
	return &Client{HttpClient: client}, nil
}

func (ghClient *Client) SendTx(transaction *transactions.Transaction) error {
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	rs, err := ghClient.HttpClient.BroadcastTxCommit(txBytes)
	if err != nil {
		return err
	}
	if rs.CheckTx.Code == ErrorCode {
		return errors.New(rs.CheckTx.Info)
	} else if rs.DeliverTx.Code == ErrorCode {
		return errors.New(rs.DeliverTx.Info)
	}
	return err
}

func (ghClient *Client) GetKey(key string) ([]byte, error) {
	rs, err := ghClient.HttpClient.ABCIQuery("key", []byte(key))
	if err != nil {
		return nil, err
	}

	if rs.Response.Code == ErrorCode {
		return nil, KeyNotFound
	}

	value, err := base64.StdEncoding.DecodeString(string(rs.Response.Value))
	if err != nil {
		return nil, err
	}

	return value, nil
}

func (ghClient *Client) GetByPrefix(prefix string) (map[string][]byte, error) {
	rs, err := ghClient.HttpClient.ABCIQuery("prefix", []byte(prefix))
	if err != nil {
		return nil, err
	}

	if rs.Response.Code == ErrorCode {
		return nil, KeyNotFound
	}

	rsValue, err := base64.StdEncoding.DecodeString(string(rs.Response.Value))
	if err != nil {
		return nil, err
	}

	values := make(map[string][]byte)
	err = json.Unmarshal(rsValue, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}

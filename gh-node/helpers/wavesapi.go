package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Gravity-Tech/gravity-core/gh-node/helpers/state"
)

const (
	GetTxPath              = "/transactions/info"
	GetUnconfirmedTxByPath = "/transactions/unconfirmed/info"
	GetStateByAddressPath  = "/addresses/data"

	WaitCount = 10
)

type Node struct {
	nodeUrl string
	apiKey  string
}

func New(nodeUrl string, apiKey string) Node {
	return Node{nodeUrl: nodeUrl, apiKey: apiKey}
}

func (node *Node) GetStateByAddressAndKey(address string, key string) (*state.State, error) {
	rsBody, _, err := sendRequest("GET", node.nodeUrl+GetStateByAddressPath+"/"+address+"?key="+key, nil, "")

	if err != nil {
		return nil, err
	}

	var states []state.State
	if err := json.Unmarshal(rsBody, &states); err != nil {
		return nil, err
	}
	if len(states) == 0 {
		return nil, nil
	}
	return &states[0], nil
}

func (node *Node) GetTxById(id string) (Transaction, error) {
	rsBody, _, err := sendRequest("GET", node.nodeUrl+GetTxPath+"/"+id, nil, "")
	if err != nil {
		return Transaction{}, err
	}

	return Unmarshal(rsBody)
}

func (node *Node) IsUnconfirmedTx(id string) (bool, error) {
	_, code, err := sendRequest("GET", node.nodeUrl+GetUnconfirmedTxByPath+"/"+id, nil, "")
	if err != nil && code != 404 {
		return true, err
	}

	return code == 200, nil
}

func (node *Node) WaitTx(id string) <-chan error {
	out := make(chan error)
	go func() {
		defer close(out)
		for i := 0; i < WaitCount; i++ {
			un, err := node.IsUnconfirmedTx(id)
			if err != nil {
				out <- err
				break
			}

			if un == false {
				tx, err := node.GetTxById(id)
				if err != nil {
					out <- err
				}
				if tx.ID == "" {
					out <- errors.New("transaction not found")
				} else {
					out <- nil
				}
				break
			}

			if i == (WaitCount - 1) {
				out <- errors.New("transaction not found")
				break
			}

			time.Sleep(time.Second)
		}
	}()
	return out
}

func sendRequest(method string, url string, rqBody []byte, apiKey string) ([]byte, int, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(rqBody))
	req.Header.Add("content-type", "application/json")
	if apiKey != "" {
		req.Header.Add("X-API-Key", apiKey)
	}

	if err != nil {
		return nil, 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		if resp != nil {
			return nil, resp.StatusCode, err
		} else {
			return nil, 520, err
		}
	}

	defer resp.Body.Close()
	rsBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return rsBody, resp.StatusCode, errors.New(string(rsBody))
	}
	return rsBody, resp.StatusCode, nil
}

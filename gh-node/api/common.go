package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

func Do(client *http.Client, url string, method Method, v interface{}, ctx context.Context) ([]byte, error) {
	var data []byte
	var err error
	if v != nil {
		data, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequestWithContext(ctx, string(method), url, bytes.NewBuffer(data))
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

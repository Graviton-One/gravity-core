package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"
)

const (
	ExtractPath   = "extract"
	InfoPath      = "info"
	AggregatePath = "aggregate"

	String DataType = "string"
	Int64  DataType = "int64"
	Base64 DataType = "base64"
)

var (
	NotFoundErr = errors.New("data not found")
)

type DataType string
type Data struct {
	Type  DataType
	Value string
}

type Client struct {
	hostUrl string
}

func New(hostUrl string) *Client {
	return &Client{
		hostUrl: hostUrl,
	}
}

func (client *Client) Extract(ctx context.Context) (*Data, error) {
	rs, err := client.do(ExtractPath, http.MethodGet, nil, ctx)
	if err != nil {
		return nil, err
	}

	var result Data
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) Aggregate(values []Data, ctx context.Context) (*Data, error) {
	rs, err := client.do(AggregatePath, http.MethodPost, values, ctx)
	if err != nil {
		return nil, err
	}
	zap.L().Sugar().Debugf("Aggragate response: %s", string(rs))
	var result Data
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) do(route string, method string, rqBody interface{}, ctx context.Context) ([]byte, error) {
	rqUrl := fmt.Sprintf("%v/%v", client.hostUrl, route)
	zap.L().Sugar().Debugf("Request URL: %s", rqUrl)
	var buf *bytes.Buffer
	var req *http.Request
	var err error
	if rqBody != nil {
		b, err := json.Marshal(&rqBody)
		if err != nil {
			return nil, err
		}
		zap.L().Sugar().Debugf("Reuest Body: %s", string(b))
		buf = bytes.NewBuffer(b)
		req, err = http.NewRequestWithContext(ctx, method, rqUrl, buf)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequestWithContext(ctx, method, rqUrl, nil)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	zap.L().Sugar().Debug("Response: ", response)
	if response.StatusCode == 404 {
		return nil, NotFoundErr
	}

	rsBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return rsBody, nil
}

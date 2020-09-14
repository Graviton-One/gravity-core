package extractor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	GetExtractedData = "extracted"
	GetExtractorInfo = "info"
	GetAggregated    = "aggregate"
)

type Client struct {
	hostUrl string
}

func New(hostUrl string) *Client {
	return &Client{
		hostUrl: hostUrl,
	}
}

func (client *Client) ExtractorInfo(ctx context.Context) (*Info, error) {
	rs, err := client.do(GetExtractorInfo, http.MethodGet, nil, ctx)
	if err != nil {
		return nil, err
	}

	var result Info
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) Aggregate(values []interface{}, ctx context.Context) (interface{}, error) {
	rs, err := client.do(GetAggregated, http.MethodPost, values, ctx)
	if err != nil {
		return nil, err
	}

	var result DataRs
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return result.Value, nil
}

func (client *Client) Extract(ctx context.Context) (interface{}, error) {
	rs, err := client.do(GetExtractedData, http.MethodGet, nil, ctx)
	if err != nil {
		return nil, err
	}

	var result DataRs
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return result.Value, nil
}

func (client *Client) do(route string, method string, rqBody interface{}, ctx context.Context) ([]byte, error) {
	rqUrl := fmt.Sprintf("%v/%v", client.hostUrl, route)

	var buf *bytes.Buffer
	var req *http.Request
	var err error
	if rqBody != nil {
		b, err := json.Marshal(&rqBody)
		if err != nil {
			return nil, err
		}
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
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	rsBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return rsBody, nil
}

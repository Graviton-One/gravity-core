package extractor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/controller"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/model"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/router"
)

type Client struct {
	hostUrl string
}

func New(hostUrl string) *Client {
	return &Client{
		hostUrl: hostUrl,
	}
}

func (client *Client) ExtractorInfo() (*model.ExtractorInfo, error) {
	rs, err := client.do(router.GetExtractorInfo, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var result model.ExtractorInfo
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *Client) Aggregate(values []interface{}) (interface{}, error) {
	rs, err := client.do(router.GetAggregated, http.MethodPost, controller.AggregationRequestBody{
		Values: values,
	})
	if err != nil {
		return nil, err
	}

	var result controller.DataRs
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return result.Value, nil
}

func (client *Client) Extract() (interface{}, error) {
	rs, err := client.do(router.GetRawData, http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	var result controller.DataRs
	err = json.Unmarshal(rs, &result)
	if err != nil {
		return nil, err
	}

	return result.Value, nil
}

func (client *Client) do(route string, method string, rqBody interface{}) ([]byte, error) {
	rqUrl := fmt.Sprintf("%v/%v", client.hostUrl, route)

	var buf *bytes.Buffer
	if rqBody != nil {
		b, err := json.Marshal(&rqBody)
		if err != nil {
			return nil, err
		}
		buf = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, rqUrl, buf)
	if err != nil {
		return nil, err
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

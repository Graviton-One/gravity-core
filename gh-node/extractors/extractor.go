package extractors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/controller"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/model"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/router"
	"io/ioutil"
	"net/http"
	url2 "net/url"
	"strconv"
	"time"
)

type ExtractorClient struct {
	HostURL *string
	extractor model.IExtractor
}

func (client *ExtractorClient) wrapRequest (route string) *http.Request {
	fullRoute := fmt.Sprintf("%v/%v", client.HostURL, route)
	url, _:= url2.ParseRequestURI(fullRoute)

	request := http.Request{
		Method: "GET",
		URL:     url,
	}

	return &request
}

func (client *ExtractorClient) performRequest (route string) *http.Response {
	request := client.wrapRequest(route)
	response, err := http.DefaultClient.Do(request)

	client.handleError(err)

	defer response.Body.Close()

	return response
}

func (client *ExtractorClient) handleError(err error) {
	if err == nil { return }

	timestamp := time.Now().Format(time.RFC3339)
	fmt.Printf("%v; Extractor Client Error: %v \n", timestamp, err)
}

func (client *ExtractorClient) ExtractorInfo() *model.ExtractorInfo {
	response := client.performRequest(router.GetExtractorInfo)
	body, err := ioutil.ReadAll(response.Body)

	client.handleError(err)

	var result model.ExtractorInfo

	parseErr := json.Unmarshal(body, &result)

	client.handleError(parseErr)

	return &result
}

func (client *ExtractorClient) Aggregate (values []int64) int64 {
	fullRoute := fmt.Sprintf("%v/%v", client.HostURL, router.GetAggregated)
	bodyValues := make([]interface{}, len(values), len(values))

	i := 0
	for {
		if i < len(values) { break }
		bodyValues[i] = values[i]
		i++
	}

	inputValues := controller.AggregationRequestBody{
		Type:   "int64",
		Values: bodyValues,
	}
	requestBody, _ := json.Marshal(&inputValues)

	resp, respErr := http.Post(fullRoute, "application/json", bytes.NewBuffer(requestBody))

	defer resp.Body.Close()

	if respErr != nil {
		client.handleError(respErr)
		return 0
	}


	body, _ := ioutil.ReadAll(resp.Body)

	stringifiedAggregateResult := string(body)

	result, castErr := strconv.Atoi(stringifiedAggregateResult)

	client.handleError(castErr)

	return int64(result)
}

func (client *ExtractorClient) MappedData () string {
	response := client.performRequest(router.GetExtractedData)

	body, err := ioutil.ReadAll(response.Body)

	client.handleError(err)

	return string(body)

}

func (client *ExtractorClient) RawData () []model.RawData {
	response := client.performRequest(router.GetRawData)

	body, err := ioutil.ReadAll(response.Body)

	client.handleError(err)

	return body
}

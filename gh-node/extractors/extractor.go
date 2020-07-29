package extractors

import (
	"encoding/json"
	"fmt"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/model"
	"github.com/Gravity-Tech/gravity-node-data-extractor/v2/router"
	"io/ioutil"
	"net/http"
	url2 "net/url"
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

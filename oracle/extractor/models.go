package extractor

type Info struct {
	Description string `json:"description"`
	DataFeedTag string `json:"datafeedtag"`
}

type DataRs struct {
	Value interface{} `json:"value"`
}

type AggregationRequestBody struct {
	Type   string        `json:"type"`
	Values []interface{} `json:"values"`
}

package extractors

type Extractor interface {
	Extract() (interface{}, error)
	Aggregate(values []interface{}) (interface{}, error)
}

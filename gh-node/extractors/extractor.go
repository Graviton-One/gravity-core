package extractors

type PriceExtractor interface {
	GetData() ([]byte, error)
	Aggregate(values [][]byte) ([]byte, error)
}

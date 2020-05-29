package extractors

type PriceExtractor interface {
	PriceNow() (float64, error)
	PriceNowByPair(pair string) (float64, error)
}

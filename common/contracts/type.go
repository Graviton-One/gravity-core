package contracts

type ExtractorType uint8

const (
	Int64Type ExtractorType = iota
	StringType
	BytesType
)

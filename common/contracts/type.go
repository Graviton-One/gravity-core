package contracts

import (
	"errors"
	"strings"
)

type ExtractorType uint8

var (
	ErrParseExtractorType = errors.New("invalid parse extractor type")
)

const (
	Int64Type ExtractorType = iota
	StringType
	BytesType
)

func ParseExtractorType(extractorType string) (ExtractorType, error) {
	switch strings.ToLower(extractorType) {
	case "int64":
		return Int64Type, nil
	case "string":
		return StringType, nil
	case "bytes":
		return BytesType, nil
	default:
		return 0, ErrParseExtractorType
	}
}

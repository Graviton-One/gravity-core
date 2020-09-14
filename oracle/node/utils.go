package node

import (
	"encoding/binary"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
)

func toBytes(value interface{}, dataType contracts.ExtractorType) []byte {
	switch dataType {
	case contracts.Int64Type:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(value.(float64)))
		return b[:]
	case contracts.StringType:
		return []byte(value.(string))
	case contracts.BytesType:
		return value.([]byte)
	}
	return nil
}

func fromBytes(value []byte, extractorType contracts.ExtractorType) interface{} {
	switch extractorType {
	case contracts.Int64Type:
		return binary.BigEndian.Uint64(value)
	case contracts.StringType:
		return string(value)
	case contracts.BytesType:
		return value
	}

	return nil
}

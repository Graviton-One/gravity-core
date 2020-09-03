package node

import (
	"encoding/binary"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
)

func toBytes(value interface{}) []byte {
	switch v := value.(type) {
	case int64:
		var b []byte
		binary.BigEndian.PutUint64(b, uint64(v))
		return b
	case string:
		return []byte(v)
	case []byte:
		return v
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

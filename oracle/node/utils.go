package node

import "encoding/binary"

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
func fromBytes(value []byte, extractorType ExtractorType) interface{} {
	switch extractorType {
	case Int64Type:
		return binary.BigEndian.Uint64(value)
	case StringType:
		return string(value)
	case BytesType:
		return value
	}

	return nil
}

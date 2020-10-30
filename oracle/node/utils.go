package node

import (
	"encoding/base64"
	"encoding/binary"
	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"strconv"
)

func toBytes(data *extractor.Data, dataType abi.ExtractorType) []byte {
	switch dataType {
	case abi.Int64Type:
		v, _ := strconv.ParseInt(data.Value, 10,64)
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(v))
		return b[:]
	case abi.StringType:
		return []byte(data.Value)
	case abi.BytesType:
		b, _ := base64.StdEncoding.DecodeString(data.Value)
		return b
	}
	return nil
}

func fromBytes(value []byte, extractorType abi.ExtractorType) *extractor.Data {
	switch extractorType {
	case abi.Int64Type:
		v := binary.BigEndian.Uint64(value)
		return &extractor.Data{
			Type:  extractor.Int64,
			Value: strconv.FormatInt(int64(v), 10),
		}
	case abi.StringType:
		return &extractor.Data{
			Type:  extractor.String,
			Value: string(value),
		}
	case abi.BytesType:
		return &extractor.Data{
			Type:  extractor.Base64,
			Value: base64.StdEncoding.EncodeToString(value),
		}
	}

	return nil
}

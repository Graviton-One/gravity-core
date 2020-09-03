package query

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/Gravity-Tech/gravity-core/common/storage"
)

type Path string

const (
	OracleByValidatorPath Path = "oraclesByValidator"
	OracleByNebulaPath    Path = "oraclesByNebula"
	BftOracleByNebulaPath Path = "bftOraclesByNebula"
	RoundHeightPath       Path = "roundHeight"
	CommitHashPath        Path = "commitHash"
	RevealPath            Path = "reveal"
	ResultPath            Path = "result"
	ResultsPath           Path = "results"
	NebulaePath           Path = "nebulae"
)

var (
	ErrInvalidPath = errors.New("invalid path")
)

func Query(store *storage.Storage, path string, rq []byte) ([]byte, error) {
	var value interface{}
	switch Path(path) {
	case OracleByValidatorPath:
		v, err := oraclesByValidator(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case OracleByNebulaPath:
		v, err := oraclesByNebula(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case RoundHeightPath:
		v, err := roundHeight(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case CommitHashPath:
		v, err := commitHash(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case RevealPath:
		v, err := reveal(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case ResultPath:
		v, err := result(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case BftOracleByNebulaPath:
		v, err := bftOraclesByNebula(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case ResultsPath:
		v, err := results(store, rq)
		if err != nil {
			return nil, err
		}
		value = v
	case NebulaePath:
		v, err := nebulae(store)
		if err != nil {
			return nil, err
		}
		value = v
	default:
		return nil, ErrInvalidPath
	}

	b, err := toBytes(value)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func toBytes(v interface{}) ([]byte, error) {
	var err error
	var b []byte
	switch v.(type) {
	case uint64:
		binary.BigEndian.PutUint64(b, v.(uint64))
	case []byte:
		b = v.([]byte)
	default:
		b, err = json.Marshal(v)
	}

	if err != nil {
		return nil, err
	}

	return b, nil
}

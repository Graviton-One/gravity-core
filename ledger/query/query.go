package query

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/config"
)

type Path string

const (
	OracleByValidatorPath      Path = "oraclesByValidator"
	OracleByNebulaPath         Path = "oraclesByNebula"
	BftOracleByNebulaPath      Path = "bftOraclesByNebula"
	RoundHeightPath            Path = "roundHeight"
	CommitHashPath             Path = "commitHash"
	RevealPath                 Path = "reveal"
	RevealsPath                Path = "reveals"
	ResultPath                 Path = "result"
	ResultsPath                Path = "results"
	NebulaePath                Path = "nebulae"
	NebulaInfoPath             Path = "nebula_info"
	LastRoundApprovedPath      Path = "lastRoundApproved"
	ConsulsPath                Path = "consuls"
	ConsulsCandidatePath       Path = "consulsCandidate"
	SignNewConsulsByConsulPath Path = "signNewConsulsByConsul"
	SignNewOraclesByConsulPath Path = "signNewOraclesByConsul"
	NebulaOraclesIndexPath     Path = "nebulaOraclesIndex"
	AllValidatorsPath          Path = "allValidators"
	ValidatorDetailsPath       Path = "validatorDetails"
	NebulaCustomParams         Path = "nebulaCustomParams"
)

var (
	ErrInvalidPath   = errors.New("invalid path")
	ErrValueNotFound = errors.New("value not found")
)

func Query(store *storage.Storage, path string, rq []byte, validatorDetails *config.ValidatorDetails) ([]byte, error) {
	var value interface{}
	var err error
	switch Path(path) {
	case OracleByValidatorPath:
		value, err = oraclesByValidator(store, rq)
	case OracleByNebulaPath:
		value, err = oraclesByNebula(store, rq)
	case RoundHeightPath:
		value, err = roundHeight(store, rq)
	case CommitHashPath:
		value, err = commitHash(store, rq)
	case RevealPath:
		value, err = reveal(store, rq)
	case RevealsPath:
		value, err = reveals(store, rq)
	case ResultPath:
		value, err = result(store, rq)
	case BftOracleByNebulaPath:
		value, err = bftOraclesByNebula(store, rq)
	case ResultsPath:
		value, err = results(store, rq)
	case NebulaePath:
		value, err = nebulae(store)
	case NebulaInfoPath:
		value, err = nebulaInfo(store, rq)
	case ConsulsPath:
		value, err = consuls(store)
	case ConsulsCandidatePath:
		value, err = consulsCandidate(store)
	case SignNewConsulsByConsulPath:
		value, err = signNewConsulsByConsul(store, rq)
	case SignNewOraclesByConsulPath:
		value, err = signNewOraclesByConsul(store, rq)
	case NebulaOraclesIndexPath:
		value, err = nebulaOraclesIndex(store, rq)
	case LastRoundApprovedPath:
		value, err = store.LastRoundApproved()
	case AllValidatorsPath:
		value, err = allValidators(store, rq)
	case ValidatorDetailsPath:
		value, err = validatorDetails.Bytes()
	case NebulaCustomParams:
		value, err = nebulaCustomParams(store, rq)
	default:
		return nil, ErrInvalidPath
	}

	if err != nil && err != storage.ErrKeyNotFound {
		return nil, err
	} else if err == storage.ErrKeyNotFound {
		return nil, ErrValueNotFound
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
		b = make([]byte, 8, 8)
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

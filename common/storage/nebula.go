package storage

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type NebulaMap map[string]NebulaInfo
type NebulaInfo struct {
	MaxPulseCountInBlock uint64
	MinScore             uint64
	ChainType            account.ChainType
	Owner                account.ConsulPubKey
}

type NebulaCustomParamsMap map[string]NebulaCustomParams
type NebulaCustomParams map[string]interface{}

func parseNebulaInfoKey(value []byte) (account.NebulaId, error) {
	hex := []byte(strings.Split(string(value), Separator)[2])
	key, err := hexutil.Decode(string(hex))
	if err != nil {
		return [32]byte{}, err
	}
	var pubKey account.NebulaId
	copy(pubKey[:], key[:])
	return pubKey, nil
}
func formNebulaInfoKey(nebulaId account.NebulaId) []byte {
	return formKey(string(NebulaInfoKey), hexutil.Encode(nebulaId[:]))
}

func (storage *Storage) Nebulae() (NebulaMap, error) {
	it := storage.txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(NebulaInfoKey)
	nebulaeInfo := make(NebulaMap)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		err := item.Value(func(v []byte) error {
			var nebulaInfo NebulaInfo
			err := json.Unmarshal(v, &nebulaInfo)
			if err != nil {
				return err
			}
			pubKey, err := parseNebulaInfoKey(k)
			if err != nil {
				return err
			}
			nebulaeInfo[pubKey.ToString(nebulaInfo.ChainType)] = nebulaInfo
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return nebulaeInfo, nil
}
func (storage *Storage) NebulaInfo(nebulaId account.NebulaId) (*NebulaInfo, error) {
	b, err := storage.getValue(formNebulaInfoKey(nebulaId))
	if err != nil {
		return nil, err
	}

	var nebulae NebulaInfo
	err = json.Unmarshal(b, &nebulae)
	if err != nil {
		return nil, err
	}

	return &nebulae, err
}

func (storage *Storage) DropNebula(nebulaId account.NebulaId) error {
	return storage.dropValue(formNebulaInfoKey(nebulaId))
}

func (storage *Storage) SetNebula(nebulaId account.NebulaId, info NebulaInfo) error {
	zap.L().Debug("Setting nebula!!!!")
	return storage.setValue(formNebulaInfoKey(nebulaId), &info)
}

func parseNebulaCustomParamsKey(value []byte) (account.NebulaId, error) {
	hex := []byte(strings.Split(string(value), Separator)[2])
	key, err := hexutil.Decode(string(hex))
	if err != nil {
		return [32]byte{}, err
	}
	var pubKey account.NebulaId
	copy(pubKey[:], key[:])
	return pubKey, nil
}
func formNebulaCustomParamsKey(nebulaId account.NebulaId) []byte {
	return formKey(string(NebulaCustomParamsKey), hexutil.Encode(nebulaId[:]))
}

func (storage *Storage) NebulaCustomParams(nebulaId account.NebulaId) (*NebulaCustomParams, error) {
	b, err := storage.getValue(formNebulaCustomParamsKey(nebulaId))
	if err != nil {
		return nil, err
	}

	var nebulae NebulaCustomParams
	err = json.Unmarshal(b, &nebulae)
	if err != nil {
		return nil, err
	}

	return &nebulae, err
}

func (storage *Storage) DropNebulaCustomParams(nebulaId account.NebulaId) error {
	return storage.dropValue(formNebulaCustomParamsKey(nebulaId))
}

func (storage *Storage) SetNebulaCustomParams(nebulaId account.NebulaId, customParams NebulaCustomParams) error {
	return storage.setValue(formNebulaCustomParamsKey(nebulaId), &customParams)
}

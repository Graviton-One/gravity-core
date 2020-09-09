package storage

import (
	"encoding/json"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type NebulaMap map[account.NebulaId]NebulaInfo
type NebulaInfo struct {
	MaxPulseCountInBlock uint64
	MinScore             uint64
	ChainType            account.ChainType
	Owner                account.ConsulPubKey
}

func (storage *Storage) Nebulae() (NebulaMap, error) {
	b, err := storage.getValue([]byte(Nebulae))
	if err == ErrKeyNotFound {
		return make(NebulaMap), nil
	}
	if err != nil {
		return nil, err
	}

	var nebulae NebulaMap
	err = json.Unmarshal(b, &nebulae)
	if err != nil {
		return nil, err
	}

	return nebulae, err
}
func (storage *Storage) SetNebula(nebulaId account.NebulaId, info NebulaInfo) error {
	nebulae, err := storage.Nebulae()
	if err != nil {
		return err
	}

	nebulae[nebulaId] = info

	return storage.setValue([]byte(Nebulae), nebulae)
}

package storage

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type OraclesByTypeMap map[account.ChainType]account.OraclesPubKey
type OraclesMap map[string]account.ChainType

func formBftOraclesByNebulaKey(nebulaId account.NebulaId) []byte {
	return formKey(string(BftOraclesByNebulaKey), hexutil.Encode(nebulaId[:]))
}
func formNebulaOraclesIndexKey(nebulaId account.NebulaId) []byte {
	return formKey(string(NebulaOraclesIndexKey), hexutil.Encode(nebulaId[:]))
}
func formSignOraclesByConsulKey(consulPubKey account.ConsulPubKey, nebulaId account.NebulaId, roundId int64) []byte {
	return formKey(string(SignOraclesResultByConsulKey), hexutil.Encode(consulPubKey[:]), hexutil.Encode(nebulaId[:]), fmt.Sprintf("%d", roundId))
}
func formOraclesByConsulKey(consulPubKey account.ConsulPubKey) []byte {
	return formKey(string(OraclesByValidatorKey), hexutil.Encode(consulPubKey[:]))
}
func formOraclesByNebulaKey(nebulaId account.NebulaId) []byte {
	return formKey(string(OraclesByNebulaKey), hexutil.Encode(nebulaId[:]))
}
func formNebulaeByOracleKey(pubKey account.OraclesPubKey) []byte {
	return formKey(string(NebulaeByOracleKey), hexutil.Encode(pubKey[:]))
}

func (storage *Storage) OraclesByNebula(nebulaId account.NebulaId) (OraclesMap, error) {
	b, err := storage.getValue(formOraclesByNebulaKey(nebulaId))
	if err != nil {
		return nil, err
	}

	var oraclesByNebula OraclesMap
	err = json.Unmarshal(b, &oraclesByNebula)
	if err != nil {
		return oraclesByNebula, err
	}

	return oraclesByNebula, err
}
func (storage *Storage) SetOraclesByNebula(nebulaAddress account.NebulaId, oracles OraclesMap) error {
	return storage.setValue(formOraclesByNebulaKey(nebulaAddress), oracles)
}

func (storage *Storage) NebulaeByOracle(pubKey account.OraclesPubKey) ([]account.NebulaId, error) {
	b, err := storage.getValue(formNebulaeByOracleKey(pubKey))
	if err != nil {
		return nil, err
	}

	var nebulae []account.NebulaId
	err = json.Unmarshal(b, &nebulae)
	if err != nil {
		return nebulae, err
	}

	return nebulae, err
}
func (storage *Storage) SetNebulaeByOracle(pubKey account.OraclesPubKey, nebulae []account.NebulaId) error {
	return storage.setValue(formNebulaeByOracleKey(pubKey), nebulae)
}

func (storage *Storage) NebulaOraclesIndex(nebulaAddress account.NebulaId) (uint64, error) {
	b, err := storage.getValue(formNebulaOraclesIndexKey(nebulaAddress))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(b), nil
}
func (storage *Storage) SetNebulaOraclesIndex(nebulaAddress account.NebulaId, index uint64) error {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], index)
	err := storage.setValue(formNebulaOraclesIndexKey(nebulaAddress), b[:])
	if err != nil {
		return err
	}

	return err
}

func (storage *Storage) OraclesByConsul(pubKey account.ConsulPubKey) (OraclesByTypeMap, error) {
	b, err := storage.getValue(formOraclesByConsulKey(pubKey))
	if err != nil {
		return nil, err
	}

	var oracles OraclesByTypeMap
	err = json.Unmarshal(b, &oracles)
	if err != nil {
		return oracles, err
	}

	return oracles, err
}
func (storage *Storage) SetOraclesByConsul(pubKey account.ConsulPubKey, oracles OraclesByTypeMap) error {
	return storage.setValue(formOraclesByConsulKey(pubKey), oracles)
}

func (storage *Storage) SignOraclesByConsul(pubKey account.ConsulPubKey, nebulaId account.NebulaId, roundId int64) ([]byte, error) {
	key := formSignOraclesByConsulKey(pubKey, nebulaId, roundId)
	item, err := storage.txn.Get(key)
	if err != nil {
		return nil, err
	}

	b, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	return b, err
}
func (storage *Storage) SetSignOracles(pubKey account.ConsulPubKey, nebulaId account.NebulaId, roundId int64, sign []byte) error {
	return storage.setValue(formSignOraclesByConsulKey(pubKey, nebulaId, roundId), sign)
}

func (storage *Storage) BftOraclesByNebula(nebulaId account.NebulaId) (OraclesMap, error) {
	b, err := storage.getValue(formBftOraclesByNebulaKey(nebulaId))
	if err != nil {
		return nil, err
	}

	var oraclesByNebula OraclesMap
	err = json.Unmarshal(b, &oraclesByNebula)
	if err != nil {
		return oraclesByNebula, err
	}

	return oraclesByNebula, err
}
func (storage *Storage) SetBftOraclesByNebula(nebulaId account.NebulaId, oracles OraclesMap) error {
	return storage.setValue(formBftOraclesByNebulaKey(nebulaId), oracles)
}

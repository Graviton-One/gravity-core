package storage

import (
	"encoding/json"

	"github.com/Gravity-Tech/proof-of-concept/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type OraclesByTypeMap map[account.ChainType]account.OraclesPubKey
type OraclesMap map[account.OraclesPubKey]bool

func formOraclesByValidatorKey(validator account.PubKey) []byte {
	return formKey(string(OraclesByValidatorKey), hexutil.Encode(validator[:]))
}

func formOraclesByNebulaKey(nebulaAddress []byte) []byte {
	return formKey(string(OraclesByNebulaKey), hexutil.Encode(nebulaAddress))
}

func (storage *Storage) OraclesByNebula(nebulaAddress []byte) (OraclesMap, error) {
	b, err := storage.getValue(formOraclesByNebulaKey(nebulaAddress))
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

func (storage *Storage) SetOraclesByNebula(nebulaAddress []byte, oracles OraclesMap) error {
	return storage.setValue(formOraclesByNebulaKey(nebulaAddress), oracles)
}

func (storage *Storage) OraclesByValidator(validator account.PubKey) (OraclesByTypeMap, error) {
	b, err := storage.getValue(formOraclesByValidatorKey(validator))
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

func (storage *Storage) SetOraclesByValidator(validator account.PubKey, oracles OraclesByTypeMap) error {
	return storage.setValue(formOraclesByValidatorKey(validator), oracles)
}

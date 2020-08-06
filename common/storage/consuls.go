package storage

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type Consul struct {
	Validator account.ValidatorPubKey
	Value     uint64
}

func formConsulSignKey(validatorAddress account.ValidatorPubKey, chainType account.ChainType, roundId int64) []byte {
	prefix := ""
	switch chainType {
	case account.Waves:
		prefix = "waves"
	case account.Ethereum:
		prefix = "ethereum"
	}
	return formKey(string(ConsulsSignKey), hexutil.Encode(validatorAddress[:]), prefix, fmt.Sprintf("%d", roundId))
}

func (storage *Storage) Consuls() ([]Consul, error) {
	var consuls []Consul

	key := []byte(ConsulsKey)
	item, err := storage.txn.Get(key)
	if err != nil {
		return nil, err
	}

	b, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &consuls)
	if err != nil {
		return nil, err
	}

	return consuls, err
}

func (storage *Storage) SetConsuls(consuls []Consul) error {
	return storage.setValue([]byte(ConsulsKey), consuls)
}

func (storage *Storage) ConsulSign(validatorAddress account.ValidatorPubKey, chainType account.ChainType, roundId int64) ([]byte, error) {
	key := formConsulSignKey(validatorAddress, chainType, roundId)
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

func (storage *Storage) SetConsulSign(validatorAddress account.ValidatorPubKey, chainType account.ChainType, roundId int64, sign []byte) error {
	return storage.setValue(formConsulSignKey(validatorAddress, chainType, roundId), sign)
}

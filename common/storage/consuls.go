package storage

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type Consul struct {
	PubKey account.ConsulPubKey
	Value  uint64
}

func formSignConsulsResultByConsulKey(pubKey account.ConsulPubKey, chainType account.ChainType, roundId int64) []byte {
	prefix := ""
	switch chainType {
	case account.Waves:
		prefix = "waves"
	case account.Ethereum:
		prefix = "ethereum"
	}
	return formKey(string(SignConsulsResultByConsulKey), hexutil.Encode(pubKey[:]), prefix, fmt.Sprintf("%d", roundId))
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

func (storage *Storage) PrevConsuls() ([]Consul, error) {
	var consuls []Consul

	key := []byte(PrevConsulsKey)
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
func (storage *Storage) SetPrevConsuls(consuls []Consul) error {
	return storage.setValue([]byte(PrevConsulsKey), consuls)
}

func (storage *Storage) SignConsulsResultByConsul(consulPubKey account.ConsulPubKey, chainType account.ChainType, roundId int64) ([]byte, error) {
	key := formSignConsulsResultByConsulKey(consulPubKey, chainType, roundId)
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
func (storage *Storage) SetSignConsulsResult(consulsPubKey account.ConsulPubKey, chainType account.ChainType, roundId int64, sign []byte) error {
	return storage.setValue(formSignConsulsResultByConsulKey(consulsPubKey, chainType, roundId), sign)
}

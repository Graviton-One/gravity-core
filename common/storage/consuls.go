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

func formSignConsulsByConsulKey(pubKey account.ConsulPubKey, chainType account.ChainType, roundId int64) []byte {
	prefix := ""
	switch chainType {
	case account.Waves:
		prefix = "waves"
	case account.Ethereum:
		prefix = "ethereum"
	case account.Binance:
		prefix = "bsc"
	case account.Fantom:
		prefix = "ftm"
	case account.Heco:
		prefix = "heco"
	case account.Ergo:
		prefix = "ergo"
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

func (storage *Storage) ConsulsCandidate() ([]Consul, error) {
	var consuls []Consul

	key := []byte(ConsulsCandidateKey)
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
func (storage *Storage) SetConsulsCandidate(consuls []Consul) error {
	return storage.setValue([]byte(ConsulsCandidateKey), consuls)
}

func (storage *Storage) SignConsulsByConsul(consulPubKey account.ConsulPubKey, chainType account.ChainType, roundId int64) ([]byte, error) {
	key := formSignConsulsByConsulKey(consulPubKey, chainType, roundId)
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
func (storage *Storage) SetSignConsuls(consulsPubKey account.ConsulPubKey, chainType account.ChainType, roundId int64, sign []byte) error {
	return storage.setValue(formSignConsulsByConsulKey(consulsPubKey, chainType, roundId), sign)
}

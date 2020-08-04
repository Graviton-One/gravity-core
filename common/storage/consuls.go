package storage

import (
	"encoding/json"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type Consul struct {
	Validator account.PubKey
	Value     uint64
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

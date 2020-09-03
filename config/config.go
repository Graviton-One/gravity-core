package config

import (
	"encoding/json"
	"io/ioutil"
)

type PrivKeys struct {
	Waves    string
	Ethereum string
	Ledger   string
}

func ParseConfig(filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, config); err != nil {
		return err
	}
	return err
}

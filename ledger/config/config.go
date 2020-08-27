package config

import (
	"encoding/json"
	"io/ioutil"
)

type AdaptorsConfig struct {
	NodeUrl                string
	GravityContractAddress string
	PrivKey                string
	Nebulae                []string
}
type Config struct {
	InitScore map[string]uint64
	Adapters  map[string]AdaptorsConfig
}

func Load(filename string) (Config, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	if err := json.Unmarshal(file, &config); err != nil {
		return Config{}, err
	}
	return config, err
}

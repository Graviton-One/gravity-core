package rpc

import (
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/tendermint/tendermint/crypto"
)

var GlobalClient *gravity.Client

type Config struct {
	Host    string
	pubKey  account.ConsulPubKey
	privKey crypto.PrivKey
	client  *gravity.Client
}

func NewConfig(host string, ghClientUrl string, privKey crypto.PrivKey) (*Config, error) {
	var ghPubKey account.ConsulPubKey
	copy(ghPubKey[:], privKey.PubKey().Bytes()[5:])

	ghClient, err := gravity.New(ghClientUrl)
	if err != nil {
		return nil, err
	}
	return &Config{
		Host:    host,
		privKey: privKey,
		pubKey:  ghPubKey,
		client:  ghClient,
	}, nil
}

func NewGlobalClient(ghClientUrl string) error {
	ghClient, err := gravity.New(ghClientUrl)
	if err != nil {
		return err
	}
	GlobalClient = ghClient
	return nil
}

package rpc

import (
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/tendermint/tendermint/crypto"
)

type RPCConfig struct {
	Host    string
	pubKey  account.ConsulPubKey
	privKey crypto.PrivKey
	client  *client.GravityClient
}
type VotesRq struct {
	votes []storage.Vote
}

func NewRPCConfig(host string, ghClientUrl string, privKey crypto.PrivKey) (*RPCConfig, error) {
	var ghPubKey account.ConsulPubKey
	copy(ghPubKey[:], privKey.PubKey().Bytes()[5:])

	ghClient, err := client.NewGravityClient(ghClientUrl)
	if err != nil {
		return nil, err
	}
	return &RPCConfig{
		Host:    host,
		privKey: privKey,
		pubKey:  ghPubKey,
		client:  ghClient,
	}, nil
}

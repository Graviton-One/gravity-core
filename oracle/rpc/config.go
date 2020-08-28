package rpc

import (
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/storage"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

type RPCConfig struct {
	Host    string
	pubKey  account.ConsulPubKey
	privKey tendermintCrypto.PrivKeyEd25519
	client  *client.GravityClient
}
type VotesRq struct {
	votes []storage.Vote
}

func NewRPCConfig(host string, ghClientUrl string, privKey []byte) (*RPCConfig, error) {
	ghPrivKey := tendermintCrypto.PrivKeyEd25519{}
	copy(ghPrivKey[:], privKey)

	var ghPubKey account.ConsulPubKey
	copy(ghPubKey[:], ghPrivKey.PubKey().Bytes()[5:])

	ghClient, err := client.NewGravityClient(ghClientUrl)
	if err != nil {
		return nil, err
	}
	return &RPCConfig{
		Host:    host,
		privKey: ghPrivKey,
		pubKey:  ghPubKey,
		client:  ghClient,
	}, nil
}

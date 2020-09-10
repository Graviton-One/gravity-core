package config

import (
	"encoding/json"
	"io/ioutil"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/ed25519"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
)

type PrivKeys struct {
	Validator    string
	TargetChains map[string]string
}

type PubKeys struct {
	Validator    string
	TargetChains map[string]string
}

func GeneratePrivKeys() (PrivKeys, PubKeys, error) {
	validatorPrivKey := ed25519.GenPrivKey()

	ethPrivKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return PrivKeys{}, PubKeys{}, err
	}

	wCrypto := wavesplatform.NewWavesCrypto()
	wSeed := wCrypto.RandomSeed()

	return PrivKeys{
			Validator: hexutil.Encode(validatorPrivKey[:]),
			TargetChains: map[string]string{
				account.Ethereum.String(): hexutil.Encode(ethCrypto.FromECDSA(ethPrivKey)),
				account.Waves.String():    string(wSeed),
			},
		},
		PubKeys{
			Validator: hexutil.Encode(validatorPrivKey.PubKey().Bytes()[5:]),
			TargetChains: map[string]string{
				account.Ethereum.String(): hexutil.Encode(ethCrypto.CompressPubkey(&ethPrivKey.PublicKey)),
				account.Waves.String():    string(wCrypto.PublicKey(wSeed)),
			},
		}, nil
}
func ParseConfig(filename string, config interface{}) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(file, config); err != nil {
		return err
	}
	return nil
}

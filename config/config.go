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

func GeneratePrivKeys() (PrivKeys, error) {
	validatorPrivKey := ed25519.GenPrivKey()

	ethPrivKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return PrivKeys{}, err
	}

	wavesGen := wavesplatform.NewWavesCrypto()
	wSeed := wavesGen.RandomSeed()

	return PrivKeys{
		Validator: hexutil.Encode(validatorPrivKey[:]),
		TargetChains: map[string]string{
			account.Ethereum.String(): hexutil.Encode(ethCrypto.FromECDSA(ethPrivKey)),
			account.Waves.String():    string(wSeed),
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
	return err
}

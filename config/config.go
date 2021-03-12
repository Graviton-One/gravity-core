package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/tendermint/tendermint/crypto/ed25519"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
)

type ChainIds map[string]account.ChainType
type Keys struct {
	Validator    Key
	TargetChains map[string]Key
	ChainIds     ChainIds
}

type Key struct {
	Address string
	PubKey  string
	PrivKey string
}

func GeneratePrivKeys() (*Keys, error) {
	validatorPrivKey := ed25519.GenPrivKey()

	ethPrivKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	wCrypto := wavesplatform.NewWavesCrypto()
	wSeed := wCrypto.RandomSeed()

	return &Keys{
		Validator: Key{
			Address: hexutil.Encode(validatorPrivKey.PubKey().Bytes()[5:]),
			PubKey:  hexutil.Encode(validatorPrivKey.PubKey().Bytes()[5:]),
			PrivKey: hexutil.Encode(validatorPrivKey[:]),
		},
		TargetChains: map[string]Key{
			account.Ethereum.String(): Key{
				Address: ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey).String(),
				PubKey:  hexutil.Encode(ethCrypto.CompressPubkey(&ethPrivKey.PublicKey)),
				PrivKey: hexutil.Encode(ethCrypto.FromECDSA(ethPrivKey)),
			},
			account.Waves.String(): Key{
				Address: string(wCrypto.AddressFromSeed(wSeed, 'S')),
				PubKey:  string(wCrypto.PublicKey(wSeed)),
				PrivKey: string(wSeed),
			},
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

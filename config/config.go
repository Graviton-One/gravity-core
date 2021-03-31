package config

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/hashicorp/vault/api"

	"github.com/ethereum/go-ethereum/common/hexutil"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/tendermint/tendermint/crypto/ed25519"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
)

type Keys struct {
	Validator    Key
	TargetChains map[string]Key
}

type Key struct {
	Address string
	PubKey  string
	PrivKey string
}
type VaultConfig struct {
	Url   string
	Token string
	Path  string
}

func LoadConfigFromVault(url string, token string, path string) (string, error) {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}
	client, err := api.NewClient(&api.Config{Address: url, HttpClient: httpClient})
	if err != nil {
		return "", err
	}

	client.SetToken(token)
	data, err := client.Logical().Read(path)
	if err != nil {
		return "", err
	}

	b, _ := json.Marshal(data.Data)
	return string(b), nil
}

func generateEthereumBasedPrivKeys() (*Key, error) {
	ethPrivKey, err := ethCrypto.GenerateKey()
	if err != nil {
		return nil, err
	}

	return &Key{
		Address: ethCrypto.PubkeyToAddress(ethPrivKey.PublicKey).String(),
		PubKey:  hexutil.Encode(ethCrypto.CompressPubkey(&ethPrivKey.PublicKey)),
		PrivKey: hexutil.Encode(ethCrypto.FromECDSA(ethPrivKey)),
	}, nil
}

func generateWavesPrivKeys(chain byte) (*Key, error) {
	wCrypto := wavesplatform.NewWavesCrypto()
	wSeed := wCrypto.RandomSeed()

	return &Key{
		Address: string(wCrypto.AddressFromSeed(wSeed, wavesplatform.WavesChainID(chain))),
		PubKey:  string(wCrypto.PublicKey(wSeed)),
		PrivKey: string(wSeed),
	}, nil
}

func GeneratePrivKeys(wavesChainID byte) (*Keys, error) {
	validatorPrivKey := ed25519.GenPrivKey()

	ethPrivKeys, err := generateEthereumBasedPrivKeys()
	if err != nil {
		return nil, err
	}

	bscPrivKeys, err := generateEthereumBasedPrivKeys()
	if err != nil {
		return nil, err
	}
	wavesPrivKeys, err := generateWavesPrivKeys(wavesChainID)
	if err != nil {
		return nil, err
	}

	return &Keys{
		Validator: Key{
			Address: hexutil.Encode(validatorPrivKey.PubKey().Bytes()[5:]),
			PubKey:  hexutil.Encode(validatorPrivKey.PubKey().Bytes()[5:]),
			PrivKey: hexutil.Encode(validatorPrivKey[:]),
		},
		TargetChains: map[string]Key{
			account.Ethereum.String(): *ethPrivKeys,
			account.Binance.String():  *bscPrivKeys,
			account.Waves.String():    *wavesPrivKeys,
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

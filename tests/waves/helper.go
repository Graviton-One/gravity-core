package tests

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"io/ioutil"
	"rh_tests/helpers"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"

	wavesClient "github.com/wavesplatform/gowaves/pkg/client"

	"github.com/wavesplatform/gowaves/pkg/crypto"
)
type Account struct {
	Address string
	// In case of waves: Secret is private key actually
	Secret  crypto.SecretKey
	PubKey  crypto.PublicKey
}

func GenerateAccountFromSeed(chainId byte, wordList string) (*Account, error) {
	seed := wavesplatform.Seed(wordList)
	wCrypto := wavesplatform.NewWavesCrypto()
	address := string(wCrypto.AddressFromSeed(seed, wavesplatform.WavesChainID(chainId)))
	seedWaves, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
	if err != nil {
		return nil, err
	}
	pubKey := wCrypto.PublicKey(seed)

	return &Account{
		Address: address,
		PubKey:  crypto.PublicKey(crypto.MustDigestFromBase58(string(pubKey))),
		Secret:  seedWaves,
	}, nil
}


func ScriptFromFile(filename string) ([]byte, error) {
	scriptBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	script, err := base64.StdEncoding.DecodeString(string(scriptBytes))
	if err != nil {
		return nil, err
	}

	return script, nil
}


type NetworkEnvironment int

const (
	Wavelet = 1e8
)

const (
	WavesStagenet NetworkEnvironment = iota
	WavesTestnet
)

func (env NetworkEnvironment) NodeURL() string {
	switch env {
	case WavesStagenet:
		return "https://nodes-stagenet.wavesnodes.com"
	}

	panic("no node url")
}

func (env NetworkEnvironment) ChainIDBytes() byte {
	return env.ChainID()[0]
}

func (env NetworkEnvironment) ChainID() string {
	switch env {
	case WavesStagenet:
		return "S"
	}

	panic("invalid chain id")
}

type WavesTestConfig struct {
	DistributorSeed string
	Environment     NetworkEnvironment
}

type WavesActor wavesplatform.Seed

func (actor WavesActor) Account(chainId byte) *Account {
	account, _ := GenerateAccountFromSeed(chainId, string(actor))
	return account
}

func (actor WavesActor) Recipient(chainId byte) proto.Recipient {
	recipient, _ := proto.NewRecipientFromString(actor.Account(chainId).Address)
	return recipient
}

type WavesActorSeedsMock struct {
	Gravity, Nebula, Subscriber WavesActor
}

func NewWavesActorsMock() WavesActorSeedsMock {
	wCrypto := wavesplatform.NewWavesCrypto()

	return WavesActorSeedsMock{
		Gravity:    WavesActor(wCrypto.RandomSeed()),
		Nebula:     WavesActor(wCrypto.RandomSeed()),
		Subscriber: WavesActor(wCrypto.RandomSeed()),
	}
}

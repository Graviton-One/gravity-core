package waves

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"

	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	RideErrorPrefix = "Error while executing account-script: "
)

type Config struct {
	LUPortScriptFile string
	NodeUrl          string
	DistributionSeed string
}

type RideErr struct {
	Message string
}

type Account struct {
	Address   string
	Secret    crypto.SecretKey
	PubKey    crypto.PublicKey
	Recipient proto.Recipient
}

func LoadConfig(filename string) (Config, error) {
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

func GenerateAddress(chainId byte) (*Account, error) {
	wCrypto := wavesplatform.NewWavesCrypto()
	seed := wCrypto.RandomSeed()
	address := string(wCrypto.AddressFromSeed(seed, wavesplatform.WavesChainID(chainId)))
	seedWaves, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
	if err != nil {
		return nil, err
	}
	pubKey := wCrypto.PublicKey(seed)

	recipient, err := proto.NewRecipientFromString(address)
	if err != nil {
		return nil, err
	}

	return &Account{
		Address:   address,
		PubKey:    crypto.PublicKey(crypto.MustDigestFromBase58(string(pubKey))),
		Secret:    seedWaves,
		Recipient: recipient,
	}, nil
}

func CheckRideError(rideErr error, msg string) error {
	body := rideErr.(*client.RequestError).Body
	var rsError RideErr
	err := json.Unmarshal([]byte(body), &rsError)
	if err != nil {
		return err
	}
	if rsError.Message != RideErrorPrefix+msg {
		return errors.New("error not found")
	}

	return nil
}

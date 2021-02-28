package tests

import (
	"context"
	"encoding/base64"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/client"
	"time"

	//"encoding/json"
	//"errors"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"io/ioutil"

	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type Account struct {
	Address string
	// In case of waves: Secret is private key actually
	Secret crypto.SecretKey
	PubKey crypto.PublicKey
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
	ctx             context.Context
	DistributorSeed string
	Environment     NetworkEnvironment
}

type WavesActor wavesplatform.Seed

func NewWavesActor() WavesActor {
	wCrypto := wavesplatform.NewWavesCrypto()

	return WavesActor(wCrypto.RandomSeed())
}

func (actor WavesActor) Account(chainId byte) *Account {
	account, _ := GenerateAccountFromSeed(chainId, string(actor))
	return account
}

func (actor WavesActor) Recipient(chainId byte) proto.Recipient {
	recipient, _ := proto.NewRecipientFromString(actor.Account(chainId).Address)
	return recipient
}

func (actor WavesActor) SecretKey() crypto.SecretKey {
	privKey, _ := crypto.NewSecretKeyFromBase58(string(actor.wcrypto().PrivateKey(wavesplatform.Seed(actor))))
	return privKey
}

func (actor WavesActor) wcrypto() wavesplatform.WavesCrypto {
	return wavesplatform.NewWavesCrypto()
}

type WavesActorSeedsMock struct {
	Gravity, Nebula, Subscriber WavesActor
}

func NewWavesActorsMock() WavesActorSeedsMock {
	return WavesActorSeedsMock{
		Gravity:    NewWavesActor(),
		Nebula:     NewWavesActor(),
		Subscriber: NewWavesActor(),
	}
}

func SignAndBroadcast(tx proto.Transaction, txID string, cfg WavesTestConfig, clientWaves *client.Client, wavesHelper helpers.ClientHelper, senderSeed crypto.SecretKey) error {
	var err error
	err = tx.Sign(cfg.Environment.ChainIDBytes(), senderSeed)
	if err != nil {
		return err
	}
	_, err = clientWaves.Transactions.Broadcast(cfg.ctx, tx)
	if err != nil {
		return err
	}

	err = <-wavesHelper.WaitTx(txID, cfg.ctx)
	if err != nil {
		return err
	}

	return nil
}

func TransferWavesTransaction(senderPubKey crypto.PublicKey, amount uint64, recipient proto.Recipient) *proto.TransferWithProofs {
	tx := &proto.TransferWithProofs{
		Type:    proto.TransferTransaction,
		Version: 1,
		Transfer: proto.Transfer{
			SenderPK:    senderPubKey,
			Fee:         0.001 * Wavelet,
			Timestamp:   client.NewTimestampFromTime(time.Now()),
			Recipient:   recipient,
			Amount:      amount,
			AmountAsset: nil,
		},
	}

	return tx
}

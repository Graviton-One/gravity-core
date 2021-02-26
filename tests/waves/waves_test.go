package tests

import (
	"flag"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"testing"
	wavesClient "github.com/wavesplatform/gowaves/pkg/client"
	"time"
)


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

type WavesActorSeedsMock struct {
	Gravity, Nebula, Subscriber wavesplatform.Seed
}

var cfg WavesTestConfig
var actorsMock WavesActorSeedsMock

var wCrypto = wavesplatform.NewWavesCrypto()

func TestMain(m *testing.M) {
	Init()
}

func Init() {
	var distributor, chainID string
	flag.StringVar(&chainID, "chain", "S", "network env")
	flag.StringVar(&distributor, "distributor", "", "waves token distributor")
	flag.Parse()

	if distributor == ""  {
		panic("distributor is invalid")
	}

	environment := WavesStagenet

	cfg = WavesTestConfig{ DistributorSeed: distributor, Environment: environment }

	actorsMock = WavesActorSeedsMock{
		Gravity:    wCrypto.RandomSeed(),
		Nebula:     wCrypto.RandomSeed(),
		Subscriber: wCrypto.RandomSeed(),
	}
}

func NebulaDistinctTest() error {

	nebulaScript, err := ScriptFromFile("./contracts/waves/nebula")

	distributionSeed, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(wavesplatform.Seed(cfg.DistributorSeed))))
	if err != nil {
		return err
	}
	//
	//transferTx := &proto.MassTransferWithProofs{
	//	Type:      proto.MassTransferTransaction,
	//	Version:   1,
	//	SenderPK:  crypto.GeneratePublicKey(distributionSeed),
	//	Fee:       5000000,
	//	Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
	//	Transfers: []proto.MassTransferEntry{
	//		{
	//			Amount:    2 * Wavelet,
	//			Recipient: actorsMock.Gravity,
	//		},
	//	},
	//	Attachment: &proto.LegacyAttachment{},
	//}


	return nil
}

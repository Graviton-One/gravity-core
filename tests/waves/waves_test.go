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

	actorsMock = NewWavesActorsMock()
}

func NebulaDeployTest(t *testing.T) error {
	nebulaScript, err := ScriptFromFile("./contracts/waves/nebula")

	distributionSeed, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(wavesplatform.Seed(cfg.DistributorSeed))))
	if err != nil {
		return err
	}

	// send waves
	// 0.1 - for script set
	// 0.1 - for data tx set
	transferTx := &proto.Transfer{
		SenderPK:  crypto.GeneratePublicKey(distributionSeed),
		Fee:       5000000,
		Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
		Recipient: actorsMock.Nebula.Recipient(cfg.Environment.ChainIDBytes()),
		Amount:    0.2 * Wavelet,
		AmountAsset: nil,
	}

	//err = transferTx.Sign(cfg.Environment.ChainIDBytes(), distributionSeed)
	//if err != nil {
	//	return nil, err
	//}
	//_, err = testConfig.Client.Transactions.Broadcast(testConfig.Ctx, massTx)
	//if err != nil {
	//	return nil, err
	//}
	//err = <-testConfig.Helper.WaitTx(massTx.ID.String(), testConfig.Ctx)
	//if err != nil {
	//	return nil, err
	//}

	return nil
}



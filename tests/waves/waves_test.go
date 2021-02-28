package tests

import (
	"context"
	"flag"
	"github.com/Gravity-Tech/gateway-deployer/waves/contracts"
	wavesDeployer "github.com/Gravity-Tech/gateway-deployer/waves/deployer"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	client "github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"strings"
	"testing"
	"time"
)

var cfg WavesTestConfig
var actorsMock WavesActorSeedsMock
var clientWaves *client.Client
var wavesHelper helpers.ClientHelper

var wCrypto = wavesplatform.NewWavesCrypto()

//var nebulaTestMockCfg []NebulaTestMockConfig
var nebulaTestMockCfg NebulaTestMockConfig

type NebulaTestMockConfig struct {
	chainId           byte

	BftCoefficient    int64          `dtx:"bft_coefficient"`
	GravityAddress    string         `dtx:"gravity_contract"`
	NebulaPubkey      string         `dtx:"contract_pubkey"`
	SubscriberAddress string         `dtx:"subscriber_address"`
	OraclesList       [5]WavesActor  `dtx:"oracles"`
}

func (nebulaMockCfg *NebulaTestMockConfig) OraclesPubKeysListDataEntry () string {
	var res []string

	for _, oracle := range nebulaMockCfg.OraclesList {
		res = append(res, oracle.Account(nebulaMockCfg.chainId).PubKey.String())
	}

	return strings.Join(res, ",")
}

func (nebulaMockCfg *NebulaTestMockConfig) DataEntries () proto.DataEntries {
	return proto.DataEntries{
		&proto.IntegerDataEntry{
			Key:   "bft_coefficient",
			Value: nebulaMockCfg.BftCoefficient,
		},
		&proto.StringDataEntry{
			Key:   "gravity_contract",
			Value: nebulaMockCfg.GravityAddress,
		},
		&proto.StringDataEntry{
			Key:   "contract_pubkey",
			Value: nebulaMockCfg.NebulaPubkey,
		},
		&proto.StringDataEntry{
			Key:   "subscriber_address",
			Value: nebulaMockCfg.SubscriberAddress,
		},
		&proto.StringDataEntry{
			Key:   "oracles",
			Value: nebulaMockCfg.OraclesPubKeysListDataEntry(),
		},
	}
}

func (nebulaMockCfg *NebulaTestMockConfig) OraclesPubKeysList () []string {
	var oraclesPubKeyList []string

	for _, mockedConsul := range nebulaMockCfg.OraclesList {
		pk := mockedConsul.Account(cfg.Environment.ChainIDBytes()).PubKey.String()
		oraclesPubKeyList = append(oraclesPubKeyList, pk)
	}

	return oraclesPubKeyList
}

func TestMain(m *testing.M) {
	Init()

	m.Run()
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

	cfg = WavesTestConfig{ ctx: context.Background(), DistributorSeed: distributor, Environment: environment }

	actorsMock = NewWavesActorsMock()

	clientWaves, _ = client.NewClient(client.Options{
		BaseUrl: cfg.Environment.NodeURL(),
		Client:  nil,
		ApiKey:  "",
	})

	wavesHelper = helpers.NewClientHelper(clientWaves)

	nebulaTestMockCfg = NebulaTestMockConfig {
		chainId: cfg.Environment.ChainIDBytes(),
		BftCoefficient: 5,
		GravityAddress: NewWavesActor().Account(cfg.Environment.ChainIDBytes()).Address,
		SubscriberAddress: NewWavesActor().Account(cfg.Environment.ChainIDBytes()).Address,
		OraclesList: [5]WavesActor{
			NewWavesActor(),
			NewWavesActor(),
			NewWavesActor(),
			NewWavesActor(),
			NewWavesActor(),
		},
	}
}

func NebulaDeployTest(t *testing.T) {
	nebulaScript, err := ScriptFromFile("./contracts/waves/nebula")

	//distributionSeed, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(wavesplatform.Seed(cfg.DistributorSeed))))
	distributionSeed := WavesActor(cfg.DistributorSeed).SecretKey()

	if err != nil {
		t.Fail()
	}

	// send waves
	// 0.1 - for script set
	// 0.1 - for data tx set
	transferTx := TransferWavesTransaction(
		crypto.GeneratePublicKey(distributionSeed),
		0.2 * Wavelet,
		actorsMock.Nebula.Recipient(cfg.Environment.ChainIDBytes()),
	)

	err = SignAndBroadcast(
		transferTx,
		transferTx.ID.String(),
		cfg,
		clientWaves,
		wavesHelper,
		distributionSeed,
	)

	for _, oracle := range nebulaTestMockCfg.OraclesList {
		// 0.2 for each oracle
		transferTx := TransferWavesTransaction(
			crypto.GeneratePublicKey(distributionSeed),
			0.2 * Wavelet,
			oracle.Recipient(cfg.Environment.ChainIDBytes()),
		)

		err = SignAndBroadcast(
			transferTx,
			transferTx.ID.String(),
			cfg,
			clientWaves,
			wavesHelper,
			distributionSeed,
		)

		if err != nil {
			t.Fail()
		}

	}

	if err != nil {
		t.Fail()
	}

	err = wavesDeployer.DeployNebulaWaves(
		clientWaves,
		wavesHelper,
		nebulaScript,
		nebulaTestMockCfg.GravityAddress,
		nebulaTestMockCfg.SubscriberAddress,
		nebulaTestMockCfg.OraclesPubKeysList(),
		nebulaTestMockCfg.BftCoefficient,
		contracts.BytesType,
		cfg.Environment.ChainIDBytes(),
		actorsMock.Nebula.SecretKey(),
		cfg.ctx,
	)

	if err != nil {
		t.Fail()
	}
}

func NebulaDataTransactionPersistTest(t *testing.T) error {
	var err error
	// for mock option in nebula mock options
	tx := &proto.DataWithProofs{
		Type:      proto.DataTransaction,
		Version:   1,
		Proofs:    nil,
		SenderPK:  crypto.PublicKey{},
		Entries:   nebulaTestMockCfg.DataEntries(),
		Fee:       0.1 * Wavelet,
		Timestamp: 0,
	}

	err = SignAndBroadcast(
		tx,
		tx.ID.String(),
		cfg,
		clientWaves,
		wavesHelper,
		actorsMock.Nebula.SecretKey(),
	)

	if err != nil {
		t.Fail()
	}

	return nil
}

func NebulaSendHashValueSucceedingTest(t *testing.T) {
	var err error
	var exampleHashBytes []byte
	exampleHash := "this is example data"

	copy(exampleHashBytes, exampleHash)

	chainID := cfg.Environment.ChainIDBytes()

	signaturesList := make([]string, 5, 5)
	for i, oracle := range nebulaTestMockCfg.OraclesList {
		signature, err := crypto.Sign(oracle.SecretKey(), exampleHashBytes)

		if err != nil {
			t.Fail()
		}

		signaturesList[i] = signature.String()
	}

	concatedSignatures := strings.Join(signaturesList, ",")

	sender := nebulaTestMockCfg.OraclesList[0]
	invokeTx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		ChainID:         cfg.Environment.ChainIDBytes(),
		SenderPK:        sender.Account(chainID).PubKey,
		ScriptRecipient: actorsMock.Nebula.Recipient(chainID),
		FunctionCall:    proto.FunctionCall{
			Name: "sendHashValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{ Value: exampleHashBytes },
				proto.StringArgument{ Value: concatedSignatures },
			},
		},
		Payments:        nil,
		Fee:             0.005 * Wavelet,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
	}

	err = SignAndBroadcast(invokeTx, invokeTx.ID.String(), cfg, clientWaves, wavesHelper, sender.SecretKey())

	if err != nil {
		t.Fail()
	}
}

func NebulaSendHashValueFailingTest(t *testing.T) {
	var err error
	var exampleHashBytes []byte
	exampleHash := "this is example data"

	copy(exampleHashBytes, exampleHash)

	chainID := cfg.Environment.ChainIDBytes()

	signaturesList := make([]string, 5, 5)
	for i, oracle := range nebulaTestMockCfg.OraclesList {
		signature, err := crypto.Sign(oracle.SecretKey(), exampleHashBytes)

		if err != nil {
			t.Fail()
		}

		signaturesList[i] = signature.String()
	}

	concatedSignatures := strings.Join(signaturesList, ",")

	sender := nebulaTestMockCfg.OraclesList[0]
	invokeTx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		ChainID:         cfg.Environment.ChainIDBytes(),
		SenderPK:        sender.Account(chainID).PubKey,
		ScriptRecipient: actorsMock.Nebula.Recipient(chainID),
		FunctionCall:    proto.FunctionCall{
			Name: "sendHashValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{ Value: []byte{'a'} },
				proto.StringArgument{ Value: concatedSignatures },
			},
		},
		Payments:        nil,
		Fee:             0.005 * Wavelet,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
	}

	err = SignAndBroadcast(invokeTx, invokeTx.ID.String(), cfg, clientWaves, wavesHelper, sender.SecretKey())

	if err == nil {
		t.Fail()
	}
}

func NebulaUpdateOraclesTest(t *testing.T) {

}
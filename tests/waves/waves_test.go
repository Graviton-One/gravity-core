package tests

import (
	"context"
	"flag"
	"fmt"
	"github.com/Gravity-Tech/gateway-deployer/waves/contracts"
	wavesDeployer "github.com/Gravity-Tech/gateway-deployer/waves/deployer"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	client "github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"os"
	"runtime/debug"
	"strings"
	"testing"
	"time"
)

var cfg WavesTestConfig
var actorsMock WavesActorSeedsMock
var clientWaves *client.Client
var wavesHelper helpers.ClientHelper

var wCrypto = wavesplatform.NewWavesCrypto()

//var nebulaTestMockCfg []NebulaTestMockConfig - validate different mocks in future
var nebulaTestMockCfg NebulaTestMockConfig

func TestMain(m *testing.M) {
	Init()

	os.Exit(m.Run())
}

func handleError(t *testing.T, err error, successMessage string) {
	if err != nil {
		t.Error(err.Error())
		debug.PrintStack()
		panic(1)
	} else {
		t.Log(successMessage)
	}
}

var testInputOption *testsInputCommand

type testsInputCommand struct {
	distributor, chainID string
}

func Init() {
	testInputOption = &testsInputCommand{}

	flag.StringVar(&testInputOption.chainID, "chain", "S", "network env")
	flag.StringVar(&testInputOption.distributor, "distributor", "", "waves token distributor")
	flag.Parse()

	if testInputOption.distributor == ""  {
		panic("distributor is invalid")
	}

	// TODO: from chain id fn should be here
	environment := WavesStagenet

	cfg = WavesTestConfig{ ctx: context.Background(), DistributorSeed: testInputOption.distributor, Environment: environment }

	actorsMock = NewWavesActorsMock()

	var err error
	clientWaves, err = client.NewClient(client.Options{
		BaseUrl: cfg.Environment.NodeURL(),
		Client:  nil,
		ApiKey:  "",
	})

	if err != nil {
		fmt.Printf("err: %v \n", err)
	}

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

func TestNebulaMockConfig(t *testing.T) {
	err := nebulaTestMockCfg.Validate()

	handleError(t, err, "nebula mock is valid")
}

func TestNebulaDeploy(t *testing.T) {


	nebulaScript, err := ScriptFromFile("../../abi/waves/nebula.abi")

	//distributionSeed, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(wavesplatform.Seed(cfg.DistributorSeed))))
	distributionSeed := WavesActor(cfg.DistributorSeed).SecretKey()

	handleError(t, err, "nebula waves script found successfully")

	// send waves
	// 0.1 - for script set
	// 0.1 - for data tx set
	transferTx := TransferWavesTransaction(
		crypto.GeneratePublicKey(distributionSeed),
		0.2 * Wavelet,
		actorsMock.Nebula.Recipient(cfg.Environment.ChainIDBytes()),
	)

	_, err = SignAndBroadcast(
		transferTx,
		cfg,
		clientWaves,
		wavesHelper,
		distributionSeed,
	)

	handleError(t, err, "successfully sent 0.2 waves to nebula")

	amountPerOracle, oracleCount := 0.1, float64(len(nebulaTestMockCfg.OraclesList))
	var oraclesMassTransferList []proto.MassTransferEntry
	for _, oracle := range nebulaTestMockCfg.OraclesList {
		// 0.2 for each oracle
		oraclesMassTransferList = append(oraclesMassTransferList, proto.MassTransferEntry{
			Recipient: oracle.Recipient(cfg.Environment.ChainIDBytes()),
			Amount:    0.1 * Wavelet,
		})
	}

	massTransferTx := &proto.MassTransferWithProofs{
		Type:       proto.MassTransferTransaction,
		Version:    1,
		SenderPK:   crypto.GeneratePublicKey(distributionSeed),
		Asset:      proto.OptionalAsset{},
		Transfers:  oraclesMassTransferList,
		Timestamp:  client.NewTimestampFromTime(time.Now()),
		Fee:        0.004 * Wavelet,
	}

	_, err = SignAndBroadcast(massTransferTx, cfg, clientWaves, wavesHelper, distributionSeed)

	handleError(t, err, fmt.Sprintf("mass transferred waves (%v * %v = %v) to oracles", amountPerOracle, oracleCount, amountPerOracle * oracleCount))

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

	handleError(t, err, "deployed nebula")
}


// TestNebulaSendHashValueSucceeding
// TestNebulaSendHashValueFailing_InvalidHash
// TestNebulaSendHashValueFailing_NotEnoughSignatures

func TestNebulaSendHashValueSucceeding(t *testing.T) {
	var err error

	t.Log("check the only succeeding behaviour of \"sendHashValue\" function")

	exampleHash := "this is example data"
	exampleHashBytes := []byte(exampleHash)

	chainID := cfg.Environment.ChainIDBytes()

	signaturesList := make([]string, 5, 5)
	for i, oracle := range nebulaTestMockCfg.OraclesList {
		signature, err := crypto.Sign(oracle.SecretKey(), exampleHashBytes)

		handleError(t, err, fmt.Sprintf("message signed by oracle #%v successfully \n", i + 1))

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

	result, err := SignAndBroadcast(invokeTx, cfg, clientWaves, wavesHelper, sender.SecretKey())

	handleError(t, err, fmt.Sprintf("success: invoked sendHashValue. tx id: %v", result.TxID))
}

func TestNebulaSendHashValueFailing_InvalidHash(t *testing.T) {
	var err error

	t.Log("ensure \"sendHashValue\" fails if signed hash differs for input hash")

	targetHash := []byte("sdgfjnsdfgjsdnfg")

	exampleHash := "this is example data"
	exampleHashBytes := []byte(exampleHash)

	copy(exampleHashBytes, exampleHash)

	chainID := cfg.Environment.ChainIDBytes()

	signaturesList := make([]string, 5, 5)
	for i, oracle := range nebulaTestMockCfg.OraclesList {
		signature, err := crypto.Sign(oracle.SecretKey(), exampleHashBytes)

		handleError(t, err, fmt.Sprintf("message signed by oracle #%v: %v \n", i + 1, signature.String()))

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
				proto.BinaryArgument{ Value: targetHash },
				proto.StringArgument{ Value: concatedSignatures },
			},
		},
		Payments:        nil,
		Fee:             0.005 * Wavelet,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
	}

	response, err := SignAndBroadcast(invokeTx, cfg, clientWaves, wavesHelper, sender.SecretKey())

	if err == nil {
		t.Error(fmt.Sprintf("error: invalid hash was accepted by nebula. tx: %v \n", response.TxID))
	} else {
		t.Log(fmt.Sprintf("success: invalid hash was rejected by nebula. response: \n %v \n", err.Error()))
	}
}

/**
 * The goal is to check that "sendHashValue" fails if total signatures count < bftCoefficient
 */
func TestNebulaSendHashValueFailing_NotEnoughSignatures(t *testing.T) {
	var err error

	t.Log("ensure \"sendHashValue\" fails if not enough oracle signatures provided")

	exampleHash := "this is example data"
	exampleHashBytes := []byte(exampleHash)

	chainID := cfg.Environment.ChainIDBytes()

	bftCoefficient := nebulaTestMockCfg.BftCoefficient

	if bftCoefficient < 1 {
		t.Error(fmt.Sprintf("error: bft_coefficient provided is less than 1. actual value: %v", bftCoefficient))
		panic(1)
	}

	signaturesList := make([]string, 5, 5)
	for i := int64(0); i < bftCoefficient - 1; i++ {
		oracle := nebulaTestMockCfg.OraclesList[i]
		signature, err := crypto.Sign(oracle.SecretKey(), exampleHashBytes)

		handleError(t, err, fmt.Sprintf("message signed by oracle #%v: %v \n", i + 1, signature.String()))

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

	response, err := SignAndBroadcast(invokeTx, cfg, clientWaves, wavesHelper, sender.SecretKey())

	if err == nil {
		t.Error(fmt.Sprintf("error: signatures count < bft_coefficient. hash was accepted by nebula. tx: %v \n", response.TxID))
	} else {
		t.Log(fmt.Sprintf("success: signatures count < bft_coefficient. hash was rejected by nebula. response: \n %v \n", err.Error()))
	}
}

/*
 * TODO: TestUpdateOraclesSucceeding -
 *   signatures count >= bft_coefficient & rest input params are valid
 */
/*
 * TODO: TestUpdateOraclesFailing_SignaturesListLengthIsLessThanFive -
 *   RIDE source code is bounded for signsList to have len() == 5, so '1,1,1,1,1' is valid.
 *   That is why '1,1,1' is succeeding case for this test
 */
/*
 * TODO: TestUpdateOraclesFailing_NotEnoughSignatures -
 *   actual signatures count is less than bft
 */

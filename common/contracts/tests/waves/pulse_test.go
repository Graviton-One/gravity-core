package waves

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
	"github.com/Gravity-Tech/gravity-core/deployer"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	BftValue    = 3
	ConsulCount = 5
	OracleCount = 5
	ChainId     = 'R'
	Wavelet     = 100000000
)

type TestPulseConfig struct {
	Helper helpers.ClientHelper
	Client *client.Client
	Ctx    context.Context

	Gravity *Account
	Nebula  *Account
	Sub     *Account

	Consuls []*Account
	Oracles []*Account
}

var config *TestPulseConfig
var tests = map[string]func(t *testing.T){
	"sendHashPositive":     testSendHashPositive,
	"sendHashInvalidSigns": testSendHashInvalidSigns,
	"sendSubPositive":      testSendSubPositive,
	"sendSubInvalidHash":   testSendSubInvalidHash,
}

func TestPulse(t *testing.T) {
	var err error
	config, err = initTests()
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range tests {
		height, _, err := config.Client.Blocks.Height(config.Ctx)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(k, v)

		err = <-config.Helper.WaitByHeight(height.Height+1, config.Ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func initTests() (*TestPulseConfig, error) {
	var testConfig TestPulseConfig
	testConfig.Ctx = context.Background()

	cfg, err := LoadConfig("config.json")
	if err != nil {
		return nil, err
	}

	wClient, err := client.NewClient(client.Options{ApiKey: "", BaseUrl: cfg.NodeUrl})
	if err != nil {
		return nil, err
	}
	testConfig.Client = wClient
	testConfig.Helper = helpers.NewClientHelper(testConfig.Client)

	testConfig.Gravity, err = GenerateAddress(ChainId)
	if err != nil {
		return nil, err
	}

	testConfig.Nebula, err = GenerateAddress(ChainId)
	if err != nil {
		return nil, err
	}

	testConfig.Sub, err = GenerateAddress(ChainId)
	if err != nil {
		return nil, err
	}

	for i := 0; i < ConsulCount; i++ {
		consul, err := GenerateAddress(ChainId)
		if err != nil {
			return nil, err
		}

		testConfig.Consuls = append(testConfig.Consuls, consul)
	}
	for i := 0; i < OracleCount; i++ {
		oracle, err := GenerateAddress(ChainId)
		if err != nil {
			return nil, err
		}

		testConfig.Oracles = append(testConfig.Consuls, oracle)
	}

	gravityScript, err := ScriptFromFile(cfg.GravityScriptFile)
	if err != nil {
		return nil, err
	}

	nebulaScript, err := ScriptFromFile(cfg.NebulaScriptFile)
	if err != nil {
		return nil, err
	}

	subScript, err := ScriptFromFile(cfg.SubMockScriptFile)
	if err != nil {
		return nil, err
	}

	wCrypto := wavesplatform.NewWavesCrypto()
	distributionSeed, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(wavesplatform.Seed(cfg.DistributionSeed))))
	if err != nil {
		return nil, err
	}

	gravityAddressRecipient, err := proto.NewRecipientFromString(testConfig.Gravity.Address)
	if err != nil {
		return nil, err
	}
	nebulaAddressRecipient, err := proto.NewRecipientFromString(testConfig.Nebula.Address)
	if err != nil {
		return nil, err
	}
	subAddressRecipient, err := proto.NewRecipientFromString(testConfig.Sub.Address)
	if err != nil {
		return nil, err
	}
	oracleRecipient, err := proto.NewRecipientFromString(testConfig.Oracles[0].Address)
	if err != nil {
		return nil, err
	}

	massTx := &proto.MassTransferWithProofs{
		Type:      proto.MassTransferTransaction,
		Version:   1,
		SenderPK:  crypto.GeneratePublicKey(distributionSeed),
		Fee:       5000000,
		Timestamp: client.NewTimestampFromTime(time.Now()),
		Transfers: []proto.MassTransferEntry{
			{
				Amount:    2 * Wavelet,
				Recipient: gravityAddressRecipient,
			},
			{
				Amount:    2 * Wavelet,
				Recipient: nebulaAddressRecipient,
			},
			{
				Amount:    2 * Wavelet,
				Recipient: subAddressRecipient,
			},
			{
				Amount:    2 * Wavelet,
				Recipient: oracleRecipient,
			},
		},
		Attachment: &proto.LegacyAttachment{},
	}
	err = massTx.Sign(ChainId, distributionSeed)
	if err != nil {
		return nil, err
	}
	_, err = testConfig.Client.Transactions.Broadcast(testConfig.Ctx, massTx)
	if err != nil {
		return nil, err
	}
	err = <-testConfig.Helper.WaitTx(massTx.ID.String(), testConfig.Ctx)
	if err != nil {
		return nil, err
	}

	var consulsString []string
	for _, v := range testConfig.Consuls {
		consulsString = append(consulsString, v.Address)
	}
	err = deployer.DeployGravityWaves(testConfig.Client, testConfig.Helper, gravityScript, consulsString, BftValue, ChainId, testConfig.Gravity.Secret, testConfig.Ctx)
	if err != nil {
		return nil, err
	}

	err = deployer.DeploySubWaves(testConfig.Client, testConfig.Helper, subScript, ChainId, testConfig.Sub.Secret, testConfig.Ctx)
	if err != nil {
		return nil, err
	}

	var oraclesString []string
	for _, v := range testConfig.Oracles {
		oraclesString = append(oraclesString, v.PubKey.String())
	}
	err = deployer.DeployNebulaWaves(testConfig.Client, testConfig.Helper, nebulaScript, testConfig.Gravity.Address,
		testConfig.Sub.Address, oraclesString, BftValue, contracts.BytesType, ChainId, testConfig.Nebula.Secret, testConfig.Ctx)
	if err != nil {
		return nil, err
	}

	return &testConfig, nil
}

func testSendHashPositive(t *testing.T) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := crypto.Keccak256(id)
	if err != nil {
		t.Fatal(err)
	}

	var signs []string
	for _, v := range config.Oracles {
		sign, err := crypto.Sign(v.Secret, hash.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		signs = append(signs, sign.String())
	}

	recipient, err := proto.NewAddressFromString(config.Nebula.Address)
	if err != nil {
		t.Fatal(err)
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Oracles[0].PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "sendHashValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: hash.Bytes(),
				},
				proto.StringArgument{
					Value: strings.Join(signs, ","),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: proto.NewRecipientFromAddress(recipient),
	}
	err = tx.Sign(ChainId, config.Oracles[0].Secret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.Client.Transactions.Broadcast(config.Ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	err = <-config.Helper.WaitTx(tx.ID.String(), config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
}
func testSendHashInvalidSigns(t *testing.T) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := crypto.Keccak256(id)
	if err != nil {
		t.Fatal(err)
	}

	var signs []string
	for i := 0; i < OracleCount; i++ {
		if i >= (BftValue - 1) {
			signs = append(signs, "")
			continue
		}
		sign, err := crypto.Sign(config.Oracles[i].Secret, hash.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		signs = append(signs, sign.String())
	}

	recipient, err := proto.NewAddressFromString(config.Nebula.Address)
	if err != nil {
		t.Fatal(err)
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Oracles[0].PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "sendHashValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: hash.Bytes(),
				},
				proto.StringArgument{
					Value: strings.Join(signs, ","),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: proto.NewRecipientFromAddress(recipient),
	}
	err = tx.Sign(ChainId, config.Oracles[0].Secret)
	if err != nil {
		t.Fatal(err)
	}

	_, err = config.Client.Transactions.Broadcast(config.Ctx, tx)
	if err == nil {
		t.Fatal("invalid signs not fail in contract")
	}

	err = CheckRideError(err, "invalid bft count")
	if err != nil {
		t.Fatal(err)
	}
}

func testSendSubPositive(t *testing.T) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := crypto.Keccak256(id)
	if err != nil {
		t.Fatal(err)
	}

	var signs []string
	for _, v := range config.Oracles {
		sign, err := crypto.Sign(v.Secret, hash.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		signs = append(signs, sign.String())
	}

	recipient, err := proto.NewAddressFromString(config.Nebula.Address)
	if err != nil {
		t.Fatal(err)
	}
	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Oracles[0].PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "sendHashValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: hash.Bytes(),
				},
				proto.StringArgument{
					Value: strings.Join(signs, ","),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: proto.NewRecipientFromAddress(recipient),
	}

	err = tx.Sign(ChainId, config.Oracles[0].Secret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.Client.Transactions.Broadcast(config.Ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	err = <-config.Helper.WaitTx(tx.ID.String(), config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	lastPulseState, _, err := config.Helper.GetStateByAddressAndKey(config.Nebula.Address, "last_pulse_id", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	lastPulseId := int64(lastPulseState.Value.(float64))
	recipient, err = proto.NewAddressFromString(config.Sub.Address)
	if err != nil {
		t.Fatal(err)
	}
	tx = &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Nebula.PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "attachValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: id,
				},
				proto.IntegerArgument{
					Value: lastPulseId,
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: proto.NewRecipientFromAddress(recipient),
	}

	err = tx.Sign(ChainId, config.Oracles[0].Secret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.Client.Transactions.Broadcast(config.Ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	err = <-config.Helper.WaitTx(tx.ID.String(), config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	subValue, _, err := config.Helper.GetStateByAddressAndKey(config.Sub.Address, fmt.Sprintf("%d", lastPulseId), config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	value, err := base64.StdEncoding.DecodeString(strings.Split(subValue.Value.(string), ":")[1])
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(value, id) {
		t.Fatal("invalid sent value")
	}
}
func testSendSubInvalidHash(t *testing.T) {
	id := make([]byte, 32)
	_, err := rand.Read(id)
	if err != nil {
		t.Fatal(err)
	}

	hash, err := crypto.Keccak256(id)
	if err != nil {
		t.Fatal(err)
	}

	var signs []string
	for _, v := range config.Oracles {
		sign, err := crypto.Sign(v.Secret, hash.Bytes())
		if err != nil {
			t.Fatal(err)
		}

		signs = append(signs, sign.String())
	}

	recipient, err := proto.NewAddressFromString(config.Nebula.Address)
	if err != nil {
		t.Fatal(err)
	}
	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Oracles[0].PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "sendHashValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: hash.Bytes(),
				},
				proto.StringArgument{
					Value: strings.Join(signs, ","),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: proto.NewRecipientFromAddress(recipient),
	}

	err = tx.Sign(ChainId, config.Oracles[0].Secret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.Client.Transactions.Broadcast(config.Ctx, tx)
	if err != nil {
		t.Fatal(err)
	}
	err = <-config.Helper.WaitTx(tx.ID.String(), config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	lastPulseId, _, err := config.Helper.GetStateByAddressAndKey(config.Nebula.Address, "last_pulse_id", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	lastPulse := int64(lastPulseId.Value.(float64))
	recipient, err = proto.NewAddressFromString(config.Sub.Address)
	if err != nil {
		t.Fatal(err)
	}
	tx = &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Nebula.PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "attachValue",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: make([]byte, 32, 32),
				},
				proto.IntegerArgument{
					Value: lastPulse,
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: proto.NewRecipientFromAddress(recipient),
	}

	err = tx.Sign(ChainId, config.Oracles[0].Secret)
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.Client.Transactions.Broadcast(config.Ctx, tx)
	if err == nil {
		t.Fatal("invalid value is sent")
	}
	err = CheckRideError(err, "invalid keccak256(value)")
	if err != nil {
		t.Fatal(err)
	}
}

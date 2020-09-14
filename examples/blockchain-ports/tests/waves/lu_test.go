package waves

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"rebbit-hole/deployer"
	"rebbit-hole/helpers"
	"testing"
	"time"

	"github.com/mr-tron/base58"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	APPROVE = 1
	UNLOCK  = 2

	NEW       = 1
	COMPLETED = 2

	ChainId = 'R'
	Wavelet = 100000000

	DefaultRqAmount = 1000
)

type TestConfig struct {
	Helper helpers.ClientHelper
	Client *client.Client
	Ctx    context.Context

	LUPort  *Account
	Tester  *Account
	AssetId string
}

var config *TestConfig
var tests = map[string]func(t *testing.T){
	"createFiveRq":    testCreateFiveRq,
	"approveMiddleRq": testApproveMiddleRq,
	"approveFirstRq":  testApproveFirstRq,
	"approveLastRq":   testApproveLastRq,
	"unlockRq":        testUnlockRq,
}

func TestLU(t *testing.T) {
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

		success := t.Run(k, v)
		if !success {
			t.Fatal("invalid run tests")
		}
		err = <-config.Helper.WaitByHeight(height.Height+1, config.Ctx)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func initTests() (*TestConfig, error) {
	var testConfig TestConfig
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

	testConfig.LUPort, err = GenerateAddress(ChainId)
	if err != nil {
		return nil, err
	}

	testConfig.Tester, err = GenerateAddress(ChainId)
	if err != nil {
		return nil, err
	}

	luPortScript, err := ScriptFromFile(cfg.LUPortScriptFile)
	if err != nil {
		return nil, err
	}

	wCrypto := wavesplatform.NewWavesCrypto()
	distributionSeed, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(wavesplatform.Seed(cfg.DistributionSeed))))
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
				Recipient: testConfig.LUPort.Recipient,
			},
			{
				Amount:    100 * Wavelet,
				Recipient: testConfig.Tester.Recipient,
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

	testConfig.AssetId, err = deployer.CreateToken(testConfig.Client, testConfig.Helper, "TEST", "TEST", 10000000000, 0, ChainId, testConfig.Tester.Secret, testConfig.Ctx)
	if err != nil {
		return nil, err
	}

	err = deployer.DeployLUPort(testConfig.Client, testConfig.Helper, luPortScript, testConfig.Tester.Address, testConfig.AssetId, ChainId, testConfig.LUPort.Secret, testConfig.Ctx)
	if err != nil {
		return nil, err
	}

	return &testConfig, nil
}

func testCreateFiveRq(t *testing.T) {
	asset, err := proto.NewOptionalAssetFromString(config.AssetId)
	if err != nil {
		t.Fatal(err)
	}

	var requests []string
	for i := 0; i < 5; i++ {
		tx := &proto.InvokeScriptWithProofs{
			Type:     proto.InvokeScriptTransaction,
			Version:  1,
			SenderPK: config.Tester.PubKey,
			ChainID:  ChainId,
			FunctionCall: proto.FunctionCall{
				Name: "createTransferWrapRq",
				Arguments: proto.Arguments{
					proto.StringArgument{
						Value: "test",
					},
				},
			},
			Payments: proto.ScriptPayments{
				proto.ScriptPayment{
					Amount: DefaultRqAmount,
					Asset:  *asset,
				},
			},
			Fee:             5000000,
			Timestamp:       client.NewTimestampFromTime(time.Now()),
			ScriptRecipient: config.LUPort.Recipient,
		}
		err := tx.Sign(ChainId, config.Tester.Secret)
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

		requests = append(requests, tx.ID.String())
	}

	firstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	first := firstRqState.Value.(string)

	lastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	last := lastRqState.Value.(string)

	if first != requests[0] {
		t.Fatal("invalid first rq")
	} else if last != requests[len(requests)-1] {
		t.Fatal("invalid last rq")
	}

	rq := first
	for i := 1; i < len(requests)-2; i++ {
		nextRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "next_rq_"+rq, config.Ctx)
		if err != nil {
			t.Fatal(err)
		}

		prevRq := rq
		rq = nextRqState.Value.(string)

		prevRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "prev_rq_"+rq, config.Ctx)
		if err != nil {
			t.Fatal(err)
		}

		if prevRqState.Value.(string) != prevRq {
			t.Fatal("invalid prev rq position")
		} else if rq != requests[i] {
			t.Fatal("invalid rq position")
		}
	}
}

func testApproveMiddleRq(t *testing.T) {
	firstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	first := firstRqState.Value.(string)

	lastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	last := lastRqState.Value.(string)

	nextRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "next_rq_"+first, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	rq := nextRqState.Value.(string)
	nextNextRqStatus, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "next_rq_"+rq, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	nextNextRq := nextNextRqStatus.Value.(string)

	var action [8]byte
	binary.BigEndian.PutUint64(action[:], APPROVE)

	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Tester.PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "attachData",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: append(action[:], crypto.MustBytesFromBase58(rq)...),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: config.LUPort.Recipient,
	}
	err = tx.Sign(ChainId, config.Tester.Secret)
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

	newFirstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	newFirst := newFirstRqState.Value.(string)

	newLastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	newLast := newLastRqState.Value.(string)

	if first != newFirst {
		t.Fatal("invalid drop approved rq (first rq dropped)")
	} else if last != newLast {
		t.Fatal("invalid drop approved rq (last rq dropped)")
	}

	newNextRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "next_rq_"+newFirst, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	if nextNextRq != newNextRqState.Value.(string) {
		t.Fatal("invalid drop approved rq (invalid next rq)")
	}

	prevRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "prev_rq_"+nextNextRq, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	if first != prevRqState.Value.(string) {
		t.Fatal("invalid drop approved rq (invalid prev rq)")
	}

	statusRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "rq_status_"+rq, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	if int64(statusRqState.Value.(float64)) != COMPLETED {
		t.Fatal("invalid changed status for rq")
	}
}
func testApproveFirstRq(t *testing.T) {
	firstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	first := firstRqState.Value.(string)

	lastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	last := lastRqState.Value.(string)

	nextRqStatus, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "next_rq_"+first, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	nextRq := nextRqStatus.Value.(string)

	var action [8]byte
	binary.BigEndian.PutUint64(action[:], APPROVE)

	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Tester.PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "attachData",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: append(action[:], crypto.MustBytesFromBase58(first)...),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: config.LUPort.Recipient,
	}
	err = tx.Sign(ChainId, config.Tester.Secret)
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

	newFirstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	newFirst := newFirstRqState.Value.(string)

	newLastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	newLast := newLastRqState.Value.(string)

	if first == newFirst {
		t.Fatal("invalid drop first approved rq")
	} else if last != newLast {
		t.Fatal("invalid drop approved rq (last rq dropped)")
	}

	if newFirst != nextRq {
		t.Fatal("invalid drop first approved rq")
	}

	statusRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "rq_status_"+first, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	if int64(statusRqState.Value.(float64)) != COMPLETED {
		t.Fatal("invalid changed status for rq")
	}
}
func testApproveLastRq(t *testing.T) {
	firstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	first := firstRqState.Value.(string)

	lastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	last := lastRqState.Value.(string)

	prevRqStatus, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "prev_rq_"+last, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	prevRq := prevRqStatus.Value.(string)

	var action [8]byte
	binary.BigEndian.PutUint64(action[:], APPROVE)

	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Tester.PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "attachData",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: append(action[:], crypto.MustBytesFromBase58(last)...),
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: config.LUPort.Recipient,
	}
	err = tx.Sign(ChainId, config.Tester.Secret)
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

	newFirstRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "first_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	newFirst := newFirstRqState.Value.(string)

	newLastRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "last_rq", config.Ctx)
	if err != nil {
		t.Fatal(err)
	}
	newLast := newLastRqState.Value.(string)

	if first != newFirst {
		t.Fatal("invalid drop first approved rq (first rq dropped)")
	} else if newLast != prevRq {
		t.Fatal("invalid drop approved rq")
	}

	statusRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "rq_status_"+last, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	if int64(statusRqState.Value.(float64)) != COMPLETED {
		t.Fatal("invalid changed status for rq")
	}
}

func testUnlockRq(t *testing.T) {
	var args []byte

	var action [8]byte
	binary.BigEndian.PutUint64(action[:], UNLOCK)

	var amount [8]byte
	binary.BigEndian.PutUint64(amount[:], DefaultRqAmount)

	var rqId [32]byte
	_, err := rand.Read(rqId[:])
	if err != nil {
		t.Fatal(err)
	}

	args = append(args, action[:]...)
	args = append(args, rqId[:]...)
	args = append(args, amount[:]...)
	args = append(args, crypto.MustBytesFromBase58(config.Tester.Address[:])...)
	tx := &proto.InvokeScriptWithProofs{
		Type:     proto.InvokeScriptTransaction,
		Version:  1,
		SenderPK: config.Tester.PubKey,
		ChainID:  ChainId,
		FunctionCall: proto.FunctionCall{
			Name: "attachData",
			Arguments: proto.Arguments{
				proto.BinaryArgument{
					Value: args,
				},
			},
		},
		Fee:             5000000,
		Timestamp:       client.NewTimestampFromTime(time.Now()),
		ScriptRecipient: config.LUPort.Recipient,
	}
	err = tx.Sign(ChainId, config.Tester.Secret)
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

	rqIdString := base58.Encode(rqId[:])

	statusRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "rq_status_"+rqIdString, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	amountRqState, _, err := config.Helper.GetStateByAddressAndKey(config.LUPort.Address, "rq_amount_"+rqIdString, config.Ctx)
	if err != nil {
		t.Fatal(err)
	}

	if int64(statusRqState.Value.(float64)) != COMPLETED {
		t.Fatal("invalid changed status for rq")
	} else if int64(amountRqState.Value.(float64)) != DefaultRqAmount {
		t.Fatal("invalid rq amount")
	}
}

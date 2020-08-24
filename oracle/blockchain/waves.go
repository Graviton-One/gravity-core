package blockchain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/client"

	wavesHelper "github.com/Gravity-Tech/gravity-core/common/helpers/waves"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/btcsuite/btcutil/base58"
	wavesClient "github.com/wavesplatform/gowaves/pkg/client"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type WavesClient struct {
	ghClient        *client.Client
	wavesClient     *wavesClient.Client
	contractAddress string
	chainID         byte
	helper          wavesHelper.ClientHelper
	nebulaId        account.NebulaId
	privKey         []byte
}

func NewWavesClient(ghClient *client.Client, nebulaId account.NebulaId, privKey []byte, nodeUrl string) (*WavesClient, error) {
	wavesClient, err := wavesClient.NewClient(wavesClient.Options{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	return &WavesClient{
		ghClient:        ghClient,
		wavesClient:     wavesClient,
		contractAddress: base58.Encode(nebulaId),
		nebulaId:        nebulaId,
		chainID:         'T',
		privKey:         privKey,
		helper:          wavesHelper.NewClientHelper(wavesClient),
	}, nil
}

func (client *WavesClient) GetHeight(ctx context.Context) (uint64, error) {
	wavesHeight, _, err := client.wavesClient.Blocks.Height(ctx)
	if err != nil {
		return 0, err
	}

	return wavesHeight.Height, nil
}

func (client *WavesClient) SendResult(tcHeight uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	state, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, fmt.Sprintf("%d", tcHeight), ctx)
	if err != nil {
		return "", err
	} else if state != nil {
		return "", nil
	}

	funcArgs := new(proto.Arguments)
	funcArgs.Append(proto.BinaryArgument{
		Value: hash,
	})
	realSignCount := 0
	var signs []string
	oracles, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, "oracles", ctx)
	if err != nil {
		return "", err
	}

	for _, oracle := range strings.Split(oracles.Value.(string), ",") {
		var pubKey account.OraclesPubKey
		copy(pubKey[:], base58.Decode(oracle))
		sign, err := client.ghClient.Result(account.Ethereum, client.nebulaId, int64(tcHeight), pubKey)
		if err != nil {
			signs = append(signs, "nil")
			continue
		}
		signs = append(signs, base58.Encode(sign))
		realSignCount++
	}
	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(signs, ","),
	})

	bft, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, "bft_coefficient", ctx)
	if err != nil {
		return "", err
	}

	if realSignCount == 0 {
		return "", nil
	}
	if realSignCount >= int(bft.Value.(float64)) {
		secret, err := wavesCrypto.NewSecretKeyFromBytes(client.privKey)

		asset, err := proto.NewOptionalAssetFromString("WAVES")
		if err != nil {
			return "", err
		}
		contract, err := proto.NewRecipientFromString(client.contractAddress)
		if err != nil {
			return "", err
		}
		tx := &proto.InvokeScriptWithProofs{
			Type:            proto.InvokeScriptTransaction,
			Version:         1,
			SenderPK:        wavesCrypto.GeneratePublicKey(secret),
			ChainID:         client.chainID,
			ScriptRecipient: contract,
			FunctionCall: proto.FunctionCall{
				Name:      "confirmData",
				Arguments: *funcArgs,
			},
			Payments:  nil,
			FeeAsset:  *asset,
			Fee:       500000,
			Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
		}

		err = tx.Sign(client.chainID, secret)
		if err != nil {
			return "", err
		}

		_, err = client.wavesClient.Transactions.Broadcast(ctx, tx)
		if err != nil {
			return "", err
		}

		fmt.Printf("Tx finilize: %s \n", tx.ID)
		return tx.ID.String(), nil
	}

	return "", nil
}

func (client *WavesClient) SendSubs(tcHeight uint64, value interface{}, ctx context.Context) error {
	state, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, fmt.Sprintf("%d", tcHeight), ctx)
	if err != nil {
		return err
	} else if state == nil {
		return nil
	}

	subContract, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, "subscriber_address", ctx)
	if err != nil {
		return err
	}

	pubKeyNebulaContract, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, "contract_pubkey", ctx)
	if err != nil {
		return err
	}

	height, _, err := client.helper.GetStateByAddressAndKey(subContract.Value.(string), fmt.Sprintf("%d", tcHeight), ctx)
	if err != nil {
		return err
	}

	if height != nil {
		return err
	}
	secret, err := wavesCrypto.NewSecretKeyFromBytes(client.privKey)

	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		return err
	}

	contract, err := proto.NewRecipientFromString(subContract.Value.(string))
	if err != nil {
		return err
	}

	pubKey, err := wavesCrypto.NewPublicKeyFromBase58(pubKeyNebulaContract.Value.(string))
	if err != nil {
		return err
	}

	nebulaType, _, err := client.helper.GetStateByAddressAndKey(client.contractAddress, "type", ctx)
	if err != nil {
		return err
	}

	args := proto.Arguments{}
	switch SubType(nebulaType.Value.(int8)) {
	case Int64:
		args.Append(
			proto.IntegerArgument{
				Value: value.(int64),
			})
	case String:
		args.Append(
			proto.StringArgument{
				Value: value.(string),
			})
	case Bytes:
		args.Append(
			proto.BinaryArgument{
				Value: value.([]byte),
			})
	}
	args.Append(
		proto.IntegerArgument{
			Value: int64(tcHeight),
		})

	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        pubKey,
		ChainID:         client.chainID,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name:      "attachData",
			Arguments: args,
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
	}

	err = tx.Sign(client.chainID, secret)
	if err != nil {
		return err
	}

	_, err = client.wavesClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return err
	}

	fmt.Printf("Sub send tx: %s \n", tx.ID.String())

	return nil
}

func (client *WavesClient) WaitTx(id string, ctx context.Context) error {
	return <-client.helper.WaitTx(id, ctx)
}

func (client *WavesClient) Sign(msg []byte) ([]byte, error) {
	secret, err := wavesCrypto.NewSecretKeyFromBytes(client.privKey)
	if err != nil {
		return nil, err
	}
	sig, err := wavesCrypto.Sign(secret, msg)
	if err != nil {
		return nil, err
	}
	return sig.Bytes(), nil
}

package blockchain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"
	ghClient "github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/oracle-node/helpers"
	"github.com/btcsuite/btcutil/base58"
	"github.com/wavesplatform/gowaves/pkg/client"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Waves struct {
	client          *client.Client
	contractAddress string
	chainID         byte
}

func NewWaves(contractAddress string, nodeUrl string, ctx context.Context) (*Waves, error) {
	wavesClient, err := client.NewClient(client.Options{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	return &Waves{
		client:          wavesClient,
		contractAddress: contractAddress,
		chainID:         'T',
	}, nil
}

func (waves *Waves) GetHeight(ctx context.Context) (uint64, error) {
	wavesHeight, _, err := waves.client.Blocks.Height(ctx)
	if err != nil {
		return 0, err
	}

	return wavesHeight.Height, nil
}

func (waves *Waves) SendResult(ghClient *ghClient.Client, tcHeight uint64, privKey []byte, nebulaId account.NebulaId, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	helperWaves := helpers.New(waves.client.GetOptions().BaseUrl, "")
	state, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, fmt.Sprintf("%d", tcHeight))
	if err != nil {
		return "", err
	}
	if state == nil {
		funcArgs := new(proto.Arguments)
		funcArgs.Append(proto.BinaryArgument{
			Value: hash,
		})
		realSignCount := 0
		var signs []string
		oracles, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "oracles")
		if err != nil {
			return "", err
		}

		for _, oracle := range strings.Split(oracles.Value.(string), ",") {
			var pubKey account.OraclesPubKey
			copy(pubKey[:], base58.Decode(oracle))
			sign, err := ghClient.Result(account.Ethereum, nebulaId, int64(tcHeight), pubKey)
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

		bft, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "bft_coefficient")
		if err != nil {
			return "", err
		}

		if realSignCount == 0 {
			return "", nil
		}
		if realSignCount >= int(bft.Value.(float64)) {
			secret, err := wavesCrypto.NewSecretKeyFromBytes(privKey)

			asset, err := proto.NewOptionalAssetFromString("WAVES")
			if err != nil {
				return "", err
			}
			contract, err := proto.NewRecipientFromString(waves.contractAddress)
			if err != nil {
				return "", err
			}
			tx := &proto.InvokeScriptWithProofs{
				Type:            proto.InvokeScriptTransaction,
				Version:         1,
				SenderPK:        wavesCrypto.GeneratePublicKey(secret),
				ChainID:         waves.chainID,
				ScriptRecipient: contract,
				FunctionCall: proto.FunctionCall{
					Name:      "confirmData",
					Arguments: *funcArgs,
				},
				Payments:  nil,
				FeeAsset:  *asset,
				Fee:       500000,
				Timestamp: client.NewTimestampFromTime(time.Now()),
			}

			err = tx.Sign(waves.chainID, secret)
			if err != nil {
				return "", err
			}

			_, err = waves.client.Transactions.Broadcast(ctx, tx)
			if err != nil {
				return "", err
			}

			fmt.Printf("Tx finilize: %s \n", tx.ID)
			return tx.ID.String(), nil
		}
	}

	return "", nil
}

func (waves *Waves) SendSubs(tcHeight uint64, privKey []byte, value interface{}, ctx context.Context) error {
	helperWaves := helpers.New(waves.client.GetOptions().BaseUrl, "")
	state, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, fmt.Sprintf("%d", tcHeight))
	if err != nil {
		return err
	}
	if state != nil {
		subContract, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "subscriber_address")
		if err != nil {
			return err
		}

		pubKeyNebulaContract, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "contract_pubkey")
		if err != nil {
			return err
		}

		height, err := helperWaves.GetStateByAddressAndKey(subContract.Value.(string), fmt.Sprintf("%d", tcHeight))
		if err != nil {
			return err
		}

		if height != nil {
			return err
		}
		secret, err := wavesCrypto.NewSecretKeyFromBytes(privKey)

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

		nebulaType, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "type")
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
			ChainID:         waves.chainID,
			ScriptRecipient: contract,
			FunctionCall: proto.FunctionCall{
				Name:      "attachData",
				Arguments: args,
			},
			Payments:  nil,
			FeeAsset:  *asset,
			Fee:       500000,
			Timestamp: client.NewTimestampFromTime(time.Now()),
		}

		err = tx.Sign(waves.chainID, secret)
		if err != nil {
			return err
		}

		_, err = waves.client.Transactions.Broadcast(ctx, tx)
		if err != nil {
			return err
		}

		fmt.Printf("Sub send tx: %s \n", tx.ID.String())

	}
	return nil
}

func (waves *Waves) WaitTx(id string) error {
	helperWaves := helpers.New(waves.client.GetOptions().BaseUrl, "")
	return <-helperWaves.WaitTx(id)
}

package blockchain

import (
	"context"
	"fmt"
	"gravity-hub/common/keys"
	"gravity-hub/gh-node/api/gravity"
	"gravity-hub/gh-node/helpers"
	"strings"
	"time"

	"github.com/btcsuite/btcutil/base58"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/wavesplatform/gowaves/pkg/proto"

	"github.com/wavesplatform/gowaves/pkg/client"
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
		chainID:         'R',
	}, nil
}

func (waves *Waves) GetHeight(ctx context.Context) (uint64, error) {
	wavesHeight, _, err := waves.client.Blocks.Height(ctx)
	if err != nil {
		return 0, err
	}

	return wavesHeight.Height, nil
}

func (waves *Waves) SendResult(tcHeight uint64, privKey []byte, nebulaId []byte, ghClient *gravity.Client, validators [][]byte, hash []byte, ctx context.Context) (string, error) {
	helperWaves := helpers.New(waves.client.GetOptions().BaseUrl, "")
	state, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, fmt.Sprintf("%d", tcHeight))
	if err != nil {
		return "", err
	}
	if state == nil {
		funcArgs := new(proto.Arguments)
		funcArgs.Append(proto.StringArgument{
			Value: base58.Encode(hash),
		})
		realSignCount := 0
		var signs []string
		oracles, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "oracles")
		if err != nil {
			return "", err
		}
		for _, oracle := range strings.Split(oracles.Value.(string), ",") {
			pubKey := base58.Decode(oracle)
			sign, err := ghClient.GetKey(keys.FormSignResultKey(nebulaId, tcHeight, pubKey))
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
func (waves *Waves) SendSubs(tcHeight uint64, privKey []byte, value uint64, ctx context.Context) error {
	helperWaves := helpers.New(waves.client.GetOptions().BaseUrl, "")
	state, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, fmt.Sprintf("%d", tcHeight))
	if err != nil {
		return err
	}
	if state == nil {
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

		tx := &proto.InvokeScriptWithProofs{
			Type:            proto.InvokeScriptTransaction,
			Version:         1,
			SenderPK:        pubKey,
			ChainID:         waves.chainID,
			ScriptRecipient: contract,
			FunctionCall: proto.FunctionCall{
				Name: "attachData",
				Arguments: proto.Arguments{
					proto.IntegerArgument{
						Value: int64(value),
					},
					proto.IntegerArgument{
						Value: int64(tcHeight),
					},
				},
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

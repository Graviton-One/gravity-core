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
	}, nil
}

func (waves *Waves) GetHeight(ctx context.Context) (uint64, error) {
	wavesHeight, _, err := waves.client.Blocks.Height(ctx)
	if err != nil {
		return 0, err
	}

	return wavesHeight.Height, nil
}

func (waves *Waves) SendResult(tcHeight uint64, privKey []byte, nebulaId []byte, ghClient *gravity.Client, validators [][]byte, hash []byte, ctx context.Context) error {
	helperWaves := helpers.New(waves.client.GetOptions().BaseUrl, "")
	state, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, fmt.Sprintf("%d", tcHeight))
	if err != nil {
		return err
	}
	if state == nil {
		funcArgs := new(proto.Arguments)
		funcArgs.Append(proto.StringArgument{
			Value: base58.Encode(hash),
		})
		bft := int(float32(len(validators)) * 0.7)
		realSignCount := 0
		var signs []string
		oracles, err := helperWaves.GetStateByAddressAndKey(waves.contractAddress, "oracles")
		if err != nil {
			return err
		}
		for _, oracle := range strings.Split(oracles.Value.(string), ",") {
			pubKey := base58.Decode(oracle)
			sign, err := ghClient.GetKey(keys.FormSignResultKey(nebulaId, tcHeight, pubKey), ctx)
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

		if realSignCount >= bft {
			secret, err := wavesCrypto.NewSecretKeyFromBytes(privKey)

			asset, err := proto.NewOptionalAssetFromString("WAVES")
			if err != nil {
				return err
			}
			contract, err := proto.NewRecipientFromString(waves.contractAddress)
			if err != nil {
				return err
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

			err = tx.Sign('T', secret)
			if err != nil {
				return err
			}

			_, err = waves.client.Transactions.Broadcast(ctx, tx)
			if err != nil {
				return err
			}

			fmt.Printf("Tx finilize: %s \n", tx.ID)
		}
	}

	return nil
}
func (waves *Waves) SendSubs(tcHeight uint64, privKey []byte, value uint64) error {
	return nil
}

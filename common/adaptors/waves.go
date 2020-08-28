package adaptors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/helpers"

	"github.com/Gravity-Tech/gravity-core/common/client"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/btcsuite/btcutil/base58"
	wavesClient "github.com/wavesplatform/gowaves/pkg/client"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type WavesAdaptor struct {
	secret wavesCrypto.SecretKey

	ghClient    *client.GravityClient
	wavesClient *wavesClient.Client
	helper      helpers.ClientHelper

	gravityContract string
	chainID         byte
}
type WavesAdapterOption func(*WavesAdaptor) error

func WithWavesGravityContract(address string) WavesAdapterOption {
	return func(h *WavesAdaptor) error {
		h.gravityContract = address
		return nil
	}
}
func WavesAdapterWithGhClient(ghClient *client.GravityClient) WavesAdapterOption {
	return func(h *WavesAdaptor) error {
		h.ghClient = ghClient
		return nil
	}
}

func NewWavesAdapter(seed []byte, nodeUrl string, opts ...WavesAdapterOption) (*WavesAdaptor, error) {
	wClient, err := wavesClient.NewClient(wavesClient.Options{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	secret, err := wavesCrypto.NewSecretKeyFromBytes(seed)
	adapter := &WavesAdaptor{
		secret:      secret,
		wavesClient: wClient,
		helper:      helpers.NewClientHelper(wClient),
	}
	for _, opt := range opts {
		err := opt(adapter)
		if err != nil {
			return nil, err
		}
	}
	return adapter, nil
}

func (adaptor *WavesAdaptor) GetHeight(ctx context.Context) (uint64, error) {
	wavesHeight, _, err := adaptor.wavesClient.Blocks.Height(ctx)
	if err != nil {
		return 0, err
	}

	return wavesHeight.Height, nil
}
func (adaptor *WavesAdaptor) WaitTx(id string, ctx context.Context) error {
	return <-adaptor.helper.WaitTx(id, ctx)
}
func (adaptor *WavesAdaptor) Sign(msg []byte) ([]byte, error) {
	sig, err := wavesCrypto.Sign(adaptor.secret, msg)
	if err != nil {
		return nil, err
	}
	return sig.Bytes(), nil
}
func (adaptor *WavesAdaptor) PubKey() account.OraclesPubKey {
	var oraclePubKey account.OraclesPubKey
	pubKey := wavesCrypto.GeneratePublicKey(adaptor.secret)
	copy(oraclePubKey[:], pubKey.Bytes())
	return oraclePubKey
}

func (adaptor *WavesAdaptor) SendDataResult(nebulaId account.NebulaId, tcHeight uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	nebulaAddress := base58.Encode(nebulaId)

	state, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, fmt.Sprintf("%d", tcHeight), ctx)
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
	oracles, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "oracles", ctx)
	if err != nil {
		return "", err
	}

	for _, oracle := range strings.Split(oracles.Value.(string), ",") {
		var pubKey account.OraclesPubKey
		copy(pubKey[:], base58.Decode(oracle))
		sign, err := adaptor.ghClient.Result(account.Ethereum, nebulaId, int64(tcHeight), pubKey)
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

	bft, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "bft_coefficient", ctx)
	if err != nil {
		return "", err
	}

	if realSignCount == 0 {
		return "", nil
	}
	if realSignCount < int(bft.Value.(float64)) {
		return "", nil
	}

	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		return "", err
	}

	contract, err := proto.NewRecipientFromString(nebulaAddress)
	if err != nil {
		return "", err
	}
	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        wavesCrypto.GeneratePublicKey(adaptor.secret),
		ChainID:         adaptor.chainID,
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

	err = tx.Sign(adaptor.chainID, adaptor.secret)
	if err != nil {
		return "", err
	}

	_, err = adaptor.wavesClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func (adaptor *WavesAdaptor) SendDataToSubs(nebulaId account.NebulaId, tcHeight uint64, value interface{}, ctx context.Context) error {
	nebulaAddress := base58.Encode(nebulaId)
	state, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, fmt.Sprintf("%d", tcHeight), ctx)
	if err != nil {
		return err
	} else if state == nil {
		return nil
	}

	subContract, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "subscriber_address", ctx)
	if err != nil {
		return err
	}

	pubKeyNebulaContract, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "contract_pubkey", ctx)
	if err != nil {
		return err
	}

	height, _, err := adaptor.helper.GetStateByAddressAndKey(subContract.Value.(string), fmt.Sprintf("%d", tcHeight), ctx)
	if err != nil {
		return err
	}

	if height != nil {
		return err
	}

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

	nebulaType, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "type", ctx)
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
		ChainID:         adaptor.chainID,
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

	err = tx.Sign(adaptor.chainID, adaptor.secret)
	if err != nil {
		return err
	}

	_, err = adaptor.wavesClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}

func (adaptor *WavesAdaptor) SendOraclesToNebula(nebulaId account.NebulaId, oracles []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error) {
	nebulaAddress := base58.Encode(nebulaId)
	lastRoundState, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "last_round_"+fmt.Sprintf("%d", round), ctx)
	if err != nil {
		return "", err
	}

	if lastRoundState != nil {
		return "", err
	}

	oracleCountState, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "oracle_count", ctx)
	if err != nil {
		return "", err
	}

	var newOracles []string
	var stringSigns []string
	for _, v := range oracles {
		newOracles = append(newOracles, base58.Encode(v.ToBytes(account.Waves)))
	}
	for _, v := range signs {
		stringSigns = append(stringSigns, base58.Encode(v))
	}

	emptyCount := oracleCountState.Value.(int) - len(newOracles)
	for i := 0; i < emptyCount; i++ {
		newOracles = append(newOracles, base58.Encode([]byte{0}))
	}

	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		return "", err
	}

	contract, err := proto.NewRecipientFromString(nebulaAddress)
	if err != nil {
		return "", err
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        wavesCrypto.GeneratePublicKey(adaptor.secret),
		ChainID:         adaptor.chainID,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name: "setSortedOracles",
			Arguments: proto.Arguments{
				proto.StringArgument{
					Value: strings.Join(newOracles, ","),
				},
				proto.StringArgument{
					Value: strings.Join(stringSigns, ","),
				},
				proto.IntegerArgument{
					Value: round,
				},
			},
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
	}

	err = tx.Sign(adaptor.chainID, adaptor.secret)
	if err != nil {
		return "", err
	}

	_, err = adaptor.wavesClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func (adaptor *WavesAdaptor) SendConsulsToGravityContract(newConsulsAddresses []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error) {
	state, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, "last_round_"+fmt.Sprintf("%d", round), ctx)
	if err != nil {
		return "", err
	}
	if state != nil {
		return "", err
	}

	var stringSigns []string
	oneSigFound := false
	for _, v := range signs {
		oneSigFound = true
		stringSigns = append(stringSigns, base58.Encode(v))
	}

	var newConsulsString []string

	for _, v := range newConsulsAddresses {
		newConsulsString = append(newConsulsString, base58.Encode(v.ToBytes(account.Waves)))
	}

	consulsCountState, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, "consuls_count", ctx)
	if err != nil {
		return "", err
	}

	emptyCount := consulsCountState.Value.(int) - len(newConsulsString)
	for i := 0; i < emptyCount; i++ {
		newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
	}
	if !oneSigFound {
		return "", nil
	}

	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		return "", err
	}

	contract, err := proto.NewRecipientFromString(adaptor.gravityContract)
	if err != nil {
		return "", err
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        wavesCrypto.GeneratePublicKey(adaptor.secret),
		ChainID:         adaptor.chainID,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name: "setConsuls",
			Arguments: proto.Arguments{
				proto.StringArgument{
					Value: strings.Join(newConsulsString, ","),
				},
				proto.StringArgument{
					Value: strings.Join(stringSigns, ","),
				},
				proto.IntegerArgument{
					Value: round,
				},
			},
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: wavesClient.NewTimestampFromTime(time.Now()),
	}

	err = tx.Sign(adaptor.chainID, adaptor.secret)
	if err != nil {
		return "", err
	}

	_, err = adaptor.wavesClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	err = <-adaptor.helper.WaitTx(tx.ID.String(), ctx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}
func (adaptor *WavesAdaptor) SignConsuls(consulsAddresses []account.OraclesPubKey) ([]byte, error) {
	return adaptor.signOraclesPubKey(consulsAddresses)
}
func (adaptor *WavesAdaptor) SignOracles(nebulaId account.NebulaId, oracles []account.OraclesPubKey) ([]byte, error) {
	return adaptor.signOraclesPubKey(oracles)
}

func (adaptor *WavesAdaptor) signOraclesPubKey(consulsAddresses []account.OraclesPubKey) ([]byte, error) {
	var stringOracles []string
	for _, v := range consulsAddresses {
		stringOracles = append(stringOracles, base58.Encode(v.ToBytes(account.Waves)))
	}

	sign, err := adaptor.Sign([]byte(strings.Join(stringOracles, ",")))
	if err != nil {
		return nil, err
	}

	return sign, nil
}

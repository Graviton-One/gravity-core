package adaptors

import (
	"context"
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"github.com/gookit/validate"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/storage"

	"github.com/Gravity-Tech/gravity-core/common/helpers"

	"github.com/Gravity-Tech/gravity-core/common/gravity"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/btcsuite/btcutil/base58"
	wclient "github.com/wavesplatform/gowaves/pkg/client"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

const (
	ConsulsCount = 5
)

type WavesAdaptor struct {
	secret crypto.SecretKey

	ghClient        *gravity.Client      `option:"ghClient"`
	wavesClient     *wclient.Client      `option:"wvClient"`
	helper          helpers.ClientHelper `option:"-"`
	gravityContract string               `option:"gravityContract"`
	chainID         byte                 `option:"chainID"`
}
type WavesAdapterOption func(*WavesAdaptor) error

func (wa *WavesAdaptor) applyOpts(opts AdapterOptions) error {
	err := validateWavesAdapterOptions(opts)
	if err != nil {
		return err
	}
	v := reflect.TypeOf(*wa)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("option")
		val, ok := opts[tag]
		if ok {
			switch tag {
			case "ghClient":
				wa.ghClient = val.(*gravity.Client)
			case "wvClient":
				wa.wavesClient = val.(*wclient.Client)
			case "gravityContract":
				wa.gravityContract = val.(string)
			case "chainID":
				wa.chainID = val.(byte)
			}
		}
	}
	return nil
}

func validateWavesAdapterOptions(opts AdapterOptions) error {
	v := validate.Map(opts)

	v.AddRule("chainID", "isByte")
	v.AddRule("ghClient", "isGhClient")
	v.AddRule("wvClient", "isWvClient")
	v.AddRule("gravityContract", "string")

	if !v.Validate() { // validate ok
		return v.Errors
	}
	return nil
}

func WithWavesGravityContract(address string) WavesAdapterOption {
	return func(h *WavesAdaptor) error {
		h.gravityContract = address
		return nil
	}
}
func WavesAdapterWithGhClient(ghClient *gravity.Client) WavesAdapterOption {
	return func(h *WavesAdaptor) error {
		h.ghClient = ghClient
		return nil
	}
}
func NewWavesAdapterByOpts(seed []byte, nodeUrl string, opts AdapterOptions) (*WavesAdaptor, error) {
	wClient, err := wclient.NewClient(wclient.Options{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	secret, err := crypto.NewSecretKeyFromBytes(seed)
	adapter := &WavesAdaptor{
		secret:      secret,
		wavesClient: wClient,
		helper:      helpers.NewClientHelper(wClient),
	}
	err = adapter.applyOpts(opts)
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

func NewWavesAdapter(seed []byte, nodeUrl string, chainId byte, opts ...WavesAdapterOption) (*WavesAdaptor, error) {
	wClient, err := wclient.NewClient(wclient.Options{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	secret, err := crypto.NewSecretKeyFromBytes(seed)
	adapter := &WavesAdaptor{
		secret:      secret,
		wavesClient: wClient,
		helper:      helpers.NewClientHelper(wClient),
		chainID:     chainId,
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
	sig, err := crypto.Sign(adaptor.secret, msg)
	if err != nil {
		return nil, err
	}
	return sig.Bytes(), nil
}
func (adaptor *WavesAdaptor) SignHash(nebulaId account.NebulaId, intervalId uint64, pulseId uint64, hash []byte) ([]byte, error) {
	return adaptor.Sign(hash)
}
func (adaptor *WavesAdaptor) PubKey() account.OraclesPubKey {
	pubKey := crypto.GeneratePublicKey(adaptor.secret)
	oraclePubKey := account.BytesToOraclePubKey(pubKey[:], account.Waves)
	return oraclePubKey
}
func (adaptor *WavesAdaptor) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
	nebulaAddress := base58.Encode(nebulaId.ToBytes(account.Waves))
	state, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "type", ctx)
	if err != nil {
		return 0, err
	}

	return abi.ExtractorType(state.Value.(float64)), nil
}

func (adaptor *WavesAdaptor) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	nebulaAddress := base58.Encode(nebulaId.ToBytes(account.Waves))

	state, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, fmt.Sprintf("data_hash_%d", pulseId), ctx)
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
		pubKey, err := account.StringToOraclePubKey(oracle, account.Waves)
		if err != nil {
			signs = append(signs, "nil")
			continue
		}

		sign, err := adaptor.ghClient.Result(account.Waves, nebulaId, int64(pulseId), pubKey)
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
		SenderPK:        crypto.GeneratePublicKey(adaptor.secret),
		ChainID:         adaptor.chainID,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name:      "sendHashValue",
			Arguments: *funcArgs,
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: wclient.NewTimestampFromTime(time.Now()),
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
func (adaptor *WavesAdaptor) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
	nebulaAddress := base58.Encode(nebulaId.ToBytes(account.Waves))
	zap.L().Sugar().Debugf("SendValueToSubs: nebulaAddress - %s", nebulaAddress)
	state, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, fmt.Sprintf("data_hash_%d", pulseId), ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	} else if state == nil {
		zap.L().Debug("SendValueToSubs: state is nil")
		return nil
	}

	subContract, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "subscriber_address", ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debug("SendValueToSubs: subcontract ", subContract)

	pubKeyNebulaContract, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "contract_pubkey", ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debug("SendValueToSubs: pubKeyNebulaContract ", pubKeyNebulaContract)

	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	contract, err := proto.NewRecipientFromString(subContract.Value.(string))
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debug("SendValueToSubs: contract ", contract)

	pubKey, err := crypto.NewPublicKeyFromBase58(pubKeyNebulaContract.Value.(string))
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	nebulaType, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "type", ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	args := proto.Arguments{}
	switch SubType(int8(nebulaType.Value.(float64))) {
	case Int64:
		v, err := strconv.ParseInt(value.Value, 10, 64)
		if err != nil {
			return err
		}
		args.Append(
			proto.IntegerArgument{
				Value: v,
			})
	case String:
		args.Append(
			proto.StringArgument{
				Value: value.Value,
			})
	case Bytes:
		v, err := base64.StdEncoding.DecodeString(value.Value)
		if err != nil {
			return err
		}
		args.Append(
			proto.BinaryArgument{
				Value: v,
			})
	}
	args.Append(
		proto.IntegerArgument{
			Value: int64(pulseId),
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
		Fee:       900000,
		Timestamp: wclient.NewTimestampFromTime(time.Now()),
	}
	zap.L().Sugar().Debug("SendValueToSubs: tx ", tx)
	err = tx.Sign(adaptor.chainID, adaptor.secret)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	zap.L().Sugar().Debug("SendValueToSubs: Broadcast ", tx)
	_, err = adaptor.wavesClient.Transactions.Broadcast(ctx, tx)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

func (adaptor *WavesAdaptor) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	defer func() {
		if r := recover(); r != nil {
			zap.L().Sugar().Error("Recovered in SignOracles", r)
		}
	}()
	nebulaAddress := base58.Encode(nebulaId.ToBytes(account.Waves))
	lastRoundState, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "last_round_"+fmt.Sprintf("%d", round), ctx)
	if err != nil {
		return "", err
	}

	if lastRoundState != nil {
		return "", err
	}

	var newOracles []string
	var stringSigns [5]string

	consulsState, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, fmt.Sprintf("consuls_%d", round), ctx)
	if err != nil {
		return "", err
	}

	consuls := strings.Split(consulsState.Value.(string), ",")
	for k, v := range signs {
		pubKey := k.ToString(account.Waves)
		index := -1

		for i, v := range consuls {
			if v == pubKey {
				index = i
				break
			}
		}

		if index == -1 {
			continue
		}

		stringSigns[index] = base58.Encode(v)
	}
	for i, v := range stringSigns {
		if v != "" {
			continue
		}

		stringSigns[i] = base58.Encode([]byte{0})
	}

	for _, v := range oracles {
		if v == nil {
			newOracles = append(newOracles, base58.Encode([]byte{1}))
			continue
		}
		newOracles = append(newOracles, base58.Encode(v.ToBytes(account.Waves)))
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
		SenderPK:        crypto.GeneratePublicKey(adaptor.secret),
		ChainID:         adaptor.chainID,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name: "updateOracles",
			Arguments: proto.Arguments{
				proto.StringArgument{
					Value: strings.Join(newOracles, ","),
				},
				proto.StringArgument{
					Value: strings.Join(stringSigns[:], ","),
				},
				proto.IntegerArgument{
					Value: round,
				},
			},
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: wclient.NewTimestampFromTime(time.Now()),
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
func (adaptor *WavesAdaptor) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	var stringSigns [5]string
	lastRoundState, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, "last_round", ctx)
	if err != nil {
		return "", err
	}

	lastRound := uint64(lastRoundState.Value.(float64))
	consulsState, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, fmt.Sprintf("consuls_%d", lastRound), ctx)
	if err != nil {
		return "", err
	}

	consuls := strings.Split(consulsState.Value.(string), ",")
	for k, v := range signs {
		pubKey := k.ToString(account.Waves)
		index := -1

		for i, v := range consuls {
			if v == pubKey {
				index = i
				break
			}
		}

		if index == -1 {
			continue
		}

		stringSigns[index] = base58.Encode(v)
	}

	for i, v := range stringSigns {
		if v != "" {
			continue
		}

		stringSigns[i] = base58.Encode([]byte{0})
	}

	var newConsulsString []string

	for _, v := range newConsulsAddresses {
		if v == nil {
			newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
			continue
		}
		newConsulsString = append(newConsulsString, base58.Encode(v.ToBytes(account.Waves)))
	}

	emptyCount := ConsulsCount - len(newConsulsString)
	for i := 0; i < emptyCount; i++ {
		newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
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
		SenderPK:        crypto.GeneratePublicKey(adaptor.secret),
		ChainID:         adaptor.chainID,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name: "updateConsuls",
			Arguments: proto.Arguments{
				proto.StringArgument{
					Value: strings.Join(newConsulsString, ","),
				},
				proto.StringArgument{
					Value: strings.Join(stringSigns[:], ","),
				},
				proto.IntegerArgument{
					Value: round,
				},
			},
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: wclient.NewTimestampFromTime(time.Now()),
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
func (adaptor *WavesAdaptor) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64, sender account.OraclesPubKey) ([]byte, error) {
	var msg []string
	for _, v := range consulsAddresses {
		if v == nil {
			msg = append(msg, base58.Encode([]byte{0}))
			continue
		}
		msg = append(msg, base58.Encode(v.ToBytes(account.Waves)))
	}
	msg = append(msg, fmt.Sprintf("%d", roundId))

	sign, err := adaptor.Sign([]byte(strings.Join(msg, ",")))
	if err != nil {
		return nil, err
	}

	return sign, err
}
func (adaptor *WavesAdaptor) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, round int64, sender account.OraclesPubKey) ([]byte, error) {
	var stringOracles []string
	for _, v := range oracles {
		if v == nil {
			stringOracles = append(stringOracles, base58.Encode([]byte{1}))
			continue
		}
		stringOracles = append(stringOracles, base58.Encode(v.ToBytes(account.Waves)))
	}

	sign, err := adaptor.Sign([]byte(strings.Join(stringOracles, ",")))
	if err != nil {
		return nil, err
	}

	return sign, err
}

func (adaptor *WavesAdaptor) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	nebulaAddress := base58.Encode(nebulaId.ToBytes(account.Waves))
	state, _, err := adaptor.helper.GetStateByAddressAndKey(nebulaAddress, "last_pulse_id", ctx)
	if err != nil && err != storage.ErrKeyNotFound {
		return 0, err
	}

	if err == storage.ErrKeyNotFound || state == nil {
		return 0, nil
	}

	return uint64(state.Value.(float64)), nil
}
func (adaptor *WavesAdaptor) LastRound(ctx context.Context) (uint64, error) {
	state, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, "last_round", ctx)
	if err != nil && err != storage.ErrKeyNotFound {
		return 0, err
	}

	if state == nil {
		return 0, err
	}

	return uint64(state.Value.(float64)), nil
}
func (adaptor *WavesAdaptor) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	state, _, err := adaptor.helper.GetStateByAddressAndKey(adaptor.gravityContract, fmt.Sprintf("consuls_%d", roundId), ctx)
	if err != nil {
		return false, err
	}
	if state == nil {
		return false, nil
	} else {
		return true, nil
	}
}

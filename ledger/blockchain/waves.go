package blockchain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/helpers/waves"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/btcsuite/btcutil/base58"
	"github.com/wavesplatform/gowaves/pkg/client"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"
)

type Waves struct {
	privKeyBytes           []byte
	secret                 wavesCrypto.SecretKey
	client                 *client.Client
	helper                 waves.ClientHelper
	gravityContractAddress string
	chainID                byte
}

func NewWaves(gravityContractAddress string, privKey []byte, nodeUrl string) (*Waves, error) {
	wavesClient, err := client.NewClient(client.Options{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	return &Waves{
		privKeyBytes:           privKey,
		helper:                 waves.NewClientHelper(wavesClient),
		client:                 wavesClient,
		gravityContractAddress: gravityContractAddress,
		chainID:                'T',
	}, nil
}

func (waves *Waves) SendOraclesToNebula(nebulaId account.NebulaId, oracles []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error) {
	nebulaAddress := base58.Encode(nebulaId)
	lastRoundState, _, err := waves.helper.GetStateByAddressAndKey(nebulaAddress, "last_round_"+fmt.Sprintf("%d", round), ctx)
	if err != nil {
		return "", err
	}

	if lastRoundState != nil {
		return "", err
	}

	oracleCountState, _, err := waves.helper.GetStateByAddressAndKey(nebulaAddress, "oracle_count", ctx)
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
		SenderPK:        wavesCrypto.GeneratePublicKey(waves.secret),
		ChainID:         waves.chainID,
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
		Timestamp: client.NewTimestampFromTime(time.Now()),
	}

	err = tx.Sign(waves.chainID, waves.secret)
	if err != nil {
		return "", err
	}

	_, err = waves.client.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}

func (waves *Waves) SendConsulsToGravityContract(newConsulsAddresses []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error) {
	state, _, err := waves.helper.GetStateByAddressAndKey(waves.gravityContractAddress, "last_round_"+fmt.Sprintf("%d", round), ctx)
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

	consulsCountState, _, err := waves.helper.GetStateByAddressAndKey(waves.gravityContractAddress, "consuls_count", ctx)
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

	contract, err := proto.NewRecipientFromString(waves.gravityContractAddress)
	if err != nil {
		return "", err
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        wavesCrypto.GeneratePublicKey(waves.secret),
		ChainID:         waves.chainID,
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
		Timestamp: client.NewTimestampFromTime(time.Now()),
	}

	err = tx.Sign(waves.chainID, waves.secret)
	if err != nil {
		return "", err
	}

	_, err = waves.client.Transactions.Broadcast(ctx, tx)
	if err != nil {
		return "", err
	}

	err = <-waves.helper.WaitTx(tx.ID.String(), ctx)
	if err != nil {
		return "", err
	}

	return tx.ID.String(), nil
}

func (waves *Waves) SignConsuls(consulsAddresses []account.OraclesPubKey) ([]byte, error) {
	return waves.sign(consulsAddresses)
}

func (waves *Waves) SignOracles(nebulaId account.NebulaId, oracles []account.OraclesPubKey) ([]byte, error) {
	return waves.sign(oracles)
}

func (waves *Waves) PubKey() []byte {
	pubKey := wavesCrypto.GeneratePublicKey(waves.secret)
	return pubKey.Bytes()
}

func (waves *Waves) sign(consulsAddresses []account.OraclesPubKey) ([]byte, error) {
	var stringOracles []string
	for _, v := range consulsAddresses {
		stringOracles = append(stringOracles, base58.Encode(v.ToBytes(account.Waves)))
	}

	sign, err := waves.signMsg([]byte(strings.Join(stringOracles, ",")))
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func (waves *Waves) signMsg(msg []byte) ([]byte, error) {
	secret, err := wavesCrypto.NewSecretKeyFromBytes(waves.privKeyBytes)
	if err != nil {
		return nil, err
	}
	sig, err := wavesCrypto.Sign(secret, msg)
	if err != nil {
		return nil, err
	}
	return sig.Bytes(), nil
}

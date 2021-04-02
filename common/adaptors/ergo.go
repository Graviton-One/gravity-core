package adaptors

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/abi/ethereum"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gookit/validate"
	wclient "github.com/wavesplatform/gowaves/pkg/client"
	"go.uber.org/zap"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

const (
	Consuls = 5
)

type ErgoAdaptor struct {
	secret 			crypto.SecretKey
	ghClient        *gravity.Client      `option:"ghClient"`
	//wavesClient     *wclient.Client      `option:"wvClient"`
	//helper          helpers.ClientHelper `option:"-"`
	gravityContract string               `option:"gravityContract"`
}

type ErgoAdapterOption func(*ErgoAdaptor) error

func (er *ErgoAdaptor) applyOpts(opts AdapterOptions) error {
	err := validateErgoAdapterOptions(opts)
	if err != nil {
		return err
	}
	v := reflect.TypeOf(*er)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("option")
		val, ok := opts[tag]
		if ok {
			switch tag {
			case "ghClient":
				er.ghClient = val.(*gravity.Client)
			case "gravityContract":
				er.gravityContract = val.(string)

			}
		}
	}
	return nil
}

func validateErgoAdapterOptions(opts AdapterOptions) error {
	v := validate.Map(opts)

	v.AddRule("ghClient", "isGhClient")
	v.AddRule("gravityContract", "string")

	if !v.Validate() { // validate ok
		return v.Errors
	}
	return nil
}

func WithErgoGravityContract(address string) ErgoAdapterOption {
	return func(h *ErgoAdaptor) error {
		h.gravityContract = address
		return nil
	}
}

func ErgoAdapterWithGhClient(ghClient *gravity.Client) ErgoAdapterOption {
	return func(h *ErgoAdaptor) error {
		h.ghClient = ghClient
		return nil
	}
}

func NewErgoAdapterByOpts(seed []byte, nodeUrl string, opts AdapterOptions) (*ErgoAdaptor, error) {
	//wClient, err := wclient.NewClient(wclient.Options{ApiKey: "", BaseUrl: nodeUrl})
	//if err != nil {
	//	return nil, err
	//}

	secret, err := crypto.NewSecretKeyFromBytes(seed)
	adapter := &ErgoAdaptor{
		secret:      secret,
	}
	err = adapter.applyOpts(opts)
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

func NewErgoAdapter(seed []byte, nodeUrl string, ctx context.Context, opts ...ErgoAdapterOption) (*ErgoAdaptor, error) {
	//wClient, err := wclient.NewClient(wclient.Options{ApiKey: "", BaseUrl: nodeUrl})
	//if err != nil {
	//	return nil, err
	//}

	secret, err := crypto.NewSecretKeyFromBytes(seed)
	if err != nil {
		return nil, err
	}
	adapter := &ErgoAdaptor{
		secret:      secret,

	}
	for _, opt := range opts {
		err := opt(adapter)
		if err != nil {
			return nil, err
		}
	}
	return adapter, nil
}


func (adaptor *ErgoAdaptor) GetHeight(ctx context.Context) (uint64, error) {
	type Response struct {
		Status  bool    `json:"name"`
		Height	uint64	`json:"pokemon_entries"`
	}
	res, err := http.NewRequest("GET", "http://127.0.0.1:9000/height", nil)
	if err != nil {
		return 0, err
	}
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)
	if !responseObject.Status {
		err = fmt.Errorf("proxy connection problem")
	}
	return responseObject.Height, err
}
func (adaptor *ErgoAdaptor) WaitTx(id string, ctx context.Context) error {
	return <-adaptor.helper.WaitTx(id, ctx)
}
func (adaptor *ErgoAdaptor) Sign(msg []byte) ([]byte, error) {
	sig, err := crypto.Sign(adaptor.secret, msg)
	if err != nil {
		return nil, err
	}
	return sig.Bytes(), nil
}
func (adaptor *ErgoAdaptor) PubKey() account.OraclesPubKey {
	pubKey := crypto.GeneratePublicKey(adaptor.secret)
	oraclePubKey := account.BytesToOraclePubKey(pubKey[:], account.Ergo)
	return oraclePubKey
}

func (adaptor *ErgoAdaptor) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
	nebula, err := ethereum.NewNebula(common.BytesToAddress(nebulaId.ToBytes(account.Ethereum)), adaptor.ethClient)
	if err != nil {
		return 0, err
	}

	exType, err := nebula.DataType(nil)
	if err != nil {
		return 0, err
	}

	return abi.ExtractorType(exType), nil
}

func (adaptor *ErgoAdaptor) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	nebula, err := ethereum.NewNebula(common.BytesToAddress(nebulaId.ToBytes(account.Ethereum)), adaptor.ethClient)
	if err != nil {
		return "", err
	}

	data, err := nebula.Pulses(nil, big.NewInt(int64(pulseId)))
	if err != nil {
		return "", err
	}

	if bytes.Equal(data.DataHash[:], make([]byte, 32, 32)) != true {
		return "", nil
	}

	bft, err := nebula.BftValue(nil)
	if err != nil {
		return "", err
	}

	realSignCount := 0

	oracles, err := nebula.GetOracles(nil)
	if err != nil {
		return "", err
	}
	var r [5][32]byte
	var s [5][32]byte
	var v [5]uint8
	for _, validator := range validators {
		pubKey, err := crypto.DecompressPubkey(validator.ToBytes(account.Ethereum))
		if err != nil {
			return "", err
		}
		validatorAddress := crypto.PubkeyToAddress(*pubKey)
		position := 0
		isExist := false
		for i, address := range oracles {
			if validatorAddress == address {
				position = i
				isExist = true
				break
			}
		}
		if !isExist {
			continue
		}

		sign, err := adaptor.ghClient.Result(account.Ethereum, nebulaId, int64(pulseId), validator)
		if err != nil {
			r[position] = [32]byte{}
			s[position] = [32]byte{}
			v[position] = byte(0)
			continue
		}
		copy(r[position][:], sign[:32])
		copy(s[position][:], sign[32:64])
		v[position] = sign[64] + 27

		realSignCount++
	}

	if realSignCount < int(bft.Uint64()) {
		return "", nil
	}

	var resultBytes32 [32]byte
	copy(resultBytes32[:], hash)

	opt := bind.NewKeyedTransactor(adaptor.privKey)

	opt.GasPrice, err = adaptor.ethClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", err
	}

	opt.GasPrice.Mul(opt.GasPrice, big.NewInt(2))
	tx, err := nebula.SendHashValue(opt, resultBytes32, v[:], r[:], s[:])
	if err != nil {
		return "", err
	}
	return tx.Hash().String(), nil
}
func (adaptor *ErgoAdaptor) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
	var err error

	nebula, err := ethereum.NewNebula(common.BytesToAddress(nebulaId.ToBytes(account.Ethereum)), adaptor.ethClient)
	if err != nil {
		return err
	}

	ids, err := nebula.GetSubscribersIds(nil)
	if err != nil {
		return err
	}

	for _, id := range ids {
		zap.L().Sugar().Debug("IDs iterate", id)
		t, err := nebula.DataType(nil)
		if err != nil {
			return err
		}

		transactOpt := bind.NewKeyedTransactor(adaptor.privKey)
		zap.L().Sugar().Debug("transactOpt is nil", transactOpt == nil)

		switch SubType(t) {
		case Int64:
			zap.L().Sugar().Debugf("SendIntValueToSubs")
			v, err := strconv.ParseInt(value.Value, 10, 64)
			if err != nil {
				return err
			}
			_, err = nebula.SendValueToSubInt(transactOpt, v, big.NewInt(int64(pulseId)), id)
			if err != nil {
				return err
			}
		case String:
			zap.L().Sugar().Debugf("SendStringValueToSubs")
			_, err = nebula.SendValueToSubString(transactOpt, value.Value, big.NewInt(int64(pulseId)), id)
			if err != nil {
				return err
			}
		case Bytes:
			//println(value.Value)
			v, err := base64.StdEncoding.DecodeString(value.Value)
			if err != nil {
				return err
			}

			_, err = nebula.SendValueToSubByte(transactOpt, v, big.NewInt(int64(pulseId)), id)
			if err != nil {
				zap.L().Error(err.Error())
				continue
			}
		}
	}
	return nil
}

func (adaptor *ErgoAdaptor) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	nebula, err := ethereum.NewNebula(common.BytesToAddress(nebulaId.ToBytes(account.Ethereum)), adaptor.ethClient)
	if err != nil {
		return "", err
	}

	lastRound, err := nebula.Rounds(nil, big.NewInt(round))
	if err != nil {
		return "", err
	}

	if lastRound {
		return "", err
	}

	var oraclesAddresses []common.Address
	for _, v := range oracles {
		if v == nil {
			oraclesAddresses = append(oraclesAddresses, common.Address{})
			continue
		}

		pubKey, err := crypto.DecompressPubkey(v.ToBytes(account.Ethereum))
		if err != nil {
			return "", err
		}
		oraclesAddresses = append(oraclesAddresses, crypto.PubkeyToAddress(*pubKey))
	}

	consuls, err := adaptor.gravityContract.GetConsuls(nil)
	if err != nil {
		return "", err
	}

	var r [5][32]byte
	var s [5][32]byte
	var v [5]uint8
	for pubKey, sign := range signs {
		index := -1
		ethPubKey, err := crypto.DecompressPubkey(pubKey.ToBytes(account.Ethereum))
		if err != nil {
			return "", err
		}
		validatorAddress := crypto.PubkeyToAddress(*ethPubKey)
		for i, v := range consuls {
			if v == validatorAddress {
				index = i
				break
			}
		}

		if index == -1 {
			continue
		}

		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])

		r[index] = bytes32R
		s[index] = bytes32S
		v[index] = sign[64:][0] + 27
	}

	tx, err := nebula.UpdateOracles(bind.NewKeyedTransactor(adaptor.privKey), oraclesAddresses, v[:], r[:], s[:], big.NewInt(round))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
func (adaptor *ErgoAdaptor) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	consuls, err := adaptor.gravityContract.GetConsuls(nil)
	if err != nil {
		return "", err
	}

	var consulsAddress []common.Address
	for _, v := range newConsulsAddresses {
		if v == nil {
			consulsAddress = append(consulsAddress, common.Address{})
			continue
		}
		pubKey, err := crypto.DecompressPubkey(v.ToBytes(account.Ethereum))
		if err != nil {
			return "", err
		}
		consulsAddress = append(consulsAddress, crypto.PubkeyToAddress(*pubKey))
	}

	var r [5][32]byte
	var s [5][32]byte
	var v [5]uint8
	for pubKey, sign := range signs {
		index := -1
		ethPubKey, err := crypto.DecompressPubkey(pubKey.ToBytes(account.Ethereum))
		if err != nil {
			return "", err
		}
		validatorAddress := crypto.PubkeyToAddress(*ethPubKey)
		for i, v := range consuls {
			if v == validatorAddress {
				index = i
				break
			}
		}

		if index == -1 {
			continue
		}

		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])

		r[index] = bytes32R
		s[index] = bytes32S
		v[index] = sign[64:][0] + 27
	}

	tx, err := adaptor.gravityContract.UpdateConsuls(bind.NewKeyedTransactor(adaptor.privKey), consulsAddress, v[:], r[:], s[:], big.NewInt(round))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
func (adaptor *ErgoAdaptor) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64) ([]byte, error) {
	var oraclesAddresses []common.Address
	for _, v := range consulsAddresses {
		if v == nil {
			oraclesAddresses = append(oraclesAddresses, common.Address{})
			continue
		}
		pubKey, err := crypto.DecompressPubkey(v.ToBytes(account.Ethereum))
		if err != nil {
			return nil, err
		}
		oraclesAddresses = append(oraclesAddresses, crypto.PubkeyToAddress(*pubKey))
	}
	hash, err := adaptor.gravityContract.HashNewConsuls(nil, oraclesAddresses, big.NewInt(roundId))
	if err != nil {
		return nil, err
	}

	sign, err := adaptor.Sign(hash[:])
	if err != nil {
		return nil, err
	}

	return sign, nil
}
func (adaptor *ErgoAdaptor) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey) ([]byte, error) {
	nebula, err := ethereum.NewNebula(common.BytesToAddress(nebulaId.ToBytes(account.Ethereum)), adaptor.ethClient)
	if err != nil {
		return nil, err
	}

	var oraclesAddresses []common.Address
	for _, v := range oracles {
		if v == nil {
			oraclesAddresses = append(oraclesAddresses, common.Address{})
			continue
		}
		pubKey, err := crypto.DecompressPubkey(v.ToBytes(account.Ethereum))
		if err != nil {
			return nil, err
		}
		oraclesAddresses = append(oraclesAddresses, crypto.PubkeyToAddress(*pubKey))
	}

	hash, err := nebula.HashNewOracles(nil, oraclesAddresses)
	if err != nil {
		return nil, err
	}

	sign, err := adaptor.Sign(hash[:])
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func (adaptor *ErgoAdaptor) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	nebula, err := ethereum.NewNebula(common.BytesToAddress(nebulaId.ToBytes(account.Ethereum)), adaptor.ethClient)
	if err != nil {
		return 0, err
	}

	lastId, err := nebula.LastPulseId(nil)
	if err != nil {
		return 0, err
	}

	return lastId.Uint64(), nil
}
func (adaptor *ErgoAdaptor) LastRound(ctx context.Context) (uint64, error) {
	lastRound, err := adaptor.gravityContract.LastRound(nil)
	if err != nil {
		return 0, err
	}

	return lastRound.Uint64(), nil
}
func (adaptor *ErgoAdaptor) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	consuls, err := adaptor.gravityContract.GetConsulsByRoundId(nil, big.NewInt(roundId))
	if err != nil {
		return false, err
	}

	if len(consuls) > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

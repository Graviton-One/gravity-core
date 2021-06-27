package adaptors

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/abi/ethereum"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

//const (
//	Int64  SubType = 0
//	String SubType = 1
//	Bytes  SubType = 2
//
//	waitTimeout = 240
//)
//
//type SubType uint8
type BinanceAdaptor struct {
	privKey *ecdsa.PrivateKey

	ghClient  *gravity.Client
	ethClient *ethclient.Client

	gravityContract *ethereum.Gravity
}
type BinanceAdapterOption func(*BinanceAdaptor) error

func WithBinanceGravityContract(address string) BinanceAdapterOption {
	return func(h *BinanceAdaptor) error {
		hexAddress, err := hexutil.Decode(address)
		if err != nil {
			return err
		}
		ethContractAddress := common.Address{}
		ethContractAddress.SetBytes(hexAddress)
		h.gravityContract, err = ethereum.NewGravity(ethContractAddress, h.ethClient)
		if err != nil {
			return err
		}

		return nil
	}
}
func BinanceAdapterWithGhClient(ghClient *gravity.Client) BinanceAdapterOption {
	return func(h *BinanceAdaptor) error {
		h.ghClient = ghClient
		return nil
	}
}

func NewBinanceAdaptor(privKey []byte, nodeUrl string, ctx context.Context, opts ...BinanceAdapterOption) (*BinanceAdaptor, error) {
	ethClient, err := ethclient.DialContext(ctx, nodeUrl)
	if err != nil {
		return nil, err
	}

	ethPrivKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
		},
		D: new(big.Int),
	}
	ethPrivKey.D.SetBytes(privKey)
	ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(privKey)

	adapter := &BinanceAdaptor{
		privKey:   ethPrivKey,
		ethClient: ethClient,
	}
	for _, opt := range opts {
		err := opt(adapter)
		if err != nil {
			return nil, err
		}
	}

	return adapter, nil
}

func (adaptor *BinanceAdaptor) GetHeight(ctx context.Context) (uint64, error) {
	tcHeightRq, err := adaptor.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return tcHeightRq.NumberU64(), nil
}
func (adaptor *BinanceAdaptor) Sign(msg []byte) ([]byte, error) {
	sig, err := crypto.Sign(msg, adaptor.privKey)
	if err != nil {
		return nil, err
	}

	return sig, nil
}
func (adaptor *BinanceAdaptor) WaitTx(id string, ctx context.Context) error {
	tx, _, err := adaptor.ethClient.TransactionByHash(ctx, common.HexToHash(id))
	if err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, adaptor.ethClient, tx)
	if err != nil {
		return err
	}

	return nil
}
func (adaptor *BinanceAdaptor) PubKey() account.OraclesPubKey {
	pubKey := crypto.CompressPubkey(&adaptor.privKey.PublicKey)
	oraclePubKey := account.BytesToOraclePubKey(pubKey[:], account.Ethereum)
	return oraclePubKey
}
func (adaptor *BinanceAdaptor) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
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

func (adaptor *BinanceAdaptor) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
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
		fmt.Printf("Exist bft count %d < min bft count (%d)", realSignCount, bft.Uint64())
		return "", nil
	}

	var resultBytes32 [32]byte
	copy(resultBytes32[:], hash)

	opt := bind.NewKeyedTransactor(adaptor.privKey)
	opt.Context = ctx
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
func (adaptor *BinanceAdaptor) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
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
		t, err := nebula.DataType(nil)
		if err != nil {
			return err
		}

		transactOpt := bind.NewKeyedTransactor(adaptor.privKey)
		transactOpt.Context = ctx
		switch SubType(t) {
		case Int64:
			v, err := strconv.ParseInt(value.Value, 10, 64)
			if err != nil {
				return err
			}
			_, err = nebula.SendValueToSubInt(transactOpt, v, big.NewInt(int64(pulseId)), id)
			if err != nil {
				return err
			}
		case String:
			_, err = nebula.SendValueToSubString(transactOpt, value.Value, big.NewInt(int64(pulseId)), id)
			if err != nil {
				return err
			}
		case Bytes:
			println(value.Value)
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

func (adaptor *BinanceAdaptor) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
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
	opts := bind.NewKeyedTransactor(adaptor.privKey)
	opts.Context = ctx
	tx, err := nebula.UpdateOracles(opts, oraclesAddresses, v[:], r[:], s[:], big.NewInt(round))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
func (adaptor *BinanceAdaptor) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
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
	opts := bind.NewKeyedTransactor(adaptor.privKey)
	opts.Context = ctx
	tx, err := adaptor.gravityContract.UpdateConsuls(opts, consulsAddress, v[:], r[:], s[:], big.NewInt(round))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}
func (adaptor *BinanceAdaptor) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64, sender account.OraclesPubKey) ([]byte, error) {
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
func (adaptor *BinanceAdaptor) SignHash(nebulaId account.NebulaId, intervalId uint64, pulseId uint64, hash []byte) ([]byte, error) {
	return adaptor.Sign(hash)
}
func (adaptor *BinanceAdaptor) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, round int64, sender account.OraclesPubKey) ([]byte, error) {
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

func (adaptor *BinanceAdaptor) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
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
func (adaptor *BinanceAdaptor) LastRound(ctx context.Context) (uint64, error) {
	lastRound, err := adaptor.gravityContract.LastRound(nil)
	if err != nil {
		return 0, err
	}

	return lastRound.Uint64(), nil
}
func (adaptor *BinanceAdaptor) RoundExist(roundId int64, ctx context.Context) (bool, error) {
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

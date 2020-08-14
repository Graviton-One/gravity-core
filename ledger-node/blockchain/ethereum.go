package blockchain

import (
	"context"
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/Gravity-Tech/gravity-core/common/contracts"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	privKey         *ecdsa.PrivateKey
	privKeyBytes    []byte
	client          *ethclient.Client
	gravityEthereum *contracts.Gravity
}

func NewEthereum(gravityContractAddress string, privKey []byte, nodeUrl string, ctx context.Context) (*Ethereum, error) {
	ethContractAddress := common.Address{}
	hexAddress, err := hexutil.Decode(gravityContractAddress)
	if err != nil {
		return nil, err
	}
	ethContractAddress.SetBytes(hexAddress)

	ethClient, err := ethclient.DialContext(ctx, nodeUrl)
	if err != nil {
		return nil, err
	}
	gravityEthereum, err := contracts.NewGravity(ethContractAddress, ethClient)
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

	return &Ethereum{
		privKey:         ethPrivKey,
		gravityEthereum: gravityEthereum,
		client:          ethClient,
		privKeyBytes:    privKey,
	}, nil
}

func (ethereum *Ethereum) SendOraclesToNebula(nebulaId account.NebulaId, oracles []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error) {
	nebula, err := contracts.NewNebula(common.BytesToAddress(nebulaId), ethereum.client)
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
		oraclesAddresses = append(oraclesAddresses, common.BytesToAddress(v.ToBytes(account.Ethereum)))
	}

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, sign := range signs {
		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])

		r = append(r, bytes32R)
		s = append(s, bytes32S)
		v = append(v, sign[64:][0]+27)
	}

	tx, err := nebula.UpdateOracles(bind.NewKeyedTransactor(ethereum.privKey), oraclesAddresses, v, r, s, big.NewInt(round))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

func (ethereum *Ethereum) SendConsulsToGravityContract(newConsulsAddresses []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error) {
	lastRound, err := ethereum.gravityEthereum.Rounds(nil, big.NewInt(round))
	if err != nil {
		return "", err
	}

	if lastRound {
		return "", nil
	}

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, sign := range signs {
		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])

		r = append(r, bytes32R)
		s = append(s, bytes32S)
		v = append(v, sign[64:][0]+27)
	}

	var consulsAddress []common.Address

	for _, v := range newConsulsAddresses {
		consulsAddress = append(consulsAddress, common.BytesToAddress(v.ToBytes(account.Ethereum)))
	}

	tx, err := ethereum.gravityEthereum.UpdateConsuls(bind.NewKeyedTransactor(ethereum.privKey), consulsAddress, v, r, s, big.NewInt(round))
	if err != nil {
		return "", nil
	}

	return tx.Hash().Hex(), nil
}

func (ethereum *Ethereum) SignConsuls(consulsAddresses []account.OraclesPubKey) ([]byte, error) {
	var oraclesAddresses []common.Address
	for _, v := range consulsAddresses {
		oraclesAddresses = append(oraclesAddresses, common.BytesToAddress(v.ToBytes(account.Ethereum)))
	}
	hash, err := ethereum.gravityEthereum.HashNewConsuls(nil, oraclesAddresses)
	if err != nil {
		return nil, err
	}

	sign, err := account.SignWithTC(ethereum.privKeyBytes, hash[:], account.Ethereum)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func (ethereum *Ethereum) SignOracles(nebulaId account.NebulaId, oracles []account.OraclesPubKey) ([]byte, error) {
	nebula, err := contracts.NewNebula(common.BytesToAddress(nebulaId), ethereum.client)
	if err != nil {
		return nil, err
	}

	var oraclesAddresses []common.Address
	for _, v := range oracles {
		oraclesAddresses = append(oraclesAddresses, common.BytesToAddress(v.ToBytes(account.Ethereum)))
	}

	hash, err := nebula.HashNewOracles(nil, oraclesAddresses)
	if err != nil {
		return nil, err
	}
	sign, err := account.SignWithTC(ethereum.privKeyBytes, hash[:], account.Ethereum)
	if err != nil {
		return nil, err
	}

	return sign, nil
}

func (ethereum *Ethereum) PubKey() []byte {
	return crypto.PubkeyToAddress(ethereum.privKey.PublicKey).Bytes()
}

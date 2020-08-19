package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/Gravity-Tech/gravity-core/common/contracts/sender"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	Int64  SubType = 0
	String SubType = 1
	Bytes  SubType = 2
)

type SubType uint8
type EthereumClient struct {
	ghClient  *client.Client
	ethClient *ethclient.Client
	nebula    *contracts.Nebula
	privKey   []byte
	nebulaId  account.NebulaId
}

func NewEthereumClient(ghClient *client.Client, nebulaId account.NebulaId, nodeUrl string, privKey []byte, ctx context.Context) (*EthereumClient, error) {
	ethContractAddress := common.Address{}
	ethContractAddress.SetBytes(nebulaId)

	ethClient, err := ethclient.DialContext(ctx, nodeUrl)
	if err != nil {
		return nil, err
	}
	nebulaContract, err := contracts.NewNebula(ethContractAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &EthereumClient{
		nebula:    nebulaContract,
		ethClient: ethClient,
		ghClient:  ghClient,
		privKey:   privKey,
		nebulaId:  nebulaId,
	}, nil
}

func (client *EthereumClient) GetHeight(ctx context.Context) (uint64, error) {
	tcHeightRq, err := client.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return tcHeightRq.NumberU64(), nil
}

func (client *EthereumClient) SendResult(tcHeight uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	data, err := client.nebula.Pulses(nil, big.NewInt(int64(tcHeight)))
	if err != nil {
		return "", err
	}

	if bytes.Equal(data[:], make([]byte, 32, 32)) == true {
		bft, err := client.nebula.BftValue(nil)
		if err != nil {
			return "", err
		}

		realSignCount := 0

		oracles, err := client.nebula.GetOracles(nil)
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

			sign, err := client.ghClient.Result(account.Ethereum, client.nebulaId, int64(tcHeight), validator)
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

		if realSignCount >= int(bft.Uint64()) {
			ethPrivKey := &ecdsa.PrivateKey{
				PublicKey: ecdsa.PublicKey{
					Curve: secp256k1.S256(),
				},
				D: new(big.Int),
			}
			ethPrivKey.D.SetBytes(client.privKey)
			ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(client.privKey)

			transactOpt := bind.NewKeyedTransactor(ethPrivKey)
			var resultBytes32 [32]byte
			copy(resultBytes32[:], hash)
			tx, err := client.nebula.SendHashValue(transactOpt, resultBytes32, v[:], r[:], s[:])
			if err != nil {
				return "", err
			}

			fmt.Printf("Tx finilize: %s \n", tx.Hash().String())

			return tx.Hash().String(), nil
		}
	}
	return "", nil
}

func (client *EthereumClient) SendSubs(tcHeight uint64, value interface{}, ctx context.Context) error {
	var err error
	ids, err := client.nebula.GetSubscribersIds(nil)
	if err != nil {
		return err
	}

	for _, id := range ids {
		ethPrivKey := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: secp256k1.S256(),
			},
			D: new(big.Int),
		}
		ethPrivKey.D.SetBytes(client.privKey)
		ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(client.privKey)

		t, err := client.nebula.DataType(nil)
		if err != nil {
			return err
		}

		transactOpt := bind.NewKeyedTransactor(ethPrivKey)
		subSenderAddress, err := client.nebula.SenderToSubs(nil)
		if err != nil {
			return err
		}

		var tx *types.Transaction
		switch SubType(t) {
		case Int64:
			subsSenderContract, err := sender.NewSubsSenderInt(subSenderAddress, client.ethClient)
			if err != nil {
				return err
			}

			tx, err = subsSenderContract.SendValueToSub(transactOpt, value.(int64), big.NewInt(int64(tcHeight)), id)
			if err != nil {
				return err
			}
		case String:
			subsSenderContract, err := sender.NewSubsSenderString(subSenderAddress, client.ethClient)
			if err != nil {
				return err
			}

			tx, err = subsSenderContract.SendValueToSub(transactOpt, value.(string), big.NewInt(int64(tcHeight)), id)
			if err != nil {
				return err
			}
		case Bytes:
			subsSenderContract, err := sender.NewSubsSenderBytes(subSenderAddress, client.ethClient)
			if err != nil {
				return err
			}

			tx, err = subsSenderContract.SendValueToSub(transactOpt, value.([]byte), big.NewInt(int64(tcHeight)), id)
			if err != nil {
				return err
			}
		}

		fmt.Printf("Sub send tx: %s \n", tx.Hash().String())
	}
	return nil
}

func (client *EthereumClient) WaitTx(id string, ctx context.Context) error {
	return nil
}

func (client *EthereumClient) Sign(msg []byte) ([]byte, error) {
	ethPrivKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
		},
		D: new(big.Int),
	}
	ethPrivKey.D.SetBytes(client.privKey)
	ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(client.privKey)

	sig, err := crypto.Sign(msg, ethPrivKey)
	if err != nil {
		return nil, err
	}

	return sig, nil
}

package blockchain

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"fmt"
	"proof-of-concept/common/contracts"
	"proof-of-concept/common/keys"
	"proof-of-concept/gh-node/api/gravity"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Ethereum struct {
	ethClient *ethclient.Client
	nebula    *contracts.Nebula
}

func NewEthereum(contractAddress string, nodeUrl string, ctx context.Context) (*Ethereum, error) {
	ethContractAddress := common.Address{}
	hexAddress, err := hexutil.Decode(contractAddress)
	if err != nil {
		return nil, err
	}
	ethContractAddress.SetBytes(hexAddress)

	ethClient, err := ethclient.DialContext(ctx, nodeUrl)
	if err != nil {
		return nil, err
	}
	nebulaContract, err := contracts.NewNebula(ethContractAddress, ethClient)
	if err != nil {
		return nil, err
	}

	return &Ethereum{
		nebula:    nebulaContract,
		ethClient: ethClient,
	}, nil
}

func (ethereum *Ethereum) GetHeight(ctx context.Context) (uint64, error) {
	tcHeightRq, err := ethereum.ethClient.BlockByNumber(ctx, nil)
	if err != nil {
		return 0, err
	}

	return tcHeightRq.NumberU64(), nil
}

func (ethereum *Ethereum) SendResult(tcHeight uint64, privKey []byte, nebulaId []byte, ghClient *gravity.Client, validators [][]byte, hash []byte, ctx context.Context) (string, error) {
	data, err := ethereum.nebula.Pulses(nil, big.NewInt(int64(tcHeight)))
	if err != nil {
		return "", err
	}

	if bytes.Equal(data[:], make([]byte, 32, 32)) == true {
		bft, err := ethereum.nebula.BftValue(nil)
		if err != nil {
			return "", err
		}

		realSignCount := 0

		oracles, err := ethereum.nebula.GetOracles(nil)
		if err != nil {
			return "", err
		}
		var r [5][32]byte
		var s [5][32]byte
		var v [5]uint8
		for _, validator := range validators {
			pubKey, err := crypto.DecompressPubkey(validator)
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

			sign, err := ghClient.GetKey(keys.FormSignResultKey(nebulaId, tcHeight, validator))
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
			ethPrivKey.D.SetBytes(privKey)
			ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(privKey)

			transactOpt := bind.NewKeyedTransactor(ethPrivKey)
			var resultBytes32 [32]byte
			copy(resultBytes32[:], hash)
			tx, err := ethereum.nebula.ConfirmData(transactOpt, resultBytes32, v[:], r[:], s[:])
			if err != nil {
				return "", err
			}

			fmt.Printf("Tx finilize: %s \n", tx.Hash().String())

			return tx.Hash().String(), nil
		}
	}
	return "", nil
}

func (ethereum *Ethereum) SendSubs(tcHeight uint64, privKey []byte, value uint64, ctx context.Context) error {
	ids, err := ethereum.nebula.GetSubscriptionIds(nil)
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
		ethPrivKey.D.SetBytes(privKey)
		ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(privKey)

		transactOpt := bind.NewKeyedTransactor(ethPrivKey)
		tx, err := ethereum.nebula.SendData(transactOpt, value, big.NewInt(int64(tcHeight)), id)
		if err != nil {
			return err
		}

		fmt.Printf("Sub send tx: %s \n", tx.Hash().String())
	}
	return nil
}

func (ethereum *Ethereum) WaitTx(id string) error {
	return nil
}

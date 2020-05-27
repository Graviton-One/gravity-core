package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gravity-hub/gh-node/api/gravity"
	"gravity-hub/gh-node/extractors"
	"gravity-hub/gh-node/keys"
	"gravity-hub/gh-node/nebula"
	"gravity-hub/gh-node/transaction"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/ethereum/go-ethereum/crypto"
)

const (
	Rounds = 4
)

type Client struct {
	pubKey     []byte
	nebulaId   []byte
	ghClient   *gravity.Client
	ethClient  *ethclient.Client
	extractor  extractors.PriceExtractor
	privKey    *ecdsa.PrivateKey
	nebula     *nebula.Nebula
	round      uint64
	timeout    int
	validators [][]byte
}

func New(privKeyBytes []byte, nebulaId []byte, contractAddress string, ghClient *gravity.Client, ethClient *ethclient.Client, extractor extractors.PriceExtractor, timeout int, ctx context.Context) (*Client, error) {
	privKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
		},
		D: new(big.Int),
	}
	privKey.D.SetBytes(privKeyBytes)
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(privKeyBytes)

	ethContractAddress := common.Address{}
	hexAddress, err := hex.DecodeString(contractAddress)
	if err != nil {
		return nil, err
	}
	ethContractAddress.SetBytes(hexAddress)

	nebulaContract, err := nebula.NewNebula(ethContractAddress, ethClient)
	if err != nil {
		return nil, err
	}
	pubKey := crypto.CompressPubkey(&privKey.PublicKey)
	validatorKey := keys.FormValidatorKey(nebulaId, pubKey)
	_, err = ghClient.GetKey(validatorKey, ctx)

	if err != nil && err != gravity.KeyNotFound {
		return nil, err
	} else if err == gravity.KeyNotFound {
		tx, err := transaction.New(pubKey, transaction.AddValidator, privKey, append(nebulaId, pubKey...))
		if err != nil {
			return nil, err
		}

		err = ghClient.SendTx(tx, ctx)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Add validator tx id: %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}

	validatorPrefix := strings.Join([]string{string(keys.ValidatorKey), hex.EncodeToString(nebulaId)}, "_")
	values, err := ghClient.GetByPrefix(validatorPrefix, ctx)
	if err != nil {
		return nil, err
	}

	var validators [][]byte
	var myRound = 0
	for k, _ := range values {
		keyParts := strings.Split(k, "_")
		validator, err := hex.DecodeString(keyParts[2])
		if err != nil {
			continue
		}

		if bytes.Compare(validator, pubKey) == 0 {
			myRound = len(validators)
		}

		validators = append(validators, validator)
	}

	return &Client{
		pubKey:     pubKey,
		nebulaId:   nebulaId,
		privKey:    privKey,
		nebula:     nebulaContract,
		ghClient:   ghClient,
		ethClient:  ethClient,
		extractor:  extractor,
		round:      uint64(myRound),
		validators: validators,
	}, nil
}

func (client *Client) Start(ctx context.Context) error {
	var lastEthHeight uint64
	var lastGHHeight int64

	commitPrice := make(map[uint64]uint64)
	commitHash := make(map[uint64][]byte)
	resultValue := make(map[uint64]uint64)
	resultHash := make(map[uint64][]byte)
	for {
		ethHeightRq, err := client.ethClient.BlockByNumber(ctx, nil)
		if err != nil {
			return err
		}

		ethHeight := ethHeightRq.NumberU64()
		if lastEthHeight != ethHeight {
			fmt.Printf("Ethereum Height: %d\n", ethHeight)
			lastEthHeight = ethHeight
		}

		block, err := client.ghClient.GetBlock(ctx)
		if err != nil {
			return err
		}

		ghHeight, err := strconv.ParseInt(block.Result.Block.Header.Height, 10, 64)
		if err != nil {
			return err
		}

		if lastGHHeight != ghHeight {
			fmt.Printf("GH Height: %d\n", ghHeight)
			lastGHHeight = ghHeight
		}

		switch ghHeight % Rounds {
		case 0:
			if _, ok := commitHash[ethHeight]; ok {
				continue
			}
			commitKey := keys.FormCommitKey(client.nebulaId, ethHeight, client.pubKey)
			_, err = client.ghClient.GetKey(commitKey, ctx)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}

			price, commit, err := client.commit(ethHeight, ctx)
			if err != nil {
				return err
			}
			commitHash[ethHeight] = commit
			commitPrice[ethHeight] = price
		case 1:
			if _, ok := commitHash[ethHeight]; !ok {
				continue
			}
			revealKey := keys.FormRevealKey(client.nebulaId, ethHeight, commitHash[ethHeight])
			_, err = client.ghClient.GetKey(revealKey, ctx)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}
			err := client.reveal(ethHeight, commitPrice[ethHeight], commitHash[ethHeight], ctx)
			if err != nil {
				return err
			}
		case 2:
			signKey := keys.FormSignResultKey(client.nebulaId, ethHeight, client.pubKey)
			_, err = client.ghClient.GetKey(signKey, ctx)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}
			value, hash, err := client.signResult(ethHeight, ctx)
			if err != nil {
				return err
			}
			resultValue[ethHeight] = value
			resultHash[ethHeight] = hash
		case 3:
			if ethHeight%uint64(len(client.validators)) != client.round {
				continue
			}
			if _, ok := resultValue[ethHeight]; !ok {
				continue
			}
			err = client.sendResult(ethHeight, resultHash[ethHeight], ctx)
			if err != nil {
				return err
			}
		}

		time.Sleep(time.Duration(client.timeout) * time.Second)
	}
}

func (client *Client) commit(ethHeight uint64, ctx context.Context) (uint64, []byte, error) {
	var commitPrice uint64
	price, err := client.extractor.PriceNow()
	if err != nil {
		return 0, nil, err
	}

	commitPrice = uint64(price * 100)

	commitPriceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(commitPriceBytes, commitPrice)
	commit := crypto.Keccak256(commitPriceBytes)

	fmt.Printf("Commit: %.2f - %s \n", float32(commitPrice)/100, hex.EncodeToString(commit[:]))
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, ethHeight)

	tx, err := transaction.New(client.pubKey, transaction.Commit, client.privKey, append(client.nebulaId, append(heightBytes, commit[:]...)...))
	if err != nil {
		return 0, nil, err
	}

	err = client.ghClient.SendTx(tx, ctx)
	if err != nil {
		return 0, nil, err
	}

	fmt.Printf("Commit txId: %s\n", tx.Id)

	return commitPrice, commit, nil
}
func (client *Client) reveal(ethHeight uint64, price uint64, hash []byte, ctx context.Context) error {
	fmt.Printf("Reveal: %.2f - %s \n", float32(price)/100, hex.EncodeToString(hash))
	commitPriceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(commitPriceBytes, price)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, ethHeight)
	var args []byte
	args = append(args, hash...)
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, commitPriceBytes...)

	tx, err := transaction.New(client.pubKey, transaction.Reveal, client.privKey, args)
	if err != nil {
		return err
	}

	err = client.ghClient.SendTx(tx, ctx)
	if err != nil {
		return err
	}
	fmt.Printf("Reveal txId: %s\n", tx.Id)

	return nil
}
func (client *Client) signResult(ethHeight uint64, ctx context.Context) (uint64, []byte, error) {
	prefix := strings.Join([]string{string(keys.RevealKey), hex.EncodeToString(client.nebulaId), fmt.Sprintf("%d", ethHeight)}, "_")
	values, err := client.ghClient.GetByPrefix(prefix, ctx)
	if err != nil {
		panic(err)
	}

	var reveals []uint64
	for _, v := range values {
		reveals = append(reveals, binary.BigEndian.Uint64(v))
	}
	var average uint64
	for _, v := range reveals {
		average += v
	}
	value := uint64(float64(average) / float64(len(reveals)))

	bytesResult := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytesResult, value)

	hash := crypto.Keccak256(bytesResult)
	msg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(hash))
	resultHash := crypto.Keccak256(append([]byte(msg), hash...))
	sign, err := crypto.Sign(resultHash, client.privKey)
	if err != nil {
		return 0, nil, err
	}

	fmt.Printf("Result hash: %s \n", hex.EncodeToString(resultHash))

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, ethHeight)
	var args []byte
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, resultHash...)
	args = append(args, sign...)

	tx, err := transaction.New(client.pubKey, transaction.SignResult, client.privKey, args)
	if err != nil {
		return 0, nil, err
	}

	err = client.ghClient.SendTx(tx, ctx)
	if err != nil {
		return 0, nil, err
	}
	fmt.Printf("Sign result txId: %s\n", tx.Id)
	return value, resultHash, nil
}
func (client *Client) sendResult(ethHeight uint64, hash []byte, ctx context.Context) error {
	data, err := client.nebula.Pulses(nil, big.NewInt(int64(ethHeight)))
	if err != nil {
		return err
	}

	if bytes.Equal(data[:], make([]byte, 32, 32)) == true {
		bft := int(float32(len(client.validators)) * 0.7)
		realSignCount := 0

		oracles, err := client.nebula.GetOracles(nil)
		if err != nil {
			return err
		}
		var r [5][32]byte
		var s [5][32]byte
		var v [5]uint8
		for _, validator := range client.validators {
			pubKey, err := crypto.DecompressPubkey(validator)
			if err != nil {
				return err
			}
			validatorAddress := crypto.PubkeyToAddress(*pubKey)
			position := 0
			for i, address := range oracles {
				if validatorAddress == address {
					position = i
					break
				}
			}

			sign, err := client.ghClient.GetKey(keys.FormSignResultKey(client.nebulaId, ethHeight, validator), ctx)
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

		if realSignCount >= bft {
			transactOpt := bind.NewKeyedTransactor(client.privKey)
			var resultBytes32 [32]byte
			copy(resultBytes32[:], hash)
			tx, err := client.nebula.ConfirmData(transactOpt, resultBytes32, v[:], r[:], s[:])
			if err != nil {
				return err
			}

			fmt.Printf("Tx finilize: %s \n", tx.Hash().String())
		}
	}
	return nil
}

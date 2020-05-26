package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gh-node/api/gravity"
	"gh-node/extractors"
	"gh-node/keys"
	"gh-node/nebula"
	"gh-node/transaction"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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
	round      int
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
	hexAddress, err := hexutil.Decode(contractAddress)
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
		round:      myRound,
		validators: validators,
	}, nil
}

func (client *Client) Start(ctx context.Context) error {
	var lastEthHeight uint64
	var lastGHHeight int64
	var lastCommitPrice string
	var commitHeight int64
	var commitHash []byte
	var resultHash []byte
	for {
		price, err := client.extractor.PriceNow()
		if err != nil {
			return err
		}

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

		if ghHeight%Rounds == 0 {
			commitKey := keys.FormCommitKey(client.nebulaId, ethHeight, client.pubKey)
			_, err := client.ghClient.GetKey(commitKey, ctx)
			if err == gravity.KeyNotFound {
				lastCommitPrice = fmt.Sprintf("%.2f", price)
				commit := crypto.Keccak256([]byte(lastCommitPrice))

				fmt.Printf("Commit: %s - %s \n", lastCommitPrice, hex.EncodeToString(commit[:]))
				heightBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(heightBytes, ethHeight)

				tx, err := transaction.New(client.pubKey, transaction.Commit, client.privKey, append(client.nebulaId, append(heightBytes, commit[:]...)...))
				if err != nil {
					return err
				}

				err = client.ghClient.SendTx(tx, ctx)
				if err != nil {
					return err
				}

				fmt.Printf("Commit txId: %s\n", tx.Id)

				commitHeight = ghHeight
				commitHash = commit[:]
			} else {
				return err
			}
		}

		if commitHeight != 0 && ghHeight%Rounds == 1 {
			revealKey := keys.FormRevealKey(client.nebulaId, ethHeight, commitHash)

			_, err = client.ghClient.GetKey(revealKey, ctx)
			if err == gravity.KeyNotFound {
				fmt.Printf("Reveal: %s - %s \n", lastCommitPrice, hex.EncodeToString(commitHash[:]))
				heightBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(heightBytes, ethHeight)
				var args []byte
				args = append(args, commitHash[:]...)
				args = append(args, client.nebulaId...)
				args = append(args, heightBytes...)
				args = append(args, lastCommitPrice...)

				tx, err := transaction.New(client.pubKey, transaction.Reveal, client.privKey, args)
				if err != nil {
					return err
				}

				err = client.ghClient.SendTx(tx, ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Reveal txId: %s\n", tx.Id)
				commitHeight = 0
			} else {
				if err != nil {
					return err
				}
			}
		}

		if ghHeight%Rounds == 2 {
			signKey := keys.FormSignResultKey(client.nebulaId, ethHeight, client.pubKey)

			_, err = client.ghClient.GetKey(signKey, ctx)
			if err == gravity.KeyNotFound {
				prefix := strings.Join([]string{string(keys.RevealKey), hex.EncodeToString(client.nebulaId), fmt.Sprintf("%d", ethHeight)}, "_")

				values, err := client.ghClient.GetByPrefix(prefix, ctx)
				if err != nil {
					panic(err)
				}

				var reveals []float64
				for _, v := range values {
					value, err := strconv.ParseFloat(string(v), 64)
					if err != nil {
						continue
					}
					reveals = append(reveals, value)
				}
				var average float64
				for _, v := range reveals {
					average += v
				}
				average = average / float64(len(reveals))
				bytesResult := make([]byte, 8)
				binary.LittleEndian.PutUint64(bytesResult, uint64(int64(average*100)))
				resultHashByte32 := crypto.Keccak256(bytesResult)
				resultHash = resultHashByte32[:]
				fmt.Printf("Result hash: %s \n", hex.EncodeToString(resultHash))
				signBytes, err := signEthMsg(resultHash, client.privKey)
				if err != nil {
					return err
				}

				heightBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(heightBytes, ethHeight)
				var args []byte
				args = append(args, client.nebulaId...)
				args = append(args, heightBytes...)
				args = append(args, resultHash...)
				args = append(args, signBytes...)

				tx, err := transaction.New(client.pubKey, transaction.SignResult, client.privKey, args)
				if err != nil {
					return err
				}

				err = client.ghClient.SendTx(tx, ctx)
				if err != nil {
					return err
				}
				fmt.Printf("Sign result txId: %s\n", tx.Id)
			}
		}

		if resultHash != nil && ghHeight%Rounds == 3 && int(ethHeight)%len(client.validators) == client.round {
			data, err := client.nebula.Pulses(nil, big.NewInt(int64(ethHeight)))
			if err != nil {
				return err
			}

			if bytes.Equal(data.DataHash[:], make([]byte, 32, 32)) == true {
				bft := int(float32(len(client.validators)) * 0.7)
				realSignCount := 0

				var r [][32]byte
				var s [][32]byte
				var v []uint8
				for i, validator := range client.validators {
					sign, err := client.ghClient.GetKey(keys.FormSignResultKey(client.nebulaId, ethHeight, validator), ctx)
					if err != nil {
						copy(r[i][:], make([]byte, 32, 32))
						copy(s[i][:], make([]byte, 32, 32))
						v[i] = byte(0)
						continue
					}
					copy(r[i][:], sign[:32])
					copy(s[i][:], sign[32:64])
					v[i] = sign[64] + 27

					realSignCount++
				}

				if realSignCount >= bft {
					transactOpt := bind.NewKeyedTransactor(client.privKey)
					var resultBytes32 [32]byte
					copy(resultBytes32[:], resultHash)
					tx, err := client.nebula.ConfirmData(transactOpt, resultBytes32, v, r, s)
					if err != nil {
						return err
					}

					fmt.Printf("Tx finilize: %s \n", tx.Hash())
				}
			}
		}
		time.Sleep(time.Duration(client.timeout) * time.Second)
	}
}
func signEthMsg(message []byte, privKey *ecdsa.PrivateKey) ([]byte, error) {
	validationMsg := "\x19Ethereum Signed Message:\n" + strconv.Itoa(len(message))
	validationHash := crypto.Keccak256(append([]byte(validationMsg), message[:]...))
	return crypto.Sign(validationHash, privKey)
}

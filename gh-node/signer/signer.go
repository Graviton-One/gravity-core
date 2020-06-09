package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"gravity-hub/common/account"
	"gravity-hub/common/keys"
	"gravity-hub/common/transactions"
	"gravity-hub/gh-node/api/gravity"
	"gravity-hub/gh-node/blockchain"
	"gravity-hub/gh-node/extractors"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/crypto"

	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
)

type Status struct {
	commitPrice uint64
	commitHash  []byte
	resultValue uint64
	resultHash  []byte
	isSent      bool
	isSentSub   bool
}
type Client struct {
	tcPubKey   []byte
	nebulaId   []byte
	tcPrivKey  []byte
	privKey    []byte
	ghClient   *gravity.Client
	extractor  extractors.PriceExtractor
	round      uint64
	timeout    int
	validators [][]byte
	chainType  account.ChainType
	blockchain blockchain.IBlockchain
}

func New(privKeyString string, tcPrivKeyString string, nebulaId []byte, chainType account.ChainType, contractAddress string, nodeUrl string, ghClient *gravity.Client, extractor extractors.PriceExtractor, timeout int, ctx context.Context) (*Client, error) {
	var pubKey []byte
	var tcPrivKey []byte
	var targetBlockchain blockchain.IBlockchain
	var err error
	switch chainType {
	case account.Ethereum:
		privKeyBytes, err := hex.DecodeString(tcPrivKeyString)
		if err != nil {
			return nil, err
		}
		ethPrivKey := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: secp256k1.S256(),
			},
			D: new(big.Int),
		}
		ethPrivKey.D.SetBytes(privKeyBytes)
		ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(privKeyBytes)
		pubKey = crypto.CompressPubkey(&ethPrivKey.PublicKey)

		targetBlockchain, err = blockchain.NewEthereum(contractAddress, nodeUrl, ctx)
		if err != nil {
			return nil, err
		}
		privKey = privKeyBytes
	case account.Waves:
		wCrypto := wavesplatform.NewWavesCrypto()
		seed := wavesplatform.Seed(tcPrivKeyString)
		secret, err := wavesCrypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
		if err != nil {
			panic(err)
		}
		key := wavesCrypto.GeneratePublicKey(secret)
		pubKey = key.Bytes()
		privKey = secret.Bytes()
		targetBlockchain, err = blockchain.NewWaves(contractAddress, nodeUrl, ctx)
		if err != nil {
			return nil, err
		}
	}

	validatorKey := keys.FormValidatorKey(nebulaId, pubKey)
	_, err = ghClient.GetKey(validatorKey, ctx)

	if err != nil && err != gravity.KeyNotFound {
		return nil, err
	} else if err == gravity.KeyNotFound {
		tx, err := transactions.New(pubKey, transactions.AddValidator, chainType, privKey, append(nebulaId, pubKey...))
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
		tcPrivKey:  privKey,
		ghClient:   ghClient,
		extractor:  extractor,
		round:      uint64(myRound),
		validators: validators,
		chainType:  chainType,
		blockchain: targetBlockchain,
		timeout:    timeout,
	}, nil
}

func (client *Client) Start(ctx context.Context) error {
	var lastGHHeight int64

	blockStatus := make(map[uint64]*Status)
	for {
		tcHeight, err := client.blockchain.GetHeight(ctx)
		if err != nil {
			return err
		}

		block, err := client.ghClient.GetBlock(ctx)
		if err != nil {
			return err
		}

		ghHeight, err := strconv.ParseInt(block.Result.Block.Header.Height, 10, 64)
		if err != nil {
			return err
		}

		blockKey := keys.FormBlockKey(client.chainType, tcHeight)
		startGhHeightBytes, err := client.ghClient.GetKey(blockKey, ctx)
		if err != nil && err != gravity.KeyNotFound {
			return err
		}
		var startGhHeight uint64
		if err == gravity.KeyNotFound {
			fmt.Printf("Target Chain Height: %d\n", tcHeight)
			tcHeightBytes := make([]byte, 8)
			ghHeightBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(tcHeightBytes, tcHeight)
			binary.BigEndian.PutUint64(ghHeightBytes, uint64(ghHeight))
			tx, err := transactions.New(client.pubKey, transactions.NewRound, client.chainType, client.privKey, append(tcHeightBytes, ghHeightBytes...))
			if err != nil {
				return err
			}

			err = client.ghClient.SendTx(tx, ctx)
			if err != nil {
				return err
			}
			startGhHeight = uint64(ghHeight)
			fmt.Printf("GH Height Round Start: %d\n", startGhHeight)
		} else {
			startGhHeight = binary.BigEndian.Uint64(startGhHeightBytes)
		}
		if lastGHHeight != ghHeight {
			fmt.Printf("GH Height: %d\n", ghHeight)
			lastGHHeight = ghHeight
		}
		switch uint64(ghHeight) {
		case startGhHeight:
			fallthrough
		case startGhHeight + 1:
			if _, ok := blockStatus[tcHeight]; ok {
				continue
			}
			commitKey := keys.FormCommitKey(client.nebulaId, tcHeight, client.pubKey)
			_, err = client.ghClient.GetKey(commitKey, ctx)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}

			price, commit, err := client.commit(tcHeight, ctx)
			if err != nil {
				return err
			}
			blockStatus[tcHeight] = &Status{
				commitHash:  commit,
				commitPrice: price,
			}
		case startGhHeight + 2:
			if _, ok := blockStatus[tcHeight]; !ok {
				continue
			}
			revealKey := keys.FormRevealKey(client.nebulaId, tcHeight, blockStatus[tcHeight].commitHash)
			_, err = client.ghClient.GetKey(revealKey, ctx)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}
			err := client.reveal(tcHeight, blockStatus[tcHeight].commitPrice, blockStatus[tcHeight].commitHash, ctx)
			if err != nil {
				return err
			}
		case startGhHeight + 3:
			if _, ok := blockStatus[tcHeight]; !ok {
				continue
			}
			signKey := keys.FormSignResultKey(client.nebulaId, tcHeight, client.pubKey)
			_, err = client.ghClient.GetKey(signKey, ctx)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}
			isSuccess, value, hash, err := client.signResult(tcHeight, ctx)
			if err != nil {
				return err
			}
			if !isSuccess {
				continue
			}
			blockStatus[tcHeight].resultValue = value
			blockStatus[tcHeight].resultHash = hash
		default:
			if _, ok := blockStatus[tcHeight]; !ok {
				continue
			}
			if tcHeight%uint64(len(client.validators)) != client.round || blockStatus[tcHeight].isSent {
				continue
			}
			if _, ok := blockStatus[tcHeight]; !ok && blockStatus[tcHeight].resultValue != 0 {
				continue
			}
			txId, err := client.blockchain.SendResult(tcHeight, client.tcPrivKey, client.nebulaId, client.ghClient, client.validators, blockStatus[tcHeight].resultHash, ctx)
			if err != nil {
				return err
			}
			blockStatus[tcHeight].isSent = true
			go func() {
				err := client.blockchain.WaitTx(txId)
				if err != nil {
					println(err.Error())
					return
				}
				for i := 0; i < 1; i++ {
					err = client.blockchain.SendSubs(tcHeight, client.tcPrivKey, blockStatus[tcHeight].resultValue, ctx)
					if err != nil {
						println(err.Error())
						time.Sleep(time.Second)
						continue
					}
					break
				}
			}()
		}

		time.Sleep(time.Duration(client.timeout) * time.Second)
	}
}

func (client *Client) commit(tcHeight uint64, ctx context.Context) (uint64, []byte, error) {
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
	binary.BigEndian.PutUint64(heightBytes, tcHeight)

	tx, err := transactions.New(client.pubKey, transactions.Commit, client.chainType, client.privKey, append(client.nebulaId, append(heightBytes, commit[:]...)...))
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
func (client *Client) reveal(tcHeight uint64, price uint64, hash []byte, ctx context.Context) error {
	fmt.Printf("Reveal: %.2f - %s \n", float32(price)/100, hex.EncodeToString(hash))
	commitPriceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(commitPriceBytes, price)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, tcHeight)
	var args []byte
	args = append(args, hash...)
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, commitPriceBytes...)

	tx, err := transactions.New(client.pubKey, transactions.Reveal, client.chainType, client.privKey, args)
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
func (client *Client) signResult(tcHeight uint64, ctx context.Context) (bool, uint64, []byte, error) {
	prefix := strings.Join([]string{string(keys.RevealKey), hex.EncodeToString(client.nebulaId), fmt.Sprintf("%d", tcHeight)}, "_")
	values, err := client.ghClient.GetByPrefix(prefix, ctx)
	if err != nil {
		panic(err)
	}

	var reveals []uint64
	for _, v := range values {
		reveals = append(reveals, binary.BigEndian.Uint64(v))
	}
	if reveals == nil {
		return false, 0, nil, nil
	}

	var average uint64
	for _, v := range reveals {
		average += v
	}
	value := uint64(float64(average) / float64(len(reveals)))

	bytesResult := make([]byte, 8)
	binary.BigEndian.PutUint64(bytesResult, value)

	hash := crypto.Keccak256(bytesResult)
	sign := account.Sign(client.tcPrivKey, hash)

	fmt.Printf("Result hash: %s \n", hex.EncodeToString(hash))

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, tcHeight)
	var args []byte
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, hash...)
	args = append(args, sign...)

	tx, err := transactions.New(client.pubKey, transactions.SignResult, client.chainType, client.privKey, args)
	if err != nil {
		return false, 0, nil, err
	}

	err = client.ghClient.SendTx(tx, ctx)
	if err != nil {
		return false, 0, nil, err
	}
	fmt.Printf("Sign result txId: %s\n", tx.Id)
	return true, value, hash, nil
}

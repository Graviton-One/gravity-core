package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gravity-hub/common/account"
	"gravity-hub/common/keys"
	"gravity-hub/common/transactions"
	"gravity-hub/gh-node/api/gravity"
	"gravity-hub/gh-node/blockchain"
	"gravity-hub/gh-node/config"
	"gravity-hub/gh-node/extractors"
	"gravity-hub/gh-node/rpc"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/crypto"
	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
)

type TCAccount struct {
	privKey []byte
	pubKey  []byte
}

type Status struct {
	commitPrice uint64
	commitHash  []byte
	resultValue uint64
	resultHash  []byte
	isSent      bool
	isSentSub   bool
}

type Client struct {
	nebulaId []byte
	TCAccount
	ghPrivKey  tendermintCrypto.PrivKeyEd25519
	ghClient   *gravity.Client
	extractor  extractors.PriceExtractor
	timeout    int
	chainType  account.ChainType
	blockchain blockchain.IBlockchain
}

func New(cfg config.Config, ctx context.Context) (*Client, error) {
	nebulaId, err := hexutil.Decode(cfg.NebulaId)
	if err != nil {
		return nil, err
	}

	chainType, err := account.ParseChainType(cfg.ChainType)
	if err != nil {
		return nil, err
	}

	ghClient, err := gravity.NewClient(cfg.GHNodeURL)
	if err != nil {
		panic(err)
	}

	var tcPubKey []byte
	var tcPrivKey []byte
	var targetBlockchain blockchain.IBlockchain

	switch chainType {
	case account.Ethereum:
		privKeyBytes, err := hexutil.Decode(cfg.TCPrivKey)
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
		tcPubKey = crypto.CompressPubkey(&ethPrivKey.PublicKey)

		targetBlockchain, err = blockchain.NewEthereum(cfg.NebulaContract, cfg.NodeUrl, ctx)
		if err != nil {
			return nil, err
		}
		tcPrivKey = privKeyBytes
	case account.Waves:
		wCrypto := wavesplatform.NewWavesCrypto()
		seed := wavesplatform.Seed(cfg.TCPrivKey)
		secret, err := wavesCrypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
		if err != nil {
			panic(err)
		}
		key := wavesCrypto.GeneratePublicKey(secret)
		tcPubKey = key.Bytes()
		tcPrivKey = secret.Bytes()
		targetBlockchain, err = blockchain.NewWaves(cfg.NebulaContract, cfg.NodeUrl, ctx)
		if err != nil {
			return nil, err
		}
	}

	ghPrivKeyBytes, err := hexutil.Decode(cfg.GHPrivKey)
	if err != nil {
		return nil, err
	}
	ghPrivKey := tendermintCrypto.GenPrivKeyFromSecret(ghPrivKeyBytes)
	ghPubKey := ghPrivKey.PubKey().Bytes()

	oracleKey := keys.FormOraclesByValidatorKey(tcPubKey)
	b, err := ghClient.GetKey(oracleKey)

	var oracles map[account.ChainType][]byte
	err = json.Unmarshal(b, &oracles)

	isFound := false
	for _, value := range oracles {
		if bytes.Equal(value, tcPubKey) {
			isFound = true
			break
		}
	}

	if err != nil && err != gravity.KeyNotFound {
		return nil, err
	} else if err == gravity.KeyNotFound || !isFound {
		tx, err := transactions.New(ghPubKey, transactions.AddOracle, chainType, ghPrivKey, append([]byte{byte(chainType)}, tcPubKey...))
		if err != nil {
			return nil, err
		}

		err = ghClient.SendTx(tx)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Add oracle tx id: %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}

	oracleNebulaKey := keys.FormOraclesByNebulaKey(nebulaId)
	item, err := ghClient.GetKey(oracleNebulaKey)
	if err != nil && err != gravity.KeyNotFound {
		return nil, err
	}

	oraclesByNebula := make(map[string]string)
	err = json.Unmarshal(item, &oraclesByNebula)

	if _, ok := oraclesByNebula[hexutil.Encode(tcPubKey)]; !ok {
		tx, err := transactions.New(ghPubKey, transactions.AddOracleInNebula, chainType, ghPrivKey, append(nebulaId, tcPubKey...))
		if err != nil {
			return nil, err
		}

		err = ghClient.SendTx(tx)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Add oracle in nebula id: %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}

	err = rpc.ListenRpcServer(rpc.ServerConfig{
		Host:      cfg.RPCHost,
		PubKey:    tcPubKey,
		PrivKey:   ghPrivKey,
		ChainType: chainType,
		GhClient:  ghClient,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		TCAccount: TCAccount{
			pubKey:  tcPubKey,
			privKey: tcPrivKey,
		},
		nebulaId:   nebulaId,
		ghPrivKey:  ghPrivKey,
		ghClient:   ghClient,
		extractor:  &extractors.BinanceExtractor{},
		chainType:  chainType,
		blockchain: targetBlockchain,
		timeout:    cfg.Timeout,
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

		info, err := client.ghClient.HttpClient.ABCIInfo()
		if err != nil {
			return err
		}

		ghHeight := info.Response.LastBlockHeight

		blockKey := keys.FormBlockKey(client.chainType, tcHeight)
		startGhHeightBytes, err := client.ghClient.GetKey(blockKey)
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
			tx, err := transactions.New(client.pubKey, transactions.NewRound, client.chainType, client.ghPrivKey, append(tcHeightBytes, ghHeightBytes...))
			if err != nil {
				return err
			}

			err = client.ghClient.SendTx(tx)
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
			_, err = client.ghClient.GetKey(commitKey)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}

			price, commit, err := client.commit(tcHeight)
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
			_, err = client.ghClient.GetKey(revealKey)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}
			err := client.reveal(tcHeight, blockStatus[tcHeight].commitPrice, blockStatus[tcHeight].commitHash)
			if err != nil {
				return err
			}
		case startGhHeight + 3:
			if _, ok := blockStatus[tcHeight]; !ok {
				continue
			}
			signKey := keys.FormSignResultKey(client.nebulaId, tcHeight, client.pubKey)
			_, err = client.ghClient.GetKey(signKey)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err != gravity.KeyNotFound {
				continue
			}
			isSuccess, value, hash, err := client.signResult(tcHeight)
			if err != nil {
				return err
			}
			if !isSuccess {
				continue
			}
			blockStatus[tcHeight].resultValue = value
			blockStatus[tcHeight].resultHash = hash
		default:
			var oracles [][]byte
			var myRound uint64

			item, err := client.ghClient.GetKey(keys.FormBftOraclesByNebulaKey(client.nebulaId))
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			oraclesMap := make(map[string]string)
			err = json.Unmarshal(item, &oraclesMap)
			if err != nil {
				return err
			}

			var count uint64
			for k, _ := range oraclesMap {
				v, err := hexutil.Decode(k)
				if err != nil {
					continue
				}
				oracles = append(oracles, v)
				if bytes.Equal(v, client.TCAccount.pubKey) {
					myRound = count
				}
				count++
			}

			if _, ok := blockStatus[tcHeight]; !ok {
				continue
			}
			if tcHeight%uint64(len(oracles)) != myRound || blockStatus[tcHeight].isSent {
				continue
			}
			if _, ok := blockStatus[tcHeight]; !ok && blockStatus[tcHeight].resultValue != 0 {
				continue
			}
			txId, err := client.blockchain.SendResult(tcHeight, client.TCAccount.privKey, client.nebulaId, client.ghClient, oracles, blockStatus[tcHeight].resultHash, ctx)
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
					err = client.blockchain.SendSubs(tcHeight, client.TCAccount.privKey, blockStatus[tcHeight].resultValue, ctx)
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

func (client *Client) commit(tcHeight uint64) (uint64, []byte, error) {
	var commitPrice uint64
	price, err := client.extractor.PriceNow()
	if err != nil {
		return 0, nil, err
	}

	commitPrice = uint64(price * 100)

	commitPriceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(commitPriceBytes, commitPrice)
	commit := crypto.Keccak256(commitPriceBytes)

	fmt.Printf("Commit: %.2f - %s \n", float32(commitPrice)/100, hexutil.Encode(commit[:]))
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, tcHeight)

	tx, err := transactions.New(client.pubKey, transactions.Commit, client.chainType, client.ghPrivKey, append(client.nebulaId, append(heightBytes, commit[:]...)...))
	if err != nil {
		return 0, nil, err
	}

	err = client.ghClient.SendTx(tx)
	if err != nil {
		return 0, nil, err
	}

	fmt.Printf("Commit txId: %s\n", tx.Id)

	return commitPrice, commit, nil
}
func (client *Client) reveal(tcHeight uint64, price uint64, hash []byte) error {
	fmt.Printf("Reveal: %.2f - %s \n", float32(price)/100, hexutil.Encode(hash))
	commitPriceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(commitPriceBytes, price)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, tcHeight)
	var args []byte
	args = append(args, hash...)
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, commitPriceBytes...)

	tx, err := transactions.New(client.pubKey, transactions.Reveal, client.chainType, client.ghPrivKey, args)
	if err != nil {
		return err
	}

	err = client.ghClient.SendTx(tx)
	if err != nil {
		return err
	}
	fmt.Printf("Reveal txId: %s\n", tx.Id)

	return nil
}
func (client *Client) signResult(tcHeight uint64) (bool, uint64, []byte, error) {
	prefix := strings.Join([]string{string(keys.RevealKey), hexutil.Encode(client.nebulaId), fmt.Sprintf("%d", tcHeight)}, "_")
	values, err := client.ghClient.GetByPrefix(prefix)
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
	sign, err := account.SignWithTCPriv(client.TCAccount.privKey, hash, client.chainType)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Result hash: %s \n", hexutil.Encode(hash))

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, tcHeight)
	var args []byte
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, hash...)
	args = append(args, sign...)

	tx, err := transactions.New(client.pubKey, transactions.SignResult, client.chainType, client.ghPrivKey, args)
	if err != nil {
		return false, 0, nil, err
	}

	err = client.ghClient.SendTx(tx)
	if err != nil {
		return false, 0, nil, err
	}
	fmt.Printf("Sign result txId: %s\n", tx.Id)
	return true, value, hash, nil
}

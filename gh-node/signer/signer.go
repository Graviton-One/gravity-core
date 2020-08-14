package signer

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/Gravity-Tech/proof-of-concept/common/account"
	"github.com/Gravity-Tech/proof-of-concept/common/keys"
	"github.com/Gravity-Tech/proof-of-concept/common/transactions"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/api/gravity"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/blockchain"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/config"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/extractor"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/rpc"

	"github.com/btcsuite/btcutil/base58"

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
	ghPubKey   []byte
	ghClient   *gravity.Client
	extractor  *extractor.Client
	timeout    int
	chainType  account.ChainType
	blockchain blockchain.IBlockchain
}

func New(cfg config.Config, ctx context.Context) (*Client, error) {
	chainType, err := account.ParseChainType(cfg.ChainType)
	if err != nil {
		return nil, err
	}

	var nebulaId []byte
	switch chainType {
	case account.Waves:
		nebulaId = base58.Decode(cfg.NebulaId)
	case account.Ethereum:
		nebulaId, err = hexutil.Decode(cfg.NebulaId)
		if err != nil {
			return nil, err
		}
	}

	ghClient, err := gravity.NewClient(cfg.GHNodeURL)
	if err != nil {
		return nil, err
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

		targetBlockchain, err = blockchain.NewEthereum(cfg.NebulaId, cfg.NodeUrl, ctx)
		if err != nil {
			return nil, err
		}
		tcPrivKey = privKeyBytes
	case account.Waves:
		wCrypto := wavesplatform.NewWavesCrypto()
		seed := wavesplatform.Seed(cfg.TCPrivKey)
		secret, err := wavesCrypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
		if err != nil {
			return nil, err
		}
		key := wavesCrypto.GeneratePublicKey(secret)
		tcPubKey = key.Bytes()
		tcPrivKey = secret.Bytes()
		targetBlockchain, err = blockchain.NewWaves(cfg.NebulaId, cfg.NodeUrl, ctx)
		if err != nil {
			return nil, err
		}
	}

	ghPrivKeyBytes, err := base64.StdEncoding.DecodeString(cfg.GHPrivKey)
	if err != nil {
		return nil, err
	}
	ghPrivKey := tendermintCrypto.PrivKeyEd25519{}
	copy(ghPrivKey[:], ghPrivKeyBytes)
	ghPubKey := ghPrivKey.PubKey().Bytes()[5:]

	oracleKey := keys.FormOraclesByValidatorKey(ghPubKey)

	rs, err := ghClient.HttpClient.ABCIQuery("key", []byte(oracleKey))
	if err != nil {
		return nil, err
	}

	var oracles map[account.ChainType][]byte
	if rs.Response.Value != nil {
		err = json.Unmarshal(rs.Response.Value, &oracles)
		if err != nil {
			return nil, err
		}
	}

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
	rs, err = ghClient.HttpClient.ABCIQuery("key", []byte(oracleNebulaKey))
	if err != nil {
		return nil, err
	}

	oraclesByNebula := make(map[string]string)
	if rs.Response.Value != nil {
		err = json.Unmarshal(rs.Response.Value, &oraclesByNebula)
		if err != nil {
			return nil, err
		}
	}

	if _, ok := oraclesByNebula[hexutil.Encode(tcPubKey)]; !ok {
		args := []byte{byte(len(nebulaId))}
		tx, err := transactions.New(ghPubKey, transactions.AddOracleInNebula, chainType, ghPrivKey, append(args, append(nebulaId, tcPubKey...)...))
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

	go rpc.ListenRpcServer(rpc.ServerConfig{
		Host:      cfg.RPCHost,
		PubKey:    tcPubKey,
		PrivKey:   ghPrivKey,
		ChainType: chainType,
		GhClient:  ghClient,
	})

	return &Client{
		TCAccount: TCAccount{
			pubKey:  tcPubKey,
			privKey: tcPrivKey,
		},
		nebulaId:   nebulaId,
		ghPrivKey:  ghPrivKey,
		ghClient:   ghClient,
		extractor:  &extractor.Client{hostURL: &cfg.ExtractorURL},
		chainType:  chainType,
		blockchain: targetBlockchain,
		timeout:    cfg.Timeout,
		ghPubKey:   ghPubKey,
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

		info, err := client.ghClient.HttpClient.Status()
		if err != nil {
			return err
		}

		ghHeight := info.SyncInfo.LatestBlockHeight

		blockKey := keys.FormBlockKey(client.chainType, tcHeight)
		rs, err := client.ghClient.HttpClient.ABCIQuery("key", []byte(blockKey))
		if err != nil {
			return err
		}

		var startGhHeight uint64
		if rs.Response.Value == nil {
			fmt.Printf("Target Chain Height: %d\n", tcHeight)
			tcHeightBytes := make([]byte, 8)
			ghHeightBytes := make([]byte, 8)
			binary.BigEndian.PutUint64(tcHeightBytes, tcHeight)
			binary.BigEndian.PutUint64(ghHeightBytes, uint64(ghHeight))
			tx, err := transactions.New(client.ghPubKey, transactions.NewRound, client.chainType, client.ghPrivKey, append(tcHeightBytes, ghHeightBytes...))
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
			startGhHeight = binary.BigEndian.Uint64(rs.Response.Value)
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
			commitKey := keys.FormCommitKey(client.nebulaId, tcHeight, client.TCAccount.pubKey)
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
			signKey := keys.FormSignResultKey(client.nebulaId, tcHeight, client.TCAccount.pubKey)
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
			if tcHeight%uint64(len(oracles)) != myRound {
				continue
			}

			signKey := keys.FormSignResultKey(client.nebulaId, tcHeight, client.TCAccount.pubKey)
			_, err = client.ghClient.GetKey(signKey)
			if err != nil && err != gravity.KeyNotFound {
				return err
			}
			if err == gravity.KeyNotFound {
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

			if txId == "" {

				go func(txId string) {
					err := client.blockchain.WaitTx(txId)
					if err != nil {
						println(err.Error())
						return
					}
					for i := 0; i < 1; i++ {
						err = client.blockchain.SendSubs(tcHeight, client.TCAccount.privKey, blockStatus[tcHeight].resultValue, ctx)
						if err != nil {
							time.Sleep(time.Second)
							continue
						}
						break
					}
				}(txId)

			}
		}

		time.Sleep(time.Duration(client.timeout) * time.Second)
	}
}

func (client *Client) commit(tcHeight uint64) (uint64, []byte, error) {
	var commitPrice uint64
	commitBytes := client.extractor.RawData()

	commit := crypto.Keccak256(commitBytes)

	fmt.Printf("Commit: %.2f - %s \n", float32(commitPrice)/100, hexutil.Encode(commit[:]))
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, tcHeight)

	args := []byte{byte(len(client.nebulaId))}
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, commit[:]...)
	args = append(args, client.TCAccount.pubKey[:]...)

	tx, err := transactions.New(client.ghPubKey, transactions.Commit, client.chainType, client.ghPrivKey, args)
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
	args = append(args, byte(len(client.nebulaId)))
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, byte(len(commitPriceBytes)))
	args = append(args, commitPriceBytes...)
	args = append(args, client.TCAccount.pubKey...)

	tx, err := transactions.New(client.ghPubKey, transactions.Reveal, client.chainType, client.ghPrivKey, args)
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
	args = append(args, byte(len(client.nebulaId)))
	args = append(args, client.nebulaId...)
	args = append(args, heightBytes...)
	args = append(args, hash...)
	args = append(args, sign...)

	tx, err := transactions.New(client.ghPubKey, transactions.SignResult, client.chainType, client.ghPrivKey, args)
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

package node

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/state"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/Gravity-Tech/gravity-core/oracle-node/blockchain"
	"github.com/Gravity-Tech/gravity-core/oracle-node/config"
	"github.com/Gravity-Tech/gravity-core/oracle-node/extractors"
	"github.com/Gravity-Tech/gravity-core/oracle-node/rpc"

	"github.com/btcsuite/btcutil/base58"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	Int64Type  ExtractorType = "int64"
	StringType ExtractorType = "string"
	BytesType  ExtractorType = "bytes"
)

type ExtractorType string

type Node struct {
	nebulaId account.NebulaId
	TCAccount
	ghPrivKey     tendermintCrypto.PrivKeyEd25519
	ghPubKey      account.ConsulPubKey
	ghClient      *client.Client
	extractor     extractors.Extractor
	timeout       int
	chainType     account.ChainType
	blockchain    blockchain.IBlockchain
	extractorType ExtractorType
}

func New(cfg config.Config, ctx context.Context) (*Node, error) {
	chainType, err := account.ParseChainType(cfg.ChainType)
	if err != nil {
		return nil, err
	}

	var nebulaId account.NebulaId
	switch chainType {
	case account.Waves:
		nebulaId = base58.Decode(cfg.NebulaId)
	case account.Ethereum:
		nebulaId, err = hexutil.Decode(cfg.NebulaId)
		if err != nil {
			return nil, err
		}
	}

	var targetBlockchain blockchain.IBlockchain
	switch chainType {
	case account.Ethereum:
		targetBlockchain, err = blockchain.NewEthereum(cfg.NebulaId, cfg.NodeUrl, ctx)
		if err != nil {
			return nil, err
		}
	case account.Waves:
		targetBlockchain, err = blockchain.NewWaves(cfg.NebulaId, cfg.NodeUrl, ctx)
		if err != nil {
			return nil, err
		}
	}

	tcPrivKey, tcPubKey, err := account.HexToPrivKey(cfg.TCPrivKey, chainType)
	if err != nil {
		return nil, err
	}

	ghPrivKeyBytes, err := base64.StdEncoding.DecodeString(cfg.GHPrivKey)
	if err != nil {
		return nil, err
	}

	ghPrivKey := tendermintCrypto.PrivKeyEd25519{}
	copy(ghPrivKey[:], ghPrivKeyBytes)

	var ghPubKey account.ConsulPubKey
	copy(ghPubKey[:], ghPrivKey.PubKey().Bytes()[5:])

	ghClient, err := client.New(cfg.GHNodeURL)
	if err != nil {
		return nil, err
	}

	go rpc.ListenRpcServer(rpc.ServerConfig{
		Host:      cfg.RPCHost,
		PubKey:    ghPubKey,
		PrivKey:   ghPrivKey,
		ChainType: chainType,
		GhClient:  ghClient,
	})

	return &Node{
		TCAccount: TCAccount{
			pubKey:  tcPubKey,
			privKey: tcPrivKey,
		},
		nebulaId:      nebulaId,
		ghPrivKey:     ghPrivKey,
		ghClient:      ghClient,
		extractor:     &extractors.BinanceExtractor{},
		chainType:     chainType,
		blockchain:    targetBlockchain,
		timeout:       cfg.Timeout,
		ghPubKey:      ghPubKey,
		extractorType: ExtractorType(cfg.ExtractorType),
	}, nil
}

func (node *Node) Init() error {
	oraclesByValidator, err := node.ghClient.OraclesByValidator(node.ghPubKey)
	if err != nil {
		return err
	}

	oracle, ok := oraclesByValidator[node.chainType]
	if !ok || oracle == node.TCAccount.pubKey {
		args := []transactions.Args{
			{
				Value: node.chainType,
			},
			{
				Value: node.TCAccount.pubKey,
			},
		}

		tx, err := transactions.New(node.ghPubKey, transactions.AddOracle, node.ghPrivKey, args)
		if err != nil {
			return err
		}

		err = node.ghClient.SendTx(tx)
		if err != nil {
			return err
		}

		fmt.Printf("Add oracle (TXID): %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}

	oraclesByNebulaKey, err := node.ghClient.OraclesByNebula(node.nebulaId, node.chainType)
	if err != nil {
		return err
	}

	_, ok = oraclesByNebulaKey[node.TCAccount.pubKey]
	if !ok {
		args := []transactions.Args{
			{
				Value: node.nebulaId,
			},
			{
				Value: node.TCAccount.pubKey,
			},
		}

		tx, err := transactions.New(node.ghPubKey, transactions.AddOracleInNebula, node.ghPrivKey, args)
		if err != nil {
			return err
		}

		err = node.ghClient.SendTx(tx)
		if err != nil {
			return err
		}

		fmt.Printf("Add oracle in nebula (TXID): %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}
	return nil
}

func (node *Node) Start(ctx context.Context) error {
	var lastLedgerHeight uint64

	roundState := make(map[uint64]*RoundState)
	for {
		tcHeight, err := node.blockchain.GetHeight(ctx)
		if err != nil {
			return err
		}

		info, err := node.ghClient.HttpClient.Status()
		if err != nil {
			return err
		}

		ledgerHeight := uint64(info.SyncInfo.LatestBlockHeight)

		roundHeight, err := node.ghClient.RoundHeight(node.chainType, ledgerHeight)
		if err != nil && err != client.ErrValueNotFound {
			return err
		}

		var startGhHeight uint64
		if err == client.ErrValueNotFound {
			fmt.Printf("Target Chain Height: %d\n", tcHeight)

			args := []transactions.Args{
				{
					Value: node.chainType,
				},
				{
					Value: tcHeight,
				},
			}

			tx, err := transactions.New(node.ghPubKey, transactions.NewRound, node.ghPrivKey, args)
			if err != nil {
				return err
			}
			err = node.ghClient.SendTx(tx)
			if err != nil {
				return err
			}

			startGhHeight = ledgerHeight
			fmt.Printf("Round Start (Height): %d\n", startGhHeight)
		} else {
			startGhHeight = roundHeight
		}

		if lastLedgerHeight != ledgerHeight {
			fmt.Printf("Ledger Height: %d\n", ledgerHeight)
			lastLedgerHeight = ledgerHeight
		}

		switch state.CalculateSubRound(ledgerHeight) {
		case state.CommitSubRound:
			if _, ok := roundState[tcHeight]; ok {
				continue
			}

			_, err := node.ghClient.CommitHash(node.chainType, node.nebulaId, int64(tcHeight), node.TCAccount.pubKey)
			if err != nil && err != client.ErrValueNotFound {
				return err
			}

			data, err := node.extractor.Extract()
			if err != nil {
				return err
			}
			commit, err := node.commit(data, tcHeight)
			if err != nil {
				return err
			}
			roundState[tcHeight] = &RoundState{
				commitHash: commit,
				data:       data,
			}
		case state.RevealSubRound:
			if _, ok := roundState[tcHeight]; !ok {
				continue
			}
			_, err := node.ghClient.Reveal(node.nebulaId, int64(tcHeight), roundState[tcHeight].commitHash)
			if err != nil && err != client.ErrValueNotFound {
				return err
			}

			err = node.reveal(tcHeight, roundState[tcHeight].data, roundState[tcHeight].commitHash)
			if err != nil {
				return err
			}
		case state.ResultSubRound:
			if _, ok := roundState[tcHeight]; !ok {
				continue
			}

			_, err := node.ghClient.Result(node.chainType, node.nebulaId, int64(tcHeight), node.TCAccount.pubKey)
			if err != nil && err != client.ErrValueNotFound {
				return err
			}

			isSuccess, value, hash, err := node.signResult(tcHeight)
			if err != nil {
				return err
			}
			if !isSuccess {
				continue
			}
			roundState[tcHeight].resultValue = value
			roundState[tcHeight].resultHash = hash
		default:
			var oracles []account.OraclesPubKey
			var myRound uint64

			oraclesMap, err := node.ghClient.BftOraclesByNebula(node.chainType, node.nebulaId)
			if err != nil {
				return err
			}
			var count uint64
			for oracle, _ := range oraclesMap {
				oracles = append(oracles, oracle)
				if node.TCAccount.pubKey == oracle {
					myRound = count
				}
				count++
			}

			if _, ok := roundState[tcHeight]; !ok {
				continue
			}
			if tcHeight%uint64(len(oracles)) != myRound {
				continue
			}

			_, err = node.ghClient.Result(node.chainType, node.nebulaId, int64(tcHeight), node.TCAccount.pubKey)
			if err == client.ErrValueNotFound {
				continue
			} else if err != nil {
				return err
			}

			if _, ok := roundState[tcHeight]; !ok && roundState[tcHeight].resultValue != 0 {
				continue
			}

			txId, err := node.blockchain.SendResult(node.ghClient, tcHeight, node.TCAccount.privKey, node.nebulaId, oracles, roundState[tcHeight].resultHash, ctx)
			if err != nil {
				return err
			}

			roundState[tcHeight].isSent = true

			if txId == "" {
				go func(txId string) {
					err := node.blockchain.WaitTx(txId)
					if err != nil {
						println(err.Error())
						return
					}
					for i := 0; i < 1; i++ {
						err = node.blockchain.SendSubs(tcHeight, node.TCAccount.privKey, roundState[tcHeight].resultValue, ctx)
						if err != nil {
							time.Sleep(time.Second)
							continue
						}
						break
					}
				}(txId)
			}
		}

		time.Sleep(time.Duration(node.timeout) * time.Second)
	}
}

func (node *Node) commit(data interface{}, tcHeight uint64) ([]byte, error) {
	dataBytes := toBytes(data)
	commit := crypto.Keccak256(dataBytes)
	fmt.Printf("Commit: %s - %s \n", hexutil.Encode(dataBytes), hexutil.Encode(commit[:]))

	args := []transactions.Args{
		{
			Value: node.nebulaId,
		},
		{
			Value: tcHeight,
		},
		{
			Value: commit,
		},
		{
			Value: node.TCAccount.pubKey,
		},
	}

	tx, err := transactions.New(node.ghPubKey, transactions.Commit, node.ghPrivKey, args)
	if err != nil {
		return nil, err
	}

	err = node.ghClient.SendTx(tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Commit txId: %s\n", tx.Id)

	return commit, nil
}
func (node *Node) reveal(tcHeight uint64, reveal interface{}, commit []byte) error {
	dataBytes := toBytes(reveal)
	fmt.Printf("Reveal: %s  - %s \n", hexutil.Encode(dataBytes), hexutil.Encode(commit))

	args := []transactions.Args{
		{
			Value: commit,
		},
		{
			Value: node.nebulaId,
		},
		{
			Value: tcHeight,
		},
		{
			Value: reveal,
		},
		{
			Value: node.TCAccount.pubKey,
		},
	}

	tx, err := transactions.New(node.ghPubKey, transactions.Reveal, node.ghPrivKey, args)
	if err != nil {
		return err
	}

	err = node.ghClient.SendTx(tx)
	if err != nil {
		return err
	}
	fmt.Printf("Reveal txId: %s\n", tx.Id)

	return nil
}
func (node *Node) signResult(tcHeight uint64) (bool, interface{}, []byte, error) {
	var values []interface{}
	bytesValues, err := node.ghClient.Results(tcHeight, node.chainType, node.nebulaId)
	if err != nil {
		return false, nil, nil, err
	}

	for _, v := range bytesValues {
		values = append(values, fromBytes(v, node.extractorType))
	}

	result, err := node.extractor.Aggregate(values)
	if err != nil {
		return false, nil, nil, err
	}

	hash := crypto.Keccak256(toBytes(result))
	sign, err := account.SignWithTC(node.TCAccount.privKey, hash, node.chainType)
	if err != nil {
		return false, nil, nil, err
	}
	fmt.Printf("Result hash: %s \n", hexutil.Encode(hash))

	args := []transactions.Args{
		{
			Value: node.nebulaId,
		},
		{
			Value: tcHeight,
		},
		{
			Value: sign,
		},
		{
			Value: byte(node.chainType),
		},
		{
			Value: node.TCAccount.pubKey,
		},
	}
	tx, err := transactions.New(node.ghPubKey, transactions.Result, node.ghPrivKey, args)
	if err != nil {
		return false, nil, nil, err
	}

	err = node.ghClient.SendTx(tx)
	if err != nil {
		return false, nil, nil, err
	}

	fmt.Printf("Sign result txId: %s\n", tx.Id)
	return true, result, hash, nil
}

func toBytes(value interface{}) []byte {
	switch v := value.(type) {
	case int64:
		var b []byte
		binary.BigEndian.PutUint64(b, uint64(v))
		return b
	case string:
		return []byte(v)
	case []byte:
		return v
	}
	return nil
}
func fromBytes(value []byte, extractorType ExtractorType) interface{} {
	switch extractorType {
	case Int64Type:
		return binary.BigEndian.Uint64(value)
	case StringType:
		return string(value)
	case BytesType:
		return value
	}

	return nil
}

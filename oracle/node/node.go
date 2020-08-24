package node

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"os"
	"time"

	"log"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/Gravity-Tech/gravity-core/common/state"

	"github.com/Gravity-Tech/gravity-core/oracle/extractor"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/Gravity-Tech/gravity-core/oracle/blockchain"
	"github.com/Gravity-Tech/gravity-core/oracle/config"
	"github.com/Gravity-Tech/gravity-core/oracle/rpc"

	"github.com/btcsuite/btcutil/base58"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ethereum/go-ethereum/crypto"
	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	errorLogger = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

type Node struct {
	nebulaId        account.NebulaId
	tcPubKey        account.OraclesPubKey
	ghPrivKey       tendermintCrypto.PrivKeyEd25519
	ghPubKey        account.ConsulPubKey
	ghClient        *client.Client
	timeout         int
	chainType       account.ChainType
	blockchain      blockchain.IBlockchainClient
	extractorClient *extractor.Client
	extractorType   contracts.ExtractorType
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

	ghClient, err := client.New(cfg.GHNodeURL)
	if err != nil {
		return nil, err
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

	var targetBlockchain blockchain.IBlockchainClient
	switch chainType {
	case account.Ethereum:
		targetBlockchain, err = blockchain.NewEthereumClient(ghClient, nebulaId, cfg.NodeUrl, tcPrivKey, ctx)
		if err != nil {
			return nil, err
		}
	case account.Waves:
		targetBlockchain, err = blockchain.NewWavesClient(ghClient, nebulaId, tcPrivKey, cfg.NodeUrl)
		if err != nil {
			return nil, err
		}
	}

	go rpc.ListenRpcServer(rpc.ServerConfig{
		Host:      cfg.RPCHost,
		PubKey:    ghPubKey,
		PrivKey:   ghPrivKey,
		ChainType: chainType,
		GhClient:  ghClient,
	})

	extractorClient := extractor.New(cfg.ExtractorUrl)

	return &Node{
		tcPubKey:        tcPubKey,
		nebulaId:        nebulaId,
		ghPrivKey:       ghPrivKey,
		ghPubKey:        ghPubKey,
		ghClient:        ghClient,
		extractorClient: extractorClient,
		chainType:       chainType,
		blockchain:      targetBlockchain,
		timeout:         cfg.Timeout,
		extractorType:   contracts.ExtractorType(cfg.ExtractorType),
	}, nil
}

func (node *Node) Init() error {
	oraclesByValidator, err := node.ghClient.OraclesByValidator(node.ghPubKey)
	if err != nil {
		return err
	}

	oracle, ok := oraclesByValidator[node.chainType]
	if !ok || oracle == node.tcPubKey {
		args := []transactions.Args{
			{
				Value: node.chainType,
			},
			{
				Value: node.tcPubKey,
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

	_, ok = oraclesByNebulaKey[node.tcPubKey]
	if !ok {
		args := []transactions.Args{
			{
				Value: node.nebulaId,
			},
			{
				Value: node.tcPubKey,
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

func (node *Node) Start(ctx context.Context) {
	var lastLedgerHeight uint64
	roundState := make(map[uint64]*RoundState)
	for {
		info, err := node.ghClient.HttpClient.Status()
		if err != nil {
			errorLogger.Print(err)
		} else {
			ledgerHeight := uint64(info.SyncInfo.LatestBlockHeight)
			if lastLedgerHeight != ledgerHeight {
				fmt.Printf("Ledger Height: %d\n", ledgerHeight)
				lastLedgerHeight = ledgerHeight
			}

			err := node.execute(ledgerHeight, roundState, ctx)
			if err != nil {
				errorLogger.Print(err)
			}
		}

		time.Sleep(time.Duration(node.timeout) * time.Second)
	}
}

func (node *Node) execute(ledgerHeight uint64, roundState map[uint64]*RoundState, ctx context.Context) error {
	tcHeight, err := node.blockchain.GetHeight(ctx)
	if err != nil {
		return err
	}

	roundHeight, err := node.ghClient.RoundHeight(node.chainType, ledgerHeight)
	if err != nil && err != client.ErrValueNotFound {
		return err
	}

	var startGhHeight uint64
	if err != client.ErrValueNotFound {
		startGhHeight = roundHeight
	} else {
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
	}

	switch state.CalculateSubRound(ledgerHeight) {
	case state.CommitSubRound:
		if _, ok := roundState[tcHeight]; ok {
			return nil
		}

		_, err := node.ghClient.CommitHash(node.chainType, node.nebulaId, int64(tcHeight), node.tcPubKey)
		if err != nil && err != client.ErrValueNotFound {
			return err
		}

		data, err := node.extractorClient.Extract(ctx)
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
			return nil
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
			return nil
		}

		_, err := node.ghClient.Result(node.chainType, node.nebulaId, int64(tcHeight), node.tcPubKey)
		if err != nil && err != client.ErrValueNotFound {
			return err
		}

		isSuccess, value, hash, err := node.signResult(tcHeight, ctx)
		if err != nil {
			return err
		}
		if !isSuccess {
			return nil
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
			if node.tcPubKey == oracle {
				myRound = count
			}
			count++
		}

		if _, ok := roundState[tcHeight]; !ok {
			return nil
		}
		if tcHeight%uint64(len(oracles)) != myRound {
			return nil
		}

		_, err = node.ghClient.Result(node.chainType, node.nebulaId, int64(tcHeight), node.tcPubKey)
		if err == client.ErrValueNotFound {
			return nil
		} else if err != nil {
			return err
		}

		if _, ok := roundState[tcHeight]; !ok && roundState[tcHeight].resultValue != 0 {
			return nil
		}

		txId, err := node.blockchain.SendResult(tcHeight, oracles, roundState[tcHeight].resultHash, ctx)
		if err != nil {
			return err
		}

		roundState[tcHeight].isSent = true

		//TODO
		if txId == "" {
			go func(txId string) {
				err := node.blockchain.WaitTx(txId, ctx)
				if err != nil {
					println(err.Error())
					return
				}
				for i := 0; i < 1; i++ {
					err = node.blockchain.SendSubs(tcHeight, roundState[tcHeight].resultValue, ctx)
					if err != nil {
						time.Sleep(time.Second)
						continue
					}
					break
				}
			}(txId)
		}
	}
	return nil
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
			Value: node.tcPubKey,
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
			Value: node.tcPubKey,
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
func (node *Node) signResult(tcHeight uint64, ctx context.Context) (bool, interface{}, []byte, error) {
	var values []interface{}
	bytesValues, err := node.ghClient.Results(tcHeight, node.chainType, node.nebulaId)
	if err != nil {
		return false, nil, nil, err
	}

	for _, v := range bytesValues {
		values = append(values, fromBytes(v, node.extractorType))
	}

	result, err := node.extractorClient.Aggregate(values, ctx)
	if err != nil {
		return false, nil, nil, err
	}

	hash := crypto.Keccak256(toBytes(result))
	sign, err := node.blockchain.Sign(hash)
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
			Value: node.tcPubKey,
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
func fromBytes(value []byte, extractorType contracts.ExtractorType) interface{} {
	switch extractorType {
	case contracts.Int64Type:
		return binary.BigEndian.Uint64(value)
	case contracts.StringType:
		return string(value)
	case contracts.BytesType:
		return value
	}

	return nil
}

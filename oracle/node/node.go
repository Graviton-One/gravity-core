package node

import (
	"context"
	"fmt"
	"os"
	"time"

	"log"

	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/Gravity-Tech/gravity-core/common/state"

	"github.com/Gravity-Tech/gravity-core/oracle/extractor"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/transactions"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	Timeout = 1
)

var (
	errorLogger = log.New(os.Stdout,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
)

type Validator struct {
	privKey tendermintCrypto.PrivKeyEd25519
	pubKey  account.ConsulPubKey
}

func NewValidator(privKey []byte) *Validator {
	validatorPrivKey := tendermintCrypto.PrivKeyEd25519{}
	copy(validatorPrivKey[:], privKey)

	var ghPubKey account.ConsulPubKey
	copy(ghPubKey[:], validatorPrivKey.PubKey().Bytes()[5:])

	return &Validator{
		privKey: validatorPrivKey,
		pubKey:  ghPubKey,
	}
}

type Extractor struct {
	*extractor.Client
	ExtractorType contracts.ExtractorType
}

type Node struct {
	nebulaId  account.NebulaId
	chainType account.ChainType

	validator     *Validator
	oraclePubKey  account.OraclesPubKey
	gravityClient *client.GravityClient

	adaptor   adaptors.IBlockchainAdaptor
	extractor *Extractor
}

func New(nebulaId account.NebulaId, chainType account.ChainType, oracleSecretKey []byte, validator *Validator, extractor *Extractor, gravityNodeUrl string, targetChainNodeUrl string, ctx context.Context) (*Node, error) {
	ghClient, err := client.NewGravityClient(gravityNodeUrl)
	if err != nil {
		return nil, err
	}

	var adaptor adaptors.IBlockchainAdaptor
	switch chainType {
	case account.Ethereum:
		adaptor, err = adaptors.NewEthereumAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.EthAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Waves:
		adaptor, err = adaptors.NewWavesAdapter(oracleSecretKey, targetChainNodeUrl, adaptors.WavesAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	}

	return &Node{
		validator:     validator,
		nebulaId:      nebulaId,
		extractor:     extractor,
		chainType:     chainType,
		adaptor:       adaptor,
		gravityClient: ghClient,
		oraclePubKey:  adaptor.PubKey(),
	}, nil
}

func (node *Node) Init() error {
	oraclesByValidator, err := node.gravityClient.OraclesByValidator(node.validator.pubKey)
	if err != nil {
		return err
	}

	oracle, ok := oraclesByValidator[node.chainType]
	if !ok || oracle == node.oraclePubKey {
		args := []transactions.Args{
			{
				Value: node.chainType,
			},
			{
				Value: node.oraclePubKey,
			},
		}

		tx, err := transactions.New(node.validator.pubKey, transactions.AddOracle, node.validator.privKey, args)
		if err != nil {
			return err
		}

		err = node.gravityClient.SendTx(tx)
		if err != nil {
			return err
		}

		fmt.Printf("Add oracle (TXID): %s\n", tx.Id)
		time.Sleep(time.Duration(5) * time.Second)
	}

	oraclesByNebulaKey, err := node.gravityClient.OraclesByNebula(node.nebulaId, node.chainType)
	if err != nil {
		return err
	}

	_, ok = oraclesByNebulaKey[node.oraclePubKey]
	if !ok {
		args := []transactions.Args{
			{
				Value: node.nebulaId,
			},
			{
				Value: node.oraclePubKey,
			},
		}

		tx, err := transactions.New(node.validator.pubKey, transactions.AddOracleInNebula, node.validator.privKey, args)
		if err != nil {
			return err
		}

		err = node.gravityClient.SendTx(tx)
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
		info, err := node.gravityClient.HttpClient.Status()
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

		time.Sleep(time.Duration(Timeout) * time.Second)
	}
}

func (node *Node) execute(ledgerHeight uint64, roundState map[uint64]*RoundState, ctx context.Context) error {
	tcHeight, err := node.adaptor.GetHeight(ctx)
	if err != nil {
		return err
	}

	roundHeight, err := node.gravityClient.RoundHeight(node.chainType, ledgerHeight)
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

		tx, err := transactions.New(node.validator.pubKey, transactions.NewRound, node.validator.privKey, args)
		if err != nil {
			return err
		}
		err = node.gravityClient.SendTx(tx)
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

		_, err := node.gravityClient.CommitHash(node.chainType, node.nebulaId, int64(tcHeight), node.oraclePubKey)
		if err != nil && err != client.ErrValueNotFound {
			return err
		}

		data, err := node.extractor.Extract(ctx)
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
		_, err := node.gravityClient.Reveal(node.nebulaId, int64(tcHeight), roundState[tcHeight].commitHash)
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

		_, err := node.gravityClient.Result(node.chainType, node.nebulaId, int64(tcHeight), node.oraclePubKey)
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

		oraclesMap, err := node.gravityClient.BftOraclesByNebula(node.chainType, node.nebulaId)
		if err != nil {
			return err
		}
		var count uint64
		for oracle, _ := range oraclesMap {
			oracles = append(oracles, oracle)
			if node.oraclePubKey == oracle {
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

		_, err = node.gravityClient.Result(node.chainType, node.nebulaId, int64(tcHeight), node.oraclePubKey)
		if err == client.ErrValueNotFound {
			return nil
		} else if err != nil {
			return err
		}

		if _, ok := roundState[tcHeight]; !ok && roundState[tcHeight].resultValue != 0 {
			return nil
		}

		txId, err := node.adaptor.SendDataResult(node.nebulaId, tcHeight, oracles, roundState[tcHeight].resultHash, ctx)
		if err != nil {
			return err
		}

		roundState[tcHeight].isSent = true

		err = node.adaptor.WaitTx(txId, ctx)
		if err != nil {
			return err
		}

		err = node.adaptor.SendDataToSubs(node.nebulaId, tcHeight, roundState[tcHeight].resultValue, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

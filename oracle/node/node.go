package node

import (
	"context"
	"errors"
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
	TimeoutMs = 200
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

	MaxPulseCountInBlock uint64
}

func New(nebulaId account.NebulaId, chainType account.ChainType, oracleSecretKey []byte, validator *Validator, extractorUrl string, gravityNodeUrl string, targetChainNodeUrl string, ctx context.Context) (*Node, error) {
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

	exType, err := adaptor.GetExtractorType(nebulaId, ctx)
	if err != nil {
		return nil, err
	}

	return &Node{
		validator: validator,
		nebulaId:  nebulaId,
		extractor: &Extractor{
			ExtractorType: exType,
			Client:        extractor.New(extractorUrl),
		},
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

	nebulae, err := node.gravityClient.Nebulae()
	if err != nil {
		return err
	}
	nebulaInfo, ok := nebulae[node.nebulaId]
	if !ok {
		return errors.New("nebula not found")
	}

	node.MaxPulseCountInBlock = nebulaInfo.MaxPulseCountInBlock
	return nil
}

func (node *Node) Start(ctx context.Context) {
	var lastLedgerHeight uint64
	var lastTcHeight uint64
	var pulseCountInBlock uint64
	roundState := new(RoundState)
	for {
		time.Sleep(time.Duration(TimeoutMs) * time.Millisecond)

		info, err := node.gravityClient.HttpClient.Status()
		if err != nil {
			errorLogger.Print(err)
			continue
		}

		ledgerHeight := uint64(info.SyncInfo.LatestBlockHeight)
		if lastLedgerHeight != ledgerHeight {
			fmt.Printf("Ledger Height: %d\n", ledgerHeight)
			lastLedgerHeight = ledgerHeight
		}

		tcHeight, err := node.adaptor.GetHeight(ctx)
		if err != nil {
			errorLogger.Print(err)
		}

		if tcHeight != lastTcHeight {
			pulseCountInBlock = 0
		}
		if pulseCountInBlock >= node.MaxPulseCountInBlock {
			continue
		}

		if state.CalculateSubRound(ledgerHeight) == state.CommitSubRound {
			roundState = new(RoundState)
		}

		err = node.execute(ledgerHeight, tcHeight, roundState, ctx)
		if err != nil {
			errorLogger.Print(err)
		}

		tcHeight = lastTcHeight
		if state.CalculateSubRound(ledgerHeight) == state.SendToTargetChain {
			pulseCountInBlock++
		}
	}
}

func (node *Node) execute(ledgerHeight uint64, tcHeight uint64, roundState *RoundState, ctx context.Context) error {
	pulseId, err := node.adaptor.LastPulseId(node.nebulaId, ctx)
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
		_, err := node.gravityClient.CommitHash(node.chainType, node.nebulaId, int64(pulseId), node.oraclePubKey)
		if err != nil && err != client.ErrValueNotFound {
			return err
		} else if err == nil {
			return nil
		}

		data, err := node.extractor.Extract(ctx)
		if err != nil {
			return err
		}

		if data == nil {
			return nil
		}

		commit, err := node.commit(data, pulseId)
		if err != nil {
			return err
		}

		roundState.commitHash = commit
		roundState.data = data
	case state.RevealSubRound:
		if roundState.commitHash == nil {
			return nil
		}
		_, err := node.gravityClient.Reveal(node.chainType, node.nebulaId, int64(pulseId), roundState.commitHash)
		if err != nil && err != client.ErrValueNotFound {
			return err
		} else if err == nil {
			return nil
		}

		err = node.reveal(pulseId, roundState.data, roundState.commitHash)
		if err != nil {
			return err
		}
	case state.ResultSubRound:
		if roundState.data == nil {
			return nil
		}

		_, err := node.gravityClient.Reveal(node.chainType, node.nebulaId, int64(pulseId), roundState.commitHash)
		if err != nil && err != client.ErrValueNotFound {
			return err
		} else if err == client.ErrValueNotFound {
			return nil
		}

		_, err = node.gravityClient.Result(node.chainType, node.nebulaId, int64(pulseId), node.oraclePubKey)
		if err != nil && err != client.ErrValueNotFound {
			return err
		} else if err == nil {
			return nil
		}

		value, hash, err := node.signResult(pulseId, ctx)
		if err != nil {
			return err
		}

		roundState.resultValue = value
		roundState.resultHash = hash
	case state.SendToTargetChain:
		var oracles []account.OraclesPubKey
		var myRound uint64

		if roundState.isSent {
			return nil
		}

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

		if tcHeight%uint64(len(oracles)) != myRound {
			return nil
		}

		if roundState.resultValue == nil {
			return nil
		}

		txId, err := node.adaptor.SendDataResult(node.nebulaId, pulseId, oracles, roundState.resultHash, ctx)
		if err != nil {
			return err
		}

		err = node.adaptor.WaitTx(txId, ctx)
		if err != nil {
			return err
		}

		roundState.isSent = true

		err = node.adaptor.SendDataToSubs(node.nebulaId, pulseId, roundState.resultValue, ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

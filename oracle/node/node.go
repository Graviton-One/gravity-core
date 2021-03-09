package node

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Gravity-Tech/gravity-core/abi"
	"go.uber.org/zap"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"log"

	"github.com/Gravity-Tech/gravity-core/common/state"

	"github.com/Gravity-Tech/gravity-core/oracle/extractor"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/transactions"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	TimeoutMs = 1000
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
	ExtractorType abi.ExtractorType
}

type Node struct {
	nebulaId  account.NebulaId
	chainType account.ChainType
	//Фактический номер сети в общем списке
	chainSelector account.ChainType
	validator     *Validator
	oraclePubKey  account.OraclesPubKey
	gravityClient *gravity.Client

	adaptor              adaptors.IBlockchainAdaptor
	extractor            *Extractor
	blocksInterval       uint64
	MaxPulseCountInBlock uint64
}

func New(nebulaId account.NebulaId, chainType string, chainSelector account.ChainType,
	chainId byte, oracleSecretKey []byte, validator *Validator,
	extractorUrl string, gravityNodeUrl string, blocksInterval uint64,
	targetChainNodeUrl string, ctx context.Context) (*Node, error) {

	ghClient, err := gravity.New(gravityNodeUrl)
	if err != nil {
		return nil, err
	}
	opts := adaptors.AdapterOptions{
		"ghClient": ghClient,
	}
	adaptor, err := adaptors.NewFactory().CreateAdaptor(chainType, oracleSecretKey, targetChainNodeUrl, ctx, opts)
	if err != nil {
		return nil, err
	}
	ct, err := account.ParseChainType(chainType)

	if err != nil {
		return nil, err
	}

	exType, err := adaptor.ValueType(nebulaId, ctx)
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
		chainType:      ct,
		chainSelector:  chainSelector,
		adaptor:        adaptor,
		gravityClient:  ghClient,
		oraclePubKey:   adaptor.PubKey(),
		blocksInterval: blocksInterval,
	}, nil
}

func (node *Node) Init() error {
	oraclesByValidator, err := node.gravityClient.OraclesByValidator(node.validator.pubKey)
	if err != nil {
		return err
	}

	oracle, ok := oraclesByValidator[node.chainSelector]
	if !ok || oracle != node.oraclePubKey {
		tx, err := transactions.New(node.validator.pubKey, transactions.AddOracle, node.validator.privKey)
		if err != nil {
			return err
		}

		tx.AddValues([]transactions.Value{
			transactions.BytesValue{
				Value: []byte{byte(node.chainSelector)},
			},
			transactions.BytesValue{
				Value: node.oraclePubKey[:],
			},
		})
		err = node.gravityClient.SendTx(tx)
		if err != nil {
			return err
		}

		zap.L().Sugar().Infof("Add oracle (TXID): %s\n", hexutil.Encode(tx.Id[:]))
		time.Sleep(time.Duration(5) * time.Second)
	}

	oraclesByNebulaKey, err := node.gravityClient.OraclesByNebula(node.nebulaId, node.chainSelector)
	if err != nil {
		return err
	}

	_, ok = oraclesByNebulaKey[node.oraclePubKey.ToString(node.chainType)]
	if !ok {
		tx, err := transactions.New(node.validator.pubKey, transactions.AddOracleInNebula, node.validator.privKey)
		if err != nil {
			return err
		}

		tx.AddValues([]transactions.Value{
			transactions.BytesValue{
				Value: node.nebulaId[:],
			},
			transactions.BytesValue{
				Value: node.oraclePubKey[:],
			},
		})

		err = node.gravityClient.SendTx(tx)
		if err != nil {
			return err
		}

		zap.L().Sugar().Infof("Add oracle in nebula (TXID): %s\n", hexutil.Encode(tx.Id[:]))
		time.Sleep(time.Duration(5) * time.Second)
	}

	nebulaInfo, err := node.gravityClient.NebulaInfo(node.nebulaId, node.chainType, node.chainSelector)
	if err == gravity.ErrValueNotFound {
		return errors.New("nebula not found")
	} else if err != nil {
		return err
	}

	node.MaxPulseCountInBlock = nebulaInfo.MaxPulseCountInBlock
	return nil
}

func (node *Node) Start(ctx context.Context) {
	var lastLedgerHeight uint64
	var lastTcHeight uint64
	var pulseCountInBlock uint64
	var lastPulseId uint64

	roundState := new(RoundState)
	for {
		time.Sleep(time.Duration(TimeoutMs) * time.Millisecond)
		if pulseCountInBlock >= node.MaxPulseCountInBlock {
			continue
		}

		newLastPulseId, err := node.adaptor.LastPulseId(node.nebulaId, ctx)
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}

		if lastPulseId != newLastPulseId {
			lastPulseId = newLastPulseId
			roundState = new(RoundState)
		}

		tcHeight, err := node.adaptor.GetHeight(ctx)
		if err != nil {
			zap.L().Error(err.Error())
		}

		if tcHeight != lastTcHeight {
			zap.L().Sugar().Infof("Tc Height: %d\n", tcHeight)
			lastTcHeight = tcHeight
			if tcHeight%node.blocksInterval == 0 {
				pulseCountInBlock = 0
				roundState = new(RoundState)
			}
		}

		oraclesMap, err := node.gravityClient.BftOraclesByNebula(node.chainType, node.chainSelector, node.nebulaId)
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}
		if _, ok := oraclesMap[node.oraclePubKey.ToString(node.chainType)]; !ok {
			continue
		}

		info, err := node.gravityClient.HttpClient.Status()
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}

		ledgerHeight := uint64(info.SyncInfo.LatestBlockHeight)
		if lastLedgerHeight != ledgerHeight {
			zap.L().Sugar().Infof("Ledger Height: %d\n", ledgerHeight)
			lastLedgerHeight = ledgerHeight
		}

		interval := (tcHeight - 2*node.blocksInterval/state.SubRoundCount) / node.blocksInterval

		fmt.Printf("Interval: %d\n", interval)
		round := state.CalculateSubRound(tcHeight, node.blocksInterval)

		fmt.Printf("Round: %d\n", round)

		err = node.execute(lastPulseId+1, round, tcHeight, interval, roundState, ctx)
		if err != nil {
			zap.L().Error(err.Error())
		}
	}
}

func (node *Node) execute(pulseId uint64, round state.SubRound, tcHeight uint64, intervalId uint64, roundState *RoundState, ctx context.Context) error {
	switch round {
	case state.CommitSubRound:
		zap.L().Debug("Commit subround")
		if roundState.commitHash != nil {
			return nil
		}
		_, err := node.gravityClient.CommitHash(node.chainType, node.chainSelector, node.nebulaId, int64(intervalId), int64(pulseId), node.oraclePubKey)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		} else if err == nil {
			return nil
		}

		data, err := node.extractor.Extract(ctx)
		if err != nil && err != extractor.NotFoundErr {
			return err
		} else if err == extractor.NotFoundErr {
			return nil
		}

		if data == nil {
			return nil
		}

		commit, err := node.commit(data, intervalId, pulseId)
		if err != nil {
			return err
		}

		roundState.commitHash = commit
		roundState.data = data
	case state.RevealSubRound:
		zap.L().Debug("Reveal subround")
		if roundState.commitHash == nil || roundState.RevealExist {
			return nil
		}
		_, err := node.gravityClient.Reveal(node.chainType, node.chainSelector, node.oraclePubKey, node.nebulaId, int64(intervalId), int64(pulseId), roundState.commitHash)
		if err != nil && err != gravity.ErrValueNotFound {
			return err
		} else if err == nil {
			return nil
		}

		err = node.reveal(intervalId, pulseId, roundState.data, roundState.commitHash)
		if err != nil {
			return err
		}
		roundState.RevealExist = true
	case state.ResultSubRound:
		zap.L().Debug("Result subround")
		if roundState.data == nil && !roundState.RevealExist {
			return nil
		}

		value, hash, err := node.signResult(intervalId, pulseId, ctx)
		if err != nil {
			return err
		}
		//TODO migrate to err
		if value == nil {
			zap.L().Sugar().Debugf("Value is nil: %t", value == nil)
			return nil
		}

		roundState.resultValue = value
		roundState.resultHash = hash
	case state.SendToTargetChain:
		zap.L().Debug("Send to target chain subround")
		var oracles []account.OraclesPubKey
		var myRound uint64

		if roundState.isSent || roundState.resultValue == nil {
			zap.L().Sugar().Debugf("roundState.isSent: %t, resultValue is nil: %t", roundState.isSent, roundState.resultValue == nil)
			return nil
		}

		oraclesMap, err := node.gravityClient.BftOraclesByNebula(node.chainType, node.chainSelector, node.nebulaId)
		if err != nil {
			zap.L().Sugar().Debugf("BFT error: %s , \n %s", err, zap.Stack("trace").String)
			return nil
		}
		if _, ok := oraclesMap[node.oraclePubKey.ToString(node.chainType)]; !ok {
			zap.L().Debug("Oracle not found")
			return nil
		}

		var count uint64
		for k, v := range oraclesMap {
			oracle, err := account.StringToOraclePubKey(k, v)
			if err != nil {
				return err
			}
			oracles = append(oracles, oracle)
			if node.oraclePubKey == oracle {
				myRound = count
			}
			count++
		}

		if len(oracles) == 0 {
			zap.L().Debug("Oracles map is empty")
			return nil
		}
		if tcHeight%uint64(len(oracles)) != myRound {
			zap.L().Debug("Len oracles != myRound")
			return nil
		}
		zap.L().Debug("Adding pulse")
		txId, err := node.adaptor.AddPulse(node.nebulaId, pulseId, oracles, roundState.resultHash, ctx)

		if err != nil {
			zap.L().Sugar().Debugf("Error: %s", err)
			return err
		}

		if txId != "" {
			err = node.adaptor.WaitTx(txId, ctx)
			if err != nil {
				return err
			}

			zap.L().Sugar().Infof("Result tx id: %s", txId)

			roundState.isSent = true
			zap.L().Debug("Sending Value to subs")
			err = node.adaptor.SendValueToSubs(node.nebulaId, pulseId, roundState.resultValue, ctx)
			if err != nil {
				return err
			}
		} else {
			fmt.Printf("Info: Tx result not sent")
		}
	}
	return nil
}

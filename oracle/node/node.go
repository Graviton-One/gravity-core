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

	validator     *Validator
	oraclePubKey  account.OraclesPubKey
	gravityClient *gravity.Client

	adaptor              adaptors.IBlockchainAdaptor
	extractor            *Extractor
	blocksInterval       uint64
	MaxPulseCountInBlock uint64
}

func New(nebulaId account.NebulaId, chainType account.ChainType,
	chainId byte, oracleSecretKey []byte, validator *Validator,
	extractorUrl string, gravityNodeUrl string, blocksInterval uint64,
	targetChainNodeUrl string, ctx context.Context, customParams map[string]interface{}) (*Node, error) {

	ghClient, err := gravity.New(gravityNodeUrl)
	if err != nil {
		return nil, err
	}

	var adaptor adaptors.IBlockchainAdaptor
	switch chainType {
	case account.Heco:
		adaptor, err = adaptors.NewHecoAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.HecoAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Fantom:
		adaptor, err = adaptors.NewFantomAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.FantomAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Avax:
		adaptor, err = adaptors.NewAvaxAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.AvaxAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Binance:
		adaptor, err = adaptors.NewBinanceAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.BinanceAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Polygon:
		adaptor, err = adaptors.NewPolygonAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.PolygonAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Ethereum:
		adaptor, err = adaptors.NewEthereumAdaptor(oracleSecretKey, targetChainNodeUrl, ctx, adaptors.EthAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Waves:
		adaptor, err = adaptors.NewWavesAdapter(oracleSecretKey, targetChainNodeUrl, chainId, adaptors.WavesAdapterWithGhClient(ghClient))
		if err != nil {
			return nil, err
		}
	case account.Solana:
		adaptor, err = adaptors.NewSolanaAdaptor(oracleSecretKey, targetChainNodeUrl, adaptors.SolanaAdapterWithGhClient(ghClient), adaptors.SolanaAdapterWithCustom(customParams))
		if err != nil {
			return nil, err
		}
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
		chainType:      chainType,
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

	oracle, ok := oraclesByValidator[node.chainType]
	if !ok || oracle != node.oraclePubKey {
		tx, err := transactions.New(node.validator.pubKey, transactions.AddOracle, node.validator.privKey)
		if err != nil {
			return err
		}

		tx.AddValues([]transactions.Value{
			transactions.BytesValue{
				Value: []byte{byte(node.chainType)},
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

	oraclesByNebulaKey, err := node.gravityClient.OraclesByNebula(node.nebulaId, node.chainType)
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

	nebulaInfo, err := node.gravityClient.NebulaInfo(node.nebulaId, node.chainType)
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

	node.gravityClient.HttpClient.WSEvents.Start()
	defer node.gravityClient.HttpClient.WSEvents.Stop()
	ch, err := node.gravityClient.HttpClient.WSEvents.Subscribe(ctx, "gravity-oracle", "tm.event='NewBlock'", 999)

	if err != nil {
		zap.L().Sugar().Debug("Subscribe Error: ", err.Error())
		panic(err)
	}
	// else {
	// 	go func() {
	// 		for {
	// 			a := <-ch
	// 			zap.L().Sugar().Debug(a)
	// 		}
	// 	}()
	// }

	roundState := new(RoundState)
	for {
		//time.Sleep(time.Duration(TimeoutMs) * time.Millisecond)
		<-ch
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
		zap.L().Sugar().Debugf("Round Loop Pulse new: %d last:%d", newLastPulseId, lastPulseId)
		tcHeight, err := node.adaptor.GetHeight(ctx)
		if err != nil {
			zap.L().Error(err.Error())
		}
		checkRound := state.CalculateSubRound(tcHeight, node.blocksInterval)
		if checkRound == state.CommitSubRound {
			if tcHeight != lastTcHeight {
				zap.L().Sugar().Infof("Tc Height: %d\n", tcHeight)
				lastTcHeight = tcHeight
				if tcHeight%node.blocksInterval == 0 {
					pulseCountInBlock = 0
					roundState = new(RoundState)
				}
			}
		}

		oraclesMap, err := node.gravityClient.BftOraclesByNebula(node.chainType, node.nebulaId)
		if err != nil {
			zap.L().Error(err.Error())
			continue
		}
		zap.L().Sugar().Debug("oracles: ", oraclesMap)
		if _, ok := oraclesMap[node.oraclePubKey.ToString(node.chainType)]; !ok {
			zap.L().Sugar().Debugf("oracle [%s] not found in map", node.oraclePubKey.ToString(node.chainType))
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

		interval := (lastTcHeight - 2*node.blocksInterval/state.SubRoundCount) / node.blocksInterval

		fmt.Printf("Interval: %d\n", interval)
		fmt.Printf("TcHeight: %d\n", tcHeight)
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
		zap.L().Sugar().Debugf("Commit subround pulseId: %d", pulseId)
		if roundState.commitHash != nil {
			return nil
		}
		_, err := node.gravityClient.CommitHash(node.chainType, node.nebulaId, int64(intervalId), int64(pulseId), node.oraclePubKey)
		if err != nil && err != gravity.ErrValueNotFound {
			zap.L().Error(err.Error())
			return err
		} else if err == nil {
			return nil
		}

		data, err := node.extractor.Extract(ctx)
		if err != nil && err != extractor.NotFoundErr {
			zap.L().Error(err.Error())
			return err
		} else if err == extractor.NotFoundErr {
			return nil
		}

		if data == nil {
			zap.L().Debug("Commit subround Extractor Data is empty")
			return nil
		}
		zap.L().Sugar().Debug("Extracted data ", data)

		commit, err := node.commit(data, intervalId, pulseId)
		if err != nil {
			return err
		}

		roundState.commitHash = commit
		roundState.data = data
		zap.L().Sugar().Debug("Commit round end ", roundState)
	case state.RevealSubRound:
		zap.L().Debug("Reveal subround")
		if roundState.commitHash == nil || roundState.RevealExist {
			zap.L().Sugar().Debugf("CommitHash is nil: %t, RevealExist: %t", roundState.commitHash == nil, roundState.RevealExist)
			return nil
		}
		_, err := node.gravityClient.Reveal(node.chainType, node.oraclePubKey, node.nebulaId, int64(intervalId), int64(pulseId), roundState.commitHash)
		if err != nil && err != gravity.ErrValueNotFound {
			zap.L().Error(err.Error())
			return err
		} else if err == nil {
			return nil
		}

		err = node.reveal(intervalId, pulseId, roundState.data, roundState.commitHash)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		roundState.RevealExist = true
		zap.L().Sugar().Debug("Reveal round end ", roundState)
	case state.ResultSubRound:
		zap.L().Debug("Result subround")
		if roundState.data == nil && !roundState.RevealExist {
			return nil
		}
		if roundState.resultValue != nil {
			zap.L().Debug("Round sign exists")
			return nil
		}
		value, hash, err := node.signResult(intervalId, pulseId, ctx)
		if err != nil {
			zap.L().Error(err.Error())
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

		oraclesMap, err := node.gravityClient.BftOraclesByNebula(node.chainType, node.nebulaId)
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
		if intervalId%uint64(len(oracles)) != myRound {
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
				zap.L().Sugar().Debugf("Error: %s", err)
				return err
			}

			zap.L().Sugar().Infof("Result tx id: %s", txId)

			roundState.isSent = true
			zap.L().Debug("Sending Value to subs")
			err = node.adaptor.SendValueToSubs(node.nebulaId, pulseId, roundState.resultValue, ctx)
			if err != nil {
				zap.L().Sugar().Debugf("Error: %s", err)
				return err
			}
		} else {
			fmt.Printf("Info: Tx result not sent")
		}
	}
	return nil
}

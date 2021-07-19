package node

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/mr-tron/base58"
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

// IntPow calculates n to the mth power. Since the result is an int, it is assumed that m is a positive power
func IntPow(n, m int) int {
	if m == 0 {
		return 1
	}
	result := n
	for i := 2; i <= m; i++ {
		result *= n
	}
	return result
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
	zap.L().Sugar().Debug("Creating oracle, nbula is ", nebulaId)
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
	firstCommitIteration := true
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
	attempts := 1
	for {
		//time.Sleep(time.Duration(TimeoutMs) * time.Millisecond)
		<-ch
		if pulseCountInBlock >= node.MaxPulseCountInBlock {
			continue
		}

		newLastPulseId, err := node.adaptor.LastPulseId(node.nebulaId, ctx)
		if err != nil {
			zap.L().Error(err.Error())
			time.Sleep(time.Millisecond * time.Duration(IntPow(2, attempts)) * 20)
			attempts = attempts + 1
			continue
		}
		attempts = 1

		if lastPulseId != newLastPulseId {
			lastPulseId = newLastPulseId
			roundState = &RoundState{
				data:        nil,
				commitHash:  []byte{},
				resultValue: nil,
				resultHash:  []byte{},
				isSent:      false,
				commitSent:  false,
				RevealExist: false,
			}
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
				if firstCommitIteration {
					pulseCountInBlock = 0
					roundState = &RoundState{
						data:        nil,
						commitHash:  []byte{},
						resultValue: nil,
						resultHash:  []byte{},
						isSent:      false,
						commitSent:  false,
						RevealExist: false,
					}
					firstCommitIteration = false
				}
			}
		} else {
			firstCommitIteration = true
		}

		zap.L().Sugar().Debugf("getting oracles for nebula [%s] chain type [%s]", base58.Encode(node.nebulaId[:]), node.chainType)
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("execute round", r)
		}
	}()

	return oracleRoundExecutor.Execute(node, &roundExecuteProps{
		PulseID:    pulseId,
		Round:      round,
		IntervalID: intervalId,
		RoundState: roundState,
		Ctx:        ctx,
	})
}

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/ledger-node/app"
	"github.com/Gravity-Tech/gravity-core/ledger-node/blockchain"
	"github.com/Gravity-Tech/gravity-core/ledger-node/scheduler"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/spf13/viper"

	"github.com/dgraph-io/badger"

	cfg "github.com/tendermint/tendermint/config"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	"github.com/tendermint/tendermint/libs/log"
	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/proxy"
)

var configFile, db string

func init() {
	flag.StringVar(&db, "db", "./.db", "Path to config.toml")
	flag.StringVar(&configFile, "config", "./data/config/config.toml", "Path to config.toml")
	flag.Parse()
}

func main() {
	flag.Parse()
	db, err := badger.Open(badger.DefaultOptions(db).WithTruncate(true))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open badger db: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	node, err := newTendermint(db, configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	node.Start()
	defer func() {
		node.Stop()
		node.Wait()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}

func newTendermint(db *badger.DB, configFile string) (*nm.Node, error) {
	//TODO refactoring

	// read config
	config := cfg.DefaultConfig()
	config.RootDir = filepath.Dir(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read config file: %w", err)
	} else if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("viper failed to unmarshal config: %w", err)
	} else if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	initScore := viper.GetString("initScore")
	initScoreMap := make(map[string]uint64)
	for _, v := range strings.Split(initScore, ",") {
		if v == "" {
			continue
		}
		elements := strings.Split(v, "@")
		value, err := strconv.ParseUint(elements[1], 10, 64)
		if err != nil {
			return nil, err
		}
		initScoreMap[elements[0]] = value
	}

	// read private validator
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	ctx := context.Background()
	ethClientHost := viper.GetString("ethNodeUrl")
	wavesClientHost := viper.GetString("wavesNodeUrl")

	contracts := viper.GetStringMapString("targetContracts")
	privKeys := viper.GetStringMapString("privKeys")
	wavesPrivKey := crypto.MustBytesFromBase58(privKeys["waves"])
	ethereumPrivKey, err := hexutil.Decode(privKeys["ethereum"])
	if err != nil {
		return nil, err
	}

	ledgerPrivKey := ed25519.PrivKeyEd25519{}
	copy(ledgerPrivKey[:], pv.Key.PrivKey.Bytes()[5:])

	ledgerPubKey := ed25519.PubKeyEd25519{}
	lPubKey, err := pv.GetPubKey()
	if err != nil {
		return nil, err
	}
	copy(ledgerPubKey[:], lPubKey.Bytes()[5:])
	ledger := &scheduler.LedgerValidator{
		PrivKey: ledgerPrivKey,
		PubKey:  account.ConsulPubKey(ledgerPubKey),
	}
	nebulae := make(map[account.ChainType][][]byte)
	nebulaeCfg := viper.GetStringMap("nebulae")
	for chainType, v := range nebulaeCfg {
		nebulaeStrings := v.([]interface{})
		var nebulaeBytes [][]byte
		for _, v := range nebulaeStrings {
			var b []byte
			switch chainType {
			case "waves":
				b = crypto.MustBytesFromBase58(v.(string))
			case "ethereum":
				b, err = hexutil.Decode(v.(string))
				if err != nil {
					continue
				}
			}

			nebulaeBytes = append(nebulaeBytes, b)
		}

		switch chainType {
		case "waves":
			nebulae[account.Waves] = nebulaeBytes
		case "ethereum":
			nebulae[account.Ethereum] = nebulaeBytes
		}
	}

	blockchains := make(map[account.ChainType]blockchain.IBlockchain)
	wavesBlockchain, err := blockchain.NewWaves(contracts["waves"], wavesPrivKey, wavesClientHost)
	if err != nil {
		return nil, err
	}
	blockchains[account.Waves] = wavesBlockchain

	ethBlockchain, err := blockchain.NewEthereum(contracts["ethereum"], ethereumPrivKey, ethClientHost, ctx)
	if err != nil {
		return nil, err
	}
	blockchains[account.Ethereum] = ethBlockchain

	s, err := scheduler.New(blockchains, config.RPC.ListenAddress, context.Background(), ledger, nebulae)
	if err != nil {
		return nil, err
	}

	application, err := app.NewGHApplication(ethClientHost, wavesClientHost, s, db, initScoreMap, ctx)
	if err != nil {
		return nil, err
	}

	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	// read node key
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load node's key: %w", err)
	}
	// create node
	node, err := nm.NewNode(
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(application),
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger)

	if err != nil {
		return nil, fmt.Errorf("failed to create new Tendermint node: %w", err)
	}

	return node, nil
}

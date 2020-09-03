package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/transactions"

	"github.com/Gravity-Tech/gravity-core/ledger/config"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/ledger/app"
	"github.com/Gravity-Tech/gravity-core/ledger/scheduler"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/tendermint/tendermint/abci/types"

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

var tendermintConfigFile, gravityConfigFile, dbPath string

func init() {
	flag.StringVar(&dbPath, "db", "./.db", "Path to config.toml")
	flag.StringVar(&tendermintConfigFile, "tendermintConfig", "./data/config/config.toml", "Path to config.toml")
	flag.StringVar(&gravityConfigFile, "gravityConfig", "./gravity.config", "Path to config.toml")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	db, err := badger.Open(badger.DefaultOptions(dbPath).WithTruncate(true))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	node, err := newNode(db, tendermintConfigFile, gravityConfigFile, ctx)
	if err != nil {
		panic(err)
	}

	err = node.Start()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := node.Stop()
		node.Wait()

		if err != nil {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}

func newNode(db *badger.DB, tendermintConfigFile string, gravityConfigFile string, ctx context.Context) (*nm.Node, error) {
	// read config
	config := cfg.DefaultConfig()
	config.RootDir = filepath.Dir(filepath.Dir(tendermintConfigFile))

	viper.SetConfigFile(tendermintConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read config file: %w", err)
	} else if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("viper failed to unmarshal config: %w", err)
	} else if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	// read private validator
	pv := privval.LoadFilePV(
		config.PrivValidatorKeyFile(),
		config.PrivValidatorStateFile(),
	)

	// create logger
	logger, err := tmflags.ParseLogLevel(config.LogLevel, log.NewTMLogger(log.NewSyncWriter(os.Stdout)), cfg.DefaultLogLevel())
	if err != nil {
		return nil, fmt.Errorf("failed to parse log level: %w", err)
	}

	// read node key
	nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
	if err != nil {
		return nil, fmt.Errorf("failed to load node's key: %w", err)
	}

	app, err := crateApp(db, pv, gravityConfigFile, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse gravity config: %w", err)
	}

	// create node
	node, err := nm.NewNode(
		config,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(app),
		nm.DefaultGenesisDocProviderFunc(config),
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(config.Instrumentation),
		logger)

	if err != nil {
		return nil, fmt.Errorf("failed to create new Tendermint node: %w", err)
	}

	return node, nil
}

func crateApp(db *badger.DB, pv *privval.FilePV, configFile string, ctx context.Context) (types.Application, error) {
	cfg, err := config.Load(configFile)
	if err != nil {
		return nil, err
	}

	ledgerPubKey := account.ConsulPubKey{}
	lPubKey, err := pv.GetPubKey()
	if err != nil {
		return nil, err
	}

	copy(ledgerPubKey[:], lPubKey.Bytes()[5:])

	ledgerValidator := &scheduler.LedgerValidator{
		PrivKey: pv.Key.PrivKey,
		PubKey:  (ledgerPubKey),
	}

	nodeUrls := make(map[account.ChainType]string)
	adaptersConfig := make(map[account.ChainType]*scheduler.AdaptorConfig)
	for k, v := range cfg.Adapters {
		chainType, err := account.ParseChainType(k)
		if err != nil {
			return nil, err
		}
		privKey, err := account.StringToPrivKey(v.PrivKey, chainType)
		if err != nil {
			return nil, err
		}

		var adaptor adaptors.IBlockchainAdaptor

		switch chainType {
		case account.Ethereum:
			adaptor, err = adaptors.NewEthereumAdaptor(privKey, v.NodeUrl, ctx, adaptors.WithEthereumGravityContract(v.GravityContractAddress))
			if err != nil {
				return nil, err
			}
		case account.Waves:
			adaptor, err = adaptors.NewWavesAdapter(privKey, v.NodeUrl, adaptors.WithWavesGravityContract(v.GravityContractAddress))
			if err != nil {
				return nil, err
			}
		}

		var nebulae []account.NebulaId
		for _, address := range v.Nebulae {
			nebulaId, err := account.StringToNebulaId(address, chainType)
			if err != nil {
				return nil, err
			}
			nebulae = append(nebulae, nebulaId)
		}
		adaptersConfig[chainType] = &scheduler.AdaptorConfig{
			IBlockchainAdaptor: adaptor,
			Nebulae:            nebulae,
		}
		nodeUrls[chainType] = v.NodeUrl

		if cfg.BootstrapUrl != nil {
			err := setOraclePubKey(*cfg.BootstrapUrl, ledgerValidator.PubKey, ledgerValidator.PrivKey, adaptor.PubKey(), chainType)
			if err != nil {
				return nil, err
			}
		}
	}

	blockScheduler, err := scheduler.New(adaptersConfig, ledgerValidator, ctx)
	if err != nil {
		return nil, err
	}

	genesis := new(app.Genesis)
	for k, v := range cfg.Genesis.InitScore {
		pubKey, err := account.HexToValidatorPubKey(k)
		if err != nil {
			return nil, err
		}
		genesis.InitScore[pubKey] = v
	}
	for k, v := range cfg.Genesis.OraclesAddressByValidator {
		validatorPubKey, err := account.HexToValidatorPubKey(k)
		if err != nil {
			return nil, err
		}
		for chainTypeString, oracle := range v {
			chainType, err := account.ParseChainType(chainTypeString)
			if err != nil {
				return nil, err
			}

			oraclePubKey, err := account.StringToOraclePubKey(oracle, chainType)
			if err != nil {
				return nil, err
			}

			genesis.OraclesAddressByValidator[validatorPubKey][chainType] = oraclePubKey
		}
	}

	application, err := app.NewGHApplication(nodeUrls[account.Ethereum], nodeUrls[account.Waves], blockScheduler, db, genesis, ctx)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func setOraclePubKey(bootstrapUrl string, pubKey account.ConsulPubKey, privKey ed25519.PrivKeyEd25519, oracle account.OraclesPubKey, chainType account.ChainType) error {
	gravityClient, err := client.NewGravityClient(bootstrapUrl)
	if err != nil {
		return err
	}

	oracles, err := gravityClient.OraclesByValidator(pubKey)
	if err != nil {
		return err
	}

	if _, ok := oracles[chainType]; ok {
		return nil
	}

	args := []transactions.Args{
		{
			Value: chainType,
		},
		{
			Value: oracle,
		},
	}

	tx, err := transactions.New(pubKey, transactions.AddOracle, privKey, args)
	if err != nil {
		return err
	}

	err = gravityClient.SendTx(tx)
	if err != nil {
		return err
	}

	return nil
}

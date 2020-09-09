package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/tendermint/tendermint/privval"

	"github.com/tendermint/tendermint/crypto"

	nm "github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/Gravity-Tech/gravity-core/rpc"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/ledger/app"
	"github.com/Gravity-Tech/gravity-core/ledger/scheduler"
	"github.com/dgraph-io/badger"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types"

	cfg "github.com/tendermint/tendermint/config"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/tendermint/tendermint/p2p"

	tOs "github.com/tendermint/tendermint/libs/os"

	"github.com/Gravity-Tech/gravity-core/config"
	"github.com/urfave/cli/v2"
)

type Network string
type ChainId string

const (
	DefaultBftOracleInNebulaCount = 5
	PrivateRPCHostFlag            = "rpc"
	NetworkFlag                   = "network"
	BootstrapUrlFlag              = "bootstrap"

	Custom Network = "custom"
	DevNet Network = "devnet"

	DevNetId ChainId = "gravity-devnet"
	CustomId ChainId = "gravity-custom"

	DefaultBootstrapUrl   = "http://127.0.0.1:5001"
	DefaultPrivateRpcHost = "http://127.0.0.1:2500"

	DefaultSeeds = "asdasd@127.0.0.1:26667"
)

var (
	DevNetConfig = config.LedgerConfig{
		Moniker:    config.DefaultMoniker,
		RPC:        cfg.DefaultRPCConfig(),
		IsFastSync: true,
		Mempool:    cfg.DefaultMempoolConfig(),
		Adapters: map[string]config.AdaptorsConfig{
			account.Ethereum.String(): {
				NodeUrl:                "http://127.0.0.1:8545",
				GravityContractAddress: "0x0000000000000",
			},
			account.Waves.String(): {
				NodeUrl:                "http://127.0.0.1:6869",
				GravityContractAddress: "0x0000000000000",
				ChainId:                "R",
			},
		},
	}
	DevNetGenesis = config.Genesis{
		ConsulsCount: 5,
		GenesisTime:  time.Unix(1599142244, 0),
		ChainID:      string(DevNetId),
		Block: types.BlockParams{
			MaxBytes:   1048576,
			MaxGas:     -1,
			TimeIotaMs: 1000,
		},
		Evidence: types.EvidenceParams{
			MaxAgeNumBlocks: 100000,
			MaxAgeDuration:  1728 * time.Second,
		},
		InitScore: map[string]uint64{
			"0x0000": 100,
		},
		OraclesAddressByValidator: map[string]map[string]string{
			"0x0000": {
				"waves":    "",
				"ethereum": "",
			},
		},
	}

	CustomNetGenesis = config.Genesis{
		GenesisTime: time.Now(),
		ChainID:     string(CustomId),
		Block: types.BlockParams{
			MaxBytes:   1048576,
			MaxGas:     -1,
			TimeIotaMs: 1000,
		},
		Evidence: types.EvidenceParams{
			MaxAgeNumBlocks: 100000,
			MaxAgeDuration:  1728 * time.Second,
		},
	}
)

var (
	LedgerCommand = &cli.Command{
		Name:        "ledger",
		Usage:       "",
		Description: "Commands to control ledger",
		Subcommands: []*cli.Command{
			{
				Name:        "init",
				Usage:       "Generate ledger config",
				Description: "",
				Action:      initLedgerConfig,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  NetworkFlag,
						Value: string(DevNet),
					},
				},
			},
			{
				Name:        "start",
				Usage:       "Start ledger node",
				Description: "",
				Action:      startLedger,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  PrivateRPCHostFlag,
						Value: DefaultPrivateRpcHost,
						Usage: "RPC server host",
					},
					&cli.StringFlag{
						Name:  BootstrapUrlFlag,
						Value: DefaultBootstrapUrl,
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  HomeFlag,
				Value: "./",
				Usage: "Home dir for gravity config and files",
			},
		},
	}
)

func initLedgerConfig(ctx *cli.Context) error {
	var err error

	home := ctx.String(HomeFlag)
	network := Network(ctx.String(NetworkFlag))

	if _, err := os.Stat(home); os.IsNotExist(err) {
		err = os.Mkdir(home, 0644)
		if err != nil {
			return err
		}
	}

	privKeysFile := path.Join(home, PrivKeysConfigFileName)
	if tOs.FileExists(privKeysFile) {
		var privKeysCfg config.PrivKeys
		err = config.ParseConfig(privKeysFile, privKeysCfg)
		if err != nil {
			return err
		}
	} else {
		privKeysCfg, pubKeys, err := config.GeneratePrivKeys()
		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(&privKeysCfg, "", " ")
		err = ioutil.WriteFile(path.Join(home, PrivKeysConfigFileName), b, 0644)
		if err != nil {
			return err
		}

		fmt.Printf("Validator PubKey: %s\n", pubKeys.Validator)
		for k, v := range pubKeys.TargetChains {
			fmt.Printf("%s PubKey: %s\n", k, v)
		}
	}

	var genesis config.Genesis
	if network == DevNet {
		genesis = DevNetGenesis
	} else {
		genesis = CustomNetGenesis
	}
	b, err := json.MarshalIndent(&genesis, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(home, GenesisFileName), b, 0644)
	if err != nil {
		return err
	}

	var ledgerConf config.LedgerConfig
	if network == DevNet {
		ledgerConf = DevNetConfig
		ledgerConf.P2P = cfg.DefaultP2PConfig()
		ledgerConf.P2P.Seeds = DefaultSeeds
	} else {
		ledgerConf = config.DefaultLedgerConfig()
	}
	b, err = json.MarshalIndent(&ledgerConf, "", " ")
	err = ioutil.WriteFile(path.Join(home, LedgerConfigFileName), b, 0644)
	if err != nil {
		return err
	}

	nodeKeyFile := path.Join(home, NodeKeyFileName)
	if !tOs.FileExists(nodeKeyFile) {
		if _, err := p2p.LoadOrGenNodeKey(nodeKeyFile); err != nil {
			return err
		}
	}

	keyStateFile := path.Join(home, LedgerKeyStateFileName)
	if !tOs.FileExists(keyStateFile) {
		keyState := privval.GenFilePV("", keyStateFile)
		keyState.LastSignState.Save()
	}

	return nil
}

func startLedger(ctx *cli.Context) error {
	home := ctx.String(HomeFlag)
	bootstrap := ctx.String(BootstrapUrlFlag)
	rpcHost := ctx.String(PrivateRPCHostFlag)

	var err error
	sysCtx := context.Background()

	dbDir := path.Join(home, DbDir)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.Mkdir(dbDir, 0644)
		if err != nil {
			return err
		}
	}
	db, err := badger.Open(badger.DefaultOptions(dbDir).WithTruncate(true))
	if err != nil {
		return err
	}
	defer db.Close()

	var privKeysCfg config.PrivKeys
	err = config.ParseConfig(path.Join(home, PrivKeysConfigFileName), &privKeysCfg)
	if err != nil {
		return err
	}

	var genesis config.Genesis
	err = config.ParseConfig(path.Join(home, GenesisFileName), &genesis)
	if err != nil {
		return err
	}

	var ledgerConf config.LedgerConfig
	err = config.ParseConfig(path.Join(home, LedgerConfigFileName), &ledgerConf)
	if err != nil {
		return err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(path.Join(home, NodeKeyFileName))
	if err != nil {
		return err
	}

	tConfig := cfg.DefaultConfig()

	tConfig.P2P = ledgerConf.P2P
	tConfig.Moniker = ledgerConf.Moniker
	tConfig.Mempool = ledgerConf.Mempool
	tConfig.FastSyncMode = ledgerConf.IsFastSync
	tConfig.RPC = ledgerConf.RPC

	tConfig.RootDir = home
	tConfig.P2P.RootDir = home
	tConfig.Consensus.RootDir = home
	tConfig.Consensus.TimeoutCommit = time.Second * 3
	tConfig.RPC.RootDir = home
	tConfig.Mempool.RootDir = home

	logger, err := tmflags.ParseLogLevel(tConfig.LogLevel, log.NewTMLogger(log.NewSyncWriter(os.Stdout)), cfg.DefaultLogLevel())
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	var ledgerPrivKey ed25519.PrivKeyEd25519
	ledgerPrivKeyBytes, err := hexutil.Decode(privKeysCfg.Validator)
	if err != nil {
		return err
	}
	copy(ledgerPrivKey[:], ledgerPrivKeyBytes)

	var ledgerPubKey account.ConsulPubKey
	ledgerPubKeyBytes := ledgerPrivKey.PubKey().Bytes()[5:]
	copy(ledgerPubKey[:], ledgerPubKeyBytes)

	ledgerValidator := &account.LedgerValidator{
		PrivKey: ledgerPrivKey,
		PubKey:  ledgerPubKey,
	}

	gravityApp, err := crateApp(db, ledgerValidator, privKeysCfg.TargetChains, ledgerConf, genesis, bootstrap, tConfig.RPC.ListenAddress, sysCtx)
	if err != nil {
		return fmt.Errorf("failed to parse gravity config: %w", err)
	}

	gravityApp.IsSync = true
	var validators []types.GenesisValidator
	for k, v := range genesis.InitScore {
		pubKey, err := account.HexToValidatorPubKey(k)
		if err != nil {
			return err
		}

		validators = append(validators, types.GenesisValidator{
			PubKey: ed25519.PubKeyEd25519(pubKey),
			Power:  int64(v),
		})
	}

	pv := privval.GenFilePV("", path.Join(home, LedgerKeyStateFileName))
	pv.Key = privval.FilePVKey{
		Address: ledgerValidator.PrivKey.PubKey().Address(),
		PubKey:  ledgerValidator.PrivKey.PubKey(),
		PrivKey: ledgerValidator.PrivKey,
	}

	node, err := nm.NewNode(
		tConfig,
		pv,
		nodeKey,
		proxy.NewLocalClientCreator(gravityApp),
		func() (*types.GenesisDoc, error) {
			return &types.GenesisDoc{
				GenesisTime: genesis.GenesisTime,
				ChainID:     genesis.ChainID,
				ConsensusParams: &types.ConsensusParams{
					Block:    genesis.Block,
					Evidence: genesis.Evidence,
					Validator: types.ValidatorParams{
						PubKeyTypes: []string{"ed25519"},
					},
				},
				Validators: validators,
			}, nil
		},
		nm.DefaultDBProvider,
		nm.DefaultMetricsProvider(tConfig.Instrumentation),
		logger)
	if err != nil {
		return nil
	}

	gravityApp.IsSync = false
	err = node.Start()
	if err != nil {
		return err
	}
	defer func() {
		err := node.Stop()
		node.Wait()

		if err != nil {
			panic(err)
		}
	}()

	rpcConfig, err := rpc.NewConfig(rpcHost, tConfig.RPC.ListenAddress, ledgerValidator.PrivKey)
	if err != nil {
		return err
	}
	go rpc.ListenRpcServer(rpcConfig)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}

func crateApp(db *badger.DB, ledgerValidator *account.LedgerValidator, privKeys map[string]string, cfg config.LedgerConfig, genesisCfg config.Genesis, bootstrap string, localHost string, ctx context.Context) (*app.GHApplication, error) {
	bAdaptors := make(map[account.ChainType]adaptors.IBlockchainAdaptor)
	for k, v := range cfg.Adapters {
		chainType, err := account.ParseChainType(k)
		if err != nil {
			return nil, err
		}

		privKey, err := account.StringToPrivKey(privKeys[k], chainType)
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
			adaptor, err = adaptors.NewWavesAdapter(privKey, v.NodeUrl, v.ChainId[0], adaptors.WithWavesGravityContract(v.GravityContractAddress))
			if err != nil {
				return nil, err
			}
		}

		bAdaptors[chainType] = adaptor
		if bootstrap != "" {
			err := setOraclePubKey(bootstrap, ledgerValidator.PubKey, ledgerValidator.PrivKey, adaptor.PubKey(), chainType)
			if err != nil {
				return nil, err
			}
		}
	}
	blockScheduler, err := scheduler.New(bAdaptors, ledgerValidator, localHost, ctx)
	if err != nil {
		return nil, err
	}

	genesis := app.Genesis{
		ConsulsCount:              genesisCfg.ConsulsCount,
		BftOracleInNebulaCount:    DefaultBftOracleInNebulaCount,
		OraclesAddressByValidator: make(map[account.ConsulPubKey][]app.OraclesAddresses),
	}

	for k, v := range genesisCfg.OraclesAddressByValidator {
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

			genesis.OraclesAddressByValidator[validatorPubKey] = append(genesis.OraclesAddressByValidator[validatorPubKey], app.OraclesAddresses{
				ChainType:     chainType,
				OraclesPubKey: oraclePubKey,
			})
		}
	}

	application, err := app.NewGHApplication(bAdaptors, blockScheduler, db, &genesis, ctx)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func setOraclePubKey(bootstrapUrl string, pubKey account.ConsulPubKey, privKey crypto.PrivKey, oracle account.OraclesPubKey, chainType account.ChainType) error {
	gravityClient, err := gravity.New(bootstrapUrl)
	if err != nil {
		return err
	}

	oracles, err := gravityClient.OraclesByValidator(pubKey)
	if err != nil && err != gravity.ErrValueNotFound {
		return err
	}

	if _, ok := oracles[chainType]; ok {
		return nil
	}

	tx, err := transactions.New(pubKey, transactions.AddOracle, privKey)
	if err != nil {
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{
			Value: []byte{byte(chainType)},
		},
		transactions.BytesValue{
			Value: oracle[:],
		},
	})

	err = gravityClient.SendTx(tx)
	if err != nil {
		return err
	}

	return nil
}

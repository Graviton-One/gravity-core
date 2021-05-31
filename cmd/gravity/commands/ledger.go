package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/tendermint/tendermint/privval"
	"go.uber.org/zap"

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

	Custom  Network = "custom"
	Mainnet Network = "mainnet"

	CustomId  ChainId = "gravity-custom"
	MainnetId ChainId = "gravity-mainnet"

	DefaultPrivateRpcHost = "127.0.0.1:2500"

	MainnetBootstrapHost   = "http://134.122.37.128:26657"
	MainnetPersistentPeers = "9e91414d53328d46d68415e1e31934ab09f69511@134.122.37.128:26656,3d8306a9e006687374f23905adc06f3517306cd7@212.111.41.159:26656,2c17ad4dcc9947342ff03be6827cfb285b6494bb@3.135.223.165:26666"
)

var (
	MainnetConfig = config.LedgerConfig{
		Moniker:    config.DefaultMoniker,
		RPC:        cfg.DefaultRPCConfig(),
		IsFastSync: true,
		Mempool:    cfg.DefaultMempoolConfig(),
		Details:    (&config.ValidatorDetails{}).DefaultNew(),
		Adapters: map[string]config.AdaptorsConfig{
			account.Binance.String(): {
				NodeUrl:                "https://bsc-dataseed4.ninicoin.io/",
				ChainType:              account.Binance.String(),
				GravityContractAddress: "0x5b875E3457ce737D42593aB5d6e5cfBF7896a27d",
			},
			account.Waves.String(): {
				NodeUrl:                "https://nodes.swop.fi",
				ChainId:                "W",
				ChainType:              "waves",
				GravityContractAddress: "3PLpMu2cAg618e7xXYHtckFJjFZksPFHoLm",
			},
			account.Ergo.String(): {
				NodeUrl: "http://10.10.10.4:9016",
				ChainType: account.Ergo.String(),
				GravityContractAddress: "",
			},
			account.Heco.String(): {
				NodeUrl:                "https://http-mainnet.hecochain.com",
				GravityContractAddress: "0x8f56C70A8d473e58b47BAc0D0f24eF630064D7ed",
			},
			account.Fantom.String(): {
				NodeUrl:                "https://rpcapi.fantom.network",
				GravityContractAddress: "0xB883418014e73228F1Ec470714802c59bB49f1eC",
			},
		},
	}
	MainnetGenesis = config.Genesis{
		ConsulsCount: 3,
		GenesisTime:  time.Unix(1614613181, 0),
		ChainID:      string(MainnetId),
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
			"0xcacb7145b7b70211ed43acd648878336915d981abea6bd0b3bdd6a4ff5dad5cf": 100,
			"0xc5597e2c40b78e4fecd862f54cc8e8c284f85006afcdc8253564e9f5c452ca9a": 100,
			"0x50fd18a5c1969a2369f778dbbca8f7a7e99e1236d81be3780255d8a3da89f9c4": 100,
		},
		OraclesAddressByValidator: map[string]map[string]string{
			"0xcacb7145b7b70211ed43acd648878336915d981abea6bd0b3bdd6a4ff5dad5cf": {
				"waves": "4ArMUAxJZ3ETB1xSBqJkdhM19TXoEuWHsWzHZqKo3rvY",
				"bsc":   "0x032405b9ef3cc5ed099ee13f8084f972cfdd6cec85835628ead918712b6a0fab65",
			},
			"0xc5597e2c40b78e4fecd862f54cc8e8c284f85006afcdc8253564e9f5c452ca9a": {
				"waves": "51yKUBQ7pxGJ1UNCgwdjKMQswYWHGq28thbdpo8gLoEK",
				"bsc":   "0x02826f06ab27fd6d1f0574c020e6c06010489c5eb07ba4c2aa0149f303f85db215",
			},
			"0x50fd18a5c1969a2369f778dbbca8f7a7e99e1236d81be3780255d8a3da89f9c4": {
				"waves": "FpsFLbAfmUqngqS54knGJe11Y68mgDYoLHT8QRY4fdYD",
				"bsc":   "0x023522b43a9820d64d6b61186efdace2d787931844c6512e3ca2253e38f5c3e522",
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
						Value: string(Custom),
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
						Value: MainnetBootstrapHost,
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

func getPublicIP() (string, error) {
	ifaces, _ := net.Interfaces()

	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		for _, addr := range addrs {

			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if strings.Contains(fmt.Sprintf("%v", addr), "/24") {
				return fmt.Sprintf("%v", ip), nil
			}

		}
	}

	return "", fmt.Errorf("not found valid ip")
}

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
		var privKeysCfg config.Keys
		err = config.ParseConfig(privKeysFile, privKeysCfg)
		if err != nil {
			return err
		}
	} else {
		var keysCfg *config.Keys
		var err error

		if network == Custom {
			keysCfg, err = config.GeneratePrivKeys('S')
		} else if network == Mainnet {
			keysCfg, err = config.GeneratePrivKeys('W')
		}

		if err != nil {
			return err
		}

		b, err := json.MarshalIndent(&keysCfg, "", " ")
		err = ioutil.WriteFile(path.Join(home, PrivKeysConfigFileName), b, 0644)
		if err != nil {
			return err
		}

		fmt.Printf("Validator PubKey: %s\n", keysCfg.Validator.PubKey)
		for k, v := range keysCfg.TargetChains {
			fmt.Printf("%s PubKey: %s\n", k, v.PubKey)
		}
	}

	var genesis config.Genesis
	if network == Mainnet {
		genesis = MainnetGenesis

		dateString := "2021-03-01T15:39:41.222470458Z"
		timeInst := time.Now()
		_ = timeInst.UnmarshalText([]byte(dateString))

		genesis.GenesisTime = timeInst
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
	if network == Mainnet {
		ledgerConf = MainnetConfig
		ledgerConf.P2P = cfg.DefaultP2PConfig()
		ledgerConf.P2P.PersistentPeers = MainnetPersistentPeers
		ledgerConf.P2P.ListenAddress = "tcp://0.0.0.0:26656"
	} else {
		ledgerConf = config.DefaultLedgerConfig()
	}

	ledgerConf.PublicIP, _ = getPublicIP()

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
	zaplog, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(zaplog)

	home := ctx.String(HomeFlag)
	bootstrap := ctx.String(BootstrapUrlFlag)
	rpcHost := ctx.String(PrivateRPCHostFlag)

	var err error
	sysCtx := context.Background()

	dbDir := path.Join(home, DbDir)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		err = os.Mkdir(dbDir, 0644)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
	}
	db, err := badger.Open(badger.DefaultOptions(dbDir).WithTruncate(true))
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	defer db.Close()

	var privKeysCfg config.Keys
	err = config.ParseConfig(path.Join(home, PrivKeysConfigFileName), &privKeysCfg)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	var genesis config.Genesis
	err = config.ParseConfig(path.Join(home, GenesisFileName), &genesis)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	var ledgerConf config.LedgerConfig
	err = config.ParseConfig(path.Join(home, LedgerConfigFileName), &ledgerConf)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	nodeKey, err := p2p.LoadOrGenNodeKey(path.Join(home, NodeKeyFileName))
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	tConfig := cfg.DefaultConfig()

	tConfig.P2P = ledgerConf.P2P
	tConfig.Moniker = ledgerConf.Moniker
	tConfig.Mempool = ledgerConf.Mempool
	tConfig.FastSyncMode = ledgerConf.IsFastSync
	tConfig.RPC = ledgerConf.RPC

	tConfig.RootDir = home
	tConfig.Consensus.RootDir = home
	tConfig.Consensus.TimeoutCommit = time.Second * 3

	logger, err := tmflags.ParseLogLevel(tConfig.LogLevel, log.NewTMLogger(log.NewSyncWriter(os.Stdout)), cfg.DefaultLogLevel())
	if err != nil {
		zap.L().Error(err.Error())
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	var ledgerPrivKey ed25519.PrivKeyEd25519
	ledgerPrivKeyBytes, err := hexutil.Decode(privKeysCfg.Validator.PrivKey)
	if err != nil {
		zap.L().Error(err.Error())
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

	gravityApp, err := createApp(db, ledgerValidator, privKeysCfg.TargetChains, ledgerConf, genesis, bootstrap, tConfig.RPC.ListenAddress, sysCtx)
	if err != nil {
		zap.L().Error(err.Error())
		return fmt.Errorf("failed to parse gravity config: %w", err)
	}

	gravityApp.IsSync = true
	var validators []types.GenesisValidator
	for k, v := range genesis.InitScore {
		pubKey, err := account.HexToValidatorPubKey(k)
		if err != nil {
			zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return err
	}

	gravityApp.IsSync = false
	err = node.Start()
	if err != nil {
		zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return err
	}
	go rpc.ListenRpcServer(rpcConfig)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c

	return nil
}

func createApp(db *badger.DB, ledgerValidator *account.LedgerValidator, privKeys map[string]config.Key, cfg config.LedgerConfig, genesisCfg config.Genesis, bootstrap string, localHost string, ctx context.Context) (*app.GHApplication, error) {
	bAdaptors := make(map[account.ChainType]adaptors.IBlockchainAdaptor)
	for k, v := range cfg.Adapters {
		chainType, err := account.ParseChainType(k)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}

		privKey, err := account.StringToPrivKey(privKeys[k].PrivKey, chainType)
		if err != nil {
			zap.L().Error(err.Error())
			return nil, err
		}

		var adaptor adaptors.IBlockchainAdaptor

		switch chainType {
		case account.Heco:
			adaptor, err = adaptors.NewHecoAdaptor(privKey, v.NodeUrl, ctx, adaptors.WithHecoGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Fantom:
			adaptor, err = adaptors.NewFantomAdaptor(privKey, v.NodeUrl, ctx, adaptors.WithFantomGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Avax:
			adaptor, err = adaptors.NewAvaxAdaptor(privKey, v.NodeUrl, ctx, adaptors.WithAvaxGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Binance:
			adaptor, err = adaptors.NewBinanceAdaptor(privKey, v.NodeUrl, ctx, adaptors.WithBinanceGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Ethereum:
			adaptor, err = adaptors.NewEthereumAdaptor(privKey, v.NodeUrl, ctx, adaptors.WithEthereumGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Waves:
			adaptor, err = adaptors.NewWavesAdapter(privKey, v.NodeUrl, v.ChainId[0], adaptors.WithWavesGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Solana:
			adaptor, err = adaptors.NewSolanaAdaptor(privKey, v.NodeUrl, adaptors.SolanaAdapterWithCustom(v.Custom))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}
		case account.Ergo:
			adaptor, err = adaptors.NewErgoAdapter(privKey, v.NodeUrl, ctx, adaptors.WithErgoGravityContract(v.GravityContractAddress))
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}

		}

		bAdaptors[chainType] = adaptor

		// if bootstrap != "" {
		// 	err = setOraclePubKey(bootstrap, ledgerValidator.PubKey, ledgerValidator.PrivKey, adaptor.PubKey(), chainType)
		// 	if err != nil {
		// 		zap.L().Error(err.Error())
		// 		return nil, err
		// 	}
		// }

	}
	blockScheduler, err := scheduler.New(bAdaptors, ledgerValidator, localHost, ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	genesis := app.Genesis{
		ConsulsCount:              genesisCfg.ConsulsCount,
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
				zap.L().Error(err.Error())
				return nil, err
			}

			oraclePubKey, err := account.StringToOraclePubKey(oracle, chainType)
			if err != nil {
				zap.L().Error(err.Error())
				return nil, err
			}

			genesis.OraclesAddressByValidator[validatorPubKey] = append(genesis.OraclesAddressByValidator[validatorPubKey], app.OraclesAddresses{
				ChainType:     chainType,
				OraclesPubKey: oraclePubKey,
			})
		}
	}

	application, err := app.NewGHApplication(bAdaptors, blockScheduler, db, &genesis, ctx, &cfg)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, err
	}

	return application, nil
}

func setOraclePubKey(bootstrapUrl string, pubKey account.ConsulPubKey, privKey crypto.PrivKey, oracle account.OraclesPubKey, chainType account.ChainType) error {
	gravityClient, err := gravity.New(bootstrapUrl)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	oracles, err := gravityClient.OraclesByValidator(pubKey)
	if err != nil && err != gravity.ErrValueNotFound {
		zap.L().Error(err.Error())
		return err
	}

	if _, ok := oracles[chainType]; ok {
		zap.L().Debug("Oracle exists")
		return nil
	}

	tx, err := transactions.New(pubKey, transactions.AddOracle, privKey)
	if err != nil {
		zap.L().Error(err.Error())
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
		zap.L().Error(err.Error())
		return err
	}

	return nil
}

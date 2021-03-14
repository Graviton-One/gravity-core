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
	DevNet  Network = "devnet"
	Mainnet Network = "mainnet"

	DevNetId ChainId = "gravity-devnet"
	CustomId ChainId = "gravity-custom"

	DefaultBootstrapUrl   = "http://104.248.255.124:26657"
	DefaultPrivateRpcHost = "127.0.0.1:2500"

	DefaultPersistentPeers = "2a0d75cc7833ad4780a1035b633c5bf4ef94ea4c@104.248.255.124:26656,32a091dfea2b4191d710d2609ca21a8abfe585ac@164.90.184.213:26656,34f38d98e78ed7965a56399998d9c1dccba24fe1@164.90.185.82:26656,c22e04514ce4ae0feb3480d03593d34e4713c86d@161.35.207.224:26656"
)

var (
	DevNetConfig = config.LedgerConfig{
		Moniker:    config.DefaultMoniker,
		RPC:        cfg.DefaultRPCConfig(),
		IsFastSync: true,
		Mempool:    cfg.DefaultMempoolConfig(),
		Details:    (&config.ValidatorDetails{}).DefaultNew(),
		Adapters: map[string]config.AdaptorsConfig{
			"ethereum": {
				NodeUrl:                "https://ropsten.infura.io/v3/598efca7168947c6a186e2f85b600be1",
				GravityContractAddress: "0x80C52beF8622cDF368Bf8AaD5ee4A78cB68E2a79",
				ChainType:              "ethereum",
			},
			"waves": {
				NodeUrl:                "https://nodes-stagenet.wavesnodes.com",
				GravityContractAddress: "3MfrQBknYJSnifUxD86yMPTSHEhgcPe3NBq",
				ChainId:                "S",
				ChainType:              "waves",
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
			"0xd7f746727e21ecf461bb8e8926a6aeb0931eb09311ed13d25a182d4f17339d1d": 100,
			"0x0e7f182e6d2a11bd5d8a34531243435e2aeaa0eed7cad3c5361a81328051fa02": 100,
			"0xe09b444f5c5f2fbdca58bdb37a2dcc90d370ff72f28a6d4b6a6ef732c44afa24": 100,
			"0xd70f6fcdac1a6f2292a330cc830db5e9041939ff79a87ff8536040b07378ca02": 100,
		},
		OraclesAddressByValidator: map[string]map[string]string{
			"0xd7f746727e21ecf461bb8e8926a6aeb0931eb09311ed13d25a182d4f17339d1d": {
				"ethereum": "0x038bf7253f2b3b78c7f8fbe856252373b0867098c6b3f7a6cabc6e73552be75697",
				"waves":    "CNVJbuJubqLyTZ99Y8wwuFgiKqoCUvpYCnHQsezE3Qgk",
			},
			"0x0e7f182e6d2a11bd5d8a34531243435e2aeaa0eed7cad3c5361a81328051fa02": {
				"ethereum": "0x03808de8b08ec39c720c04e7699783f1abefff809afc2a8f7e60e9dd59f039ffa8",
				"waves":    "4QWcFszF3shvhReiU26Sj8Te2QqgsfsreEgiTQNeTgB5",
			},
			"0xe09b444f5c5f2fbdca58bdb37a2dcc90d370ff72f28a6d4b6a6ef732c44afa24": {
				"ethereum": "0x0298644b29e125b1293446b3d5f5b6feb12eaf2e3245df08fe74682fe0ddce5c60",
				"waves":    "CcdpQmNU9qc1uKyr2mmkYNiyadvQ3VHrcYCtLqDfrR9a",
			},
			"0xd70f6fcdac1a6f2292a330cc830db5e9041939ff79a87ff8536040b07378ca02": {
				"ethereum": "0x026a6444ca6ad63e3fda46481d125f8fee07b9a5b5131a12393a654800956856b8",
				"waves":    "E5gz7aTwjjbCbFMYmstvcvb6NvoZyZcQSt1wp68qMpBg",
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

		if network == DevNet {
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
		ledgerConf.P2P.PersistentPeers = DefaultPersistentPeers
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

	var privKeysCfg config.Keys
	err = config.ParseConfig(path.Join(home, PrivKeysConfigFileName), &privKeysCfg)
	if err != nil {
		return err
	}
	account.ChainMapper.Assign(privKeysCfg.ChainIds)
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
	tConfig.Consensus.RootDir = home
	tConfig.Consensus.TimeoutCommit = time.Second * 3

	logger, err := tmflags.ParseLogLevel(tConfig.LogLevel, log.NewTMLogger(log.NewSyncWriter(os.Stdout)), cfg.DefaultLogLevel())
	if err != nil {
		return fmt.Errorf("failed to parse log level: %w", err)
	}

	var ledgerPrivKey ed25519.PrivKeyEd25519
	ledgerPrivKeyBytes, err := hexutil.Decode(privKeysCfg.Validator.PrivKey)
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

	gravityApp, err := createApp(db, ledgerValidator, privKeysCfg, ledgerConf, genesis, bootstrap, tConfig.RPC.ListenAddress, sysCtx)
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
		return err
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

func createApp(db *badger.DB, ledgerValidator *account.LedgerValidator, privKeys config.Keys, cfg config.LedgerConfig, genesisCfg config.Genesis, bootstrap string, localHost string, ctx context.Context) (*app.GHApplication, error) {

	bAdaptors := make(map[account.ChainType]adaptors.IBlockchainAdaptor)
	for k, v := range cfg.Adapters {
		cid, err := account.ChainMapper.ToByte(k)
		if err != nil {
			fmt.Printf("Chaintype '%s' error: %s", k, err)
			continue
		}
		account.ChainMapper.ApendAdaptor(cid, v.ChainType)

		chainType, err := account.ParseChainType(k)
		if err != nil {
			return nil, err
		}

		privKey, err := account.StringToPrivKey(privKeys.TargetChains[k].PrivKey, chainType)
		if err != nil {
			return nil, err
		}
		chain := byte(0)
		if len(v.ChainId) > 0 {
			chain = byte(v.ChainId[0])
		}
		opts := adaptors.AdapterOptions{
			"gravityContract": v.GravityContractAddress,
			"chainID":         chain,
		}
		adaptor, err := adaptors.NewFactory().CreateAdaptor(v.ChainType, privKey, v.NodeUrl, ctx, opts)
		if err != nil {
			return nil, err
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

	application, err := app.NewGHApplication(bAdaptors, blockScheduler, db, &genesis, ctx, &cfg)
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

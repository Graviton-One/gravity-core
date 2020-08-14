package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/Gravity-Tech/proof-of-concept/common/account"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/helpers"
	"github.com/Gravity-Tech/proof-of-concept/ledger-node/app"
	"github.com/Gravity-Tech/proof-of-concept/ledger-node/scheduler"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wavesplatform/gowaves/pkg/crypto"

	"github.com/spf13/viper"
	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/dgraph-io/badger"

	"github.com/ethereum/go-ethereum/ethclient"

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
	// read config
	config := cfg.DefaultConfig()
	config.RootDir = filepath.Dir(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper failed to read config file: %w", err)
	}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("viper failed to unmarshal config: %w", err)
	}
	if err := config.ValidateBasic(); err != nil {
		return nil, fmt.Errorf("config is invalid: %w", err)
	}

	ctx := context.Background()
	ethClient, err := ethclient.DialContext(ctx, viper.GetString("ethNodeUrl"))
	if err != nil {
		return nil, err
	}

	wavesClient, err := client.NewClient(client.Options{ApiKey: "", BaseUrl: viper.GetString("wavesNodeUrl")})
	if err != nil {
		return nil, err
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

	contracts := viper.GetStringMapString("targetContracts")
	privKeys := viper.GetStringMapString("privKeys")

	wavesPrivKey := crypto.MustBytesFromBase58(privKeys["waves"])

	ethereumPrivKey, err := hexutil.Decode(privKeys["ethereum"])
	if err != nil {
		return nil, err
	}
	wavesConf := &scheduler.WavesConf{
		PrivKey: wavesPrivKey,
		Client:  wavesClient,
		Helper:  helpers.New(wavesClient.GetOptions().BaseUrl, ""),
		ChainId: 'T',
	}

	ethPrivKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: secp256k1.S256(),
		},
		D: new(big.Int),
	}
	ethPrivKey.D.SetBytes(ethereumPrivKey)
	ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(ethereumPrivKey)

	ethConf := &scheduler.EthereumConf{
		PrivKey:      ethPrivKey,
		PrivKeyBytes: ethereumPrivKey,
		Client:       ethClient,
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
		PubKey:  ledgerPubKey,
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

	s, err := scheduler.New(wavesConf, ethConf, config.RPC.ListenAddress, context.Background(), ledger, nebulae, contracts["waves"], contracts["ethereum"])
	application := app.NewGHApplication(ethClient, wavesClient, s, db, initScoreMap, ctx)

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

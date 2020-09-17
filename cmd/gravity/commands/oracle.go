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

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/config"
	"github.com/Gravity-Tech/gravity-core/oracle/node"
	"github.com/urfave/cli/v2"
)

const (
	ConfigFlag = "config"

	DefaultNebulaeDir = "nebulae"
)

var (
	OracleCommand = &cli.Command{
		Name:        "oracle",
		Usage:       "",
		Description: "Commands to control oracles",
		Subcommands: []*cli.Command{
			{
				Name:        "init",
				Usage:       "Init oracle node config",
				Description: "",
				Action:      initOracleConfig,
				ArgsUsage:   "<nebulaId> <chainType> <gravityNodeUrl> <targetChainUrl> <extractorUrl>",
			},
			{
				Name:        "start",
				Usage:       "Start oracle node",
				Description: "",
				Action:      startOracle,
				ArgsUsage:   "<nebulaId>",
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

func initOracleConfig(ctx *cli.Context) error {
	home := ctx.String(HomeFlag)
	args := ctx.Args()

	if _, err := os.Stat(home); os.IsNotExist(err) {
		err = os.Mkdir(home, 0644)
		if err != nil {
			return err
		}
	}

	nebulaId := args.Get(0)
	chainTypeStr := args.Get(1)
	gravityUrl := args.Get(2)
	targetChainUrl := args.Get(3)
	extractorUrl := args.Get(4)

	cfg := config.OracleConfig{
		TargetChainNodeUrl: targetChainUrl,
		GravityNodeUrl:     gravityUrl,
		ChainId:            "R",
		ChainType:          chainTypeStr,
		ExtractorUrl:       extractorUrl,
	}
	b, err := json.MarshalIndent(&cfg, "", " ")
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Join(home, DefaultNebulaeDir)); os.IsNotExist(err) {
		err = os.Mkdir(path.Join(home, DefaultNebulaeDir), 0644)
		if err != nil {
			return err
		}
	}
	return ioutil.WriteFile(path.Join(home, DefaultNebulaeDir, fmt.Sprintf("%s.json", nebulaId)), b, 0644)
}

func startOracle(ctx *cli.Context) error {
	home := ctx.String(HomeFlag)
	nebulaIdStr := ctx.Args().First()

	var cfg config.OracleConfig
	err := config.ParseConfig(path.Join(home, DefaultNebulaeDir, fmt.Sprintf("%s.json", nebulaIdStr)), &cfg)
	if err != nil {
		return err
	}

	var privKeysCfg config.Keys
	err = config.ParseConfig(path.Join(home, PrivKeysConfigFileName), &privKeysCfg)
	if err != nil {
		return err
	}

	validatorPrivKey, err := hexutil.Decode(privKeysCfg.Validator.PrivKey)
	if err != nil {
		return err
	}

	chainType, err := account.ParseChainType(cfg.ChainType)
	if err != nil {
		return err
	}

	nebulaId, err := account.StringToNebulaId(nebulaIdStr, chainType)
	if err != nil {
		return err
	}

	oracleSecretKey, err := account.StringToPrivKey(privKeysCfg.TargetChains[chainType.String()].PrivKey, chainType)
	if err != nil {
		return err
	}

	sysCtx := context.Background()
	oracleNode, err := node.New(
		nebulaId,
		chainType,
		cfg.ChainId[0],
		oracleSecretKey,
		node.NewValidator(validatorPrivKey),
		cfg.ExtractorUrl,
		cfg.GravityNodeUrl,
		cfg.TargetChainNodeUrl,
		sysCtx)

	if err != nil {
		return err
	}

	err = oracleNode.Init()
	if err != nil {
		return err
	}

	go oracleNode.Start(sysCtx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	return nil
}

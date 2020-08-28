package main

import (
	"context"
	"encoding/base64"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gravity-Tech/gravity-core/common/contracts"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/Gravity-Tech/gravity-core/oracle/extractor"

	"github.com/Gravity-Tech/gravity-core/oracle/config"
	"github.com/Gravity-Tech/gravity-core/oracle/node"
	"github.com/Gravity-Tech/gravity-core/oracle/rpc"
)

const (
	DefaultConfigFileName = "config.json"
)

var confFileName string

func init() {
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	err := start(ctx)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}

func start(ctx context.Context) error {
	cfg, err := config.Load(confFileName)
	if err != nil {
		return err
	}

	validatorPrivKey, err := base64.StdEncoding.DecodeString(cfg.Secrets.ValidatorPrivKey)
	if err != nil {
		return err
	}

	chainType, err := account.ParseChainType(cfg.ChainType)
	if err != nil {
		return err
	}

	nebula, err := account.StringToNebulaId(cfg.NebulaId, chainType)
	if err != nil {
		return err
	}

	oracleSecretKey, err := account.StringToPrivKey(cfg.Secrets.OracleSecretKey, chainType)
	if err != nil {
		return err
	}

	exType, err := contracts.ParseExtractorType(cfg.Extractor.ExtractorType)
	if err != nil {
		return err
	}

	oracleNode, err := node.New(
		nebula,
		chainType,
		oracleSecretKey,
		node.NewValidator(validatorPrivKey),
		&node.Extractor{
			ExtractorType: exType,
			Client:        extractor.New(cfg.Extractor.ExtractorUrl),
		},
		cfg.GravityNodeUrl,
		cfg.TargetChainNodeUrl,
		ctx)

	if err != nil {
		return err
	}

	rpcConfig, err := rpc.NewRPCConfig(cfg.RPCHost, cfg.RPCHost, validatorPrivKey)
	if err != nil {
		return err
	}

	err = oracleNode.Init()
	if err != nil {
		return err
	}

	go oracleNode.Start(ctx)
	go rpc.ListenRpcServer(rpcConfig)

	return nil
}

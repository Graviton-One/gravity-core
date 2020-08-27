package main

import (
	"context"
	"encoding/base64"
	"flag"
	"os"
	"os/signal"
	"syscall"

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

	oracleNode, err := node.New(cfg, ctx)
	if err != nil {
		return err
	}

	privKey, err := base64.StdEncoding.DecodeString(cfg.GHPrivKey)
	if err != nil {
		return err
	}

	rpcConfig, err := rpc.NewRPCConfig(cfg.RPCHost, cfg.GHNodeURL, privKey)
	if err != nil {
		return err
	}

	go oracleNode.Start(ctx)
	go rpc.ListenRpcServer(rpcConfig)

	return nil
}

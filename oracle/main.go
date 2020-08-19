package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/Gravity-Tech/gravity-core/oracle/config"
	"github.com/Gravity-Tech/gravity-core/oracle/node"
)

const (
	DefaultConfigFileName = "config.json"
)

func main() {
	ctx := context.Background()
	oracleNode, err := newNode(ctx)
	if err != nil {
		panic(err)
	}

	go oracleNode.Start(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	os.Exit(0)
}

func newNode(ctx context.Context) (*node.Node, error) {
	var confFileName string
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.Parse()

	cfg, err := config.Load(confFileName)
	if err != nil {
		return nil, err
	}

	oracleNode, err := node.New(cfg, ctx)
	if err != nil {
		return nil, err
	}

	return oracleNode, nil
}

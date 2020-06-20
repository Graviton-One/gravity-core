package main

import (
	"context"
	"flag"
	"fmt"
	"gravity-hub/common/account"
	"gravity-hub/gh-node/api/gravity"
	"gravity-hub/gh-node/config"
	"gravity-hub/gh-node/extractors"
	"gravity-hub/gh-node/signer"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	DefaultConfigFileName = "config.json"
)

func logErr(err error) {
	if err == nil {
		return
	}

	fmt.Printf("Error: %s\n", err.Error())
}
func main() {
	var confFileName string
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.Load(confFileName)
	if err != nil {
		panic(err)
	}

	ghClient := gravity.NewClient(cfg.GHNodeURL)
	nebulaId, err := hexutil.Decode(cfg.NebulaId)
	if err != nil {
		panic(err)
	}

	chainType, err := account.ParseChainType(cfg.ChainType)
	if err != nil {
		panic(err)
	}

	client, err := signer.New(cfg.GHPrivKey, cfg.TCPrivKey, nebulaId, chainType, cfg.NebulaContract, cfg.NodeUrl, ghClient, &extractors.BinanceExtractor{}, cfg.Timeout, ctx)
	if err != nil {
		panic(err)
	}

	for {
		err = client.Start(ctx)
		if err != nil {
			fmt.Printf("Error:%s\n\n", err)
		}

		time.Sleep(time.Duration(cfg.Timeout) * time.Second)
	}
}

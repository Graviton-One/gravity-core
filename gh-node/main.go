package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"gh-node/api/gravity"
	"gh-node/config"
	"gh-node/extractors"
	"gh-node/signer"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"

	"golang.org/x/net/context"
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
	var confFileName, privKeyString string
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.StringVar(&privKeyString, "key", "", "set key")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.Load(confFileName)
	if err != nil {
		panic(err)
	}

	ghClient := gravity.NewClient(cfg.GHNodeURL)
	nebulaId, err := hex.DecodeString(cfg.NebulaId)
	if err != nil {
		panic(err)
	}

	privBytes, err := hexutil.Decode(privKeyString)
	if err != nil {
		panic(err)
	}

	ethClient, err := ethclient.DialContext(ctx, cfg.EthNodeUrl)
	if err != nil {
		panic(err)
	}

	client, err := signer.New(privBytes, nebulaId, cfg.NebulaContract, ghClient, ethClient, &extractors.BinanceExtractor{}, cfg.Timeout, ctx)
	if err != nil {
		panic(err)
	}

	for {
		err = client.Start(ctx)
		if err != nil {
			fmt.Printf("Error:%s\n", err.Error())
		}
	}
}

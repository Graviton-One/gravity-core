package main

import (
	"context"
	"flag"
	"fmt"
	"proof-of-concept/gh-node/config"
	"proof-of-concept/gh-node/signer"
	"time"
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

	client, err := signer.New(cfg, ctx)

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

package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/Gravity-Tech/gravity-core/oracle-node/config"
	"github.com/Gravity-Tech/gravity-core/oracle-node/node"
)

const (
	DefaultConfigFileName = "config.json"
)

func main() {
	var confFileName string
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.Parse()

	ctx := context.Background()
	cfg, err := config.Load(confFileName)
	if err != nil {
		panic(err)
	}

	client, err := node.New(cfg, ctx)
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

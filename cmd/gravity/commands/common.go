package commands

import (
	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
)

const (
	HomeFlag     = "home"
	LogLevelFlag = "loglevel"

	DbDir                  = "db"
	PrivKeysConfigFileName = "privKey.json"
	GenesisFileName        = "genesis.json"
	LedgerConfigFileName   = "config.json"
	NodeKeyFileName        = "node_key.json"
	LedgerKeyStateFileName = "key_state.json"
)

func InitLogger(ctx *cli.Context) (*zap.Logger, error) {
	level := ctx.String(LogLevelFlag)
	logger := &zap.Logger{}
	var err error
	switch level {
	case "development":
		logger, err = zap.NewDevelopment()
	case "production":
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, err
	}

	return logger, nil
}

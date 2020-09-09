package config

import (
	"github.com/Gravity-Tech/gravity-core/common/account"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	DefaultMoniker = "robot"
)

type AdaptorsConfig struct {
	NodeUrl                string
	ChainId                string
	GravityContractAddress string
}

type LedgerConfig struct {
	Moniker    string
	IsFastSync bool
	Mempool    *cfg.MempoolConfig
	RPC        *cfg.RPCConfig
	P2P        *cfg.P2PConfig

	Adapters map[string]AdaptorsConfig
}

func DefaultLedgerConfig() LedgerConfig {
	return LedgerConfig{
		Moniker:    DefaultMoniker,
		IsFastSync: true,
		Mempool:    cfg.DefaultMempoolConfig(),
		RPC:        cfg.DefaultRPCConfig(),
		P2P:        cfg.DefaultP2PConfig(),
		Adapters: map[string]AdaptorsConfig{
			account.Ethereum.String(): {
				NodeUrl:                "",
				GravityContractAddress: "",
			},
			account.Waves.String(): {
				NodeUrl:                "",
				GravityContractAddress: "",
			},
		},
	}
}

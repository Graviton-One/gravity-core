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

type ValidatorDetails struct {
	Name, Description, JoinedAt string
	// Misc
	AvatarURL, Website string
}

func (validatorDetails *ValidatorDetails) DefaultNew() *ValidatorDetails {
	return &ValidatorDetails{
		Name: "Gravity Node", Description: "", JoinedAt: "",
		AvatarURL: "", Website: "",
	}
}

type LedgerConfig struct {
	Moniker    string
	IsFastSync bool
	Mempool    *cfg.MempoolConfig
	RPC        *cfg.RPCConfig
	P2P        *cfg.P2PConfig

	Details    *ValidatorDetails

	Adapters map[string]AdaptorsConfig
}

func DefaultLedgerConfig() LedgerConfig {
	return LedgerConfig{
		Moniker:    DefaultMoniker,
		IsFastSync: true,
		Mempool:    cfg.DefaultMempoolConfig(),
		RPC:        cfg.DefaultRPCConfig(),
		P2P:        cfg.DefaultP2PConfig(),
		Details:    ValidatorDetails{}.DefaultNew(),
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

package config

import (
	"encoding/json"

	"github.com/Gravity-Tech/gravity-core/common/account"
	cfg "github.com/tendermint/tendermint/config"
)

const (
	DefaultMoniker = "robot"
)

type AdaptorsConfig struct {
	NodeUrl                string
	ChainId                string
	ChainType              string
	GravityContractAddress string
	Custom                 map[string]interface{} `json:"custom,optional"`
}

type ValidatorDetails struct {
	Name, Description, JoinedAt string
	// Misc
	AvatarURL, Website string
}

func (validatorDetails *ValidatorDetails) Bytes() ([]byte, error) {
	res, err := json.Marshal(validatorDetails)

	if err != nil {
		return make([]byte, 0), err
	}

	return res, nil
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

	Details  *ValidatorDetails
	PublicIP string

	Adapters map[string]AdaptorsConfig
}

func DefaultLedgerConfig() LedgerConfig {
	return LedgerConfig{
		Moniker:    DefaultMoniker,
		IsFastSync: true,
		Mempool:    cfg.DefaultMempoolConfig(),
		RPC:        cfg.DefaultRPCConfig(),
		P2P:        cfg.DefaultP2PConfig(),
		Details:    (&ValidatorDetails{}).DefaultNew(),
		Adapters: map[string]AdaptorsConfig{
			account.Ethereum.String(): {
				NodeUrl:                "",
				GravityContractAddress: "",
			},
			account.Waves.String(): {
				NodeUrl:                "",
				GravityContractAddress: "",
			},
			account.Ergo.String(): {
				NodeUrl:                "",
				GravityContractAddress: "",
			},
		},
	}
}

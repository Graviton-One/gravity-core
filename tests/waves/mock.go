package tests

import (
	"fmt"
	"github.com/wavesplatform/gowaves/pkg/proto"
	"strings"
)

type NebulaTestMockConfig struct {
	chainId           byte

	BftCoefficient    int64          `dtx:"bft_coefficient"`
	GravityAddress    string         `dtx:"gravity_contract"`
	NebulaPubkey      string         `dtx:"contract_pubkey"`
	SubscriberAddress string         `dtx:"subscriber_address"`
	OraclesList       [5]WavesActor  `dtx:"oracles"`
}

func (nebulaMockCfg *NebulaTestMockConfig) Validate() error {
	if nebulaMockCfg.BftCoefficient < 1 {
		return fmt.Errorf("bft_coefficient is less than 1")
	}

	if nebulaMockCfg.chainId != 'S' {
		return fmt.Errorf("only waves stagenet is supported")
	}

	if nebulaMockCfg.GravityAddress == "" {
		return fmt.Errorf("field \"gravity_contract\" cannot be empty")
	}
	if nebulaMockCfg.NebulaPubkey == "" {
		return fmt.Errorf("field \"contract_pubkey\" cannot be empty")
	}
	if nebulaMockCfg.SubscriberAddress == "" {
		return fmt.Errorf("field \"subscriber_address\" cannot be empty")
	}

	return nil
}

func (nebulaMockCfg *NebulaTestMockConfig) OraclesPubKeysListDataEntry () string {
	var res []string

	for _, oracle := range nebulaMockCfg.OraclesList {
		res = append(res, oracle.Account(nebulaMockCfg.chainId).PubKey.String())
	}

	return strings.Join(res, ",")
}

func (nebulaMockCfg *NebulaTestMockConfig) DataEntries () proto.DataEntries {
	return proto.DataEntries{
		&proto.IntegerDataEntry{
			Key:   "bft_coefficient",
			Value: nebulaMockCfg.BftCoefficient,
		},
		&proto.StringDataEntry{
			Key:   "gravity_contract",
			Value: nebulaMockCfg.GravityAddress,
		},
		&proto.StringDataEntry{
			Key:   "contract_pubkey",
			Value: nebulaMockCfg.NebulaPubkey,
		},
		&proto.StringDataEntry{
			Key:   "subscriber_address",
			Value: nebulaMockCfg.SubscriberAddress,
		},
		&proto.StringDataEntry{
			Key:   "oracles",
			Value: nebulaMockCfg.OraclesPubKeysListDataEntry(),
		},
	}
}

func (nebulaMockCfg *NebulaTestMockConfig) OraclesPubKeysList () []string {
	var oraclesPubKeyList []string

	for _, mockedConsul := range nebulaMockCfg.OraclesList {
		pk := mockedConsul.Account(cfg.Environment.ChainIDBytes()).PubKey.String()
		oraclesPubKeyList = append(oraclesPubKeyList, pk)
	}

	return oraclesPubKeyList
}

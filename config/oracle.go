package config

type OracleConfig struct {
	TargetChainNodeUrl string
	ChainId            string
	GravityNodeUrl     string
	ChainType          string
	ExtractorUrl       string
	BlocksInterval     uint64
	Custom             map[string]interface{}
}

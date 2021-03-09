package config

type OracleConfig struct {
	TargetChainNodeUrl string
	ChainId            string
	GravityNodeUrl     string
	ChainType          string
	ChainName          string
	ExtractorUrl       string
	BlocksInterval     uint64
}

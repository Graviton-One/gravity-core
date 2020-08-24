package deployer

import (
	"flag"
)

const (
	DefaultConfigFileName = "config.json"
)
const (
	Ethereum Chain = "ethereum"
	Waves    Chain = "waves"

	Gravity ContractType = "gravity"
	Nebula  ContractType = "nebula"
)

type Chain string
type ContractType string

func main() {
	var confFileName, chain, contractType string
	flag.StringVar(&confFileName, "config", DefaultConfigFileName, "set config path")
	flag.StringVar(&chain, "chain", "", "set contract chain type")
	flag.StringVar(&contractType, "type", "", "set contract type")
	flag.StringVar(&confFileName, "args", "", "set args")
	flag.Parse()

	switch Chain(chain) {
	case Ethereum:
		switch ContractType(contractType) {
		case Gravity:
		case Nebula:
		}
	case Waves:
		switch ContractType(contractType) {
		case Gravity:
		case Nebula:
		}
	}
}

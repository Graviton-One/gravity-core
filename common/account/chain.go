package account

import (
	"errors"
	"strings"
)

type ChainType byte

const (
	Ethereum ChainType = iota
	Waves
	Binance
	Heco
	Fantom
	Avax
	Ergo
)

var (
	ErrInvalidChainType = errors.New("invalid chain type")
	ErrParseChainType   = errors.New("invalid parse chain type")
)

func ParseChainType(chainType string) (ChainType, error) {
	switch strings.ToLower(chainType) {
	case "heco":
		return Heco, nil
	case "bsc":
		return Binance, nil
	case "ethereum":
		return Ethereum, nil
	case "ftm":
		return Fantom, nil
	case "avax":
		return Avax, nil
	case "waves":
		return Waves, nil
	default:
		return 0, ErrParseChainType
	}
}
func (ch ChainType) String() string {
	switch ch {
	case Ethereum:
		return "ethereum"
	case Waves:
		return "waves"
	case Binance:
		return "bsc"
	case Heco:
		return "heco"
	case Fantom:
		return "ftm"
	case Avax:
		return "avax"
	case Ergo:
		return "erg"
	default:
		return "ethereum"
	}
}

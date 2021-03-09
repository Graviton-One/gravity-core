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
)

var (
	ErrInvalidChainType = errors.New("invalid chain type")
	ErrParseChainType   = errors.New("invalid parse chain type")
)

func ParseChainType(chainType string) (ChainType, error) {
	switch strings.ToLower(chainType) {
	case "bsc":
		return Binance, nil
	case "ethereum":
		return Ethereum, nil
	case "waves":
		return Waves, nil
	default:
		// i, err := strconv.ParseInt(chainType, 10, 8) //Is third parameter not honored??
		// if err == nil {
		// 	return ChainType(i), nil
		// }
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
	default:
		return "ethereum"
		//return fmt.Sprintf("%d", ch)
	}
}

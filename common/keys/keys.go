package keys

import (
	"fmt"
	"strings"

	"github.com/Gravity-Tech/proof-of-concept/common/account"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Key string

const (
	Separator string = "_"

	ConsulsKey           Key = "consuls"
	PrevConsulsKey       Key = "prev_consuls"
	ConsulsSignKey       Key = "consuls_sing"
	OraclesSignNebulaKey Key = "oracles_sign"

	OraclesByNebulaKey    Key = "oracles_nebula"
	BftOraclesByNebulaKey Key = "bft_oracles_nebula"
	OraclesByValidatorKey Key = "oracles"
	ValidatorByOracleKey  Key = "validator"

	BlockKey      Key = "block"
	VoteKey       Key = "vote"
	ScoreKey      Key = "score"
	CommitKey     Key = "commit"
	RevealKey     Key = "reveal"
	SignResultKey Key = "signResult"
)

func FormPrevConsulsKey() string {
	return string(PrevConsulsKey)
}

func FormConsulsSignKey(validatorAddress []byte, chainType account.ChainType, roundId int64) string {
	prefix := ""
	switch chainType {
	case account.Waves:
		prefix = "waves"
	case account.Ethereum:
		prefix = "ethereum"
	}
	return strings.Join([]string{string(ConsulsSignKey), hexutil.Encode(validatorAddress), prefix, fmt.Sprintf("%d", roundId)}, Separator)
}

func FormOraclesSignNebulaKey(validatorAddress []byte, nebulaId []byte, roundId int64) string {
	return strings.Join([]string{string(OraclesSignNebulaKey), hexutil.Encode(validatorAddress), hexutil.Encode(nebulaId), fmt.Sprintf("%d", roundId)}, Separator)
}

func FormBftOraclesByNebulaKey(nebulaId []byte) string {
	return strings.Join([]string{string(BftOraclesByNebulaKey), hexutil.Encode(nebulaId)}, Separator)
}

func formKey(args []string) []byte {
	return []byte(strings.Join(args, Separator))
}

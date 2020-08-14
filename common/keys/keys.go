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

func FormConsulsKey() string {
	return string(ConsulsKey)
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

func FormOraclesByNebulaKey(nebulaId []byte) string {
	return strings.Join([]string{string(OraclesByNebulaKey), hexutil.Encode(nebulaId)}, Separator)
}

func FormBftOraclesByNebulaKey(nebulaId []byte) string {
	return strings.Join([]string{string(BftOraclesByNebulaKey), hexutil.Encode(nebulaId)}, Separator)
}

func FormOraclesByValidatorKey(validatorAddress []byte) string {
	return strings.Join([]string{string(OraclesByValidatorKey), hexutil.Encode(validatorAddress)}, Separator)
}

func FormValidatorByOracleKey(oraclePubKey []byte) string {
	return strings.Join([]string{string(ValidatorByOracleKey), hexutil.Encode(oraclePubKey)}, Separator)
}

func FormBlockKey(chainType account.ChainType, block uint64) string {
	var blockchain string
	switch chainType {
	case account.Ethereum:
		blockchain = "ethereum"
	case account.Waves:
		blockchain = "waves"
	}
	return strings.Join([]string{string(BlockKey), blockchain, fmt.Sprintf("%d", block)}, Separator)
}

func FormVoteKey(oraclePubKey []byte) string {
	return strings.Join([]string{string(VoteKey), hexutil.Encode(oraclePubKey)}, Separator)
}

func FormScoreKey(validatorAddress []byte) string {
	return strings.Join([]string{string(ScoreKey), hexutil.Encode(validatorAddress)}, Separator)
}

func FormCommitKey(nebulaAddress []byte, block uint64, oraclePubKey []byte) string {
	return strings.Join([]string{string(CommitKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(oraclePubKey)}, Separator)
}

func FormRevealKey(nebulaAddress []byte, block uint64, commitHash []byte) string {
	return strings.Join([]string{string(RevealKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(commitHash)}, Separator)
}

func FormSignResultKey(nebulaAddress []byte, block uint64, oraclePubKey []byte) string {
	return strings.Join([]string{string(SignResultKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(oraclePubKey)}, Separator)
}

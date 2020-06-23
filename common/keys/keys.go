package keys

import (
	"fmt"
	"gravity-hub/common/account"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Key string

const (
	Separator string = "_"

	NebulaeByValidatorKey   Key = "nebulaByValidators"
	ValidatorKey            Key = "validator"
	VoteKey                 Key = "vote"
	SignNebulaValidatorsKey Key = "signNebulaValidators"
	ScoreKey                Key = "score"
	CommitKey               Key = "commit"
	RevealKey               Key = "reveal"
	SignResultKey           Key = "signResult"
	BlockKey                Key = "block"
)

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

func FormScoreKey(validatorAddress []byte) string {
	return strings.Join([]string{string(ScoreKey), hexutil.Encode(validatorAddress)}, Separator)
}

func FormVoteKey(validatorAddress []byte) string {
	return strings.Join([]string{string(VoteKey), hexutil.Encode(validatorAddress)}, Separator)
}

func FormNebulaeByValidatorKey(validatorAddress []byte) string {
	return strings.Join([]string{string(NebulaeByValidatorKey), hexutil.Encode(validatorAddress)}, Separator)
}

func FormValidatorKey(nebulaAddress []byte, validatorAddress []byte) string {
	return strings.Join([]string{string(ValidatorKey), hexutil.Encode(nebulaAddress), hexutil.Encode(validatorAddress)}, Separator)
}

func FormCommitKey(nebulaAddress []byte, block uint64, validatorAddress []byte) string {
	return strings.Join([]string{string(CommitKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(validatorAddress)}, Separator)
}

func FormRevealKey(nebulaAddress []byte, block uint64, commitHash []byte) string {
	return strings.Join([]string{string(RevealKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(commitHash)}, Separator)
}

func FormSignResultKey(nebulaAddress []byte, block uint64, validatorAddress []byte) string {
	return strings.Join([]string{string(SignResultKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(validatorAddress)}, Separator)
}

func FormSignScoreValidatorsKey(validatorAddress []byte) string {
	return strings.Join([]string{string(SignNebulaValidatorsKey), hexutil.Encode(validatorAddress)}, Separator)
}

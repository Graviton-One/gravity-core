package keys

import (
	"fmt"
	"gravity-hub/common/account"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Key string

const (
	ValidatorKey  Key = "validator"
	VoteKey       Key = "vote"
	CommitKey     Key = "commit"
	RevealKey     Key = "reveal"
	SignResultKey Key = "signResult"
	BlockKey      Key = "block"
)

func FormBlockKey(chainType account.ChainType, block uint64) string {
	var blockchain string
	switch chainType {
	case account.Ethereum:
		blockchain = "ethereum"
	case account.Waves:
		blockchain = "waves"
	}
	return strings.Join([]string{string(BlockKey), blockchain, fmt.Sprintf("%d", block)}, "_")
}

func FormVoteKey(validatorAddress []byte) string {
	return strings.Join([]string{string(VoteKey), hexutil.Encode(validatorAddress)}, "_")
}

func FormValidatorKey(nebulaAddress []byte, validatorAddress []byte) string {
	return strings.Join([]string{string(ValidatorKey), hexutil.Encode(nebulaAddress), hexutil.Encode(validatorAddress)}, "_")
}

func FormCommitKey(nebulaAddress []byte, block uint64, validatorAddress []byte) string {
	return strings.Join([]string{string(CommitKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(validatorAddress)}, "_")
}

func FormRevealKey(nebulaAddress []byte, block uint64, commitHash []byte) string {
	return strings.Join([]string{string(RevealKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(commitHash)}, "_")
}

func FormSignResultKey(nebulaAddress []byte, block uint64, validatorAddress []byte) string {
	return strings.Join([]string{string(SignResultKey), hexutil.Encode(nebulaAddress), fmt.Sprintf("%d", block), hexutil.Encode(validatorAddress)}, "_")
}

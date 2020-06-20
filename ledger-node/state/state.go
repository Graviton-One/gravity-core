package state

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Key string

const (
	ValidatorKey  Key = "validator"
	CommitKey     Key = "commit"
	RevealKey     Key = "reveal"
	SignResultKey Key = "signResult"
)

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

package keys

import (
	"encoding/hex"
	"fmt"
	"gravity-hub/common/account"
	"strings"
)

type Key string

const (
	ValidatorKey  Key = "validator"
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

func FormValidatorKey(nebulaAddress []byte, validatorAddress []byte) string {
	return strings.Join([]string{string(ValidatorKey), hex.EncodeToString(nebulaAddress), hex.EncodeToString(validatorAddress)}, "_")
}

func FormCommitKey(nebulaAddress []byte, block uint64, validatorAddress []byte) string {
	return strings.Join([]string{string(CommitKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", block), hex.EncodeToString(validatorAddress)}, "_")
}

func FormRevealKey(nebulaAddress []byte, block uint64, commitHash []byte) string {
	return strings.Join([]string{string(RevealKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", block), hex.EncodeToString(commitHash)}, "_")
}

func FormSignResultKey(nebulaAddress []byte, block uint64, validatorAddress []byte) string {
	return strings.Join([]string{string(SignResultKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", block), hex.EncodeToString(validatorAddress)}, "_")
}

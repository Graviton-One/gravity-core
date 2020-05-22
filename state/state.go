package state

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type Key string

const (
	ValidatorKey  Key = "validator"
	CommitKey     Key = "commit"
	RevealKey     Key = "reveal"
	SignResultKey Key = "signResult"
	ResultKey     Key = "result"
)

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

func FormResultKey(nebulaAddress []byte, block uint64) string {
	return strings.Join([]string{string(ResultKey), hex.EncodeToString(nebulaAddress), fmt.Sprintf("%d", block)}, "_")
}

package hashing

import (
	"crypto/sha256"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/ethereum/go-ethereum/crypto"
)


func WrappedKeccak256(input []byte, chain account.ChainType) []byte {
	if chain != account.Solana {
		return crypto.Keccak256(input[:])
	}
	
	digest := sha256.Sum256(input[:])
	return digest[:]
}
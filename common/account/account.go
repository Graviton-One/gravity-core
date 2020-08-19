package account

import (
	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"

	_ "github.com/tendermint/tendermint/crypto/ed25519"
)

func Sign(privKey tendermintCrypto.PrivKeyEd25519, msg []byte) ([]byte, error) {
	return privKey.Sign(msg)
}

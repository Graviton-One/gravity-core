package account

import (
	"github.com/tendermint/tendermint/crypto"
)

func Sign(privKey crypto.PrivKey, msg []byte) ([]byte, error) {
	return privKey.Sign(msg)
}

type LedgerValidator struct {
	PrivKey crypto.PrivKey
	PubKey  ConsulPubKey
}

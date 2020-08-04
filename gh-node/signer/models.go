package signer

import "github.com/Gravity-Tech/gravity-core/common/account"

type TCAccount struct {
	privKey []byte
	pubKey  account.OraclesPubKey
}

type RoundState struct {
	data        []byte
	commitHash  []byte
	resultValue []byte
	resultHash  []byte
	isSent      bool
	isSentSub   bool
}

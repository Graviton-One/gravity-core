package node

import "github.com/Gravity-Tech/gravity-core/common/account"

type TCAccount struct {
	privKey []byte
	pubKey  account.OraclesPubKey
}

type RoundState struct {
	data        interface{}
	commitHash  []byte
	resultValue interface{}
	resultHash  []byte
	isSent      bool
	isSentSub   bool
}

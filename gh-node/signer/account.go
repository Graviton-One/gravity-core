package signer

type TCAccount struct {
	privKey []byte
	pubKey  []byte
}

type State struct {
	commitPrice uint64
	commitHash  []byte
	resultValue uint64
	resultHash  []byte
	isSent      bool
	isSentSub   bool
}

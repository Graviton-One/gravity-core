package node

type RoundState struct {
	data        interface{}
	commitHash  []byte
	resultValue interface{}
	resultHash  []byte
	isSent      bool
	RevealExist bool
}

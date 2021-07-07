package node

import "github.com/Gravity-Tech/gravity-core/oracle/extractor"

type RoundState struct {
	data        *extractor.Data
	commitHash  []byte
	resultValue *extractor.Data
	resultHash  []byte
	isSent      bool
	RevealExist bool
	commitSent  bool
}

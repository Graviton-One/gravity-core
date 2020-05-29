package blockchain

import (
	"context"
	"gravity-hub/gh-node/api/gravity"
)

type IBlockchain interface {
	GetHeight(ctx context.Context) (uint64, error)
	SendResult(tcHeight uint64, privKey []byte, nebulaId []byte, ghClient *gravity.Client, validators [][]byte, hash []byte, ctx context.Context) error
	SendSubs(tcHeight uint64, privKey []byte, value uint64) error
}

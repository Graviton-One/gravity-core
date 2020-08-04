package blockchain

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/gh-node/api/gravity"
)

type IBlockchain interface {
	GetHeight(ctx context.Context) (uint64, error)
	SendResult(tcHeight uint64, privKey []byte, nebulaId []byte, ghClient *gravity.Client, validators [][]byte, hash []byte, ctx context.Context) (string, error)
	SendSubs(tcHeight uint64, privKey []byte, value uint64, ctx context.Context) error
	WaitTx(id string) error
}

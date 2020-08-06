package blockchain

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/common/client"
)

type IBlockchain interface {
	GetHeight(ctx context.Context) (uint64, error)
	SendResult(tcHeight uint64, privKey []byte, nebulaId []byte, ghClient *client.Client, validators [][]byte, hash []byte, ctx context.Context) (string, error)
	SendSubs(tcHeight uint64, privKey []byte, value uint64, ctx context.Context) error
	WaitTx(id string) error
}

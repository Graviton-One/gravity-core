package blockchain

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type IBlockchainClient interface {
	GetHeight(ctx context.Context) (uint64, error)
	SendResult(tcHeight uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error)
	SendSubs(tcHeight uint64, value interface{}, ctx context.Context) error
	WaitTx(id string, ctx context.Context) error
	Sign(msg []byte) ([]byte, error)
}

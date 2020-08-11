package blockchain

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/Gravity-Tech/gravity-core/common/client"
)

type IBlockchain interface {
	GetHeight(ctx context.Context) (uint64, error)
	SendResult(ghClient *client.Client, tcHeight uint64, privKey []byte, nebulaId account.NebulaId, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error)
	SendSubs(tcHeight uint64, privKey []byte, value interface{}, ctx context.Context) error
	WaitTx(id string) error
}

package blockchain

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type IBlockchain interface {
	SendOraclesToNebula(nebulaId account.NebulaId, oracles []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error)
	SendConsulsToGravityContract(newConsulsAddresses []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error)
	SignConsuls(consulsAddresses []account.OraclesPubKey) ([]byte, error)
	SignOracles(nebulaId account.NebulaId, oracles []account.OraclesPubKey) ([]byte, error)
	PubKey() []byte
}

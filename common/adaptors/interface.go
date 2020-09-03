package adaptors

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/common/contracts"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type IBlockchainAdaptor interface {
	GetHeight(ctx context.Context) (uint64, error)
	WaitTx(id string, ctx context.Context) error
	Sign(msg []byte) ([]byte, error)
	PubKey() account.OraclesPubKey
	GetExtractorType(nebulaId account.NebulaId, ctx context.Context) (contracts.ExtractorType, error)

	SendDataResult(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error)
	SendDataToSubs(nebulaId account.NebulaId, pulseId uint64, value interface{}, ctx context.Context) error

	SendOraclesToNebula(nebulaId account.NebulaId, oracles []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error)
	SendConsulsToGravityContract(newConsulsAddresses []account.OraclesPubKey, signs [][]byte, round int64, ctx context.Context) (string, error)
	SignConsuls(consulsAddresses []account.OraclesPubKey) ([]byte, error)
	SignOracles(nebulaId account.NebulaId, oracles []account.OraclesPubKey) ([]byte, error)

	LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error)
}

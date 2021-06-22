package adaptors

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"

	"github.com/Gravity-Tech/gravity-core/common/account"
)

type IBlockchainAdaptor interface {
	GetHeight(ctx context.Context) (uint64, error)
	WaitTx(id string, ctx context.Context) error
	Sign(msg []byte) ([]byte, error)
	PubKey() account.OraclesPubKey
	ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error)

	AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error)
	SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error

	SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error)
	SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error)
	SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64, sender account.OraclesPubKey) ([]byte, error)
	SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, round int64, sender account.OraclesPubKey) ([]byte, error)
	SignHash(nebulaId account.NebulaId, intervalId uint64, pulseId uint64, hash []byte) ([]byte, error)
	LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error)
	LastRound(ctx context.Context) (uint64, error)
	RoundExist(roundId int64, ctx context.Context) (bool, error)
}

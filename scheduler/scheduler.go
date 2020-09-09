package rpc

import (
	"context"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/adaptors"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
)

type Scheduler struct {
	Host     string
	Adaptors map[account.ChainType]adaptors.IBlockchainAdaptor
	Ledger   *account.LedgerValidator
	ctx      context.Context
	client   *gravity.Client
}
type ConsulInfo struct {
	ConsulIndex int
	TotalCount  int
	IsConsul    bool
}

func New(host string, adaptors map[account.ChainType]adaptors.IBlockchainAdaptor, ledger *account.LedgerValidator, localHost string, ctx context.Context) (*Scheduler, error) {
	client, err := gravity.New(localHost)
	if err != nil {
		return nil, err
	}
	return &Scheduler{
		Host:     host,
		Ledger:   ledger,
		Adaptors: adaptors,
		ctx:      ctx,
		client:   client,
	}, nil
}

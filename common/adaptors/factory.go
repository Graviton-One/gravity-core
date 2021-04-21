package adaptors

import (
	"context"
	"fmt"
	"github.com/Gravity-Tech/gravity-core/common/helpers"

	"github.com/Gravity-Tech/gravity-core/abi/ethereum"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gookit/validate"
	wclient "github.com/wavesplatform/gowaves/pkg/client"
)

//AdapterOptions - map of custom adaptor creating options
type AdapterOptions map[string]interface{}

//Factory - abstract factory struct
type Factory struct {
}

func isGhClientValidator(val interface{}) bool {
	switch val.(type) {
	case *gravity.Client:
		return true
	default:
		return false
	}
}
func isByteValidator(val interface{}) bool {
	switch val.(type) {
	case byte:
		return true
	default:
		return false
	}
}
func isWvClientValidator(val interface{}) bool {
	switch val.(type) {
	case *wclient.Client:
		return true
	default:
		return false
	}
}
func isErgClientValidator(val interface{}) bool {
	switch val.(type) {
	case *helpers.ErgClient:
		return true
	default:
		return false
	}
}
func isEthClientValidator(val interface{}) bool {
	switch val.(type) {
	case *ethclient.Client:
		return true
	default:
		return false
	}
}

func isEthGravityContractValidator(val interface{}) bool {
	switch val.(type) {
	case *ethereum.Gravity:
		return true
	default:
		return false
	}
}

func NewFactory() *Factory {
	validate.AddValidator("isGhClient", isGhClientValidator)
	validate.AddValidator("isByte", isByteValidator)
	validate.AddValidator("isWvClient", isWvClientValidator)
	validate.AddValidator("isErgClient", isErgClientValidator)
	validate.AddValidator("isEthClient", isEthClientValidator)
	validate.AddValidator("isEthGravityContract", isEthGravityContractValidator)

	return &Factory{}
}

//CreateAdaptor - factory function
func (f *Factory) CreateAdaptor(name string, oracleSecretKey []byte, targetChainNodeUrl string, ctx context.Context, opts AdapterOptions) (IBlockchainAdaptor, error) {
	switch name {
	case "waves":
		return NewWavesAdapterByOpts(oracleSecretKey, targetChainNodeUrl, opts)
	case "ethereum":
		return NewEthereumsAdapterByOpts(oracleSecretKey, targetChainNodeUrl, ctx, opts)
	case "ergo":
		return NewErgoAdapterByOpts(oracleSecretKey, targetChainNodeUrl, opts)
	}
	return nil, fmt.Errorf("Unknown adaptor name %s", name)
}

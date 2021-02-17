package adaptors

import (
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/gookit/validate"
	wclient "github.com/wavesplatform/gowaves/pkg/client"
)

//Options - map of custom adaptor creating options
type AdapterOptions map[string]interface{}

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

func NewFactory() *Factory {
	validate.AddValidator("isGhClient", isGhClientValidator)
	validate.AddValidator("isByte", isByteValidator)
	validate.AddValidator("isWvClient", isWvClientValidator)
	return &Factory{}
}

//NewAdaptor - factory function
func (f *Factory) CreateAdaptor(name string, oracleSecretKey []byte, targetChainNodeUrl string, opts AdapterOptions) (IBlockchainAdaptor, error) {
	switch name {
	case "waves":
		return NewWavesAdapterByOpts(oracleSecretKey, targetChainNodeUrl, opts)
	}
	return nil, fmt.Errorf("Unknown adaptor name %s", name)
}

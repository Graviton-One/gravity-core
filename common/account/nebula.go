package account

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type NebulaId []byte

func StringToNebulaId(address string, chainType ChainType) (NebulaId, error) {
	var nebula NebulaId
	var err error
	switch chainType {
	case Ethereum:
		nebula, err = hexutil.Decode(address)
		if err != nil {
			return nil, err
		}
	case Waves:
		nebula = crypto.MustBytesFromBase58(address)
	}

	return nebula, nil
}

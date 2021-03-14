package account

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wavesplatform/gowaves/pkg/crypto"
	"go.uber.org/zap"
)

const (
	NebulaIdLength        = 32
	EthereumAddressLength = 20
	BSCAddressLength      = 20
	WavesAddressLength    = 26
)

type NebulaId [NebulaIdLength]byte

func StringToNebulaId(address string, chainType ChainType) (NebulaId, error) {
	// cid, err := ChainMapper.ToType(byte(chainType))
	// if err != nil {
	// 	return NebulaId{}, err
	// }
	// ctype := ChainType(cid)
	var nebula NebulaId

	switch chainType {
	case Ethereum, Binance:
		nebulaBytes, err := hexutil.Decode(address)
		if err != nil {
			return NebulaId{}, err
		}
		nebula = BytesToNebulaId(nebulaBytes)
	case Waves:
		nebulaBytes := crypto.MustBytesFromBase58(address)
		nebula = BytesToNebulaId(nebulaBytes)
	}

	return nebula, nil
}
func BytesToNebulaId(value []byte) NebulaId {
	var idBytes []byte
	var id NebulaId
	if len(value) < NebulaIdLength {
		idBytes = append(idBytes, make([]byte, NebulaIdLength-len(value), NebulaIdLength-len(value))...)
	}
	idBytes = append(idBytes, value...)
	copy(id[:], idBytes)

	return id
}

func (id NebulaId) ToString(chainType ChainType) string {
	zap.L().Sugar().Debugf("chainType: %d", chainType)
	cid, err := ChainMapper.ToType(byte(chainType))
	if err != nil {
		zap.L().Error(err.Error())
		return ""
	}
	ctype := ChainType(cid)

	nebula := id.ToBytes(ctype)
	switch ctype {
	case Ethereum, Binance:
		return hexutil.Encode(nebula[:])
	case Waves:
		return base58.Encode(nebula[:])
	}

	return ""
}
func (id NebulaId) ToBytes(chainType ChainType) []byte {
	// cid, err := ChainMapper.ToType(byte(chainType))
	// if err != nil {
	// 	zap.L().Error(err.Error())
	// 	return nil
	// }
	// ctype := ChainType(cid)

	switch chainType {
	case Binance:
		return id[NebulaIdLength-BSCAddressLength:]
	case Ethereum:
		return id[NebulaIdLength-EthereumAddressLength:]
	case Waves:
		return id[NebulaIdLength-WavesAddressLength:]
	}

	return nil
}

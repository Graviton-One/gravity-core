package account

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

const (
	NebulaIdLength        = 32
	EthereumAddressLength = 20
	BSCAddressLength      = 20
	WavesAddressLength    = 26
)

type NebulaId [NebulaIdLength]byte

func StringToNebulaId(address string, chainType ChainType) (NebulaId, error) {
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
	nebula := id.ToBytes(chainType)
	switch chainType {
	case Ethereum, Binance:
		return hexutil.Encode(nebula[:])
	case Waves:
		return base58.Encode(nebula[:])
	}

	return ""
}
func (id NebulaId) ToBytes(chainType ChainType) []byte {
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

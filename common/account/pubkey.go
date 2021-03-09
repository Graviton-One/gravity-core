package account

import (
	"fmt"

	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/ed25519"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type ConsulPubKey ed25519.PubKeyEd25519
type OraclesPubKey [33]byte

func StringToPrivKey(value string, chain ChainType) ([]byte, error) {
	var privKey []byte
	var err error
	cType, err := ChainMapper.ToType(byte(chain))
	if err != nil {
		return nil, err
	}
	switch ChainType(cType) {
	case Ethereum, Binance:
		privKey, err = hexutil.Decode(value)
		if err != nil {
			return nil, err
		}
	case Waves:
		wCrypto := wavesplatform.NewWavesCrypto()
		seed := wavesplatform.Seed(value)
		secret, err := crypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
		if err != nil {
			return nil, err
		}
		privKey = secret.Bytes()
	}

	return privKey, nil
}

func BytesToOraclePubKey(value []byte, chain ChainType) OraclesPubKey {
	var pubKey OraclesPubKey
	cType, err := ChainMapper.ToType(byte(chain))
	if err != nil {
		return pubKey
	}

	switch ChainType(cType) {
	case Ethereum, Binance:
		copy(pubKey[:], value[0:33])
	case Waves:
		copy(pubKey[:], append([]byte{0}, value[0:32]...))
	}
	return pubKey
}

func (pubKey *OraclesPubKey) ToBytes(chain ChainType) []byte {
	var v []byte
	cType, err := ChainMapper.ToType(byte(chain))
	if err != nil {
		return v
	}
	switch ChainType(cType) {
	case Ethereum, Binance:
		v = pubKey[:33]
	case Waves:
		v = pubKey[1:33]
	}
	return v
}
func (pubKey *OraclesPubKey) ToString(chain ChainType) string {
	b := pubKey.ToBytes(chain)
	cType, err := ChainMapper.ToType(byte(chain))
	if err != nil {
		fmt.Printf("Error converting OraclePubkey to String")
	}

	switch ChainType(cType) {
	case Ethereum, Binance:
		return hexutil.Encode(b)
	case Waves:
		return base58.Encode(b)
	}

	return ""
}

func StringToOraclePubKey(value string, chain ChainType) (OraclesPubKey, error) {
	var pubKey []byte
	var err error
	cType, err := ChainMapper.ToType(byte(chain))
	if err != nil {
		return OraclesPubKey{}, err
	}

	switch ChainType(cType) {
	case Ethereum, Binance:
		pubKey, err = hexutil.Decode(value)
		if err != nil {
			return [33]byte{}, err
		}
	case Waves:
		wPubKey, err := crypto.NewPublicKeyFromBase58(value)
		pubKey = wPubKey[:]
		if err != nil {
			return [33]byte{}, err
		}
	}
	return BytesToOraclePubKey(pubKey, chain), nil
}

func HexToValidatorPubKey(hex string) (ConsulPubKey, error) {
	b, err := hexutil.Decode(hex)
	if err != nil {
		return ConsulPubKey{}, err
	}
	pubKey := ConsulPubKey{}
	copy(pubKey[:], b)
	return pubKey, nil
}

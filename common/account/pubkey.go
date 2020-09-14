package account

import (
	"github.com/btcsuite/btcutil/base58"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/ed25519"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type ConsulPubKey ed25519.PubKeyEd25519
type OraclesPubKey [33]byte

func StringToPrivKey(value string, chainType ChainType) ([]byte, error) {
	var privKey []byte
	var err error
	switch chainType {
	case Ethereum:
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

func BytesToOraclePubKey(value []byte, chainType ChainType) OraclesPubKey {
	var pubKey OraclesPubKey
	switch chainType {
	case Ethereum:
		copy(pubKey[:], value[0:33])
	case Waves:
		copy(pubKey[:], append([]byte{0}, value[0:32]...))
	}
	return pubKey
}

func (pubKey *OraclesPubKey) ToBytes(chainType ChainType) []byte {
	var v []byte
	switch chainType {
	case Ethereum:
		v = pubKey[:33]
	case Waves:
		v = pubKey[1:33]
	}
	return v
}
func (pubKey *OraclesPubKey) ToString(chainType ChainType) string {
	b := pubKey.ToBytes(chainType)
	switch chainType {
	case Ethereum:
		return hexutil.Encode(b)
	case Waves:
		return base58.Encode(b)
	}

	return ""
}

func StringToOraclePubKey(value string, chainType ChainType) (OraclesPubKey, error) {
	var pubKey []byte
	var err error
	switch chainType {
	case Ethereum:
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
	return BytesToOraclePubKey(pubKey, chainType), nil
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

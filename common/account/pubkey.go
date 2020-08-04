package account

import (
	"crypto/ecdsa"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	wavesplatform "github.com/wavesplatform/go-lib-crypto"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
)

type PubKey [32]byte        //TODO length
type OraclesPubKey [33]byte //TODO length

func HexToPrivKey(value string, chainType ChainType) (privKey []byte, pubKey OraclesPubKey, err error) {
	var pubKeyBytes []byte
	switch chainType {
	case Ethereum:
		privKeyBytes, err := hexutil.Decode(value)
		if err != nil {
			return nil, OraclesPubKey{}, err
		}
		ethPrivKey := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: secp256k1.S256(),
			},
			D: new(big.Int),
		}
		ethPrivKey.D.SetBytes(privKeyBytes)
		ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(privKeyBytes)
		pubKeyBytes = crypto.CompressPubkey(&ethPrivKey.PublicKey)
		privKey = privKeyBytes
	case Waves:
		wCrypto := wavesplatform.NewWavesCrypto()
		seed := wavesplatform.Seed(value)
		secret, err := wavesCrypto.NewSecretKeyFromBase58(string(wCrypto.PrivateKey(seed)))
		if err != nil {
			return nil, OraclesPubKey{}, err
		}
		key := wavesCrypto.GeneratePublicKey(secret)
		pubKeyBytes = key.Bytes()
		privKey = secret.Bytes()
	}

	pubKey = BytesToOraclePubKey(pubKeyBytes, chainType)
	return
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

func HexToPubKey(hex string) (PubKey, error) {
	b, err := hexutil.Decode(hex)
	if err != nil {
		return PubKey{}, err
	}
	pubKey := PubKey{}
	copy(pubKey[:], b)
	return pubKey, nil
}

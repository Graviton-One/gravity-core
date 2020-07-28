package account

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	_ "github.com/tendermint/tendermint/crypto/ed25519"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
)

type ChainType byte
type PubKey [32]byte        //TODO length
type OraclesPubKey [33]byte //TODO length

const (
	Ethereum ChainType = iota
	Waves
)

var (
	ErrInvalidChainType = errors.New("invalid chain type")
	ErrParseChainType   = errors.New("invalid parse chain type")
)

func HexToPubKey(hex string) (PubKey, error) {
	b, err := hexutil.Decode(hex)
	if err != nil {
		return PubKey{}, err
	}
	pubKey := PubKey{}
	copy(pubKey[:], b)
	return pubKey, nil
}

func ParseChainType(chainType string) (ChainType, error) {
	switch strings.ToLower(chainType) {
	case "ethereum":
		return Ethereum, nil
	case "waves":
		return Waves, nil
	default:
		return 0, ErrParseChainType
	}
}
func (ch ChainType) String() string {
	switch ch {
	case Ethereum:
		return "ethereum"
	case Waves:
		return "waves"
	default:
		return "ethereum"
	}
}
func Sign(privKey tendermintCrypto.PrivKeyEd25519, msg []byte) ([]byte, error) {
	return privKey.Sign(msg)
}

func SignWithTC(privKey []byte, msg []byte, chainType ChainType) ([]byte, error) {
	switch chainType {
	case Ethereum:
		ethPrivKey := &ecdsa.PrivateKey{
			PublicKey: ecdsa.PublicKey{
				Curve: secp256k1.S256(),
			},
			D: new(big.Int),
		}
		ethPrivKey.D.SetBytes(privKey)
		ethPrivKey.PublicKey.X, ethPrivKey.PublicKey.Y = ethPrivKey.PublicKey.Curve.ScalarBaseMult(privKey)

		sig, err := crypto.Sign(msg, ethPrivKey)
		if err != nil {
			return nil, err
		}

		return sig, nil
	case Waves:
		secret, err := wavesCrypto.NewSecretKeyFromBytes(privKey)
		if err != nil {
			return nil, err
		}
		sig, err := wavesCrypto.Sign(secret, msg)
		if err != nil {
			return nil, err
		}
		return sig.Bytes(), nil
	default:
		return nil, ErrInvalidChainType
	}
}

func ValidateTCSign(pubKey OraclesPubKey, msg []byte, sign []byte, chainType ChainType) bool {
	switch chainType {
	case Ethereum:
		return crypto.VerifySignature(pubKey[:], msg, sign[0:64])
	case Waves:
		var wavesPubKey wavesCrypto.PublicKey
		copy(wavesPubKey[:], pubKey[:])

		var signWaves wavesCrypto.Signature
		copy(signWaves[:], sign[:])
		return wavesCrypto.Verify(wavesPubKey, signWaves, msg)
	default:
		return false
	}
}

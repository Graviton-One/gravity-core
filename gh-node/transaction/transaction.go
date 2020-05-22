package transaction

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/wavesplatform/gowaves/pkg/crypto"
)

type TxFunc string

const (
	Commit       TxFunc = "commit"
	Reveal       TxFunc = "reveal"
	AddValidator TxFunc = "addValidator"
	SignResult   TxFunc = "signResult"
)

type Transaction struct {
	Id           string //[HashKeySize]byte
	SenderPubKey string // [PublicKeySize]byte
	Signature    string // [SignatureSize]byte
	Func         TxFunc //[]byte
	Args         string //Args
}

func New(pubKey []byte, funcName TxFunc, secret crypto.SecretKey, args []byte) (*Transaction, error) {
	tx := &Transaction{
		SenderPubKey: hex.EncodeToString(pubKey),
		Args:         hex.EncodeToString(args),
		Func:         funcName,
	}
	tx.Hash()
	err := tx.Sign(secret)
	if err != nil {
		panic(err)
	}

	return tx, err
}

func (tx *Transaction) MarshalBytesWithoutSig() []byte {
	var result []byte
	result = append(result, tx.Id[:]...)
	result = append(result, tx.SenderPubKey[:]...)
	result = append(result, tx.Func...)
	result = append(result, tx.Args...)
	return result
}

func UnmarshalJson(data []byte) (*Transaction, error) {
	tx := new(Transaction)
	err := json.Unmarshal(data, tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tx *Transaction) Hash() {
	hash := sha256.Sum256(tx.MarshalBytesWithoutSig())
	tx.Id = hex.EncodeToString(hash[:])
}

func (tx *Transaction) Sign(key crypto.SecretKey) error {
	sig, err := crypto.Sign(key, tx.MarshalBytesWithoutSig())
	if err != nil {
		return err
	}
	tx.Signature = hex.EncodeToString(sig.Bytes())

	return nil
}

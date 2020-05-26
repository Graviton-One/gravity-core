package transaction

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"

	"github.com/ethereum/go-ethereum/crypto"
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

func New(pubKey []byte, funcName TxFunc, privKey *ecdsa.PrivateKey, args []byte) (*Transaction, error) {
	tx := &Transaction{
		SenderPubKey: hex.EncodeToString(pubKey),
		Args:         hex.EncodeToString(args),
		Func:         funcName,
	}
	tx.Hash()
	err := tx.Sign(privKey)
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
	hash := crypto.Keccak256(tx.MarshalBytesWithoutSig())
	tx.Id = hex.EncodeToString(hash[:])
}

func (tx *Transaction) Sign(privKey *ecdsa.PrivateKey) error {
	sig, err := crypto.Sign(tx.MarshalBytesWithoutSig(), privKey)
	if err != nil {
		return err
	}
	tx.Signature = hex.EncodeToString(sig)

	return nil
}

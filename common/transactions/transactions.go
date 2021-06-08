package transactions

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/ethereum/go-ethereum/crypto"
	tCrypto "github.com/tendermint/tendermint/crypto"
	_ "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	Commit                 TxFunc = "commit"
	Reveal                 TxFunc = "reveal"
	AddOracle              TxFunc = "addOracle"
	AddOracleInNebula      TxFunc = "addOracleInNebula"
	Result                 TxFunc = "result"
	NewRound               TxFunc = "newRound" //TODO: Legacy / not used
	Vote                   TxFunc = "vote"
	AddNebula              TxFunc = "setNebula"
	DropNebula             TxFunc = "dropNebula"
	SignNewConsuls         TxFunc = "signNewConsuls"
	SignNewOracles         TxFunc = "signNewOracles"
	ApproveLastRound       TxFunc = "approveLastRound"
	SetSolanaRecentBlock   TxFunc = "setSolanaRecentBlock"
	SetNebulaCustomParams  TxFunc = "setNebulaCustomParams"
	DropNebulaCustomParams TxFunc = "dropNebulaCustomParams"

	String Type = "string"
	Int    Type = "int"
	Bytes  Type = "bytes"
)

type ID [32]byte
type TxFunc string
type Type string
type Arg struct {
	Type  Type
	Value []byte
}
type Value interface{}
type StringValue struct {
	Value string
}
type IntValue struct {
	Value int64
}
type BytesValue struct {
	Value []byte
}

type Transaction struct {
	Id           ID
	SenderPubKey account.ConsulPubKey
	Signature    [72]byte
	Func         TxFunc
	Timestamp    uint64
	Args         []Arg
}

func New(pubKey account.ConsulPubKey, funcName TxFunc, privKey tCrypto.PrivKey) (*Transaction, error) {
	tx := &Transaction{
		SenderPubKey: pubKey,
		Func:         funcName,
		Timestamp:    uint64(time.Now().Unix()),
	}
	tx.Hash()

	err := tx.Sign(privKey)
	if err != nil {
		return tx, err
	}

	return tx, nil
}

func (tx *Transaction) Hash() {
	tx.Id = ID(crypto.Keccak256Hash(tx.Bytes()))
}

func (tx *Transaction) Sign(privKey tCrypto.PrivKey) error {
	sign, err := account.Sign(privKey, tx.Id.Bytes())
	if err != nil {
		return err
	}

	copy(tx.Signature[:], sign)
	return nil
}

func (tx *Transaction) Bytes() []byte {
	var result []byte
	result = append(result, tx.Id[:]...)
	result = append(result, tx.SenderPubKey[:]...)
	result = append(result, tx.Func...)

	for _, v := range tx.Args {
		result = append(result, v.Value...)
	}

	var b [8]byte
	binary.BigEndian.PutUint64(b[:], tx.Timestamp)
	result = append(result, b[:]...)

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

func (id ID) Bytes() []byte {
	return id[:]
}

func (tx *Transaction) AddValue(value Value) {
	var b []byte
	var t Type
	switch value.(type) {
	case StringValue:
		b = []byte(value.(StringValue).Value)
		t = String
	case IntValue:
		var bInt [8]byte
		binary.BigEndian.PutUint64(bInt[:], uint64(value.(IntValue).Value))
		b = bInt[:]
		t = Int
	case BytesValue:
		b = value.(BytesValue).Value
		t = Bytes
	}
	tx.Args = append(tx.Args, Arg{
		Type:  t,
		Value: b,
	})
}

func (tx *Transaction) AddValues(values []Value) {
	for _, value := range values {
		tx.AddValue(value)
	}
}

func (tx *Transaction) Value(index int) interface{} {
	v := tx.Args[index]

	switch v.Type {
	case String:
		return string(v.Value)
	case Int:
		return int64(binary.BigEndian.Uint64(v.Value))
	case Bytes:
		return v.Value
	}

	return nil
}

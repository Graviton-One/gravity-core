package transactions

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/ethereum/go-ethereum/crypto"

	_ "github.com/tendermint/tendermint/crypto/ed25519"
	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

const (
	Commit            TxFunc = "commit"
	Reveal            TxFunc = "reveal"
	AddOracle         TxFunc = "addOracle"
	AddOracleInNebula TxFunc = "addOracleInNebula"
	Result            TxFunc = "result"
	NewRound          TxFunc = "newRound"
	Vote              TxFunc = "vote"

	StringType ArgType = "string"
	IntType    ArgType = "int"
	BinaryType ArgType = "binary"
)

type ID [32]byte
type TxFunc string
type ArgType string
type Args struct {
	Value interface{}
}

type Transaction struct {
	Id           ID
	SenderPubKey account.ValidatorPubKey
	Signature    [72]byte
	Func         TxFunc
	Timestamp    uint64
	Args         []Args
}

func New(pubKey account.ValidatorPubKey, funcName TxFunc, privKey tendermintCrypto.PrivKeyEd25519, args []Args) (*Transaction, error) {
	tx := &Transaction{
		SenderPubKey: pubKey,
		Args:         args,
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

func (tx *Transaction) Sign(privKey tendermintCrypto.PrivKeyEd25519) error {
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
		var value []byte
		switch v.Value.(type) {
		case string:
			value = []byte(v.Value.(string))
		case int64:
			var b [8]byte
			binary.BigEndian.PutUint64(b[:], tx.Timestamp)
			value = b[:]
		case []byte:
			value = v.Value.([]byte)
		}
		result = append(result, value...)
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

package helpers

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/fatih/structs"
	"github.com/mitchellh/mapstructure"
)

type TxType int
type ContractFunc string

const (
	MassTransfer TxType = 11
	InvokeScript TxType = 16

	MassTransferReferralAttachment = "sp"
)

type Transaction struct {
	Sender          string   `structs:"sender"`
	SenderPublicKey string   `structs:"senderPublicKey"`
	Fee             int      `structs:"fee"`
	Type            TxType   `structs:"type"`
	Version         int      `structs:"version"`
	Proofs          []string `structs:"proofs"`
	ID              string   `structs:"id"`
	Timestamp       int64    `structs:"timestamp"`
	Height          int      `structs:"height"`
	Attachment      string   `structs:"attachment"`

	InvokeScriptBody *InvokeScriptBody `structs:"-"`
}

func NewTransaction(txType TxType, sender string) Transaction {
	return Transaction{
		Type:      txType,
		Version:   1,
		Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		Sender:    sender,
	}
}

func (tx Transaction) JsonMap() (map[string]interface{}, error) {
	txMap := structs.Map(tx)
	var bodyMap map[string]interface{}
	switch tx.Type {
	case InvokeScript:
		bodyMap = structs.Map(tx.InvokeScriptBody)
	default:
		errors.New("invalid tx type")

	}
	for k, v := range bodyMap {
		txMap[k] = v
	}
	return txMap, nil
}

func Parse(json map[string]interface{}) (Transaction, error) {
	tx := Transaction{}
	err := mapstructure.Decode(json, &tx)
	if err != nil {
		return tx, err
	}

	switch tx.Type {
	case InvokeScript:
		tx.InvokeScriptBody = &InvokeScriptBody{}
		err = mapstructure.Decode(json, tx.InvokeScriptBody)
	default:
		errors.New("invalid tx type")
	}

	return tx, err
}

func (tx Transaction) Marshal() ([]byte, error) {
	jsonMap, err := tx.JsonMap()
	if err != nil {
		return nil, err
	}
	return json.Marshal(jsonMap)
}

func Unmarshal(data []byte) (Transaction, error) {
	jsonMap := make(map[string]interface{})
	if err := json.Unmarshal(data, &jsonMap); err != nil {
		return Transaction{}, err
	}
	tx, err := Parse(jsonMap)

	return tx, err
}

func GroupByHeightAndFunc(txs []Transaction) map[int]map[ContractFunc][]Transaction {
	groupedTxs := make(map[int]map[ContractFunc][]Transaction)

	for _, v := range txs {
		if v.InvokeScriptBody == nil {
			continue
		}
		_, ok := groupedTxs[v.Height]
		if !ok {
			groupedTxs[v.Height] = make(map[ContractFunc][]Transaction)
		}
		groupedTxs[v.Height][v.InvokeScriptBody.Call.Function] = append(groupedTxs[v.Height][v.InvokeScriptBody.Call.Function], v)
	}
	return groupedTxs
}

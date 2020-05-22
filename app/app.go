package app

import (
	"encoding/json"
	"errors"
	"gravity-hub/transaction"

	"github.com/dgraph-io/badger"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

const (
	Success  uint32 = 0
	Error    uint32 = 1
	NotFound uint32 = 404
)

type GHApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn
}

var _ abcitypes.Application = (*GHApplication)(nil)

func NewGHApplication(db *badger.DB) *GHApplication {
	return &GHApplication{
		db: db,
	}
}

func (GHApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (GHApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *GHApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	err := app.isValid(req.Tx)
	if err != nil {
		return abcitypes.ResponseDeliverTx{Code: Error}
	}

	tx, _ := transaction.UnmarshalJson(req.Tx)
	err = tx.SetState(app.currentBatch)
	if err != nil {
		//TODO
		panic(err)
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *GHApplication) isValid(txBytes []byte) error {
	tx, err := transaction.UnmarshalJson(txBytes)
	if err != nil {
		return errors.New("invalid parse tx")
	}

	err = tx.IsValid(app.db)
	if err != nil {
		return err
	}
	return nil
}

func (app *GHApplication) CheckTx(req abcitypes.RequestCheckTx) abcitypes.ResponseCheckTx {
	err := app.isValid(req.Tx)
	if err != nil {
		return abcitypes.ResponseCheckTx{Code: Error, Info: err.Error()}
	}
	return abcitypes.ResponseCheckTx{Code: Success}
}

func (app *GHApplication) Commit() abcitypes.ResponseCommit {
	app.currentBatch.Commit()
	return abcitypes.ResponseCommit{Data: []byte{}}
}

func (app *GHApplication) Query(reqQuery abcitypes.RequestQuery) (resQuery abcitypes.ResponseQuery) {
	var err error
	switch reqQuery.Path {
	case "key":
		resQuery.Key = reqQuery.Data
		err = app.db.View(func(txn *badger.Txn) error {
			item, err := txn.Get(reqQuery.Data)
			if err != nil && err != badger.ErrKeyNotFound {
				return err
			}
			if err == badger.ErrKeyNotFound {
				resQuery.Info = "does not exist"
				resQuery.Code = Error
				return nil
			}

			return item.Value(func(val []byte) error {
				resQuery.Info = "exists"
				resQuery.Value = val
				return nil
			})
		})
	case "prefix":
		err = app.db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()
			prefix := reqQuery.Data

			values := make(map[string][]byte)
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()
				k := item.Key()
				err := item.Value(func(v []byte) error {
					values[string(k)] = v
					return nil
				})
				if err != nil {
					return err
				}
			}
			result, err := json.Marshal(&values)
			if err != nil {
				return err
			}

			resQuery.Log = "exists"
			resQuery.Value = result

			return nil
		})
	}

	if err != nil {
		panic(err)
	}
	return
}

func (GHApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	return abcitypes.ResponseInitChain{}
}

func (app *GHApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (GHApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	return abcitypes.ResponseEndBlock{}
}

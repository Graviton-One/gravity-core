package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Gravity-Tech/proof-of-concept/common/account"

	"github.com/Gravity-Tech/proof-of-concept/common/storage"
	"github.com/Gravity-Tech/proof-of-concept/common/transactions"
	"github.com/Gravity-Tech/proof-of-concept/ledger-node/scheduler"

	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/dgraph-io/badger"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

const (
	Success uint32 = 0
	Error   uint32 = 500

	ValidatorCount = 5
)

type GHApplication struct {
	db          *badger.DB
	storage     *storage.Storage
	ethClient   *ethclient.Client
	wavesClient *client.Client
	scheduler   *scheduler.Scheduler
	ctx         context.Context
	initScores  map[string]uint64
}

var _ abcitypes.Application = (*GHApplication)(nil)

func NewGHApplication(ethClient *ethclient.Client, wavesClient *client.Client, scheduler *scheduler.Scheduler, db *badger.DB, initScores map[string]uint64, ctx context.Context) *GHApplication {
	return &GHApplication{
		db:          db,
		ethClient:   ethClient,
		wavesClient: wavesClient,
		scheduler:   scheduler,
		ctx:         ctx,
		initScores:  initScores,
	}
}

func (app *GHApplication) Info(req abcitypes.RequestInfo) abcitypes.ResponseInfo {
	return abcitypes.ResponseInfo{}
}

func (app *GHApplication) SetOption(req abcitypes.RequestSetOption) abcitypes.ResponseSetOption {
	return abcitypes.ResponseSetOption{}
}

func (app *GHApplication) DeliverTx(req abcitypes.RequestDeliverTx) abcitypes.ResponseDeliverTx {
	err := app.isValid(req.Tx)
	if err != nil {
		return abcitypes.ResponseDeliverTx{Code: Error}
	}

	tx, _ := transactions.UnmarshalJson(req.Tx)
	err = tx.SetState(app.storage)
	if err != nil {
		return abcitypes.ResponseDeliverTx{Code: Error}
	}

	return abcitypes.ResponseDeliverTx{Code: 0}
}

func (app *GHApplication) isValid(txBytes []byte) error {
	tx, err := transactions.UnmarshalJson(txBytes)
	if err != nil {
		return errors.New("invalid parse tx")
	}

	err = tx.IsValid(app.ethClient, app.wavesClient, app.db, app.ctx)
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
	err := app.storage.Commit()
	if err != nil {
		panic(err)
	}
	return abcitypes.ResponseCommit{}
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
		resQuery.Info = "invalid request"
		resQuery.Code = Error
	}
	return
}

func (app *GHApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	app.storage.NewTransaction(app.db)
	for key, value := range app.initScores {
		validatorPubKey, err := account.HexToPubKey(key)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}

		err = app.storage.SetScore(validatorPubKey, value)
		if err != nil {
			panic(err)
		}
	}

	for _, value := range req.Validators {
		validatorPubKey := account.PubKey(value.PubKey.GetData())
		err := app.storage.SetScore(validatorPubKey, uint64(value.Power))
		if err != nil {
			panic(err)
		}
	}

	err := app.storage.Commit()
	if err != nil {
		panic(err)
	}

	return abcitypes.ResponseInitChain{}
}

func (app *GHApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.storage.NewTransaction(app.db)

	err := app.scheduler.HandleBlock(req.Header.Height, app.storage)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
	}

	return abcitypes.ResponseBeginBlock{}
}

func (app *GHApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	consuls, err := app.storage.GetConsuls()
	if err != nil {
		panic(err)
	}

	var newValidators []abcitypes.ValidatorUpdate
	for i := 0; i < ValidatorCount && i < len(consuls); i++ {
		if consuls[i].Value == 0 {
			continue
		}

		pubKey := abcitypes.PubKey{
			Type: "ed25519",
			Data: consuls[i].Validator,
		}

		newValidators = append(newValidators, abcitypes.ValidatorUpdate{
			PubKey: pubKey,
			Power:  int64(consuls[i].Value),
		})
	}

	if len(newValidators) > 0 {
		return abcitypes.ResponseEndBlock{ValidatorUpdates: newValidators}
	} else {
		return abcitypes.ResponseEndBlock{}
	}
}

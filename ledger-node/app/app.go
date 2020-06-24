package app

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"gravity-hub/common/keys"
	"gravity-hub/common/transactions"
	"gravity-hub/ledger-node/scheduler"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/dgraph-io/badger"
	abcitypes "github.com/tendermint/tendermint/abci/types"
)

const (
	Success uint32 = 0
	Error   uint32 = 500

	ValidatorCount = 5 // TODO: config
)

type GHApplication struct {
	db           *badger.DB
	currentBatch *badger.Txn
	ethClient    *ethclient.Client
	wavesClient  *client.Client
	scheduler    *scheduler.Scheduler
	ctx          context.Context
	initScores   map[string]uint64
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
	err = tx.SetState(app.currentBatch)
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
		resQuery.Info = "invalid request"
		resQuery.Code = Error
	}
	return
}

func (app *GHApplication) InitChain(req abcitypes.RequestInitChain) abcitypes.ResponseInitChain {
	for key, value := range app.initScores {
		var scoreBytes []byte
		binary.BigEndian.PutUint64(scoreBytes, value)

		validatorAddress, err := hexutil.Decode(key)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
		err = app.currentBatch.Set([]byte(keys.FormScoreKey(validatorAddress)), scoreBytes)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}

	return abcitypes.ResponseInitChain{}
}

func (app *GHApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	err := app.scheduler.HandleBlock(req.Header.Height, app.currentBatch)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}

	app.currentBatch = app.db.NewTransaction(true)
	return abcitypes.ResponseBeginBlock{}
}

func (app *GHApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	it := app.currentBatch.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(keys.ScoreKey)

	type Scores struct {
		Validator string
		Value     uint64
	}

	var scores []Scores

	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		err := item.Value(func(v []byte) error {
			validator := strings.Split(keys.Separator, string(k))[1]

			scores = append(scores, Scores{
				Value:     binary.BigEndian.Uint64(v),
				Validator: validator,
			})
			return nil
		})
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}

	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Value < scores[j].Value
	})

	var newValidators []abcitypes.ValidatorUpdate

	for i := 0; i < ValidatorCount; i++ {
		pubKeyBytes, err := hexutil.Decode(scores[i].Validator)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return abcitypes.ResponseEndBlock{}
		}

		pubKey := abcitypes.PubKey{}
		err = pubKey.Unmarshal(pubKeyBytes)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return abcitypes.ResponseEndBlock{}
		}

		newValidators = append(newValidators, abcitypes.ValidatorUpdate{
			PubKey: pubKey,
			Power:  int64(scores[i].Value),
		})
	}

	return abcitypes.ResponseEndBlock{ValidatorUpdates: newValidators}
}

package app

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"proof-of-concept/common/keys"
	"proof-of-concept/common/transactions"
	"proof-of-concept/ledger-node/scheduler"

	"github.com/ethereum/go-ethereum/common/hexutil"

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
	app.currentBatch = app.db.NewTransaction(true)
	for key, value := range app.initScores {
		var scoreBytes [8]byte
		binary.BigEndian.PutUint64(scoreBytes[:], value)

		validatorAddress, err := hexutil.Decode(key)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
		err = app.currentBatch.Set([]byte(keys.FormScoreKey(validatorAddress)), scoreBytes[:])
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}

	for _, value := range req.Validators {
		var scoreBytes [8]byte
		binary.BigEndian.PutUint64(scoreBytes[:], uint64(value.Power))

		validatorPubKey := value.PubKey.GetData()
		err := app.currentBatch.Set([]byte(keys.FormScoreKey(validatorPubKey)), scoreBytes[:])
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
		}
	}
	err := app.currentBatch.Commit()
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return abcitypes.ResponseInitChain{}
}

func (app *GHApplication) BeginBlock(req abcitypes.RequestBeginBlock) abcitypes.ResponseBeginBlock {
	app.currentBatch = app.db.NewTransaction(true)

	err := app.scheduler.HandleBlock(req.Header.Height, app.currentBatch)
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
	}

	return abcitypes.ResponseBeginBlock{}
}

func (app *GHApplication) EndBlock(req abcitypes.RequestEndBlock) abcitypes.ResponseEndBlock {
	var scores []scheduler.Scores
	app.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(keys.FormConsulsKey()))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if err == badger.ErrKeyNotFound {
			return nil
		}

		b, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &scores)
		if err != nil {
			return err
		}

		return nil
	})

	var newValidators []abcitypes.ValidatorUpdate

	for i := 0; i < ValidatorCount && i < len(scores); i++ {
		if scores[i].Value == 0 {
			continue
		}
		pubKeyBytes, err := hexutil.Decode(scores[i].Validator)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			return abcitypes.ResponseEndBlock{}
		}

		pubKey := abcitypes.PubKey{
			Type: "ed25519",
			Data: pubKeyBytes,
		}

		newValidators = append(newValidators, abcitypes.ValidatorUpdate{
			PubKey: pubKey,
			Power:  int64(scores[i].Value),
		})
	}

	if len(newValidators) > 0 {
		return abcitypes.ResponseEndBlock{ValidatorUpdates: newValidators}
	} else {
		return abcitypes.ResponseEndBlock{}
	}

}

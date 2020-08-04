package storage

import (
	"encoding/json"
	"strings"

	"github.com/dgraph-io/badger"
)

const (
	Separator string = "_"

	ConsulsKey           Key = "consuls"
	LastHeightKey        Key = "last_height"
	PrevConsulsKey       Key = "prev_consuls"
	ConsulsSignKey       Key = "consuls_sing"
	OraclesSignNebulaKey Key = "oracles_sign"

	OraclesByNebulaKey    Key = "oracles_nebula"
	BftOraclesByNebulaKey Key = "bft_oracles_nebula"
	OraclesByValidatorKey Key = "oracles"

	BlockKey      Key = "block"
	VoteKey       Key = "vote"
	ScoreKey      Key = "score"
	CommitKey     Key = "commit"
	RevealKey     Key = "reveal"
	SignResultKey Key = "signResult"
)

var (
	ErrKeyNotFound = badger.ErrKeyNotFound
)

type Key string
type Storage struct {
	txn *badger.Txn
}

func formKey(args ...string) []byte {
	return []byte(strings.Join(args, Separator))
}

func (storage *Storage) getValue(key []byte) ([]byte, error) {
	item, err := storage.txn.Get(key)
	if err != nil {
		return nil, err
	}

	b, err := item.ValueCopy(nil)
	if err != nil {
		return nil, err

	}

	return b, nil
}

func (storage *Storage) setValue(key []byte, jsonOrBytes interface{}) error {
	var b []byte
	var err error

	b, ok := jsonOrBytes.([]byte)
	if !ok {
		b, err = json.Marshal(jsonOrBytes)
		if err != nil {
			return err
		}
	}

	return storage.txn.Set(key, b)
}

func (storage *Storage) NewTransaction(db *badger.DB) {
	storage.txn = db.NewTransaction(true)
}

func (storage *Storage) Commit() error {
	return storage.txn.Commit()
}

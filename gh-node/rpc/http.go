package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/Gravity-Tech/proof-of-concept/common/account"
	"github.com/Gravity-Tech/proof-of-concept/common/transactions"
	"github.com/Gravity-Tech/proof-of-concept/gh-node/api/gravity"
	"github.com/Gravity-Tech/proof-of-concept/score-calculator/models"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

type ServerConfig struct {
	Host      string
	PubKey    []byte
	PrivKey   tendermintCrypto.PrivKeyEd25519
	ChainType account.ChainType
	GhClient  *gravity.Client
}

var cfg ServerConfig

func ListenRpcServer(config ServerConfig) error {
	cfg = config
	http.HandleFunc("/vote", vote)
	err := http.ListenAndServe(cfg.Host, nil)
	return err
}

func vote(w http.ResponseWriter, r *http.Request) {
	var request []models.Vote
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := transactions.New(cfg.PubKey, transactions.Vote, cfg.ChainType, cfg.PrivKey, b)
	err = cfg.GhClient.SendTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

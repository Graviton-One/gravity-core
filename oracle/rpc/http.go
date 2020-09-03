package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

type ServerConfig struct {
	Host      string
	PubKey    account.ConsulPubKey
	PrivKey   tendermintCrypto.PrivKeyEd25519
	ChainType account.ChainType
	GhClient  *client.Client
}
type VotesRq struct {
	votes []storage.Vote
}

var cfg ServerConfig

func ListenRpcServer(config ServerConfig) error {
	cfg = config
	http.HandleFunc("/vote", vote)
	err := http.ListenAndServe(cfg.Host, nil)
	return err
}

func vote(w http.ResponseWriter, r *http.Request) {
	var request VotesRq
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := json.Marshal(request.votes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := transactions.New(cfg.PubKey, transactions.Vote, cfg.PrivKey, []transactions.Args{
		{Value: b},
	})
	err = cfg.GhClient.SendTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

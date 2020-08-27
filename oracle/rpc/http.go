package rpc

import (
	"encoding/json"
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

var cfg *RPCConfig

func ListenRpcServer(config *RPCConfig) error {
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

	tx, err := transactions.New(cfg.pubKey, transactions.Vote, cfg.privKey, []transactions.Args{
		{Value: b},
	})
	err = cfg.client.SendTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

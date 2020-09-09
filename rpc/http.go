package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

var cfg *Config

func ListenRpcServer(config *Config) {
	cfg = config
	http.HandleFunc("/vote", vote)
	err := http.ListenAndServe(cfg.Host, nil)
	if err != nil {
		fmt.Printf("Error Private RPC: %s", err.Error())
	}
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

	tx, err := transactions.New(cfg.pubKey, transactions.Vote, cfg.privKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tx.AddValue(transactions.BytesValue{Value: b})
	err = cfg.client.SendTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

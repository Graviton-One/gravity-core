package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/Gravity-Tech/gravity-core/common/storage"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

var cfg *Config

type VotesRq struct {
	Votes []VoteRq
}
type VoteRq struct {
	PubKey string
	Score  uint64
}

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
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tx, err := transactions.New(cfg.pubKey, transactions.Vote, cfg.privKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var votes []storage.Vote
	for _, v := range request.Votes {
		pubKey, err := account.HexToValidatorPubKey(v.PubKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		votes = append(votes, storage.Vote{
			PubKey: pubKey,
			Score:  v.Score,
		})
	}
	b, err := json.Marshal(votes)
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

	return
}

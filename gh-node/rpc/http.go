package rpc

import (
	"encoding/json"
	"gravity-hub/common/account"
	"gravity-hub/common/transactions"
	"gravity-hub/gh-node/api/gravity"
	"gravity-hub/gravity-score-calculator/models"
	"net/http"

	tendermintCrypto "github.com/tendermint/tendermint/crypto/ed25519"
)

//TODO: to struct
var pubKey []byte
var privKey tendermintCrypto.PrivKeyEd25519
var chainType account.ChainType
var ghClient gravity.Client

func ListenRpcServer(host string, newPubKey []byte, newPrivKey tendermintCrypto.PrivKeyEd25519, newChainType account.ChainType, newGhClient gravity.Client) error {
	pubKey = newPubKey
	privKey = newPrivKey
	chainType = newChainType
	ghClient = newGhClient
	http.HandleFunc("/vote", vote)
	err := http.ListenAndServe(host, nil)
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

	tx, err := transactions.New(pubKey, transactions.Vote, chainType, privKey, b)
	err = ghClient.SendTx(tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

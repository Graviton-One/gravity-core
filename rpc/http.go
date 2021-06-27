package rpc

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/storage"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
)

var cfg *Config

type AddNebulaRq struct {
	NebulaId             string
	ChainType            string
	MaxPulseCountInBlock uint64
	MinScore             uint64
}
type DropNebulaRq struct {
	NebulaId  string
	ChainType string
}
type SetNebulaCustomParamsRq struct {
	NebulaId  string
	ChainType string
	Params    storage.NebulaCustomParams
}
type DropNebulaCustomParamsRq struct {
	NebulaId  string
	ChainType string
}

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
	http.HandleFunc("/setNebula", func(rp http.ResponseWriter, rq *http.Request) {
		nebulaHandler(rp, rq, addNebula)
	})
	http.HandleFunc("/dropNebula", func(rp http.ResponseWriter, rq *http.Request) {
		nebulaHandler(rp, rq, dropNebula)
	})
	http.HandleFunc("/setNebulaCustomParams", func(rp http.ResponseWriter, rq *http.Request) {
		nebulaHandler(rp, rq, setNebulaCustomParams)
	})
	http.HandleFunc("/dropNebulaCustomParams", func(rp http.ResponseWriter, rq *http.Request) {
		nebulaHandler(rp, rq, dropNebulaCustomParams)
	})
	http.HandleFunc("/listNebulas", func(rp http.ResponseWriter, rq *http.Request) {
		listNebulasHandler(rp, rq)
	})

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

func nebulaHandler(w http.ResponseWriter, r *http.Request, action func(rq *http.Request) error) {
	err := action(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	return
}

func dropNebula(r *http.Request) error {
	var request DropNebulaRq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return err
	}

	tx, err := transactions.New(cfg.pubKey, transactions.DropNebula, cfg.privKey)
	if err != nil {
		return err
	}

	chainType, err := account.ParseChainType(request.ChainType)
	if err != nil {
		return err
	}
	nebulaId, err := account.StringToNebulaId(request.NebulaId, chainType)
	if err != nil {
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{Value: nebulaId[:]},
	})
	err = cfg.client.SendTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func addNebula(r *http.Request) error {
	var request AddNebulaRq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return err
	}

	tx, err := transactions.New(cfg.pubKey, transactions.AddNebula, cfg.privKey)
	if err != nil {
		return err
	}

	chainType, err := account.ParseChainType(request.ChainType)
	if err != nil {
		return err
	}
	nebulaId, err := account.StringToNebulaId(request.NebulaId, chainType)
	if err != nil {
		return err
	}
	zap.L().Sugar().Debugf("Try to add nebula [%s]", request.NebulaId)
	nebulaInfo := storage.NebulaInfo{
		MaxPulseCountInBlock: request.MaxPulseCountInBlock,
		MinScore:             request.MinScore,
		ChainType:            chainType,
		Owner:                cfg.pubKey,
	}
	b, err := json.Marshal(nebulaInfo)
	if err != nil {
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{Value: nebulaId[:]},
		transactions.BytesValue{Value: b},
	})
	err = cfg.client.SendTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func setNebulaCustomParams(r *http.Request) error {
	var request SetNebulaCustomParamsRq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return err
	}

	tx, err := transactions.New(cfg.pubKey, transactions.SetNebulaCustomParams, cfg.privKey)
	if err != nil {
		return err
	}

	chainType, err := account.ParseChainType(request.ChainType)
	if err != nil {
		return err
	}
	nebulaId, err := account.StringToNebulaId(request.NebulaId, chainType)
	if err != nil {
		return err
	}

	nebulaCustomParams := request.Params
	b, err := json.Marshal(nebulaCustomParams)
	if err != nil {
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{Value: nebulaId[:]},
		transactions.BytesValue{Value: b},
	})
	err = cfg.client.SendTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func dropNebulaCustomParams(r *http.Request) error {
	var request DropNebulaCustomParamsRq
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return err
	}

	tx, err := transactions.New(cfg.pubKey, transactions.DropNebulaCustomParams, cfg.privKey)
	if err != nil {
		return err
	}

	chainType, err := account.ParseChainType(request.ChainType)
	if err != nil {
		return err
	}
	nebulaId, err := account.StringToNebulaId(request.NebulaId, chainType)
	if err != nil {
		return err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{Value: nebulaId[:]},
	})
	err = cfg.client.SendTx(tx)
	if err != nil {
		return err
	}

	return nil
}

func listNebulasHandler(w http.ResponseWriter, r *http.Request) {
	list, err := cfg.client.Nebulae()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"score-calculator/models"
	"score-calculator/trustgraph"
)

func main() {
	args := os.Args[1:]
	println("Start server")
	http.HandleFunc("/api/calculate", handleCalculateScore)
	http.ListenAndServe(args[0], nil)
}

func handleCalculateScore(w http.ResponseWriter, r *http.Request) {
	request := models.Request{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	group := trustgraph.NewGroup()
	actors := make(map[string]int)
	for i, v := range request.Actors {
		actors[v.Name] = i
		group.InitialTrust(i, v.InitScore)
	}
	for k, v := range request.Votes {
		for _, scoreV := range v {
			group.Add(actors[k], actors[scoreV.Target], scoreV.Score)
		}
	}

	out := group.Compute()

	response := models.Response{}
	response.Score = make(map[string]float32)
	for i, v := range out {
		response.Score[request.Actors[i].Name] = v
	}

	result, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprintf(w, string(result))
}

package gravity

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/ledger/query"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func TestNebulaInfo(t *testing.T) {
	rq := query.ByNebulaRq{
		ChainType:     6,
		NebulaAddress: "4VL4hsSPPNdqP5ajXinJ3L434uycugBxYaJiJ2Zv4FPo",
	}
	var rqi interface{}
	rqi = rq

	var err error
	var b []byte
	b, ok := rqi.([]byte)
	if !ok {
		b, err = json.Marshal(rq)
		if err != nil {
			log.Print("KARAMBA 1")
			log.Print(err)
			t.FailNow()
		}
	}
	client, err := rpchttp.New("http://localhost:26657", "/websocket")
	if err != nil {
		log.Print("KARAMBA 2")
		log.Print(err)
		t.FailNow()
	}

	rs, err := client.ABCIQuery(string("nebula_info"), b)
	if err != nil {
		log.Print("KARAMBA 3")
		log.Print(err)
		t.FailNow()
	}

	nebulaCustomParams := storage.NebulaCustomParams{}
	log.Print("RESULT")
	log.Print(rs.Response.Value)
	err = json.Unmarshal(rs.Response.Value, &nebulaCustomParams)
	if err != nil {
		log.Print("KARAMBA 4")
		log.Print(err)
		t.FailNow()
	}

	log.Print("KARAMBA 5")
	log.Print(nebulaCustomParams)
	t.FailNow()
}

func TestNebulaCustomParams(t *testing.T) {
	rq := query.ByNebulaRq{
		ChainType:     6,
		NebulaAddress: "4VL4hsSPPNdqP5ajXinJ3L434uycugBxYaJiJ2Zv4FPo",
	}
	var rqi interface{}
	rqi = rq

	var err error
	var b []byte
	b, ok := rqi.([]byte)
	if !ok {
		b, err = json.Marshal(rq)
		if err != nil {
			log.Print("KARAMBA 1")
			log.Print(err)
			t.FailNow()
		}
	}
	client, err := rpchttp.New("http://localhost:26657", "/websocket")
	if err != nil {
		log.Print("KARAMBA 2")
		log.Print(err)
		t.FailNow()
	}

	rs, err := client.ABCIQuery(string("nebulaCustomParams"), b)
	if err != nil {
		log.Print("KARAMBA 3")
		log.Print(err)
		t.FailNow()
	}

	nebulaCustomParams := storage.NebulaCustomParams{}
	log.Print("RESULT")
	log.Print(rs.Response.Value)
	err = json.Unmarshal(rs.Response.Value, &nebulaCustomParams)
	if err != nil {
		log.Print("KARAMBA 4")
		log.Print(err)
		t.FailNow()
	}

	log.Print("KARAMBA 5")
	log.Print(nebulaCustomParams)
	t.FailNow()
}

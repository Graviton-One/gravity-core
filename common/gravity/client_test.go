package gravity

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Gravity-Tech/gravity-core/ledger/query"
)

func TestClient_do(t *testing.T) {
	rq := query.ByNebulaRq{
		ChainType:     6,
		NebulaAddress: "4VL4hsSPPNdqP5ajXinJ3L434uycugBxYaJiJ2Zv4FPo",
	}
	var rqi interface{}

	rqi = rq
	var err error
	b, ok := rqi.([]byte)
	if !ok {
		b, err = json.Marshal(rq)
		if err != nil {
			t.FailNow()
		}
	}
	log.Print(b)
	t.FailNow()
}

package gravity

import (
	"log"
	"testing"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/btcsuite/btcutil/base58"
)

func TestClient_do(t *testing.T) {
	// rq := query.ByNebulaRq{
	// 	ChainType:     6,
	// 	NebulaAddress: "4VL4hsSPPNdqP5ajXinJ3L434uycugBxYaJiJ2Zv4FPo",
	// }
	ghClient, err := gravity.New("http://localhost:26657")

	if err != nil {
		log.Print(err)
		t.FailNow()
	}
	b := base58.Decode("4VL4hsSPPNdqP5ajXinJ3L434uycugBxYaJiJ2Zv4FPo")
	nid := account.BytesToNebulaId(b)
	r, err := ghClient.NebulaCustomParams(6, nid)
	if err != nil {
		log.Print(err)
		t.FailNow()
	}

	log.Print(r)
	t.FailNow()
}

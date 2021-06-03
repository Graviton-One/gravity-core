package instructions

import (
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/near/borsh-go"
	"github.com/portto/solana-go-sdk/types"
)

func TestSerdeSerializer(t *testing.T) {
	a := types.NewAccount()
	c := NewNebulaContract()
	// c.RoundsDict[123] = 1
	// c.RoundsDict[12] = 0
	// c.RoundsDict[14] = 1
	c.LastRound = 13
	c.MultisigAccount = a.PublicKey
	data, err := borsh.Serialize(c)

	if err != nil {
		log.Fatalf("msgpack SerializeToBytes: %v", err)
	}
	fmt.Print(hex.Dump(data))

	b := NewNebulaContract()
	err = borsh.Deserialize(&b, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, b) {
		t.Error(c, b)
	}
}

type A struct {
	X uint64
	Y string
	B map[int]byte
	Z string `borsh_skip:"true"` // will skip this field when serializing/deserializing
}

func TestBorshSerializer(t *testing.T) {
	x := A{
		X: 3301,
		Y: "liber primus",
		B: make(map[int]byte),
	}
	x.B[3] = 1
	data, err := borsh.Serialize(x)
	log.Print(data)
	if err != nil {
		t.Error(err)
	}
	y := new(A)
	err = borsh.Deserialize(y, data)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(x, *y) {
		t.Error(x, y)
	}
}

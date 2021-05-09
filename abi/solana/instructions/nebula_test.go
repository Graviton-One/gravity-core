package instructions

import (
	"encoding/hex"
	"fmt"
	"log"
	"testing"
)

func TestSerdeSerializer(t *testing.T) {
	c := NewNebulaContract()
	c.RoundsDict[123] = true
	c.RoundsDict[12] = false

	myserializer := NewNebulaSerializer(6000)

	err := c.Serialize(myserializer)

	if err != nil {
		log.Fatalf("msgpack SerializeToBytes: %v", err)
	}
	fmt.Print(hex.Dump(myserializer.GetBytes()))

	mydeserializer := NewNebulaDeserializer(myserializer.GetBytes(), 6000)

	n, err := DeserializeNebulaContract(mydeserializer)
	if err != nil {
		log.Fatalf("msgpack SerializeToBytes: %v", err)
	}
	fmt.Println("")
	fmt.Print(n.RoundsDict)
	t.FailNow()
}

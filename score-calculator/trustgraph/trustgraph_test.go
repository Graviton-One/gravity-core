package trustgraph

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	g := NewGroup()
	g.Add(1, 2, 1)
	g.Add(1, 3, .5)
	g.Add(2, 1, 1)
	g.Add(2, 3, .5)
	g.Add(3, 1, 1)
	g.Add(3, 2, 1)

	g.InitialTrust(1, 1)

	out := g.Compute()

	if out[1] < 0.975 {
		t.Error("Trust in node 1 should be closer to 1.00")
	}

	if out[2] < 0.93 {
		t.Error("Trust in node 2 should be closer to 1.00")
	}
	if out[3] < 0.4 || out[3] > 0.6 {
		t.Error("Trust in node 3 should be closer to 0.50")
	}
}

func TestRand(t *testing.T) {
	peers := 200
	rand.Seed(time.Now().UTC().UnixNano())
	g := NewGroup()

	//randomly set actual trust values for peers
	actualTrust := make([]float32, peers)
	for i := 0; i < peers; i++ {
		actualTrust[i] = rand.Float32()
	}

	// peer0 is set to and granted 100% trust
	actualTrust[0] = 1
	g.InitialTrust(0, 1)

	// set 30% of trust values to +/- 10% of actual trust
	for i := 0; i < peers; i++ {
		for j := 0; j < peers; j++ {
			if rand.Float32() > .7 {
				g.Add(i, j, randNorm(actualTrust[j]))
			}
		}
	}

	// compute trust
	out := g.Compute()

	// find RMS error
	e := float32(0)
	for i := 0; i < peers; i++ {
		x := actualTrust[i] - out[i]
		e += x * x
	}
	e = float32(math.Sqrt(float64(e / float32(peers))))

	if e > .2 {
		t.Error("RMS Error should be less than 20% for a 30% full trust grid of 200 nodes")
	}
}

// randNorm takes a float and returns a value within +/- 10%,
// without going over 1
func randNorm(x float32) float32 {
	r := rand.Float32()*.2 + .9
	x *= r
	if x > 1 {
		return 1
	}
	return x
}

func TestRangeError(t *testing.T) {
	g := NewGroup()

	err := g.Add(1, 2, 1.1)
	if err.Error() != "Trust amount cannot be greater than 1" {
		t.Error("Expected error")
	}

	err = g.Add(1, 2, -1)
	if err.Error() != "Trust amount cannot be less than 0" {
		t.Error("Expected error less than 0 error")
	}

	err = g.Add(1, 2, 1)
	if err != nil {
		t.Error("Did not expected error")
	}

	err = g.Add(1, 2, 0)
	if err != nil {
		t.Error("Did not expected error")
	}

	err = g.Add(1, 2, 0.5)
	if err != nil {
		t.Error("Did not expected error")
	}
}

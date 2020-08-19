// Package trustGraph is based on EigenTrust
// http://nlp.stanford.edu/pubs/eigentrust.pdf
package trustgraph

import (
	"errors"
)

// Group represents a group of peers. Peers need to be given unique, int IDs.
// Certainty represents the threshold of RMS change at which the algorithm will
// escape. Max is the maximum number of loos the algorithm will perform before
// escaping (regardless of certainty). These default to 0.001 and 200
// respectivly and generally don't need to be changed.
type Group struct {
	trustGrid    map[int]map[int]float32
	initialTrust map[int]float32
	Certainty    float32
	Max          int
	Alpha        float32
}

// NewGroup is the constructor for Group.
func NewGroup() Group {
	return Group{
		trustGrid:    map[int]map[int]float32{},
		initialTrust: map[int]float32{},
		Certainty:    0.001,
		Max:          200,
		Alpha:        1,
	}
}

// Add will add or override a trust relationship. The first arg is the peer who
// is extending trust, the second arg is the peer being trusted (by the peer
// in the first arg). The 3rd arg is the amount of trust, which must be
func (g Group) Add(truster, trusted int, amount float32) (err error) {
	err = float32InRange(amount)
	if err == nil {
		a, ok := g.trustGrid[truster]
		if !ok {
			a = map[int]float32{}
			g.trustGrid[truster] = a
		}
		a[trusted] = amount
	}
	return
}

// InitialTrust sets the vaulues used to seed the calculation as well as the
// corrective factor used by Alpha.
func (g Group) InitialTrust(trusted int, amount float32) (err error) {
	err = float32InRange(amount)
	if err == nil {
		g.initialTrust[trusted] = amount
	}
	return
}

// float32InRange is a helper to check that a value is 0.0 <= x <= 1.0
func float32InRange(x float32) error {
	if x < 0 {
		return errors.New("Trust amount cannot be less than 0")
	}
	if x > 1 {
		return errors.New("Trust amount cannot be greater than 1")
	}
	return nil
}

// Compute will approximate the trustworthyness of each peer from the
// information known of how much peers trust eachother.
// It wil loop, upto g.Max times or until the average difference between
// iterations is less than g.Certainty.
func (g Group) Compute() map[int]float32 {
	if len(g.initialTrust) == 0 {
		return map[int]float32{}
	}
	t0 := g.initialTrust //trust map for previous iteration

	for i := 0; i < g.Max; i++ {
		t1 := *g.computeIteration(&t0) // trust map for current iteration
		d := avgD(&t0, &t1)
		t0 = t1
		if d < g.Certainty {
			break
		}
	}

	return t0
}

// computeIteration is broken out of Compute to aid comprehension. It is the
// inner loop of Compute. It loops over every value in t (the current trust map)
// and looks up how much trust that peer extends to every other peer. The
// product of the direct trust and indirect trust
func (g Group) computeIteration(t0 *map[int]float32) *map[int]float32 {

	t1 := map[int]float32{}
	for truster, directTrust := range *t0 {
		for trusted, indirectTrust := range g.trustGrid[truster] {
			if trusted != truster {
				t1[trusted] += directTrust * indirectTrust
			}
		}
	}

	// normalize the trust values
	// in the EigenTrust paper, this was not done every step, but I prefer to
	// Not doing it means the diff (d) needs to be normalized in
	// proportion to the values (because they increase with every iteration)
	highestTrust := float32(0)
	for _, v := range t1 {
		if v > highestTrust {
			highestTrust = v
		}
	}
	//Todo handle highestTrust == 0
	for i, v := range t1 {
		t1[i] = (v/highestTrust)*g.Alpha + (1-g.Alpha)*g.initialTrust[i]
	}

	return &t1
}

// abs is helper to take abs of float32
func abs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

// avgD is helper to compare 2 maps of float32s and return the average
// difference between them
func avgD(t0, t1 *map[int]float32) float32 {
	d := float32(0)
	for i, v := range *t1 {
		d += abs(v - (*t0)[i])
	}
	d = d / float32(len(*t0))
	return d
}

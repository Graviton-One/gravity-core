package score

import (
	"math/rand"
	"testing"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"

	"github.com/Gravity-Tech/gravity-core/common/storage"
)

func TestUInt64ToFloat32Score(t *testing.T) {
	generator := rand.New(rand.NewSource(time.Now().Unix()))
	score := uint64(generator.Intn(100))
	v := UInt64ToFloat32Score(score)
	if uint64(v*Accuracy) != score {
		t.Error("invalid convert uint to float score")
	}
}
func TestFloat32ToUInt64Score(t *testing.T) {
	generator := rand.New(rand.NewSource(time.Now().Unix()))
	score := float32(generator.Intn(100)) / 100
	v := Float32ToUInt64Score(score)
	if float32(v)/Accuracy != score {
		t.Error("invalid convert float to uint score")
	}
}

func TestCalculateDropValidator(t *testing.T) {
	consuls := []account.ConsulPubKey{
		account.ConsulPubKey([32]byte{0}),
		account.ConsulPubKey([32]byte{1}),
		account.ConsulPubKey([32]byte{2}),
		account.ConsulPubKey([32]byte{3}),
		account.ConsulPubKey([32]byte{4}),
	}

	initScores := storage.ScoresByConsulMap{
		consuls[0]: Accuracy,
		consuls[1]: Accuracy,
		consuls[2]: Accuracy,
		consuls[3]: Accuracy,
		consuls[4]: Accuracy,
	}

	votes := storage.VoteByConsulMap{
		consuls[0]: []storage.Vote{
			{consuls[1], Accuracy},
			{consuls[2], Accuracy},
			{consuls[3], Accuracy},
			{consuls[4], 0},
		},
		consuls[1]: []storage.Vote{
			{consuls[0], Accuracy},
			{consuls[2], Accuracy},
			{consuls[3], Accuracy},
			{consuls[4], 0},
		},
		consuls[2]: []storage.Vote{
			{consuls[0], Accuracy},
			{consuls[1], Accuracy},
			{consuls[3], Accuracy},
			{consuls[4], 0},
		},
		consuls[3]: []storage.Vote{
			{consuls[0], Accuracy},
			{consuls[1], Accuracy},
			{consuls[2], Accuracy},
			{consuls[4], 0},
		},
		consuls[4]: []storage.Vote{
			{consuls[0], Accuracy},
			{consuls[1], Accuracy},
			{consuls[2], Accuracy},
			{consuls[3], Accuracy},
		},
	}

	score, err := Calculate(initScores, votes)
	if err != nil {
		t.Error(err)
	}

	if score[consuls[0]] != Accuracy {
		t.Error("invalid consul #1 score")
	} else if score[consuls[1]] != Accuracy {
		t.Error("invalid consul #2 score")
	} else if score[consuls[2]] != Accuracy {
		t.Error("invalid consul #3 score")
	} else if score[consuls[3]] != Accuracy {
		t.Error("invalid consul #4 score")
	} else if score[consuls[4]] != 0 {
		t.Error("invalid consul #5 score")
	}
}

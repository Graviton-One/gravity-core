package scheduler

import (
	"encoding/binary"
	"encoding/json"
	score_calculator "gravity-hub/common/api/score-calculator"
	"gravity-hub/common/keys"
	"gravity-hub/common/score"
	"gravity-hub/score-calculator/models"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/dgraph-io/badger"
)

const (
	CalculateScoreInterval = 20
)

type Scheduler struct {
	scoreClient *score_calculator.Client
}

func New(scoreClient *score_calculator.Client) *Scheduler {
	return &Scheduler{
		scoreClient: scoreClient,
	}
}

func (scheduler *Scheduler) HandleBlock(height int64, txn *badger.Txn) error {
	if height%CalculateScoreInterval == 0 {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(keys.VoteKey)

		values := make(map[string][]byte)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				values[string(k)] = v
				return nil
			})
			if err != nil {
				return err
			}
		}

		voteMap := make(map[string][]models.Vote)
		var actors []models.Actor
		for key, value := range values {
			var voteByValidator []models.Vote

			err := json.Unmarshal(value, &voteByValidator)
			if err != nil {
				return err
			}

			validator := strings.Split(key, keys.Separator)[1]
			voteMap[validator] = voteByValidator
			validatorAddress, err := hexutil.Decode(validator)
			if err != nil {
				return err
			}
			var initScore float32
			item, err := txn.Get([]byte(keys.FormScoreKey(validatorAddress)))
			if err != nil && err != badger.ErrKeyNotFound {
				return err
			} else if err == badger.ErrKeyNotFound {
				initScore = 0

			}
			scoreUInt, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			initScore = score.UInt64ToFloat32Score(binary.BigEndian.Uint64(scoreUInt))
			actors = append(actors, models.Actor{Name: validator, InitScore: initScore})
		}

		rs, err := scheduler.scoreClient.Calculate(actors, voteMap)
		if err != nil {
			return err
		}

		for key, value := range rs.Score {
			initScore := score.Float32ToUInt64Score(value)
			validatorAddress, err := hexutil.Decode(key)
			if err != nil {
				return err
			}
			var scoreBytes []byte
			binary.BigEndian.PutUint64(scoreBytes, initScore)
			err = txn.Set([]byte(keys.FormScoreKey(validatorAddress)), scoreBytes)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

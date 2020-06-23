package scheduler

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gravity-hub/common/contracts"
	"gravity-hub/common/keys"
	"gravity-hub/common/score"
	score_calculator "gravity-hub/score-calculator"
	"gravity-hub/score-calculator/models"
	"math/big"
	"strings"
	"time"

	"github.com/mr-tron/base58"
	wavesCrypto "github.com/wavesplatform/gowaves/pkg/crypto"
	"github.com/wavesplatform/gowaves/pkg/proto"

	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/wavesplatform/gowaves/pkg/client"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/dgraph-io/badger"
)

const (
	CalculateScoreInterval = 20
)

type Scheduler struct {
	wavesPrivKey     []byte
	validatorPrivKey ed25519.PrivKeyEd25519
	validatorPubKey  ed25519.PubKeyEd25519
	validatorsPubKey []ed25519.PubKeyEd25519
	wavesClient      *client.Client
	ethereum         *ethclient.Client
	gravityContract  *contracts.Gravity
	gravityWaves     string
	gravityEthereum  string
	validatorCount   int
}

func New(wavesClient *client.Client, ethereum *ethclient.Client, gravityWaves string, gravityEthereum string) (*Scheduler, error) {
	address := common.HexToAddress(gravityEthereum)
	gravityContract, err := contracts.NewGravity(address, nil)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		wavesClient:     wavesClient,
		ethereum:        ethereum,
		gravityEthereum: gravityEthereum,
		gravityWaves:    gravityWaves,
		gravityContract: gravityContract,
	}, nil
}

func (scheduler *Scheduler) HandleBlock(height int64, txn *badger.Txn) error {
	if height%CalculateScoreInterval == 0 {
		votes, err := scheduler.getVotes(txn)
		if err != nil {
			return err
		}

		scoresValue, err := scheduler.calculate(votes, txn)
		if err != nil {
			return err
		}

		err = scheduler.saveScores(scoresValue, txn)
		if err != nil {
			return err
		}

		for key, value := range scoresValue {
			scoreValue := score.Float32ToUInt64Score(value)
			pubKey, err := hexutil.Decode(key)
			if err != nil {
				return err
			}

			if scoreValue <= 0 {
				err = scheduler.dropValidator(pubKey, txn)
				if err != nil {
					return err
				}
			}
		}
	}

	err := scheduler.sendValidatorScoresEthereum(txn)
	if err != nil {
		return err
	}

	err = scheduler.sendValidatorScoresWaves(txn)
	if err != nil {
		return err
	}

	return nil
}

func (scheduler *Scheduler) sendValidatorScoresEthereum(txn *badger.Txn) error {
	var addresses []common.Address
	var newScores []*big.Int

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, value := range scheduler.validatorsPubKey {
		signKey := keys.FormSignScoreValidatorsKey(value.Bytes())
		item, err := txn.Get([]byte(signKey))
		if err != nil {
			return err
		}
		sign, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		r = append(r, bytes32R)

		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])
		s = append(s, bytes32S)
		v = append(v, sign[64:][0])
	}

	tx, err := scheduler.gravityContract.UpdateScores(nil, addresses, newScores, v, r, s)
	if err != nil {
		return err
	}

	fmt.Printf("Finilize tx id: %s", tx.Hash().String())
	return nil
}

func (scheduler *Scheduler) sendValidatorScoresWaves(txn *badger.Txn) error {
	var addresses []common.Address
	var newScores []uint64

	var addressesString []string
	var newScoresString []string
	for _, value := range addresses {
		addressesString = append(addressesString, value.String())
	}
	for _, value := range newScores {
		newScoresString = append(newScoresString, fmt.Sprintf("%d", value))
	}
	funcArgs := new(proto.Arguments)
	var signs []string
	for _, value := range scheduler.validatorsPubKey {
		signKey := keys.FormSignScoreValidatorsKey(value.Bytes())
		item, err := txn.Get([]byte(signKey))
		if err != nil {
			return err
		}
		sign, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}
		signs = append(signs, base58.Encode(sign))
	}

	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(signs, ","),
	})
	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(signs, ","),
	})
	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(signs, ","),
	})

	secret, err := wavesCrypto.NewSecretKeyFromBytes(scheduler.wavesPrivKey)

	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		return err
	}
	contract, err := proto.NewRecipientFromString(scheduler.gravityWaves)
	if err != nil {
		return err
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        wavesCrypto.GeneratePublicKey(secret),
		ChainID:         'R', //TODO config
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name:      "setSortedOracles",
			Arguments: *funcArgs,
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: client.NewTimestampFromTime(time.Now()),
	}

	//TODO chainid to cinfig
	err = tx.Sign('R', secret)
	if err != nil {
		return err
	}

	_, err = scheduler.wavesClient.Transactions.Broadcast(nil, tx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx waves finilize: %s \n", tx.ID)
	return nil

}

func (scheduler *Scheduler) dropValidator(pubKey []byte, txn *badger.Txn) error {
	key := []byte(keys.FormNebulaeByValidatorKey(pubKey))
	var nebulae []string
	item, err := txn.Get(key)
	if err != nil {
		return err
	}

	value, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}

	err = json.Unmarshal(value, &nebulae)
	if err != nil {
		return err
	}

	for _, v := range nebulae {
		nebulaAddress, err := hexutil.Decode(v)
		if err != nil {
			return err
		}

		key := []byte(keys.FormValidatorKey(nebulaAddress, pubKey))
		err = txn.Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (scheduler *Scheduler) getVotes(txn *badger.Txn) (map[string][]byte, error) {
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(keys.VoteKey)
	votes := make(map[string][]byte)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		err := item.Value(func(v []byte) error {
			votes[string(k)] = v
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return votes, nil
}

func (scheduler *Scheduler) calculate(votes map[string][]byte, txn *badger.Txn) (map[string]float32, error) {
	voteMap := make(map[string][]models.Vote)
	var actors []models.Actor
	for key, value := range votes {
		var voteByValidator []models.Vote

		err := json.Unmarshal(value, &voteByValidator)
		if err != nil {
			return nil, err
		}

		validator := strings.Split(key, keys.Separator)[1]
		voteMap[validator] = voteByValidator
		validatorAddress, err := hexutil.Decode(validator)
		if err != nil {
			return nil, err
		}
		var initScore float32
		item, err := txn.Get([]byte(keys.FormScoreKey(validatorAddress)))
		if err != nil && err != badger.ErrKeyNotFound {
			return nil, err
		} else if err == badger.ErrKeyNotFound {
			initScore = 0
		}
		scoreUInt, err := item.ValueCopy(nil)
		if err != nil {
			return nil, err
		}

		initScore = score.UInt64ToFloat32Score(binary.BigEndian.Uint64(scoreUInt))
		actors = append(actors, models.Actor{Name: validator, InitScore: initScore})
	}

	scores, err := score_calculator.Calculate(actors, voteMap)
	if err != nil {
		return nil, err
	}

	return scores, nil
}

func (scheduler *Scheduler) saveScores(scoresValue map[string]float32, txn *badger.Txn) error {
	for key, value := range scoresValue {
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

	return nil
}

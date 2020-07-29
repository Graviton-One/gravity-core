package scheduler

import (
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/account"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/contracts"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/keys"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/score"
	"github.com/Gravity-Hub-Org/proof-of-concept/common/transactions"
	"github.com/Gravity-Hub-Org/proof-of-concept/gh-node/api/gravity"
	"github.com/Gravity-Hub-Org/proof-of-concept/gh-node/helpers"
	score_calculator "github.com/Gravity-Hub-Org/proof-of-concept/score-calculator"
	"github.com/Gravity-Hub-Org/proof-of-concept/score-calculator/models"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/btcsuite/btcutil/base58"

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
	OraclesCount           = 5
)

type Scores struct {
	Validator string
	Value     uint64
}
type LedgerValidator struct {
	PrivKey ed25519.PrivKeyEd25519
	PubKey  ed25519.PubKeyEd25519
}
type WavesConf struct {
	PrivKey []byte
	Client  *client.Client
	ChainId byte
	Helper  helpers.Node
}
type EthereumConf struct {
	PrivKey      *ecdsa.PrivateKey
	PrivKeyBytes []byte
	Client       *ethclient.Client
}
type Scheduler struct {
	Ledger          *LedgerValidator
	Waves           *WavesConf
	Ethereum        *EthereumConf
	GhNode          string
	gravityEthereum *contracts.Gravity
	gravityWaves    string
	ctx             context.Context
	nebulae         map[account.ChainType][][]byte
}

func New(waves *WavesConf, ethereum *EthereumConf, ghNode string, ctx context.Context, ledger *LedgerValidator, nebulae map[account.ChainType][][]byte, gravityWaves string, gravityEthereum string) (*Scheduler, error) {
	address := common.HexToAddress(gravityEthereum)

	gravityContract, err := contracts.NewGravity(address, ethereum.Client)
	if err != nil {
		return nil, err
	}

	return &Scheduler{
		Ledger:          ledger,
		Waves:           waves,
		Ethereum:        ethereum,
		ctx:             ctx,
		GhNode:          strings.Replace(ghNode, "tcp", "http", 1),
		gravityEthereum: gravityContract,
		gravityWaves:    gravityWaves,
		nebulae:         nebulae,
	}, nil
}

func (scheduler *Scheduler) HandleBlock(height int64, txn *badger.Txn) error {
	go scheduler.setPrivKeys(txn)
	roundId := height / CalculateScoreInterval
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
	} else if height%CalculateScoreInterval < CalculateScoreInterval/2 {
		err := scheduler.signResult(roundId, OraclesCount, scheduler.nebulae[account.Waves], account.Waves, txn)
		if err != nil {
			return err
		}

		err = scheduler.signResult(roundId, OraclesCount, scheduler.nebulae[account.Ethereum], account.Ethereum, txn)
		if err != nil {
			return err
		}
	} else if height%CalculateScoreInterval > CalculateScoreInterval/2 {
		err := scheduler.sendConsulsWaves(roundId, txn)
		if err != nil {
			return err
		}

		err = scheduler.sendConsulsEthereum(roundId, txn)
		if err != nil {
			return err
		}

		for _, v := range scheduler.nebulae[account.Waves] {
		err := scheduler.sendOraclesWaves(v, roundId, txn)
			if err != nil {
				continue
			}
		}

		for _, v := range scheduler.nebulae[account.Ethereum] {
			err := scheduler.sendOraclesEthereum(v, roundId, txn)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (scheduler *Scheduler) setPrivKeys(txn *badger.Txn) error {
	item, err := txn.Get([]byte(keys.FormOraclesByValidatorKey(scheduler.Ledger.PubKey[:])))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	var oracles map[account.ChainType][]byte
	if item != nil {
		b, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &oracles)
		if err != nil {
			return err
		}
	}

	var tx *transactions.Transaction
	if _, ok := oracles[account.Waves]; !ok {
		args := []byte{byte(account.Waves)}
		s, err := wavesCrypto.NewSecretKeyFromBytes(scheduler.Waves.PrivKey)
		if err != nil {
			return err
		}

		pubKey := wavesCrypto.GeneratePublicKey(s)
		tx, err = transactions.New(scheduler.Ledger.PubKey[:], transactions.AddOracle, account.Waves, scheduler.Ledger.PrivKey, append(args, pubKey.Bytes()...))
		if err != nil {
			return err
		}
	}
	if tx != nil {
		client, err := gravity.NewClient(scheduler.GhNode)
		if err != nil {
			panic(err)
		}
		_, err = client.HttpClient.NetInfo()
		if err != nil {
			return nil
		}

		err = client.SendTx(tx)
		if err != nil {
			return err
		}
	}

	tx = nil
	if _, ok := oracles[account.Ethereum]; !ok {
		args := []byte{byte(account.Ethereum)}

		tx, err = transactions.New(scheduler.Ledger.PubKey[:], transactions.AddOracle, account.Ethereum,
			scheduler.Ledger.PrivKey, append(args, crypto.PubkeyToAddress(scheduler.Ethereum.PrivKey.PublicKey).Bytes()...))
		if err != nil {
			return err
		}
	}
	if tx != nil {
		client, err := gravity.NewClient(scheduler.GhNode)
		if err != nil {
			panic(err)
		}
		_, err = client.HttpClient.NetInfo()
		if err != nil {
			return nil
		}

		err = client.SendTx(tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (scheduler *Scheduler) signResult(roundId int64, validatorCount int, nebulaeIds [][]byte, chainType account.ChainType, txn *badger.Txn) error {
	consulsByRoundKey := keys.FormConsulsSignKey(scheduler.Ledger.PubKey[:], chainType, roundId)
	_, err := txn.Get([]byte(consulsByRoundKey))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}

	var scores []Scores
	prefix := []byte(keys.ScoreKey)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		err := item.Value(func(v []byte) error {
			validator := strings.Split(string(k), keys.Separator)[1]

			scores = append(scores, Scores{
				Value:     binary.BigEndian.Uint64(v),
				Validator: validator,
			})
			return nil
		})
		if err != nil {
			return err
		}
	}
	sort.SliceStable(scores, func(i, j int) bool {
		return scores[i].Value > scores[j].Value
	})

	if err == badger.ErrKeyNotFound {
		var newConsuls []Scores
		for i := 0; i < len(scores); i++ {
			newConsuls = append(newConsuls, scores[i])
			if len(newConsuls) >= validatorCount {
				break
			}
		}

		key := []byte(keys.FormConsulsKey())
		b, err := json.Marshal(newConsuls)
		if err != nil {
			return err
		}
		err = txn.Set(key, b)
		if err != nil {
			return err
		}

		prevConsuls, err := txn.Get(key)
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}

		if prevConsuls != nil {
			b, err := prevConsuls.ValueCopy(nil)
			key = []byte(keys.FormPrevConsulsKey())
			err = txn.Set(key, b)
			if err != nil {
				return err
			}
		}

		var sign []byte
		switch chainType {
		case account.Ethereum:
			var validators []common.Address
			for _, v := range newConsuls {
				b, err := hexutil.Decode(v.Validator)
				if err != nil {
					return err
				}

				item, err := txn.Get([]byte(keys.FormOraclesByValidatorKey(b)))
				if err != nil && err != badger.ErrKeyNotFound {
					return err
				}

				if item == nil {
					continue
				}
				b, err = item.ValueCopy(nil)
				if err != nil {
					return err
				}

				var oracles map[account.ChainType][]byte
				err = json.Unmarshal(b, &oracles)
				if err != nil {
					return err
				}

				validators = append(validators, common.BytesToAddress(oracles[account.Ethereum]))
			}

			hash, err := scheduler.gravityEthereum.HashNewConsuls(nil, validators)
			if err != nil {
				return err
			}
			sign, err = account.SignWithTCPriv(scheduler.Ethereum.PrivKeyBytes, hash[:], account.Ethereum)
			if err != nil {
				return err
			}
		case account.Waves:
			var validators []string
			for _, v := range newConsuls {
				b, err := hexutil.Decode(v.Validator)
				if err != nil {
					return err
				}

				item, err := txn.Get([]byte(keys.FormOraclesByValidatorKey(b)))
				if err != nil && err != badger.ErrKeyNotFound {
					return err
				}

				if item == nil {
					continue
				}

				b, err = item.ValueCopy(nil)
				if err != nil {
					return err
				}

				var oracles map[account.ChainType][]byte
				err = json.Unmarshal(b, &oracles)
				if err != nil {
					return err
				}

				validators = append(validators, base58.Encode(oracles[account.Waves]))
			}
			sign, err = account.SignWithTCPriv(scheduler.Waves.PrivKey, []byte(strings.Join(validators, ",")), account.Waves)
			if err != nil {
				return err
			}
		}

		err = txn.Set([]byte(consulsByRoundKey), sign)
		if err != nil {
			return err
		}

	}

	for _, nebulaId := range nebulaeIds {
		oraclesByRoundKey := keys.FormOraclesSignNebulaKey(scheduler.Ledger.PubKey[:], nebulaId, roundId)
		_, err := txn.Get([]byte(oraclesByRoundKey))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		} else if err == nil {
			continue
		}

		key := []byte(keys.FormOraclesByNebulaKey(nebulaId))
		item, err := txn.Get(key)
		if err != nil {
			continue
		}

		oraclesByNebula := make(map[string]string)
		if item != nil {
			b, err := item.ValueCopy(nil)
			if err != nil {
				continue
			}
			err = json.Unmarshal(b, &oraclesByNebula)
			if err != nil {
				continue
			}
		}

		oraclesByValidators := make(map[string]string)
		//Revert map
		for k, v := range oraclesByNebula {
			oraclesByValidators[v] = k
		}

		var newOracles []string
		newOraclesMap := make(map[string]string)
		for i := 0; i < len(scores); i++ {
			v, ok := oraclesByValidators[scores[i].Validator]
			if !ok {
				continue
			}

			newOracles = append(newOracles, v)
			newOraclesMap[v] = scores[i].Validator
			if len(newOracles) >= OraclesCount {
				break
			}
		}

		b, err := json.Marshal(&newOraclesMap)
		if err != nil {
			continue
		}

		err = txn.Set([]byte(keys.FormBftOraclesByNebulaKey(nebulaId)), b)
		if err != nil {
			continue
		}

		var sign []byte
		switch chainType {
		case account.Ethereum:
			nebula, err := contracts.NewNebula(common.BytesToAddress(nebulaId), scheduler.Ethereum.Client)
			if err != nil {
				return err
			}
			var validators []common.Address
			for _, v := range newOracles {
				validators = append(validators, common.HexToAddress(v))
			}
			hash, err := nebula.HashNewOracles(nil, validators)
			if err != nil {
				return err
			}
			sign, err = account.SignWithTCPriv(scheduler.Ethereum.PrivKeyBytes, hash[:], account.Ethereum)
			if err != nil {
				return err
			}
		case account.Waves:
			sign, err = account.SignWithTCPriv(scheduler.Waves.PrivKey, []byte(strings.Join(newOracles, ",")), account.Waves)
			if err != nil {
				return err
			}
		}

		err = txn.Set([]byte(oraclesByRoundKey), sign)
		if err != nil {
			return err
		}
	}

	return nil
}

func (scheduler *Scheduler) sendConsulsWaves(round int64, txn *badger.Txn) error {
	state, err := scheduler.Waves.Helper.GetStateByAddressAndKey(scheduler.gravityWaves, "last_round_"+fmt.Sprintf("%d", round))
	if err != nil {
		return err
	}
	if state != nil {
		return nil
	}

	item, err := txn.Get([]byte(keys.FormPrevConsulsKey()))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if err == badger.ErrKeyNotFound {
		return nil
	}
	b, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var consuls []Scores
	err = json.Unmarshal(b, &consuls)
	if err != nil {
		return err
	}

	isOneFound := false

	var signs []string
	for _, v := range consuls {
		validator, err := hexutil.Decode(v.Validator)
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte(keys.FormConsulsSignKey(validator, account.Waves, round)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if item == nil {
			signs = append(signs, base58.Encode([]byte{0}))
			continue
		}

		sign, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		isOneFound = true
		signs = append(signs, base58.Encode(sign))
	}

	var newConsulsString []string
	item, err = txn.Get([]byte(keys.FormConsulsKey()))
	if err != nil {
		return err
	}
	b, err = item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var newConsuls []Scores
	err = json.Unmarshal(b, &newConsuls)
	if err != nil {
		return err
	}

	for _, v := range newConsuls {
		validator, err := hexutil.Decode(v.Validator)
		if err != nil {
			return err
		}

		item, err := txn.Get([]byte(keys.FormOraclesByValidatorKey(validator)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if item == nil {
			newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
			continue
		}

		oracles := make(map[account.ChainType][]byte)
		if err != badger.ErrKeyNotFound {
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			err = json.Unmarshal(value, &oracles)
			if err != nil {
				return err
			}
		}

		newConsulsString = append(newConsulsString, base58.Encode(oracles[account.Waves]))
	}

	emptyCount := OraclesCount - len(signs)
	for i := 0; i < emptyCount; i++ {
		signs = append(signs, base58.Encode([]byte{0}))
	}

	emptyCount = OraclesCount - len(newConsulsString)
	for i := 0; i < emptyCount; i++ {
		newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
	}
	funcArgs := new(proto.Arguments)

	if !isOneFound {
		return nil
	}
	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(newConsulsString, ","),
	})
	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(signs, ","),
	})
	funcArgs.Append(proto.IntegerArgument{
		Value: round,
	})

	secret, err := wavesCrypto.NewSecretKeyFromBytes(scheduler.Waves.PrivKey)
	if err != nil {
		return err
	}
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
		ChainID:         scheduler.Waves.ChainId,
		ScriptRecipient: contract,
		FunctionCall: proto.FunctionCall{
			Name:      "setConsuls",
			Arguments: *funcArgs,
		},
		Payments:  nil,
		FeeAsset:  *asset,
		Fee:       500000,
		Timestamp: client.NewTimestampFromTime(time.Now()),
	}

	err = tx.Sign(scheduler.Waves.ChainId, secret)
	if err != nil {
		return err
	}

	_, err = scheduler.Waves.Client.Transactions.Broadcast(scheduler.ctx, tx)
	if err != nil {
		return err
	}

	err = <-scheduler.Waves.Helper.WaitTx(tx.ID.String())
	if err != nil {
		return err
	}

	fmt.Printf("Tx waves consuls update: %s \n", tx.ID)
	return nil

}
func (scheduler *Scheduler) sendConsulsEthereum(round int64, txn *badger.Txn) error {
	lastRound, err := scheduler.gravityEthereum.Rounds(nil, big.NewInt(round))
	if err != nil {
		return err
	}

	if lastRound {
		return nil
	}

	item, err := txn.Get([]byte(keys.FormPrevConsulsKey()))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if err == badger.ErrKeyNotFound {
		return nil
	}
	b, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var consuls []Scores
	err = json.Unmarshal(b, &consuls)
	if err != nil {
		return err
	}

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, value := range consuls {
		validator, err := hexutil.Decode(value.Validator)
		if err != nil {
			return err
		}

		item, err := txn.Get([]byte(keys.FormOraclesByValidatorKey(validator)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if item == nil {
			continue
		}
		oracles := make(map[account.ChainType][]byte)
		if err != badger.ErrKeyNotFound {
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			err = json.Unmarshal(value, &oracles)
			if err != nil {
				return err
			}
		}

		item, err = txn.Get([]byte(keys.FormConsulsSignKey(validator, account.Ethereum, round)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if item == nil {
			r = append(r, [32]byte{})
			s = append(s, [32]byte{})
			v = append(v, 0)
			continue
		}

		sign, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])

		r = append(r, bytes32R)
		s = append(s, bytes32S)
		v = append(v, sign[64:][0]+27)
	}

	var consulsAddress []common.Address
	item, err = txn.Get([]byte(keys.FormConsulsKey()))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if err == badger.ErrKeyNotFound {
		return nil
	}
	b, err = item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var newConsuls []Scores
	err = json.Unmarshal(b, &newConsuls)
	if err != nil {
		return err
	}

	for _, value := range newConsuls {
		validator, err := hexutil.Decode(value.Validator)
		if err != nil {
			return err
		}

		item, err := txn.Get([]byte(keys.FormOraclesByValidatorKey(validator)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if item == nil {
			continue
		}
		oracles := make(map[account.ChainType][]byte)
		if err != badger.ErrKeyNotFound {
			value, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			err = json.Unmarshal(value, &oracles)
			if err != nil {
				return err
			}
		}

		consulsAddress = append(consulsAddress, common.BytesToAddress(oracles[account.Ethereum]))
	}

	tx, err := scheduler.gravityEthereum.UpdateConsuls(bind.NewKeyedTransactor(scheduler.Ethereum.PrivKey), consulsAddress, v, r, s, big.NewInt(round))
	if err != nil {
		return nil
	}
	fmt.Printf("Tx ethereum consuls update: %s \n", tx.Hash().Hex())
	return nil

}

func (scheduler *Scheduler) sendOraclesWaves(nebulaId []byte, round int64, txn *badger.Txn) error {
	contractAddress := base58.Encode(nebulaId)
	state, err := scheduler.Waves.Helper.GetStateByAddressAndKey(contractAddress, "last_round_"+fmt.Sprintf("%d", round))
	if err != nil {
		return err
	}

	if state != nil {
		return nil
	}

	var newOracles []string
	item, err := txn.Get([]byte(keys.FormOraclesByNebulaKey(nebulaId)))
	if err != nil {
		return err
	}
	b, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var oracles map[string]string
	err = json.Unmarshal(b, &oracles)
	if err != nil {
		return err
	}

	for k, _ := range oracles {
		v, err := hexutil.Decode(k)
		if err != nil {
			return err
		}

		newOracles = append(newOracles, base58.Encode(v))
	}

	emptyCount := OraclesCount - len(newOracles)
	for i := 0; i < emptyCount; i++ {
		newOracles = append(newOracles, base58.Encode([]byte{0}))
	}

	item, err = txn.Get([]byte(keys.FormPrevConsulsKey()))
	if err != nil && err != badger.ErrKeyNotFound {
		return err
	}
	if err == badger.ErrKeyNotFound {
		return nil
	}

	b, err = item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var consuls []Scores
	err = json.Unmarshal(b, &consuls)
	if err != nil {
		return err
	}

	funcArgs := new(proto.Arguments)
	var signs []string

	for _, v := range consuls {
		address, err := hexutil.Decode(v.Validator)
		if err != nil {
			return err
		}

		item, err := txn.Get([]byte(keys.FormOraclesSignNebulaKey(address, nebulaId, round)))
		if err != nil {
			signs = append(signs, base58.Encode([]byte{0}))
			continue
		}

		sign, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		signs = append(signs, base58.Encode(sign))
	}

	emptyCount = OraclesCount - len(signs)
	for i := 0; i < emptyCount; i++ {
		signs = append(signs, base58.Encode([]byte{0}))
	}

	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(newOracles, ","),
	})
	funcArgs.Append(proto.StringArgument{
		Value: strings.Join(signs, ","),
	})
	funcArgs.Append(proto.IntegerArgument{
		Value: round,
	})

	secret, err := wavesCrypto.NewSecretKeyFromBytes(scheduler.Waves.PrivKey)
	asset, err := proto.NewOptionalAssetFromString("WAVES")
	if err != nil {
		return err
	}

	contract, err := proto.NewRecipientFromString(contractAddress)
	if err != nil {
		return err
	}

	tx := &proto.InvokeScriptWithProofs{
		Type:            proto.InvokeScriptTransaction,
		Version:         1,
		SenderPK:        wavesCrypto.GeneratePublicKey(secret),
		ChainID:         scheduler.Waves.ChainId,
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

	err = tx.Sign(scheduler.Waves.ChainId, secret)
	if err != nil {
		return err
	}

	_, err = scheduler.Waves.Client.Transactions.Broadcast(scheduler.ctx, tx)
	if err != nil {
		return err
	}

	fmt.Printf("Tx waves nebula (%s) oracles update: %s \n", contractAddress, tx.ID.String())
	return nil

}
func (scheduler *Scheduler) sendOraclesEthereum(nebulaId []byte, round int64, txn *badger.Txn) error {
	nebula, err := contracts.NewNebula(common.BytesToAddress(nebulaId), scheduler.Ethereum.Client)
	if err != nil {
		return err
	}

	lastRound, err := nebula.Rounds(nil, big.NewInt(round))
	if err != nil {
		return err
	}

	if lastRound {
		return nil
	}

	item, err := txn.Get([]byte(keys.FormConsulsKey()))
	if err != nil {
		return err
	}
	b, err := item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var consuls []Scores
	err = json.Unmarshal(b, &consuls)
	if err != nil {
		return err
	}

	item, err = txn.Get([]byte(keys.FormOraclesByNebulaKey(nebulaId)))
	if err != nil {
		return err
	}
	b, err = item.ValueCopy(nil)
	if err != nil {
		return err
	}

	var oracles map[string]string
	err = json.Unmarshal(b, &oracles)
	if err != nil {
		return err
	}

	var oraclesAddresses []common.Address
	for k, _ := range oracles {
		oraclesAddresses = append(oraclesAddresses, common.HexToAddress(k))
	}

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, value := range consuls {
		address, err := hexutil.Decode(value.Validator)
		if err != nil {
			return err
		}

		item, err := txn.Get([]byte(keys.FormOraclesSignNebulaKey(address, nebulaId, round)))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		}
		if item == nil {
			r = append(r, [32]byte{})
			s = append(s, [32]byte{})
			v = append(v, 0)
			continue
		}

		sign, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		var bytes32R [32]byte
		copy(bytes32R[:], sign[:32])
		var bytes32S [32]byte
		copy(bytes32S[:], sign[32:64])

		r = append(r, bytes32R)
		s = append(s, bytes32S)
		v = append(v, sign[64:][0]+27)
	}

	tx, err := nebula.UpdateOracles(bind.NewKeyedTransactor(scheduler.Ethereum.PrivKey), oraclesAddresses, v, r, s, big.NewInt(round))
	if err != nil {
		return err
	}

	fmt.Printf("Tx ethereum nebula (%s) oracles update: %s \n", common.BytesToAddress(nebulaId).Hex(), tx.Hash().Hex())
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
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	prefix := []byte(keys.ScoreKey)
	for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
		item := it.Item()
		k := item.Key()
		err := item.Value(func(v []byte) error {
			scoreUInt := v

			initScore := score.UInt64ToFloat32Score(binary.BigEndian.Uint64(scoreUInt))
			actors = append(actors, models.Actor{Name: strings.Split(string(k), keys.Separator)[1], InitScore: initScore})

			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	for key, value := range votes {
		var voteByValidator []models.Vote

		err := json.Unmarshal(value, &voteByValidator)
		if err != nil {
			return nil, err
		}

		validator := strings.Split(key, keys.Separator)[1]
		voteMap[validator] = voteByValidator
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
		var scoreBytes [8]byte
		binary.BigEndian.PutUint64(scoreBytes[:], initScore)
		err = txn.Set([]byte(keys.FormScoreKey(validatorAddress)), scoreBytes[:])
		if err != nil {
			return err
		}
	}

	return nil
}

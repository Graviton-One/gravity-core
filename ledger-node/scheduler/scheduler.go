package scheduler

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/contracts"
	"github.com/Gravity-Tech/gravity-core/common/keys"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/Gravity-Tech/gravity-core/oracle-node/helpers"
	calculator "github.com/Gravity-Tech/gravity-core/score-calculator"

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
	ValidatorCount         = 5
)

type Score struct {
	Validator account.ValidatorPubKey
	Value     uint64
}
type LedgerValidator struct {
	PrivKey ed25519.PrivKeyEd25519
	PubKey  account.ValidatorPubKey
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

func (scheduler *Scheduler) HandleBlock(height int64, store *storage.Storage) error {
	go scheduler.setPrivKeys(store, account.Ethereum) //TODO: refactoring
	go scheduler.setPrivKeys(store, account.Waves)    //TODO: refactoring

	roundId := height / CalculateScoreInterval
	if height%CalculateScoreInterval == 0 {
		newScores, err := scheduler.calculate(store)
		if err != nil {
			return err
		}

		for k, v := range newScores {
			store.SetScore(k, v)
		}

	} else if height%CalculateScoreInterval < CalculateScoreInterval/2 {
		err := scheduler.signResult(roundId, scheduler.nebulae[account.Waves], account.Waves, txn)
		if err != nil {
			return err
		}

		err = scheduler.signResult(roundId, scheduler.nebulae[account.Ethereum], account.Ethereum, txn)
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

func (scheduler *Scheduler) setPrivKeys(storage *storage.Storage, chainType account.ChainType) error {
	//TODO: refactoring
	oracles, err := storage.OraclesByValidator(scheduler.Ledger.PubKey)
	if err != nil {
		return err
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
	} else if _, ok := oracles[account.Ethereum]; !ok {
		args := []byte{byte(account.Ethereum)}

		tx, err = transactions.New(scheduler.Ledger.PubKey[:], transactions.AddOracle, account.Ethereum,
			scheduler.Ledger.PrivKey, append(args, crypto.PubkeyToAddress(scheduler.Ethereum.PrivKey.PublicKey).Bytes()...))
		if err != nil {
			return err
		}
	}

	if tx != nil {
		client, err := client.New(scheduler.GhNode)
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

func (scheduler *Scheduler) signResult(roundId int64, nebulaeIds [][]byte, chainType account.ChainType, store *storage.Storage) error {
	_, err := store.ConsulSign(scheduler.Ledger.PubKey, chainType, roundId)
	if err == badger.ErrKeyNotFound {
		return nil
	} else if err != nil {
		return err
	}

	scores, err := store.Scores()
	if err != nil {
		return err
	}

	var sortedScores []storage.Consul
	for k, v := range scores {
		sortedScores = append(sortedScores, storage.Consul{
			Validator: k,
			Value:     v,
		})
	}

	sort.SliceStable(sortedScores, func(i, j int) bool {
		return sortedScores[i].Value > sortedScores[j].Value
	})

	var sign []byte
	var newConsuls []storage.Consul
	for _, v := range sortedScores {
		newConsuls = append(newConsuls, v)
		if len(newConsuls) >= ValidatorCount {
			break
		}
	}

	switch chainType {
	case account.Ethereum:
		var validators []common.Address
		for _, v := range newConsuls {
			oracles, err := store.OraclesByValidator(v.Validator)
			if err != nil {
				return err
			}

			oraclePubKey := oracles[account.Ethereum]
			validators = append(validators, common.BytesToAddress(oraclePubKey.ToBytes(account.Ethereum)))
		}

		hash, err := scheduler.gravityEthereum.HashNewConsuls(nil, validators)
		if err != nil {
			return err
		}

		sign, err = account.SignWithTC(scheduler.Ethereum.PrivKeyBytes, hash[:], account.Ethereum)
		if err != nil {
			return err
		}
	case account.Waves:
		var validators []string
		for _, v := range newConsuls {
			oracles, err := store.OraclesByValidator(v.Validator)
			if err != nil {
				return err
			}

			oraclePubKey := oracles[account.Waves]
			validators = append(validators, base58.Encode(oraclePubKey.ToBytes(account.Waves)))
		}

		sign, err = account.SignWithTC(scheduler.Waves.PrivKey, []byte(strings.Join(validators, ",")), account.Waves)
		if err != nil {
			return err
		}
	}

	err = store.SetConsuls(newConsuls)
	if err != nil {
		return err
	}
	store.Set

	for _, nebulaId := range nebulaeIds {
		oraclesByRoundKey := keys.FormOraclesSignNebulaKey(scheduler.Ledger.PubKey[:], nebulaId, roundId)
		_, err := txn.Get([]byte(oraclesByRoundKey))
		if err != nil && err != badger.ErrKeyNotFound {
			return err
		} else if err == nil {
			continue
		}

		oraclesByNebula, err := storage.OraclesByNebula(nebulaId)
		if err != nil {
			continue
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

func (scheduler *Scheduler) calculate(store *storage.Storage) (storage.ScoresByValidatorMap, error) {
	voteMap, err := store.Votes()
	if err != nil {
		return nil, err
	}

	scores, err := store.Scores()
	if err != nil {
		return nil, err
	}

	newScores, err := calculator.Calculate(scores, voteMap)
	if err != nil {
		return nil, err
	}

	return newScores, nil
}

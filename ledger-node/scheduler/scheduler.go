package scheduler

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"time"

	"github.com/Gravity-Tech/gravity-core/common/account"
	ghClient "github.com/Gravity-Tech/gravity-core/common/client"
	"github.com/Gravity-Tech/gravity-core/common/contracts"
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

	"github.com/dgraph-io/badger"
)

const (
	CalculateScoreInterval = 20
	ValidatorCount         = 5
	OracleCount            = 5
)

type Score struct {
	PubKey account.ConsulPubKey
	Value  uint64
}
type LedgerValidator struct {
	PrivKey ed25519.PrivKeyEd25519
	PubKey  account.ConsulPubKey
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
		if err := scheduler.calculateScores(store); err != nil {
			return err
		}
	} else if height%CalculateScoreInterval < CalculateScoreInterval/2 {
		err := scheduler.signConsulsResult(roundId, account.Waves, store)
		if err != nil {
			return err
		}

		err = scheduler.signConsulsResult(roundId, account.Ethereum, store)
		if err != nil {
			return err
		}

		for _, v := range scheduler.nebulae[account.Waves] {
			err := scheduler.signOracleResultByNebula(roundId, v, account.Waves, store)
			if err != nil {
				continue
			}
		}
		for _, v := range scheduler.nebulae[account.Ethereum] {
			err := scheduler.signOracleResultByNebula(roundId, v, account.Ethereum, store)
			if err != nil {
				continue
			}
		}
	} else if height%CalculateScoreInterval > CalculateScoreInterval/2 {
		err := scheduler.sendConsulsWaves(roundId, store)
		if err != nil {
			return err
		}

		err = scheduler.sendConsulsEthereum(roundId, store)
		if err != nil {
			return err
		}

		for _, v := range scheduler.nebulae[account.Waves] {
			err := scheduler.sendOraclesWaves(v, roundId, store)
			if err != nil {
				continue
			}
		}

		for _, v := range scheduler.nebulae[account.Ethereum] {
			err := scheduler.sendOraclesEthereum(v, roundId, store)
			if err != nil {
				continue
			}
		}
	}

	return nil
}

func (scheduler *Scheduler) setPrivKeys(storage *storage.Storage, chainType account.ChainType) error {
	oracles, err := storage.OraclesByConsul(scheduler.Ledger.PubKey)
	if err != nil {
		return err
	}

	var tx *transactions.Transaction
	switch chainType {
	case account.Ethereum:
		if _, ok := oracles[account.Ethereum]; !ok {
			args := []transactions.Args{
				{
					Value: account.Ethereum,
				},
				{
					Value: crypto.PubkeyToAddress(scheduler.Ethereum.PrivKey.PublicKey).Bytes(),
				},
			}
			tx, err = transactions.New(scheduler.Ledger.PubKey, transactions.AddOracle, scheduler.Ledger.PrivKey, args)
			if err != nil {
				return err
			}
		}
	case account.Waves:
		if _, ok := oracles[account.Waves]; !ok {
			s, err := wavesCrypto.NewSecretKeyFromBytes(scheduler.Waves.PrivKey)
			if err != nil {
				return err
			}

			pubKey := wavesCrypto.GeneratePublicKey(s)
			args := []transactions.Args{
				{
					Value: account.Waves,
				},
				{
					Value: pubKey[:],
				},
			}
			tx, err = transactions.New(scheduler.Ledger.PubKey, transactions.AddOracle, scheduler.Ledger.PrivKey, args)
			if err != nil {
				return err
			}
		}
	}

	if tx != nil {
		ghClient, err := ghClient.New(scheduler.GhNode)
		if err != nil {
			panic(err)
		}
		_, err = ghClient.HttpClient.NetInfo()
		if err != nil {
			return nil
		}

		err = ghClient.SendTx(tx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (scheduler *Scheduler) signConsulsResult(roundId int64, chainType account.ChainType, store *storage.Storage) error {
	_, err := store.SignConsulsResultByConsul(scheduler.Ledger.PubKey, chainType, roundId)
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
			PubKey: k,
			Value:  v,
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
			oracles, err := store.OraclesByConsul(v.PubKey)
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
			oracles, err := store.OraclesByConsul(v.PubKey)
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

	err = store.SetSignConsulsResult(scheduler.Ledger.PubKey, chainType, roundId, sign)
	if err != nil {
		return err
	}

	return nil
}
func (scheduler *Scheduler) signOracleResultByNebula(roundId int64, nebulaId []byte, chainType account.ChainType, store *storage.Storage) error {
	_, err := store.SignOraclesResultByConsul(scheduler.Ledger.PubKey, nebulaId, roundId)
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	} else if err == nil {
		return nil
	}

	oraclesByNebula, err := store.OraclesByNebula(nebulaId)
	if err != nil {
		return err
	}

	lastIndex, err := store.NebulaOraclesIndex()
	if err != nil {
		return err
	}

	var newOracles []account.OraclesPubKey
	var oracles []account.OraclesPubKey
	newOraclesMap := make(storage.OraclesMap)

	for k, _ := range oraclesByNebula {
		oracles = append(oracles, k)
	}

	newIndex := lastIndex + 1
	if newIndex >= uint64(len(oracles)) {
		newIndex = 0
	}

	if newIndex+OracleCount > uint64(len(oracles)) {
		newOracles = oracles[newIndex:]
		newOracles = append(newOracles, newOracles[:OracleCount-len(newOracles)]...)
	} else {
		newOracles = oracles[newIndex : newIndex+OracleCount]
	}

	for _, v := range newOracles {
		newOraclesMap[v] = true
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
			validators = append(validators, common.BytesToAddress(v.ToBytes(account.Ethereum)))
		}
		hash, err := nebula.HashNewOracles(nil, validators)
		if err != nil {
			return err
		}
		sign, err = account.SignWithTC(scheduler.Ethereum.PrivKeyBytes, hash[:], account.Ethereum)
		if err != nil {
			return err
		}
	case account.Waves:
		var oraclesString []string
		for _, v := range newOracles {
			oraclesString = append(oraclesString, base58.Encode(v.ToBytes(account.Waves)))
		}

		sign, err = account.SignWithTC(scheduler.Waves.PrivKey, []byte(strings.Join(oraclesString, ",")), account.Waves)
		if err != nil {
			return err
		}
	}

	err = store.SetSignOraclesResult(scheduler.Ledger.PubKey, nebulaId, roundId, sign)
	if err != nil {
		return err
	}

	err = store.SetBftOraclesByNebula(nebulaId, newOraclesMap)
	if err != nil {
		return err
	}

	return nil
}

func (scheduler *Scheduler) sendConsulsWaves(round int64, store *storage.Storage) error {
	state, err := scheduler.Waves.Helper.GetStateByAddressAndKey(scheduler.gravityWaves, "last_round_"+fmt.Sprintf("%d", round))
	if err != nil {
		return err
	}
	if state != nil {
		return nil
	}

	prevConsuls, err := store.PrevConsuls()
	if err != nil {
		return err
	}

	oneSigFound := false
	var signs []string
	for _, v := range prevConsuls {
		sign, err := store.SignConsulsResultByConsul(v.PubKey, account.Waves, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			signs = append(signs, base58.Encode([]byte{0}))
			continue
		}

		oneSigFound = true
		signs = append(signs, base58.Encode(sign))
	}

	var newConsulsString []string
	newConsuls, err := store.Consuls()
	if err != nil {
		return err
	}

	for _, v := range newConsuls {
		oracles, err := store.OraclesByConsul(v.PubKey)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
			continue
		}

		oraclePubKey := oracles[account.Waves]
		newConsulsString = append(newConsulsString, base58.Encode(oraclePubKey.ToBytes(account.Waves)))
	}

	emptyCount := OracleCount - len(signs)
	for i := 0; i < emptyCount; i++ {
		signs = append(signs, base58.Encode([]byte{0}))
	}

	emptyCount = OracleCount - len(newConsulsString)
	for i := 0; i < emptyCount; i++ {
		newConsulsString = append(newConsulsString, base58.Encode([]byte{0}))
	}
	funcArgs := new(proto.Arguments)

	if !oneSigFound {
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
func (scheduler *Scheduler) sendConsulsEthereum(round int64, store *storage.Storage) error {
	lastRound, err := scheduler.gravityEthereum.Rounds(nil, big.NewInt(round))
	if err != nil {
		return err
	}

	if lastRound {
		return nil
	}

	prevConsuls, err := store.PrevConsuls()
	if err != nil && err != storage.ErrKeyNotFound {
		return err
	}

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, value := range prevConsuls {
		sign, err := store.SignConsulsResultByConsul(value.PubKey, account.Ethereum, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}
		if err == storage.ErrKeyNotFound {
			r = append(r, [32]byte{})
			s = append(s, [32]byte{})
			v = append(v, 0)
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
	newConsuls, err := store.Consuls()
	if err != nil {
		return err
	}

	for _, value := range newConsuls {
		oracles, err := store.OraclesByConsul(value.PubKey)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}

		if err == storage.ErrKeyNotFound {
			continue
		}

		oraclePubKey := oracles[account.Ethereum]
		consulsAddress = append(consulsAddress, common.BytesToAddress(oraclePubKey.ToBytes(account.Ethereum)))
	}

	tx, err := scheduler.gravityEthereum.UpdateConsuls(bind.NewKeyedTransactor(scheduler.Ethereum.PrivKey), consulsAddress, v, r, s, big.NewInt(round))
	if err != nil {
		return nil
	}

	fmt.Printf("Tx ethereum consuls update: %s \n", tx.Hash().Hex())
	return nil

}

func (scheduler *Scheduler) sendOraclesWaves(nebulaId []byte, round int64, store *storage.Storage) error {
	contractAddress := base58.Encode(nebulaId)
	state, err := scheduler.Waves.Helper.GetStateByAddressAndKey(contractAddress, "last_round_"+fmt.Sprintf("%d", round))
	if err != nil {
		return err
	}
	if state != nil {
		return nil
	}

	var newOracles []string
	oracles, err := store.OraclesByNebula(nebulaId)

	for k, _ := range oracles {
		newOracles = append(newOracles, base58.Encode(k.ToBytes(account.Waves)))
	}

	emptyCount := OracleCount - len(newOracles)
	for i := 0; i < emptyCount; i++ {
		newOracles = append(newOracles, base58.Encode([]byte{0}))
	}

	prevConsuls, err := store.PrevConsuls()
	if err != nil {
		return err
	}

	funcArgs := new(proto.Arguments)
	var signs []string

	for _, v := range prevConsuls {
		sign, err := store.SignOraclesResultByConsul(v.PubKey, nebulaId, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}
		if err == storage.ErrKeyNotFound {
			signs = append(signs, base58.Encode([]byte{0}))
			continue
		}

		signs = append(signs, base58.Encode(sign))
	}

	emptyCount = OracleCount - len(signs)
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
func (scheduler *Scheduler) sendOraclesEthereum(nebulaId []byte, round int64, store *storage.Storage) error {
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

	consuls, err := store.Consuls()
	if err != nil {
		return err
	}

	oracles, err := store.OraclesByNebula(nebulaId)
	if err != nil {
		return err
	}

	var oraclesAddresses []common.Address
	for k, _ := range oracles {
		oraclesAddresses = append(oraclesAddresses, common.BytesToAddress(k.ToBytes(account.Ethereum)))
	}

	var r [][32]byte
	var s [][32]byte
	var v []uint8
	for _, value := range consuls {
		sign, err := store.SignOraclesResultByConsul(value.PubKey, nebulaId, round)
		if err != nil && err != storage.ErrKeyNotFound {
			return err
		}
		if err == storage.ErrKeyNotFound {
			r = append(r, [32]byte{})
			s = append(s, [32]byte{})
			v = append(v, 0)
			continue
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

func (scheduler *Scheduler) calculateScores(store *storage.Storage) error {
	voteMap, err := store.Votes()
	if err != nil {
		return err
	}

	scores, err := store.Scores()
	if err != nil {
		return err
	}

	newScores, err := calculator.Calculate(scores, voteMap)
	if err != nil {
		return err
	}

	for k, v := range newScores {
		err := store.SetScore(k, v)
		if err != nil {
			return err
		}

		if v <= 0 {
			oracles, err := store.OraclesByConsul(k)
			if err != nil {
				return err
			}

			for _, oracle := range oracles {
				nebulae, err := store.NebulaeByOracle(oracle)
				if err != nil {
					return err
				}

				for _, nebulaId := range nebulae {
					oracles, err := store.OraclesByNebula(nebulaId)
					if err != nil {
						return err
					}

					delete(oracles, oracle)
					err = store.SetOraclesByNebula(nebulaId, oracles)
					if err != nil {
						return err
					}
				}

				err = store.SetNebulaeByOracle(oracle, []account.NebulaId{})
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

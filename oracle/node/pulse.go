package node

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/hashing"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func (node *Node) invokeCommitTx(data *extractor.Data, tcHeight uint64, pulseId uint64) ([]byte, error) {
	dataBytes := toBytes(data, node.extractor.ExtractorType)
	zap.L().Sugar().Debugf("Extractor data type: %d", node.extractor.ExtractorType)

	commit := hashing.WrappedKeccak256(dataBytes, node.chainType)
	fmt.Printf("Commit: %s - %s \n", hexutil.Encode(dataBytes), hexutil.Encode(commit[:]))

	tx, err := transactions.New(node.validator.pubKey, transactions.Commit, node.validator.privKey)
	if err != nil {
		return nil, err
	}

	tx.AddValues([]transactions.Value{
		transactions.BytesValue{
			Value: node.nebulaId[:],
		},
		transactions.IntValue{
			Value: int64(pulseId),
		},
		transactions.IntValue{
			Value: int64(tcHeight),
		},
		transactions.BytesValue{
			Value: commit,
		},
		transactions.BytesValue{
			Value: node.oraclePubKey[:],
		},
	})

	err = node.gravityClient.SendTx(tx)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Commit txId: %s\n", hexutil.Encode(tx.Id[:]))

	return commit, nil
}

func (node *Node) invokeRevealTx(tcHeight uint64, pulseId uint64, reveal *extractor.Data, commit []byte) error {
	dataBytes := toBytes(reveal, node.extractor.ExtractorType)
	fmt.Printf("Reveal: %s  - %s \n", hexutil.Encode(dataBytes), hexutil.Encode(commit))
	println(base64.StdEncoding.EncodeToString(dataBytes))
	tx, err := transactions.New(node.validator.pubKey, transactions.Reveal, node.validator.privKey)
	if err != nil {
		return err
	}
	tx.AddValues([]transactions.Value{
		transactions.BytesValue{
			Value: commit,
		},
		transactions.BytesValue{
			Value: node.nebulaId[:],
		},
		transactions.IntValue{
			Value: int64(pulseId),
		},
		transactions.IntValue{
			Value: int64(tcHeight),
		},
		transactions.BytesValue{
			Value: dataBytes,
		},
		transactions.BytesValue{
			Value: node.oraclePubKey[:],
		},
		transactions.IntValue{
			Value: int64(node.chainType),
		},
	})

	err = node.gravityClient.SendTx(tx)
	if err != nil {
		return err
	}
	fmt.Printf("Reveal txId: %s\n", hexutil.Encode(tx.Id[:]))

	return nil
}

func (node *Node) signRoundResult(intervalId uint64, pulseId uint64, ctx context.Context) (*extractor.Data, []byte, error) {
	zap.L().Sugar().Debugf("signResults: interval: %d, pulseId: %d", intervalId, pulseId)
	var values []extractor.Data
	zap.L().Sugar().Debugf("gravity Reveals: chaintype: %d, pulseId: %d NebulaId: %s", node.chainType, pulseId, node.nebulaId.ToString(node.chainType))
	bytesValues, err := node.gravityClient.Reveals(node.chainType, node.nebulaId, int64(intervalId), int64(pulseId))
	if err != nil {
		return nil, nil, err
	}
	zap.L().Sugar().Debugf("signResults: value len: %d", len(bytesValues))
	for _, v := range bytesValues {
		zap.L().Sugar().Debugf("signResults: decoding: %s", v)
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			zap.L().Sugar().Error(err.Error())
			continue
		}
		values = append(values, *fromBytes(b, node.extractor.ExtractorType))
	}
	zap.L().Sugar().Debug("Decoded values: ", values)
	if len(values) == 0 {
		return nil, nil, nil
	}
	result, err := node.extractor.Aggregate(values, ctx)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, nil, err
	}

	hash := hashing.WrappedKeccak256(toBytes(result, node.extractor.ExtractorType), node.chainType)

	sign, err := node.adaptor.SignHash(node.nebulaId, intervalId, pulseId, hash)
	if err != nil {
		zap.L().Error(err.Error())
		return nil, nil, err
	}
	zap.L().Sugar().Infof("Result hash: %s \n", hexutil.Encode(hash))

	tx, err := transactions.New(node.validator.pubKey, transactions.Result, node.validator.privKey)
	if err != nil {
		return nil, nil, err
	}
	tx.AddValues([]transactions.Value{
		transactions.BytesValue{
			Value: node.nebulaId[:],
		},
		transactions.IntValue{
			Value: int64(pulseId),
		},
		transactions.BytesValue{
			Value: sign,
		},
		transactions.BytesValue{
			Value: []byte{byte(node.chainType)},
		},
		transactions.BytesValue{
			Value: node.oraclePubKey[:],
		},
		transactions.IntValue{
			Value: int64(node.chainType),
		},
	})

	err = node.gravityClient.SendTx(tx)
	if err != nil {
		return nil, nil, err
	}

	zap.L().Sugar().Infof("Sign result txId: %s\n", hexutil.Encode(tx.Id[:]))
	return result, hash, nil
}

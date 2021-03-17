package node

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"go.uber.org/zap"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (node *Node) commit(data *extractor.Data, tcHeight uint64, pulseId uint64) ([]byte, error) {
	dataBytes := toBytes(data, node.extractor.ExtractorType)
	commit := crypto.Keccak256(dataBytes)
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
func (node *Node) reveal(tcHeight uint64, pulseId uint64, reveal *extractor.Data, commit []byte) error {
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
	})

	err = node.gravityClient.SendTx(tx)
	if err != nil {
		return err
	}
	fmt.Printf("Reveal txId: %s\n", hexutil.Encode(tx.Id[:]))

	return nil
}
func (node *Node) signResult(tcHeight uint64, pulseId uint64, ctx context.Context) (*extractor.Data, []byte, error) {
	zap.L().Sugar().Debugf("signResults: height: %d, pulseId: %d", tcHeight, pulseId)
	var values []extractor.Data
	zap.L().Sugar().Debugf("gravity Reveals: chaintype: %d, pulseId: %d", node.chainType, pulseId, node.nebulaId.ToString(node.chainType))
	bytesValues, err := node.gravityClient.Reveals(node.chainType, node.nebulaId, int64(tcHeight), int64(pulseId))
	if err != nil {
		return nil, nil, err
	}
	zap.L().Sugar().Debugf("signResults: value len: %d", len(bytesValues))
	for _, v := range bytesValues {
		zap.L().Sugar().Debugf("signResults: decoding: ", v)
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			continue
		}
		values = append(values, *fromBytes(b, node.extractor.ExtractorType))
	}

	if len(values) == 0 {
		return nil, nil, nil
	}
	result, err := node.extractor.Aggregate(values, ctx)
	if err != nil {
		return nil, nil, err
	}

	hash := crypto.Keccak256(toBytes(result, node.extractor.ExtractorType))
	sign, err := node.adaptor.Sign(hash)
	if err != nil {
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
	})

	err = node.gravityClient.SendTx(tx)
	if err != nil {
		return nil, nil, err
	}

	zap.L().Sugar().Infof("Sign result txId: %s\n", hexutil.Encode(tx.Id[:]))
	return result, hash, nil
}

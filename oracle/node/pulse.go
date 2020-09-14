package node

import (
	"context"
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func (node *Node) commit(data interface{}, tcHeight uint64, pulseId uint64) ([]byte, error) {
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
func (node *Node) reveal(tcHeight uint64, pulseId uint64, reveal interface{}, commit []byte) error {
	dataBytes := toBytes(reveal, node.extractor.ExtractorType)
	fmt.Printf("Reveal: %s  - %s \n", hexutil.Encode(dataBytes), hexutil.Encode(commit))

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
func (node *Node) signResult(tcHeight uint64, pulseId uint64, ctx context.Context) (interface{}, []byte, error) {
	var values []interface{}
	bytesValues, err := node.gravityClient.Reveals(node.chainType, node.nebulaId, int64(tcHeight), int64(pulseId))
	if err != nil {
		return nil, nil, err
	}

	for _, v := range bytesValues {
		values = append(values, fromBytes(v, node.extractor.ExtractorType))
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
	fmt.Printf("Result hash: %s \n", hexutil.Encode(hash))

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

	fmt.Printf("Sign result txId: %s\n", hexutil.Encode(tx.Id[:]))
	return result, hash, nil
}

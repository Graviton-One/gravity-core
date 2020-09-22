package node

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"

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
	var values []extractor.Data
	bytesValues, err := node.gravityClient.Reveals(node.chainType, node.nebulaId, int64(tcHeight), int64(pulseId))
	if err != nil {
		return nil, nil, err
	}

	for _, v := range bytesValues {
		b, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			continue
		}
		values = append(values, *fromBytes(b, node.extractor.ExtractorType))
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

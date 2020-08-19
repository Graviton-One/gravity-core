package client

import (
	"encoding/binary"
	"encoding/json"
	"errors"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/storage"
	"github.com/Gravity-Tech/gravity-core/common/transactions"
	"github.com/Gravity-Tech/gravity-core/ledger/query"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	InternalServerErrCode = 500
	NotFoundCode          = 404
)

var (
	ErrValueNotFound  = errors.New("value not found")
	ErrInternalServer = errors.New("internal server error")
)

type Client struct {
	Host       string
	HttpClient *rpchttp.HTTP
}

func New(host string) (*Client, error) {
	client, err := rpchttp.New(host, "/websocket")
	if err != nil {
		return nil, err
	}
	return &Client{Host: host, HttpClient: client}, nil
}

func (ghClient *Client) SendTx(transaction *transactions.Transaction) error {
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		return err
	}

	rs, err := ghClient.HttpClient.BroadcastTxCommit(txBytes)
	if err != nil {
		return err
	}
	if rs.CheckTx.Code == InternalServerErrCode {
		return errors.New(rs.CheckTx.Info)
	} else if rs.DeliverTx.Code == InternalServerErrCode {
		return errors.New(rs.DeliverTx.Info)
	}
	return err
}

func (client *Client) OraclesByValidator(pubKey account.ConsulPubKey) (storage.OraclesByTypeMap, error) {
	rq := query.ByValidatorRq{
		PubKey: hexutil.Encode(pubKey[:]),
	}

	rs, err := client.do(query.OracleByValidatorPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	oracles := make(storage.OraclesByTypeMap)
	if err == ErrValueNotFound {
		return oracles, nil
	}

	err = json.Unmarshal(rs, &oracles)
	if err != nil {
		return nil, err
	}

	return oracles, nil
}

func (client *Client) OraclesByNebula(nebulaId account.NebulaId, chainType account.ChainType) (storage.OraclesMap, error) {
	rq := query.ByNebulaRq{
		ChainType:     chainType,
		NebulaAddress: hexutil.Encode(nebulaId),
	}

	rs, err := client.do(query.OracleByNebulaPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	oracles := make(storage.OraclesMap)
	if err == ErrValueNotFound {
		return oracles, nil
	}

	err = json.Unmarshal(rs, &oracles)
	if err != nil {
		return nil, err
	}

	return oracles, nil
}

func (client *Client) BftOraclesByNebula(chainType account.ChainType, nebulaId account.NebulaId) (storage.OraclesMap, error) {
	rq := query.ByNebulaRq{
		ChainType:     chainType,
		NebulaAddress: hexutil.Encode(nebulaId),
	}

	rs, err := client.do(query.OracleByNebulaPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	oracles := make(storage.OraclesMap)
	if err == ErrValueNotFound {
		return oracles, nil
	}

	err = json.Unmarshal(rs, &oracles)
	if err != nil {
		return nil, err
	}

	return oracles, nil
}
func (client *Client) Results(height uint64, chainType account.ChainType, nebulaId account.NebulaId) ([][]byte, error) {
	rq := query.ResultsRq{
		Height:        height,
		ChainType:     chainType,
		NebulaAddress: hexutil.Encode(nebulaId),
	}

	rs, err := client.do(query.OracleByNebulaPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	var oracles [][]byte
	if err == ErrValueNotFound {
		return oracles, nil
	}

	err = json.Unmarshal(rs, &oracles)
	if err != nil {
		return nil, err
	}

	return oracles, nil
}

func (client *Client) RoundHeight(chainType account.ChainType, ledgerHeight uint64) (uint64, error) {
	rq := query.RoundHeightRq{
		ChainType:    chainType,
		LedgerHeight: ledgerHeight,
	}

	rs, err := client.do(query.RoundHeightPath, rq)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(rs), nil
}
func (client *Client) CommitHash(chainType account.ChainType, nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	rq := query.CommitHashRq{
		ChainType:     chainType,
		NebulaAddress: hexutil.Encode(nebulaAddress),
		Height:        height,
		OraclePubKey:  hexutil.Encode(oraclePubKey[:]),
	}

	rs, err := client.do(query.CommitHashPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	return rs, nil
}
func (client *Client) Reveal(nebulaAddress []byte, height int64, commitHash []byte) ([]byte, error) {
	rq := query.RevealRq{
		NebulaAddress: hexutil.Encode(nebulaAddress),
		Height:        height,
		CommitHash:    hexutil.Encode(commitHash),
	}

	rs, err := client.do(query.RevealPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	return rs, nil
}
func (client *Client) Result(chainType account.ChainType, nebulaAddress []byte, height int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	rq := query.ResultRq{
		ChainType:     chainType,
		NebulaAddress: hexutil.Encode(nebulaAddress),
		Height:        height,
		OraclePubKey:  hexutil.Encode(oraclePubKey[:]),
	}

	rs, err := client.do(query.ResultPath, rq)
	if err != nil || err != ErrValueNotFound {
		return nil, err
	}

	return rs, nil
}

func (client *Client) do(path query.Path, rq interface{}) ([]byte, error) {
	var err error
	b, ok := rq.([]byte)
	if !ok {
		b, err = json.Marshal(rq)
		if err != nil {
			return nil, err
		}
	}

	rs, err := client.HttpClient.ABCIQuery(string(path), b)
	if err != nil {
		return nil, err
	} else if rs.Response.Code == InternalServerErrCode {
		return nil, ErrInternalServer
	} else if rs.Response.Code == NotFoundCode {
		return nil, ErrValueNotFound
	}

	return rs.Response.Value, err
}

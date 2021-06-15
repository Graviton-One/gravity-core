package gravity

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
	"go.uber.org/zap"
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

func (client *Client) SendTx(transaction *transactions.Transaction) error {
	txBytes, err := json.Marshal(transaction)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}

	rs, err := client.HttpClient.BroadcastTxCommit(txBytes)
	if err != nil {
		zap.L().Error(err.Error())
		return err
	}
	if rs.CheckTx.Code == InternalServerErrCode {
		zap.L().Sugar().Error("Check error ", rs.CheckTx.Code)
		return errors.New(rs.CheckTx.Info)
	} else if rs.DeliverTx.Code == InternalServerErrCode {
		zap.L().Sugar().Error("Deliver error ", rs.CheckTx.Code)
		return errors.New(rs.DeliverTx.Info)
	}
	return err
}

func (client *Client) OraclesByValidator(pubKey account.ConsulPubKey) (storage.OraclesByTypeMap, error) {
	rq := query.ByValidatorRq{
		PubKey: hexutil.Encode(pubKey[:]),
	}

	rs, err := client.do(query.OracleByValidatorPath, rq)
	if err != nil || err == ErrValueNotFound {
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
		NebulaAddress: nebulaId.ToString(chainType),
	}

	rs, err := client.do(query.OracleByNebulaPath, rq)
	if err != nil && err != ErrValueNotFound {
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
		NebulaAddress: nebulaId.ToString(chainType),
	}

	rs, err := client.do(query.BftOracleByNebulaPath, rq)
	if err != nil && err != ErrValueNotFound {
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
func (client *Client) Results(height uint64, chainType account.ChainType, nebulaId account.NebulaId) ([]string, error) {
	rq := query.ResultsRq{
		Height:        height,
		ChainType:     chainType,
		NebulaAddress: nebulaId.ToString(chainType),
	}

	rs, err := client.do(query.ResultsPath, rq)
	if err != nil && err != ErrValueNotFound {
		return nil, err
	}

	var oracles []string
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

func (client *Client) LastRoundApproved() (uint64, error) {
	rs, err := client.do(query.LastRoundApprovedPath, nil)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(rs), nil
}
func (client *Client) CommitHash(chainType account.ChainType, nebulaId account.NebulaId, height int64, pulseId int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	rq := query.CommitHashRq{
		ChainType:     chainType,
		NebulaAddress: nebulaId.ToString(chainType),
		Height:        height,
		PulseId:       pulseId,
		OraclePubKey:  oraclePubKey.ToString(chainType),
	}

	rs, err := client.do(query.CommitHashPath, rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
func (client *Client) Reveal(chainType account.ChainType, oraclePubKey account.OraclesPubKey, nebulaId account.NebulaId, height int64, pulseId int64, commitHash []byte) ([]byte, error) {
	rq := query.RevealRq{
		ChainType:     chainType,
		NebulaAddress: nebulaId.ToString(chainType),
		Height:        height,
		PulseId:       pulseId,
		OraclePubKey:  oraclePubKey.ToString(chainType),
		CommitHash:    hexutil.Encode(commitHash),
	}

	rs, err := client.do(query.RevealPath, rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
func (client *Client) Reveals(chainType account.ChainType, nebulaId account.NebulaId, height int64, pulseId int64) ([]string, error) {
	rq := query.RevealRq{
		ChainType:     chainType,
		NebulaAddress: nebulaId.ToString(chainType),
		Height:        height,
		PulseId:       pulseId,
	}

	rs, err := client.do(query.RevealsPath, rq)
	if err != nil {
		return nil, err
	}
	zap.L().Sugar().Debugf("Reveals res: %s", string(rs))
	var reveals []string
	if err == ErrValueNotFound {
		return reveals, nil
	}

	err = json.Unmarshal(rs, &reveals)
	if err != nil {
		return nil, err
	}

	return reveals, nil
}
func (client *Client) Result(chainType account.ChainType, nebulaId account.NebulaId, height int64, oraclePubKey account.OraclesPubKey) ([]byte, error) {
	rq := query.ResultRq{
		ChainType:     chainType,
		NebulaAddress: nebulaId.ToString(chainType),
		Height:        height,
		OraclePubKey:  oraclePubKey.ToString(chainType),
	}

	rs, err := client.do(query.ResultPath, rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
func (client *Client) NebulaInfo(id account.NebulaId, chainType account.ChainType) (*storage.NebulaInfo, error) {
	rq := query.ByNebulaRq{
		ChainType:     chainType,
		NebulaAddress: id.ToString(chainType),
	}

	rs, err := client.do(query.NebulaInfoPath, rq)
	if err != nil {
		return nil, err
	}
	var nebulaInfo storage.NebulaInfo
	err = json.Unmarshal(rs, &nebulaInfo)
	if err != nil {
		return nil, err
	}

	return &nebulaInfo, nil
}
func (client *Client) Nebulae() (storage.NebulaMap, error) {
	rs, err := client.do(query.NebulaePath, nil)
	if err != nil && err != ErrValueNotFound {
		return nil, err
	}

	nebulae := make(storage.NebulaMap)
	if err == ErrValueNotFound {
		return nebulae, nil
	}

	err = json.Unmarshal(rs, &nebulae)
	if err != nil {
		return nil, err
	}

	return nebulae, nil
}
func (client *Client) Consuls() ([]storage.Consul, error) {
	rs, err := client.do(query.ConsulsPath, nil)
	if err != nil && err != ErrValueNotFound {
		return nil, err
	}

	var consuls []storage.Consul
	if err == ErrValueNotFound {
		return consuls, nil
	}

	err = json.Unmarshal(rs, &consuls)
	if err != nil {
		return nil, err
	}

	return consuls, nil
}
func (client *Client) ConsulsCandidate() ([]storage.Consul, error) {
	rs, err := client.do(query.ConsulsCandidatePath, nil)
	if err != nil && err != ErrValueNotFound {
		return nil, err
	}

	var consuls []storage.Consul
	if err == ErrValueNotFound {
		return consuls, nil
	}

	err = json.Unmarshal(rs, &consuls)
	if err != nil {
		return nil, err
	}

	return consuls, nil
}
func (client *Client) SignNewConsulsByConsul(pubKey account.ConsulPubKey, chainId account.ChainType, roundId int64) ([]byte, error) {
	rq := query.SignByConsulRq{
		ConsulPubKey: hexutil.Encode(pubKey[:]),
		ChainType:    chainId,
		RoundId:      roundId,
	}

	rs, err := client.do(query.SignNewConsulsByConsulPath, rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
func (client *Client) SignNewOraclesByConsul(pubKey account.ConsulPubKey, chainId account.ChainType, nebulaId account.NebulaId, roundId int64) ([]byte, error) {
	rq := query.SignByConsulRq{
		ConsulPubKey: hexutil.Encode(pubKey[:]),
		ChainType:    chainId,
		RoundId:      roundId,
		NebulaId:     nebulaId.ToString(chainId),
	}

	rs, err := client.do(query.SignNewOraclesByConsulPath, rq)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
func (client *Client) NebulaOraclesIndex(chainId account.ChainType, nebulaId account.NebulaId) (uint64, error) {
	rq := query.ByNebulaRq{
		ChainType:     chainId,
		NebulaAddress: nebulaId.ToString(chainId),
	}

	rs, err := client.do(query.NebulaOraclesIndexPath, rq)
	if err != nil && err != ErrValueNotFound {
		return 0, err
	}

	return binary.BigEndian.Uint64(rs), nil
}

func (client *Client) NebulaCustomParams(id account.NebulaId, chainType account.ChainType) (storage.NebulaCustomParams, error) {
	rq := query.ByNebulaRq{
		ChainType:     chainType,
		NebulaAddress: id.ToString(chainType),
	}
	rs, err := client.do(query.NebulaCustomParams, rq)
	if err != nil && err != ErrValueNotFound {
		return nil, err
	}

	nebulaCustomParams := storage.NebulaCustomParams{}

	err = json.Unmarshal(rs, &nebulaCustomParams)
	if err != nil {
		return nil, err
	}

	return nebulaCustomParams, nil
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

	return rs.Response.Value, nil
}

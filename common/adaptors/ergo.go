package adaptors

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Gravity-Tech/gravity-core/abi"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"github.com/gookit/validate"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	crypto "crypto/ed25519"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
)

const (
	ConsulsNumber = 5
)

type ErgoAdaptor struct {
	secret crypto.PrivateKey

	ergoClient      *helpers.ErgClient `option:"ergClient"`
	ghClient        *gravity.Client    `option:"ghClient"`
	gravityContract string             `option:"gravityContract"`
}

type ErgoAdapterOption func(*ErgoAdaptor) error

func (adaptor *ErgoAdaptor) applyOpts(opts AdapterOptions) error {
	err := validateErgoAdapterOptions(opts)
	if err != nil {
		return err
	}
	v := reflect.TypeOf(*adaptor)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("option")
		val, ok := opts[tag]
		if ok {
			switch tag {
			case "ghClient":
				adaptor.ghClient = val.(*gravity.Client)
			case "ergClient":
				adaptor.ergoClient = val.(*helpers.ErgClient)
			case "gravityContract":
				adaptor.gravityContract = val.(string)

			}
		}
	}
	return nil
}

func validateErgoAdapterOptions(opts AdapterOptions) error {
	v := validate.Map(opts)
	v.AddRule("ghClient", "isGhClient")
	v.AddRule("ergClient", "isErgClient")
	v.AddRule("gravityContract", "string")

	if !v.Validate() { // validate ok
		return v.Errors
	}
	return nil
}

func WithErgoGravityContract(address string) ErgoAdapterOption {
	return func(h *ErgoAdaptor) error {
		h.gravityContract = address
		return nil
	}
}

func ErgoAdapterWithGhClient(ghClient *gravity.Client) ErgoAdapterOption {
	return func(h *ErgoAdaptor) error {
		h.ghClient = ghClient
		return nil
	}
}

func NewErgoAdapterByOpts(seed []byte, nodeUrl string, ctx context.Context, opts AdapterOptions) (*ErgoAdaptor, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", nodeUrl, nil)
	if err != nil {
		return nil, err
	}
	client, err := helpers.NewClient(helpers.ErgOptions{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	_, err = client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}
	secret := crypto.NewKeyFromSeed(seed)
	adapter := &ErgoAdaptor{
		secret:     secret,
		ergoClient: client,
	}
	err = adapter.applyOpts(opts)
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

func NewErgoAdapter(seed []byte, nodeUrl string, ctx context.Context, opts ...ErgoAdapterOption) (*ErgoAdaptor, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", nodeUrl, nil)
	if err != nil {
		return nil, err
	}
	client, err := helpers.NewClient(helpers.ErgOptions{ApiKey: "", BaseUrl: nodeUrl})
	if err != nil {
		return nil, err
	}

	_, err = client.Do(ctx, req, nil)
	if err != nil {
		return nil, err
	}

	secret := crypto.NewKeyFromSeed(seed)
	er := &ErgoAdaptor{
		ergoClient: client,
		secret:     secret,
	}
	for _, opt := range opts {
		err := opt(er)
		if err != nil {
			return nil, err
		}
	}
	return er, nil
}

func (adaptor *ErgoAdaptor) WaitTx(id string, ctx context.Context) error {
	type Response struct {
		Status  bool `json:"success"`
		Confirm int  `json:"numConfirmations"`
	}
	out := make(chan error)
	const TxWaitCount = 10
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "/numConfirmations")
	if err != nil {
		out <- err
	}
	go func() {
		defer close(out)
		for i := 0; i <= TxWaitCount; i++ {
			req, err := http.NewRequest("GET", url.String()+"/:"+id, nil)
			if err != nil {
				out <- err
				break
			}
			response := new(Response)
			_, err = adaptor.ergoClient.Do(ctx, req, response)
			if err != nil {
				out <- err
				break
			}

			if response.Confirm == -1 {
				_, err = adaptor.ergoClient.Do(ctx, req, response)
				if err != nil {
					out <- err
					break
				}

				if response.Confirm == -1 {
					out <- errors.New("tx not found")
					break
				} else {
					break
				}
			}

			if TxWaitCount == i {
				out <- errors.New("tx not found")
				break
			}
			time.Sleep(time.Second)
		}
	}()
	return <-out
}

func (adaptor *ErgoAdaptor) GetHeight(ctx context.Context) (uint64, error) {
	type Response struct {
		Status bool   `json:"success"`
		Height uint64 `json:"height"`
	}
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "/height")
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return 0, err
	}
	response := new(Response)
	_, err = adaptor.ergoClient.Do(ctx, req, response)
	if err != nil {
		return 0, err
	}

	return response.Height, nil
}

func (adaptor *ErgoAdaptor) Sign(msg []byte) ([]byte, error) {
	type Sign struct {
		A string
		Z string
	}
	type Response struct {
		Status bool `json:"success"`
		Signed Sign `json:"signed"`
	}
	values := map[string]string{"msg": hex.EncodeToString(msg), "sk": hex.EncodeToString(adaptor.secret)}
	jsonValue, _ := json.Marshal(values)
	res, err := http.Post(adaptor.ergoClient.Options.BaseUrl+"/sign", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	response, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseObject Response
	err = json.Unmarshal(response, &responseObject)
	if err != nil {
		return nil, err
	}

	if !responseObject.Status {
		err = fmt.Errorf("proxy connection problem")
		return nil, err
	}
	fmt.Println(responseObject.Signed)
	return []byte(responseObject.Signed.A + responseObject.Signed.Z), nil
}

func (adaptor *ErgoAdaptor) PubKey() account.OraclesPubKey {
	pubKey := adaptor.secret.Public()
	pk, _ := pubKey.(crypto.PublicKey)
	oraclePubKey := account.BytesToOraclePubKey(pk, account.Ergo)
	return oraclePubKey
}

// ValueType TODO: returns datatype from extractor (type of target)
func (adaptor *ErgoAdaptor) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
	panic("implement me")
}

func (adaptor *ErgoAdaptor) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	type Oracle struct {
		State   bool     `json:"state"`
		Oracles []string `json:"oracles"`
		Bft     int      `json:"bft"`
	}
	type Result struct {
		Success  bool   `json:"success"`
		Response Oracle `json:"response"`
	}
	type Sign struct {
		a []string
		z []string
	}
	type Data struct {
		Signs Sign   `json:"signs"`
		Hash  []byte `json:"hashData"`
	}
	type Tx struct {
		Success bool   `json:"success"`
		TxId    string `json:"txId"`
	}
	var oracles []string
	var signsA []string
	var signsZ []string
	realSignCount := 0

	// Get oracles and bftValue
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/getPreAddPulseInfo")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String()+"/:"+strconv.FormatUint(pulseId, 10), nil)
	if err != nil {
		return "", err
	}
	result := new(Result)
	_, err = adaptor.ergoClient.Do(ctx, req, result)
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", errors.New("can't get oracles")
	} else if result.Success && !result.Response.State {
		return "", errors.New("wrong pulseID")
	} else {
		oracles = result.Response.Oracles
	}

	// Iterate over oracles and get signs
	for _, oracle := range oracles {
		pubKey, err := account.StringToOraclePubKey(oracle, account.Ergo)
		if err != nil {
			signsA = append(signsA, hex.EncodeToString([]byte{0}))
			signsZ = append(signsZ, hex.EncodeToString([]byte{0}))
			continue
		}
		sign, err := adaptor.ghClient.Result(account.Ergo, nebulaId, int64(pulseId), pubKey)

		if err != nil {
			signsA = append(signsA, hex.EncodeToString([]byte{0}))
			signsZ = append(signsZ, hex.EncodeToString([]byte{0}))
			continue
		}
		signsA = append(signsA, string(sign[:66]))
		signsZ = append(signsZ, string(sign[66:]))
		realSignCount++
	}

	// Check realSignCount with bftValue before sending data
	if realSignCount == 0 {
		return "", nil
	}
	if realSignCount < result.Response.Bft {
		return "", nil
	}

	// Send oracleSigns to be verified by contract in proxy side and get txId
	data, err := json.Marshal(Data{Signs: Sign{a: signsA, z: signsZ}, Hash: hash})
	url, err = helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/addPulse")
	if err != nil {
		return "", err
	}
	req, err = http.NewRequestWithContext(ctx, "POST", url.String(), bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	tx := new(Tx)
	_, err = adaptor.ergoClient.Do(ctx, req, tx)
	if err != nil {
		return "", err
	}

	return tx.TxId, nil
}

func (adaptor *ErgoAdaptor) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
	return nil
}

func (adaptor *ErgoAdaptor) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	var signsA [5]string
	var signsZ [5]string
	type Tx struct {
		Success bool   `json:"success"`
		TxId    string `json:"txId"`
	}
	type Consuls struct {
		Success bool     `json:"success"`
		consuls []string `json:"consuls"`
	}
	type Sign struct {
		a [5]string
		z [5]string
	}
	type Data struct {
		newOracles []string `json:"newOracles"`
		Signs	Sign	`json:"signs"`
	}

	lastRound, err := adaptor.LastRound(ctx)
	if err != nil {
		return "", err
	}
	if uint64(round) <= lastRound{
		return "", errors.New("this is not a new round")
	}

	var consuls []string
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/getConsuls")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	result := new(Consuls)
	_, err = adaptor.ergoClient.Do(ctx, req, result)
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", errors.New("can't get consuls")
	} else {
		consuls = result.consuls
	}

	for k, sign := range signs {
		pubKey := k.ToString(account.Ergo)
		index := -1

		for i, v := range consuls {
			if v == pubKey {
				index = i
				break
			}
		}

		if index == -1 {
			continue
		}
		signsA[index] = string(sign[:66])
		signsZ[index] = string(sign[66:])
	}

	for i, v := range signsA {
		if v != "" {
			continue
		}

		signsA[i] = hex.EncodeToString([]byte{0})
		signsZ[i] = hex.EncodeToString([]byte{0})
	}

	var newOracles []string

	for _, v := range oracles {
		if v == nil {
			newOracles = append(newOracles, hex.EncodeToString([]byte{0}))
			continue
		}
		newOracles = append(newOracles, hex.EncodeToString(v.ToBytes(account.Ergo)))
	}

	url, err = helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/updateOracles")
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(Data{newOracles: newOracles, Signs: Sign{a: signsA , z: signsZ} })
	req, err = http.NewRequestWithContext(ctx, "POST", url.String(), bytes.NewBuffer(data))
	tx := new(Tx)
	_, err = adaptor.ergoClient.Do(ctx, req, tx)
	if err != nil {
		return "", err
	}

	return tx.TxId, nil
}

func (adaptor *ErgoAdaptor) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	var signsA [5]string
	var signsZ [5]string
	type Tx struct {
		Success bool   `json:"success"`
		TxId    string `json:"txId"`
	}
	type Consuls struct {
		Success bool     `json:"success"`
		consuls []string `json:"consuls"`
	}
	type Sign struct {
		a [5]string
		z [5]string
	}
	type Data struct {
		newConsuls []string `json:"newConsuls"`
		Signs	Sign	`json:"signs"`
	}

	lastRound, err := adaptor.LastRound(ctx)
	if err != nil {
		return "", err
	}
	if uint64(round) <= lastRound{
		return "", errors.New("this is not a new round")
	}

	var consuls []string
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/getConsuls")
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	result := new(Consuls)
	_, err = adaptor.ergoClient.Do(ctx, req, result)
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", errors.New("can't get consuls")
	} else {
		consuls = result.consuls
	}

	for k, sign := range signs {
		pubKey := k.ToString(account.Ergo)
		index := -1

		for i, v := range consuls {
			if v == pubKey {
				index = i
				break
			}
		}

		if index == -1 {
			continue
		}
		signsA[index] = string(sign[:66])
		signsZ[index] = string(sign[66:])
	}

	for i, v := range signsA {
		if v != "" {
			continue
		}

		signsA[i] = hex.EncodeToString([]byte{0})
		signsZ[i] = hex.EncodeToString([]byte{0})
	}

	var newConsulsString []string

	for _, v := range newConsulsAddresses {
		if v == nil {
			newConsulsString = append(newConsulsString, hex.EncodeToString([]byte{0}))
			continue
		}
		newConsulsString = append(newConsulsString, hex.EncodeToString(v.ToBytes(account.Ergo)))
	}

	emptyCount := ConsulsNumber - len(newConsulsString)
	for i := 0; i < emptyCount; i++ {
		newConsulsString = append(newConsulsString, hex.EncodeToString([]byte{0}))
	}

	url, err = helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/updateConsuls")
	if err != nil {
		return "", err
	}
	data, err := json.Marshal(Data{newConsuls: newConsulsString, Signs: Sign{a: signsA , z: signsZ} })
	req, err = http.NewRequestWithContext(ctx, "POST", url.String(), bytes.NewBuffer(data))
	tx := new(Tx)
	_, err = adaptor.ergoClient.Do(ctx, req, tx)
	if err != nil {
		return "", err
	}

	return tx.TxId, nil
}

func (adaptor *ErgoAdaptor) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64) ([]byte, error) {
	var msg []string
	for _, v := range consulsAddresses {
		if v == nil {
			msg = append(msg, hex.EncodeToString([]byte{0}))
			continue
		}
		msg = append(msg, hex.EncodeToString(v.ToBytes(account.Ergo)))
	}
	msg = append(msg, fmt.Sprintf("%d", roundId))

	sign, err := adaptor.Sign([]byte(strings.Join(msg, ",")))
	if err != nil {
		return nil, err
	}

	return sign, err
}

func (adaptor *ErgoAdaptor) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey) ([]byte, error) {
	var stringOracles []string
	for _, v := range oracles {
		if v == nil {
			stringOracles = append(stringOracles, hex.EncodeToString([]byte{1}))
			continue
		}
		stringOracles = append(stringOracles, hex.EncodeToString(v.ToBytes(account.Ergo)))
	}

	sign, err := adaptor.Sign([]byte(strings.Join(stringOracles, ",")))
	if err != nil {
		return nil, err
	}

	return sign, err
}

func (adaptor *ErgoAdaptor) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	type Result struct {
		Success bool   `json:"success"`
		PulseId string `json:"pulse_id"`
	}
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/getLastPulseId")
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return 0, err
	}
	result := new(Result)
	_, err = adaptor.ergoClient.Do(ctx, req, result)
	if err != nil {
		return 0, err
	}
	if !result.Success {
		return 0, errors.New("can't get lastPulseId")
	}
	pulseId, _ := strconv.ParseUint(result.PulseId, 10, 64)
	return pulseId, nil
}

func (adaptor *ErgoAdaptor) LastRound(ctx context.Context) (uint64, error) {
	type Result struct {
		Success   bool   `json:"success"`
		LastRound string `json:"lastRound"`
	}
	url, err := helpers.JoinUrl(adaptor.ergoClient.Options.BaseUrl, "adaptor/lastRound")
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequestWithContext(ctx, "GET", url.String(), nil)
	if err != nil {
		return 0, err
	}
	result := new(Result)
	_, err = adaptor.ergoClient.Do(ctx, req, result)
	if err != nil {
		return 0, err
	}
	if !result.Success {
		return 0, errors.New("can't get lastRound")
	}
	lastRound, _ := strconv.ParseUint(result.LastRound, 10, 64)
	return lastRound, nil
}

func (adaptor *ErgoAdaptor) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	lastRound, err := adaptor.LastRound(ctx)
	if err != nil {
		return false, err
	}
	if uint64(roundId) > lastRound{
		return false, nil
	}
	return true, nil

}

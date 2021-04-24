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
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	crypto "crypto/ed25519"
	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/helpers"
)

const (
	Consuls = 5
)

type ErgoAdaptor struct {
	secret crypto.PrivateKey

	ergoClient      *helpers.ErgClient `option:"ergClient"`
	ghClient        *gravity.Client    `option:"ghClient"`
	gravityContract string             `option:"gravityContract"`
}

type ErgoAdapterOption func(*ErgoAdaptor) error

func (er *ErgoAdaptor) applyOpts(opts AdapterOptions) error {
	err := validateErgoAdapterOptions(opts)
	if err != nil {
		return err
	}
	v := reflect.TypeOf(*er)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		tag := field.Tag.Get("option")
		val, ok := opts[tag]
		if ok {
			switch tag {
			case "ghClient":
				er.ghClient = val.(*gravity.Client)
			case "ergClient":
				er.ergoClient = val.(*helpers.ErgClient)
			case "gravityContract":
				er.gravityContract = val.(string)

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

	client.Do(ctx, req, nil)
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

	client.Do(ctx, req, nil)

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

func (er *ErgoAdaptor) WaitTx(id string, ctx context.Context) error {
	type Response struct {
		Status  bool `json:"success"`
		Confirm int  `json:"numConfirmations"`
	}
	out := make(chan error)
	const TxWaitCount = 10
	url, err := helpers.JoinUrl(er.ergoClient.Options.BaseUrl, "/numConfirmations")
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
			_, err = er.ergoClient.Do(ctx, req, response)
			if err != nil {
				out <- err
				break
			}

			if response.Confirm == -1 {
				_, err = er.ergoClient.Do(ctx, req, response)
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

func (er *ErgoAdaptor) GetHeight(ctx context.Context) (uint64, error) {
	type Response struct {
		Status bool   `json:"success"`
		Height uint64 `json:"height"`
	}
	url, err := helpers.JoinUrl(er.ergoClient.Options.BaseUrl, "/height")
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return 0, err
	}
	response := new(Response)
	_, err = er.ergoClient.Do(ctx, req, response)
	if err != nil {
		return 0, err
	}

	return response.Height, nil
}

func (er *ErgoAdaptor) Sign(msg []byte) ([]byte, error) {
	type Sign struct {
		A string
		Z string
	}
	type Response struct {
		Status bool `json:"success"`
		Signed Sign `json:"signed"`
	}
	values := map[string]string{"msg": hex.EncodeToString(msg), "sk": hex.EncodeToString(er.secret)}
	jsonValue, _ := json.Marshal(values)
	res, err := http.Post(er.ergoClient.Options.BaseUrl+"/sign", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, err
	}
	response, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var responseObject Response
	json.Unmarshal(response, &responseObject)
	if !responseObject.Status {
		err = fmt.Errorf("proxy connection problem")
		return nil, err
	}
	fmt.Println(responseObject.Signed)
	return []byte(responseObject.Signed.A[:24] + responseObject.Signed.Z[:32]), nil
}

func (er *ErgoAdaptor) PubKey() account.OraclesPubKey {
	pubKey := er.secret.Public()
	pk, _ := pubKey.(crypto.PublicKey)
	oraclePubKey := account.BytesToOraclePubKey(pk, account.Ergo)
	return oraclePubKey
}

func (er *ErgoAdaptor) ValueType(nebulaId account.NebulaId, ctx context.Context) (abi.ExtractorType, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) AddPulse(nebulaId account.NebulaId, pulseId uint64, validators []account.OraclesPubKey, hash []byte, ctx context.Context) (string, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) SendValueToSubs(nebulaId account.NebulaId, pulseId uint64, value *extractor.Data, ctx context.Context) error {
	panic("implement me")
}

func (er *ErgoAdaptor) SetOraclesToNebula(nebulaId account.NebulaId, oracles []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) SendConsulsToGravityContract(newConsulsAddresses []*account.OraclesPubKey, signs map[account.OraclesPubKey][]byte, round int64, ctx context.Context) (string, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) SignConsuls(consulsAddresses []*account.OraclesPubKey, roundId int64) ([]byte, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) SignOracles(nebulaId account.NebulaId, oracles []*account.OraclesPubKey) ([]byte, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) LastPulseId(nebulaId account.NebulaId, ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) LastRound(ctx context.Context) (uint64, error) {
	panic("implement me")
}

func (er *ErgoAdaptor) RoundExist(roundId int64, ctx context.Context) (bool, error) {
	panic("implement me")
}

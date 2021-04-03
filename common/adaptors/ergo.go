package adaptors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gookit/validate"
	"io/ioutil"
	"net/http"
	"reflect"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/wavesplatform/gowaves/pkg/crypto"
)

const (
	Consuls = 5
)

type ErgoAdaptor struct {
	secret 			crypto.SecretKey
	ghClient        *gravity.Client      `option:"ghClient"`
	//wavesClient     *wclient.Client      `option:"wvClient"`
	//helper          helpers.ClientHelper `option:"-"`
	gravityContract string               `option:"gravityContract"`
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

func NewErgoAdapterByOpts(seed []byte, nodeUrl string, opts AdapterOptions) (*ErgoAdaptor, error) {
	//wClient, err := wclient.NewClient(wclient.Options{ApiKey: "", BaseUrl: nodeUrl})
	//if err != nil {
	//	return nil, err
	//}

	secret, err := crypto.NewSecretKeyFromBytes(seed)
	adapter := &ErgoAdaptor{
		secret:      secret,
	}
	err = adapter.applyOpts(opts)
	if err != nil {
		return nil, err
	}

	return adapter, nil
}

func NewErgoAdapter(seed []byte, nodeUrl string, ctx context.Context, opts ...ErgoAdapterOption) (*ErgoAdaptor, error) {
	//wClient, err := wclient.NewClient(wclient.Options{ApiKey: "", BaseUrl: nodeUrl})
	//if err != nil {
	//	return nil, err
	//}

	secret, err := crypto.NewSecretKeyFromBytes(seed)
	if err != nil {
		return nil, err
	}
	adapter := &ErgoAdaptor{
		secret:      secret,

	}
	for _, opt := range opts {
		err := opt(adapter)
		if err != nil {
			return nil, err
		}
	}
	return adapter, nil
}


func (adaptor *ErgoAdaptor) GetHeight(ctx context.Context) (uint64, error) {
	type Response struct {
		Status  bool    `json:"name"`
		Height	uint64	`json:"pokemon_entries"`
	}
	res, err := http.NewRequest("GET", "http://127.0.0.1:9000/height", nil)
	if err != nil {
		return 0, err
	}
	responseData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}
	var responseObject Response
	json.Unmarshal(responseData, &responseObject)
	if !responseObject.Status {
		err = fmt.Errorf("proxy connection problem")
	}
	return responseObject.Height, err
}
func (adaptor *ErgoAdaptor) Sign(msg []byte) ([]byte, error) {
	sig, err := crypto.Sign(adaptor.secret, msg)
	if err != nil {
		return nil, err
	}
	return sig.Bytes(), nil
}
func (adaptor *ErgoAdaptor) PubKey() account.OraclesPubKey {
	pubKey := crypto.GeneratePublicKey(adaptor.secret)
	oraclePubKey := account.BytesToOraclePubKey(pubKey[:], account.Ergo)
	return oraclePubKey
}

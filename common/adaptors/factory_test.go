package adaptors

import (
	"context"
	"reflect"
	"testing"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/ethereum/go-ethereum/ethclient"
	wclient "github.com/wavesplatform/gowaves/pkg/client"
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func EqualExcept(a interface{}, b interface{}, ExceptFields []string) bool {
	val := reflect.ValueOf(a).Elem()
	otherFields := reflect.Indirect(reflect.ValueOf(b))

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		if stringInSlice(typeField.Name, ExceptFields) {
			continue
		}

		value := val.Field(i)
		otherValue := otherFields.FieldByName(typeField.Name)

		if value.Interface() != otherValue.Interface() {
			return false
		}
	}
	return true
}
func TestNewFactory(t *testing.T) {
	tests := []struct {
		name string
		want *Factory
	}{
		{name: "create factory", want: &Factory{}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFactory(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFactory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ghClientValidator(t *testing.T) {
	type args struct {
		val interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "check valid type", args: args{val: &gravity.Client{}}, want: true},
		{name: "check fail type", args: args{val: 4}, want: false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isGhClientValidator(tt.args.val); got != tt.want {
				t.Errorf("ghClientValidator() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFactory_CreateAdaptor(t *testing.T) {
	ghClient := gravity.Client{}
	wvClient := wclient.Client{}
	ethClient := ethclient.Client{}
	wanted := WavesAdaptor{
		ghClient:        &ghClient,
		wavesClient:     &wvClient,
		chainID:         1,
		gravityContract: "contract",
	}
	validOpts := AdapterOptions{
		"ghClient":        &ghClient,
		"wvClient":        &wvClient,
		"gravityContract": "contract",
		"chainID":         byte(1)}
	ethValidOpts := AdapterOptions{
		"ghClient":        &ghClient,
		"ethClient":       &ethClient,
		"gravityContract": "0x90C52beF8733cDF368Bf8AaD5ee4A78cB68E85"}
	f := NewFactory()
	type args struct {
		name               string
		oracleSecretKey    []byte
		targetChainNodeUrl string
		opts               AdapterOptions
	}
	tests := []struct {
		name    string
		f       *Factory
		args    args
		want    IBlockchainAdaptor
		wantErr bool
	}{

		{name: "Test valid options", f: f, args: args{name: "waves", oracleSecretKey: []byte("key"), targetChainNodeUrl: "ws:url", opts: validOpts}, wantErr: false, want: &wanted},
		{name: "Test incorrect adaptor name", f: f, args: args{name: "wavesddd", oracleSecretKey: []byte("key"), targetChainNodeUrl: "ws:url", opts: validOpts}, wantErr: true},
		{name: "Test invalid options", f: f, args: args{name: "waves", oracleSecretKey: []byte("key"), targetChainNodeUrl: "ws:url", opts: AdapterOptions{
			"ghClient":        &ghClient,
			"wvClient":        5,
			"gravityContract": "contract",
			"chainID":         byte(1)}}, wantErr: true},
		{name: "Test ethereum valid options", f: f, args: args{name: "ethereum", oracleSecretKey: []byte("key"), targetChainNodeUrl: "https://ropsten.infura.io/v3/598efca7168947c6a186e2f85b600be1", opts: ethValidOpts}, wantErr: false, want: &wanted},
		{name: "Test ethereum invalid options", f: f, args: args{name: "ethereum", oracleSecretKey: []byte("key"), targetChainNodeUrl: "https://ropsten.infura.io/v3/598efca7168947c6a186e2f85b600be1", opts: AdapterOptions{
			"ghClient":        &ghClient,
			"wvClient":        5,
			"gravityContract": "contract",
			"chainID":         byte(1)}}, wantErr: true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Factory{}
			_, err := f.CreateAdaptor(tt.args.name, tt.args.oracleSecretKey, tt.args.targetChainNodeUrl, context.Background(), tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Factory.CreateAdaptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

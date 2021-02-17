package adaptors

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
	wclient "github.com/wavesplatform/gowaves/pkg/client"
)

func TestOne(t *testing.T) {
	type S struct {
		F string `species:"gopher" color:"blue"`
	}

	s := S{}
	st := reflect.TypeOf(s)
	field := st.Field(0)
	fmt.Println(field.Tag.Get("color"), field.Tag.Get("species"))
	t.FailNow()
}

func TestWavesAdaptor_applyOpts(t *testing.T) {
	f := NewFactory()
	if f != nil {
	}

	ghClient := gravity.Client{}
	wvClient := wclient.Client{}
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

	type args struct {
		opts AdapterOptions
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    WavesAdaptor
	}{
		{name: "Test valid options", args: args{opts: validOpts}, wantErr: false, want: wanted},
		{name: "Test invaid ghClient", args: args{opts: AdapterOptions{
			"ghClient":        6,
			"wvClient":        &wvClient,
			"gravityContract": "contract",
			"chainID":         1}}, wantErr: true},
		{name: "Test invaid wClient", args: args{opts: AdapterOptions{
			"ghClient":        &ghClient,
			"wvClient":        5,
			"gravityContract": "contract",
			"chainID":         1}}, wantErr: true},
		{name: "Test invaid chainID", args: args{opts: AdapterOptions{
			"ghClient":        &ghClient,
			"wvClient":        5,
			"gravityContract": "contract",
			"chainID":         "karamba"}}, wantErr: true},
		{name: "Test invaid gravityContract", args: args{opts: AdapterOptions{
			"ghClient":        &ghClient,
			"wvClient":        5,
			"gravityContract": 56,
			"chainID":         1}}, wantErr: true},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wa := &WavesAdaptor{}
			if err := wa.applyOpts(tt.args.opts); (err != nil) != tt.wantErr {
				t.Errorf("WavesAdaptor.applyOpts() error = %v, wantErr %v", err, tt.wantErr)
			} else if !tt.wantErr {
				if *wa != tt.want {
					t.Errorf("ghClientValidator() = %v, want %v", *wa, tt.want)
				}
			}
		})
	}
}

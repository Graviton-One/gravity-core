package adaptors

import (
	"reflect"
	"testing"

	"github.com/Gravity-Tech/gravity-core/common/gravity"
)

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

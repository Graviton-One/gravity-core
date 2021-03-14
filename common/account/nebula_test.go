package account

import (
	"testing"
)

func initMapper() {
	ChainMapper.Assign(map[string]ChainType{"ethereum": 0, "heco": 3})
	ChainMapper.ApendAdaptor(3, "ethereum")
}

func TestNebulaId_ToString(t *testing.T) {
	initMapper()
	type args struct {
		chainType ChainType
	}
	tests := []struct {
		name string
		id   NebulaId
		args args
		want string
	}{
		{name: "test 1", id: NebulaId{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 236, 197, 129, 50, 32, 60, 168, 223, 237, 90, 125, 47, 37, 42, 170, 140, 77, 22, 32, 239},
			args: args{chainType: ChainType(3)}, want: "",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.id.ToString(tt.args.chainType); got != tt.want {
				t.Errorf("NebulaId.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

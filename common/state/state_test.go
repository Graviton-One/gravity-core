package state

import "testing"

func TestCalculateSubRound(t *testing.T) {
	type args struct {
		tcHeight       uint64
		blocksInterval uint64
	}
	tests := []struct {
		name string
		args args
		want SubRound
	}{

		{"test1", args{2788569, 20}, 0},
		{"test2", args{2788570, 20}, 1},
		{"test3", args{2788571, 20}, 1},
		{"test4", args{2788572, 20}, 1},

		{"test5", args{2788579, 20}, 1},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateSubRound(tt.args.tcHeight, tt.args.blocksInterval); got != tt.want {
				t.Errorf("CalculateSubRound() = %v, want %v", got, tt.want)
			}
		})
	}
}

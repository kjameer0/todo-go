package keypressinterface

import (
	"reflect"
	"testing"
)

func Test_generateMatrix(t *testing.T) {
	type args struct {
		rows  int
		cols  int
		items []string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "2x2 matrix filled",
			args: args{
				rows:  2,
				cols:  2,
				items: []string{"A", "B", "C", "D"},
			},
			want: [][]string{
				{"A", "B"},
				{"C", "D"},
			},
		},
		{
			name: "3x3 matrix with empty spaces",
			args: args{
				rows:  3,
				cols:  3,
				items: []string{"X", "Y", "Z"},
			},
			want: [][]string{
				{"X", "Y", "Z"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := generateMatrix(tt.args.rows, tt.args.cols, tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateMatrix() = %v, want %v", got, tt.want)
			}
		})
	}
}

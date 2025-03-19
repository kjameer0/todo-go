package keypressinterface

import (
	"reflect"
	"testing"
)

func Test_generateRows(t *testing.T) {
	type args struct {
		items       []string
		windowWidth int
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "single row",
			args: args{
				windowWidth: 80,
				items:       []string{"A", "B", "C", "D"},
			},
			want: [][]string{
				{"A", "B", "C", "D"},
			},
		},
		{
			name: "truncated long input",
			args: args{
				windowWidth: 80,
				items:       []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			},
			want: [][]string{
				{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..."},
			},
		},
		{
			name: "truncated long input with multiple lines",
			args: args{
				windowWidth: 80,
				items:       []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAaaa"},
			},
			want: [][]string{
				{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..."},
				{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..."},
			},
		},
		{
			name: "interspersed long and short items",
			args: args{
				windowWidth: 80,
				items: []string{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", "take out trash",
					"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaAaaa", "eat food", "workout"},
			},
			want: [][]string{
				{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..."},
				{"take out trash"},
				{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa..."},
				{"eat food", "workout"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generateRows(tt.args.items, tt.args.windowWidth); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

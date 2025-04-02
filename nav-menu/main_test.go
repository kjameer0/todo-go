package navmenu

import (
	"reflect"
	"testing"
)

func Test_createLookupPairings(t *testing.T) {
	type args struct {
		items []
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createLookupPairings(tt.args.items); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createLookupPairings() = %v, want %v", got, tt.want)
			}
		})
	}
}

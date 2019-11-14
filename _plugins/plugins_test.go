package plugins

import (
	"reflect"
	"testing"

	"github.com/tylerconlee/slab/config"
)

func TestLoadPlugins(t *testing.T) {
	type args struct {
		c config.Config
	}
	tests := []struct {
		name  string
		args  args
		wantP Plugins
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotP := LoadPlugins(tt.args.c); !reflect.DeepEqual(gotP, tt.wantP) {
				t.Errorf("LoadPlugins() = %v, want %v", gotP, tt.wantP)
			}
		})
	}
}

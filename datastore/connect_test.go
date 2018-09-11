package datastore

import (
	"testing"
)

func TestSave(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name       string
		args       args
		wantResult bool
	}{
		{
			name: "Save string",
			args: args{
				key:   "key",
				value: "value",
			},
			wantResult: true,
		},
	}
	// start Redis client

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RedisConnect(0)
			if gotResult := Save(tt.args.key, tt.args.value); gotResult != tt.wantResult {
				t.Errorf("Save() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

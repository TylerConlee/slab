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
	RedisConnect(8080)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := Save(tt.args.key, tt.args.value); gotResult != tt.wantResult {
				t.Errorf("Save() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestLoad(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name: "Load string",
			args: args{
				key: "key",
			},
			wantResult: "value",
		},
	}
	RedisConnect(8080)
	s := Save("key2", "value2")
	if s {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if gotResult := Load(tt.args.key); gotResult != tt.wantResult {
					t.Errorf("Load() = %v, want %v", gotResult, tt.wantResult)
				}
			})
		}
	} else {
		t.Errorf("Save failed on setup")
	}
}

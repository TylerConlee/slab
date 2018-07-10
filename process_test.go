package main

import (
	"testing"
	"time"
)

func TestRunTimer(t *testing.T) {
	type args struct {
		interval time.Duration
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunTimer(tt.args.interval)
		})
	}
}

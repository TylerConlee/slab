package main

import (
	"reflect"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_startServer(t *testing.T) {
	tests := []struct {
		name string
		want *Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := startServer(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("startServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_shutdown(t *testing.T) {
	type args struct {
		ticker *time.Ticker
		s      *Server
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdown(tt.args.ticker, tt.args.s)
		})
	}
}

func Test_keyCheck(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := keyCheck(); got != tt.want {
				t.Errorf("keyCheck() = %v, want %v", got, tt.want)
			}
		})
	}
}

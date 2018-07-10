package main

import (
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestServer_NewRouter(t *testing.T) {
	type fields struct {
		Server *http.Server
		Router *mux.Router
		Info   *ServerInfo
		Uptime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		want   *mux.Router
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server: tt.fields.Server,
				Router: tt.fields.Router,
				Info:   tt.fields.Info,
				Uptime: tt.fields.Uptime,
			}
			if got := s.NewRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Server.NewRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestServer_Callback(t *testing.T) {
	type fields struct {
		Server *http.Server
		Router *mux.Router
		Info   *ServerInfo
		Uptime time.Time
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server: tt.fields.Server,
				Router: tt.fields.Router,
				Info:   tt.fields.Info,
				Uptime: tt.fields.Uptime,
			}
			s.Callback(tt.args.w, tt.args.r)
		})
	}
}

func TestServer_Index(t *testing.T) {
	type fields struct {
		Server *http.Server
		Router *mux.Router
		Info   *ServerInfo
		Uptime time.Time
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				Server: tt.fields.Server,
				Router: tt.fields.Router,
				Info:   tt.fields.Info,
				Uptime: tt.fields.Uptime,
			}
			s.Index(tt.args.w, tt.args.r)
		})
	}
}

func TestWriteJSON(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		info   interface{}
		status int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			WriteJSON(tt.args.w, tt.args.info, tt.args.status)
		})
	}
}

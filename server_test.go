package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestServer_StartServer(t *testing.T) {
	type fields struct {
		Server *http.Server
		Router *mux.Router
		Info   *ServerInfo
		Uptime time.Time
	}
	tests := []struct {
		name   string
		fields fields
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
			s.StartServer()
		})
	}
}

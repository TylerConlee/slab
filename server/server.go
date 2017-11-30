package slack

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Server is an overarching type that contains the router, the server
// and any information displayed on the index page
type Server struct {
	Server *http.Server
	Router *mux.Router
	Info   *ServerInfo
	Uptime time.Time
}

// ServerInfo is used on the index status page which shows the Server name and
// the version of the bot
type ServerInfo struct {
	Server  string `json:"server"`
	Version string `json:"version"`
}

// ServerStatus contains the metadata from ServerInfo as well as the uptime for
// the bot
type ServerStatus struct {
	ServerInfo

	Uptime string `json:"uptime"`
}

// StartServer loads a http.Server and starts the Slack monitor
func (s *Server) StartServer() {
	s.Router = s.NewRouter()

	log.Info("HTTP server ready")
	go StartSlack()
	http.ListenAndServe(":8080", s.Router)

}

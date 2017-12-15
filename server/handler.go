package server

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	logging "github.com/op/go-logging"
	sl "github.com/tylerconlee/slab/slack"
	"github.com/tylerconlee/slack"
)

// log adds a logger for the `api` package
var log = logging.MustGetLogger("server")
var OnCall string

// NewRouter builds a new mux Router instance with the routes that
// Slack uses to handle callbacks, and the index status page
func (s *Server) NewRouter() *mux.Router {
	log.Debug("Building Router")
	r := mux.NewRouter()
	r.HandleFunc("/slack", s.Callback).Methods("POST")
	r.HandleFunc("/", s.Index).Methods("GET")
	return r
}

// SetOncall is a handler that handles the callback from a Slack action.
// TODO: Expand this to route to multiple callbacks, allowing for more
// functionality
func (s *Server) Callback(w http.ResponseWriter, r *http.Request) {
	payload := &slack.AttachmentActionCallback{}
	log.Debug(payload)
	SetTriager(payload, r)
}

// Index is a handler that outputs the server metadata and uptime in JSON form.
// TODO: add additional metadata for a more useful status page
func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	var status = ServerStatus{
		ServerInfo: ServerInfo{
			Server:  s.Info.Server,
			Version: s.Info.Version,
		},
		Uptime: time.Now().Sub(s.Uptime).String(),
	}

	WriteJSON(w, &status, http.StatusOK)
}

// WriteJSON takes the ResponseWriter, a generic structure of data and a status
// code and outputs it in JSON
func WriteJSON(w http.ResponseWriter, info interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Cache-Control", "no-store")

	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(info); err != nil {
		log.Debug("Failed to write JSON")
	} else {
		log.Debug(json.Marshal(info))
	}
}

func SetTriager(payload *slack.AttachmentActionCallback, r *http.Request) {
	err := json.Unmarshal([]byte(r.PostFormValue("payload")), payload)
	if err != nil {
		log.Critical("Unable to parse JSON for callback payload")
		os.Exit(1)
	}

	if len(payload.Actions) == 0 {
		return
	}
	log.Debug("Parsing action for callback")
	if sl.VerifyUser(payload.Actions[0].SelectedOptions[0].Value) {
		sl.OnCall = payload.Actions[0].SelectedOptions[0].Value
		OnCall = sl.OnCall
		sl.ChatUpdate(payload.Channel.ID, payload.MessageTs, payload.Actions[0].SelectedOptions[0].Value)
	}
}

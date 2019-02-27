package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/nlopes/slack"
	sl "github.com/tylerconlee/slab/slack"
)

// NewRouter builds a new mux Router instance with the routes that
// Slack uses to handle callbacks, and the index status page
func (s *Server) NewRouter() *mux.Router {
	log.Info("Building router", map[string]interface{}{
		"module": "server",
	})
	r := mux.NewRouter()
	r.HandleFunc("/slack", s.Callback).Methods("POST")
	r.HandleFunc("/api", s.Index).Methods("GET")
	return r
}

// Callback is a handler that handles the callback from a Slack action.
func (s *Server) Callback(w http.ResponseWriter, r *http.Request) {
	payload := &slack.AttachmentActionCallback{}
	err := json.Unmarshal([]byte(r.PostFormValue("payload")), payload)
	if err != nil {
		log.Error("Error unmarshaling callback payload.", map[string]interface{}{
			"module": "server",
			"error":  err,
		})
	}
	log.Info("Callback received.", map[string]interface{}{
		"module":   "server",
		"callback": payload.CallbackID,
	})
	switch payload.CallbackID {
	case "sla":
		if payload.Actions[0].Value == "ack" {
			sl.AcknowledgeSLA(payload)
		} else {
			sl.MoreInfoSLA(payload)
		}
	case "newticket":
		sl.AcknowledgeNewTicket(payload)
	case "triage_set":
		sl.SetTriager(payload)
	case "cfgwiz":
		log.Info("Config wizard step detected", map[string]interface{}{
			"module": "server",
			"step":   payload.Actions[0].Value,
		})
		sl.AddChannel(payload.Channel.ID, 1)
		if sl.ChannelSelect {
			sl.AddChannel(payload.Actions[0].SelectedOptions[0].Value, 2)
		}
		switch {
		case payload.Actions[0].Value == "start":
			sl.NextStep("start")
		case payload.Actions[0].Value == "view":
			sl.ViewConfig()
		case strings.Contains(payload.Actions[0].Value, "channel"):
			sl.NextStep(strings.Trim(payload.Actions[0].Value, "channel"))
		default:
			sl.NextStep(payload.Actions[0].SelectedOptions[0].Value)
		}
	}
	return
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
		log.Error("Failed to write JSON", map[string]interface{}{
			"module": "server",
			"error":  err,
		})
	} else {
		j, _ := json.Marshal(info)
		log.Debug("JSON write complete", map[string]interface{}{
			"module": "server",
			"json":   j,
		})
	}
}

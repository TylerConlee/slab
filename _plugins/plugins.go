package plugins

import (
	"strings"

	"github.com/nlopes/slack"
)

var Commands map[string]func([]string) ([]slack.Attachment, string)
var Send map[string]func(*Plugins, string)

func init() {
	Commands = make(map[string]func([]string) ([]slack.Attachment, string))
	Send = make(map[string]func(*Plugins, string))
}

// Plugins contains a list of all available plugins
type Plugins struct {
	Twilio    Twilio
	PagerDuty PagerDuty
}

// PagerDuty contains the connection details for the PagerDuty API:
// https://v2.developer.pagerduty.com/docs/rest-api
type PagerDuty struct {
	Email   string
	Key     string
	Enabled bool
}

// SendDispatcher receives the message from the process loop and checks which
// plugins are enabled and sends the appropriate notifications through them.
func (p *Plugins) SendDispatcher(message string) {
	log.Info("Plugins reached.", map[string]interface{}{
		"module": "plugin",
		"plugin": p,
	})
	for _, function := range Send {
		function(p, message)
	}
}

// ParsePluginCommand is ran before the core Slab commands are parsed to determine if a given command
// is related to a plugin or not. If it is, it returns a message and attachment which is then sent in
// slack/commands.go
func ParsePluginCommand(text string, user *slack.User) (message string, attachments []slack.Attachment) {
	t := strings.Fields(text)
	if len(t) > 1 {
		for command, function := range Commands {
			if t[1] == command {
				function(t)
			} else {
				attachments = []slack.Attachment{}
				message = ""

			}
		}
	}
	return
}

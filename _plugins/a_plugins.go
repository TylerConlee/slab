package plugins

import (
	"strings"

	"github.com/nlopes/slack"
	l "github.com/tylerconlee/slab/log"
)

// Commands is a map of all of the text commands needed to trigger individual plugins
var Commands map[string]func([]string) ([]slack.Attachment, string)

// Send is a map of every function used to send plugin messages
var Send map[string]func(*Plugins, string)
var log = l.Log

func init() {
	Commands = make(map[string]func([]string) ([]slack.Attachment, string))
	Send = make(map[string]func(*Plugins, string))
}

// Plugins contains a list of all available plugins
type Plugins struct {
	Twilio    Twilio
	PagerDuty PagerDuty
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
		log.Info("Plugin command received", map[string]interface{}{
			"module":  "plugin",
			"command": t,
		})
		for command, function := range Commands {
			if t[1] == command {
				attachments, message = function(t)
			} else {
				attachments = []slack.Attachment{}
				message = ""

			}
		}
	}
	return message, attachments
}

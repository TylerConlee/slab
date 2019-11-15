package plugins

import (
	"strings"

	"github.com/nlopes/slack"
)

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
	if TwilioPhone == "" {
		log.Info("To phone number for Twilio not set.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
	}
	if TwilioFrom == "" {
		log.Info("From phone number for Twilio not set.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
	}
	if (TwilioEnabled) && (TwilioPhone != "") {
		log.Info("Plugin loaded. Sending Twilio message.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
		p.SendTwilio(message)
	}
}

// ParsePluginCommand is ran before the core Slab commands are parsed to determine if a given command
// is related to a plugin or not. If it is, it returns a message and attachment which is then sent in
// slack/commands.go
func ParsePluginCommand(text string, user *slack.User) (message string, attachments []slack.Attachment) {
	t := strings.Fields(text)
	if len(t) > 1 {
		switch t[1] {
		case "twilio":
			t := strings.Fields(text)
			if len(t) > 1 {
				p := LoadPlugins()
				switch t[2] {
				case "set":
					if len(t) > 3 {
						s := TwilioSet(t[3])
						attachments = []slack.Attachment{s}
						message = "Plugin message"
					}
				case "unset":
					s := TwilioUnset()
					attachments = []slack.Attachment{s}
					message = "Plugin message"
				case "configure":
					if len(t) > 3 {
						s := TwilioConfigure(t[3])
						attachments = []slack.Attachment{s}
						message = "Plugin message"
					}
				case "status":
					s := p.TwilioStatus()
					attachments = []slack.Attachment{s}
					message = "Plugin status"
				case "enable":
					p.EnableTwilio()
					a := slack.Attachment{
						Title: "Twilio Plugin",
						Fields: []slack.AttachmentField{
							slack.AttachmentField{
								Title: "Enabled",
								Value: ":white_check_mark:",
							},
						},
					}
					attachments = []slack.Attachment{a}
					message = "Plugin Twilio has been updated"

				case "disable":
					p.DisableTwilio()
					a := slack.Attachment{
						Title: "Twilio Plugin",
						Fields: []slack.AttachmentField{
							slack.AttachmentField{
								Title: "Enabled",
								Value: ":x:",
							},
						},
					}
					attachments = []slack.Attachment{a}
					message = "Plugin Twilio has been updated"
				}
			}
		default:
			attachments = []slack.Attachment{}
			message = ""
		}

	}
	return
}

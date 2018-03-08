package slack

import (
	"strings"

	"github.com/tylerconlee/slab/plugins"
	"github.com/tylerconlee/slack"
)

// parseCommand takes the message that mentions the bot user and identifies
// what the user is asking for.
func parseCommand(text string, user string) {
	t := strings.Fields(text)
	switch t[1] {
	case "set":
		SetMessage()
	case "diag":
		DiagMessage(user)
	case "whois":
		WhoIsMessage()
	case "status":
		StatusMessage()
	case "help":
		HelpMessage()
	case "unset":
		UnsetMessage()
	case "twilio":
		switch t[2] {
		case "set":
			plugins.TwilioSet(t[3])
		case "unset":
			plugins.TwilioUnset()
		case "status":
			plugins.TwilioStatus()
		case "enable":
			plugins.EnableTwilio()
			a := slack.Attachment{
				Title: "Twilio Plugin",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Enabled",
						Value: ":white_check_mark:",
					},
				},
			}
			SendMessage("Plugin Twilio has been updated", a)

		case "disable":
			plugins.DisableTwilio()
			a := slack.Attachment{
				Title: "Twilio Plugin",
				Fields: []slack.AttachmentField{
					slack.AttachmentField{
						Title: "Enabled",
						Value: ":x:",
					},
				},
			}
			SendMessage("Plugin Twilio has been updated", a)
		}

	}

}

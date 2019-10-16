package slack

import (
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	"github.com/tylerconlee/slab/plugins"
)

// parseCommand takes the message that mentions the bot user and identifies
// what the user is asking for.
func parseCommand(text string, user *slack.User) {
	var attachment slack.Attachment
	t := strings.Fields(text)
	if len(t) > 1 {
		switch t[1] {
		case "set":

			attachment = SetMessage()
			attachments := []slack.Attachment{attachment}
			SendMessage("", c.Slack.ChannelID, attachments)
		case "diag":

			DiagMessage(user)
		case "whois":
			attachment = WhoIsMessage(user)
			attachments := []slack.Attachment{attachment}
			SendMessage("", c.Slack.ChannelID, attachments)
		case "status":
			StatusMessage(user)
		case "help":
			HelpMessage(user)
		case "unset":
			attachment = UnsetMessage(user)
			attachments := []slack.Attachment{attachment}
			SendMessage("Triager has been reset. Please use `@slab set` to set Triager.", c.Slack.ChannelID, attachments)

		case "tag":
			switch t[2] {
			case "create":
				CreateTagMessage(user)
			case "list":
				ListTagMessage(user)
			case "update":
				if len(t) > 3 {
					_, err := strconv.Atoi(t[3])
					if err != nil {
						UnknownCommandMessage(text, user.ID)
					}
					UpdateTagMessage(user, t[3])
				}
			case "delete":
				if len(t) > 3 {
					_, err := strconv.Atoi(t[3])
					if err != nil {
						UnknownCommandMessage(text, user.ID)
					}
					DeleteTagMessage(user, t[3])
				}
			}
		case "twilio":
			p := plugins.LoadPlugins(c)
			switch t[2] {
			case "set":
				if len(t) > 3 {
					s := plugins.TwilioSet(t[3])
					attachments := []slack.Attachment{s}
					SendMessage("Plugin message", c.Slack.ChannelID, attachments)
				}
			case "unset":
				s := plugins.TwilioUnset()
				attachments := []slack.Attachment{s}
				SendMessage("Plugin message", c.Slack.ChannelID, attachments)
			case "configure":
				if len(t) > 3 {
					s := plugins.TwilioConfigure(t[3])
					attachments := []slack.Attachment{s}
					SendMessage("Plugin message", c.Slack.ChannelID, attachments)
				}
			case "status":
				s := p.TwilioStatus()
				attachments := []slack.Attachment{s}
				SendMessage("Plugin status", c.Slack.ChannelID, attachments)
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
				attachments := []slack.Attachment{a}
				SendMessage("Plugin Twilio has been updated", c.Slack.ChannelID, attachments)

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
				attachments := []slack.Attachment{a}
				SendMessage("Plugin Twilio has been updated", c.Slack.ChannelID, attachments)
			}

		default:
			UnknownCommandMessage(text, user.ID)
		}
	}

}

func parseDMCommand(text string, user string) {
	t := strings.ToLower(text)
	switch t {
	case "start config":
		StartWizard(user)
		ConfirmWizard()
	default:
		UnknownCommandMessage(text, user)
	}

}

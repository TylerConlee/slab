package slack

import (
	"strconv"
	"strings"

	"github.com/nlopes/slack"
)

// parseCommand takes the message that mentions the bot user and identifies
// what the user is asking for.
func parseCommand(text string, user *slack.User) {
	var attachment slack.Attachment
	t := strings.Fields(text)
	if len(t) > 1 {
		// Send to plugins first to run through
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

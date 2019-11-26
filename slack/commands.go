package slack

import (
	"strconv"
	"strings"

	"github.com/nlopes/slack"
	plugins "github.com/tylerconlee/slab/_plugins"
)

// parseCommand takes the message that mentions the bot user and identifies
// what the user is asking for.
func parseCommand(text string, user *slack.User) {
	var attachment slack.Attachment
	attachments := []slack.Attachment{}
	message := ""
	t := strings.Fields(text)
	if len(t) > 1 {
		log.Info("Command identified", map[string]interface{}{
			"module":  "slack",
			"command": t,
		})

		message, attachments = plugins.ParsePluginCommand(text, user)

		if message == "" {
			switch t[1] {
			case "set":
				attachment = SetMessage()
				attachments = []slack.Attachment{attachment}
			case "diag":

				DiagMessage(user)
			case "whois":
				attachment = WhoIsMessage(user)
				attachments = []slack.Attachment{attachment}
			case "status":
				attachments = StatusMessage(user)
			case "help":
				message = HelpMessage(user)
			case "history":
				attachments = HistoryMessage(user)
			case "unset":
				attachment = UnsetMessage(user)
				attachments = []slack.Attachment{attachment}
				message = "Triager has been reset. Please use `@slab set` to set Triager."

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
		} else {
			log.Info("Plugin command processed", map[string]interface{}{
				"module":  "slack",
				"command": t,
				"message": message,
			})
		}

		SendMessage(message, c.Slack.ChannelID, attachments)
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

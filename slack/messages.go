package slack

import (
	"fmt"
	"os"
	"strings"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/zendesk"

	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slack"
)

var (
	c   = config.LoadConfig()
	log = logging.MustGetLogger("slack")
)

func SendMessage(attachment slack.Attachment, message string) {
	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{attachment}
	params.LinkNames = 1
	// Send a message to the given channel with pretext and the parameters
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// Log message if succesfully sent.
	log.Infof("Message successfully sent to channel %s at %s", channelID, timestamp)
}

// SetMessage creates and sends a message to Slack with a menu attachment,
// allowing users to set the OnCall staff member.
func SetMessage() {
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "triage_set",
		// Show the current OnCall member
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Currently OnCall",
				Value: fmt.Sprintf("<@%s>", OnCall),
			},
		},
		// Show a dropdown of all users to select new OnCall target
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:       "triage_select",
				Text:       "Select Team Member",
				Type:       "select",
				Style:      "primary",
				DataSource: "users",
			},
		},
	}
	SendMessage(attachment, "...")
}

// WhoIsMessage creates and sends a Slack message that sends out the value of
// OnCall.
func WhoIsMessage() {
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "triage_whois",
		// Show the current OnCall member
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Currently OnCall",
				Value: fmt.Sprintf("<@%s>", OnCall),
			},
		},
	}
	SendMessage(attachment, "...")
}

// Send sends off the SLA notification to Slack using the configured API key
func SLAMessage(n string, ticket zendesk.ActiveTicket) {
	description := ticket.Description
	if len(ticket.Description) > 100 {
		description = description[0:100] + "..."
	}
	url := fmt.Sprintf("%s/agent/tickets/%d", c.Zendesk.URL, ticket.ID)
	attachment := slack.Attachment{
		// Uncomment the following part to send a field too
		Title:      ticket.Subject,
		TitleLink:  url,
		CallbackID: "sla",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Description",
				Value: description,
			},
			slack.AttachmentField{
				Title: "Priority",
				Value: strings.Title(ticket.Priority.(string)),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Created At",
				Value: ticket.CreatedAt.String(),
				Short: true,
			},
		},
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "ack_sla",
				Text:  "Acknowledge",
				Type:  "button",
				Style: "primary",
				Confirm: &slack.ConfirmationField{
					Text:        "Are you sure?",
					OkText:      "Take it",
					DismissText: "Leave it",
				},
			},
		},
	}
	SendMessage(attachment, n)
}

// ChatUpdate takes a channel ID, a timestamp and message text
// and updated the message in the given Slack channel at the given
// timestamp with the given message text. Currently, it also updates the
// attachment specifically for the Set message output.
func ChatUpdate(
	payload *slack.AttachmentActionCallback,
	attachment slack.Attachment,
) {

	params := slack.PostMessageParameters{}

	params.Attachments = []slack.Attachment{attachment}
	// Send an update to the given channel with pretext and the parameters
	channelID, timestamp, t, err := api.UpdateMessageWithParams(
		payload.Channel.ID,
		payload.OriginalMessage.Timestamp,
		payload.OriginalMessage.Text,
		params,
	)
	log.Debug(channelID, timestamp, t, err)
}

// parseCommand takes the message that mentions the bot user and identifies
// what the user is asking for.
func parseCommand(text string) {
	t := strings.Fields(text)
	switch t[1] {
	case "set":
		SetMessage()
	case "whois":
		WhoIsMessage()
	}

}

// verifyUser takes a User ID string and runs the Slack GetUserInfo request. If
// the user exists, the function returns true.
func VerifyUser(user string) bool {
	_, err := api.GetUserInfo(user)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	return true
}

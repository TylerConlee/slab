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

// SetMessage creates and sends a message to Slack with a menu attachment,
// allowing users to set the OnCall staff member.
func SetMessage() {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "oncall_dropdown",
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
				Name:       "oncall_select",
				Text:       "Select Team Member",
				Type:       "select",
				Style:      "primary",
				DataSource: "users",
			},
		},
	}

	// Add the attachment to the parameters of a new message
	params.Attachments = []slack.Attachment{attachment}

	// Send a message to the given channel with pretext and the parameters
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, "...", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// Log message if succesfully sent.
	log.Infof("Set message successfully sent to channel %s at %s", channelID, timestamp)
}

// WhoIsMessage creates and sends a Slack message that sends out the value of
// OnCall.
func WhoIsMessage() {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "oncall_dropdown",
		// Show the current OnCall member
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Currently OnCall",
				Value: fmt.Sprintf("<@%s>", OnCall),
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	// Send a message to the given channel with pretext and the parameters
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, "...", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// Send a message to the given channel with pretext and the parameters
	log.Infof("WhoIs message successfully sent to channel %s at %s", channelID, timestamp)
}

// Send sends off the SLA notification to Slack using the configured API key
func SendSLAMessage(n string, ticket zendesk.ActiveTicket) {

	color := "warning"

	if strings.Contains(n, "@sup") {
		color = "danger"
	}
	description := ticket.Description
	if len(ticket.Description) > 100 {
		description = description[0:100] + "..."
	}
	params := slack.PostMessageParameters{}
	url := fmt.Sprintf("%s/agent/tickets/%d", c.Zendesk.URL, ticket.ID)
	attachment := slack.Attachment{
		Color: color,
		// Uncomment the following part to send a field too
		Title:     ticket.Subject,
		TitleLink: url,
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
	}
	params.LinkNames = 1
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, n, params)
	if err != nil {
		log.Criticalf("%s\n", err)
		os.Exit(1)
	}
	log.Debugf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

// ChatUpdate takes a channel ID, a timestamp and message text
// and updated the message in the given Slack channel at the given
// timestamp with the given message text. Currently, it also updates the
// attachment specifically for the Set message output.
func ChatUpdate(channel string, ts string, text string) {
	t := fmt.Sprintf("Updated OnCall person to <@%s>", text)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "oncall_dropdown",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	params.Attachments = []slack.Attachment{attachment}
	text = "..."
	// Send an update to the given channel with pretext and the parameters
	channelID, timestamp, t, err := api.UpdateMessageWithParams(channel, ts, text, params)
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

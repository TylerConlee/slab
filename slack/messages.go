package slack

import (
	"fmt"
	"os"
	"strings"
	"time"

	logging "github.com/op/go-logging"

	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slack"
)

var (
	c       = config.LoadConfig()
	log     = logging.MustGetLogger("slack")
	uptime  time.Time
	version string
)

// SendMessage takes an attachment and message and composes a message to be
// sent to the configured Slack channel ID
func SendMessage(message string, attachment slack.Attachment) {
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
// allowing users to set the triager staff member.
func SetMessage() {
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the triager here.",
		CallbackID: "triage_set",
		// Show the current triager
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Triager",
				Value: fmt.Sprintf("<@%s>", Triager),
			},
		},

		// Show a dropdown of all users to select new Triager target
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
	SendMessage("...", attachment)
}

// WhoIsMessage creates and sends a Slack message that sends out the value of
// Triager.
func WhoIsMessage() {
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the triager here.",
		CallbackID: "triage_whois",
		// Show the current Triager member
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Triager",
				Value: fmt.Sprintf("<@%s>", Triager),
			},
		},
	}
	SendMessage("...", attachment)
}

// Ticket represents an individual ticket to be used in SLAMessage and
// NewTicketMessage
type Ticket struct {
	ID          int
	Subject     string
	SLA         []interface{}
	Tags        []string
	Level       string
	Priority    interface{}
	CreatedAt   time.Time
	Description string
}

// SLAMessage sends off the SLA notification to Slack using the configured API key
func SLAMessage(n string, ticket Ticket) {
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
	SendMessage(n, attachment)
}

// NewTicketMessage takes a slice of tickets that have been created in the last
// loop interval and sends the IDs and links to the tickets to the user
// currently set as triager.
func NewTicketMessage(tickets []Ticket) {
	params := slack.PostMessageParameters{}
	for _, ticket := range tickets {
		attachment := slack.Attachment{
			Title: ticket.Subject,
			TitleLink: fmt.Sprintf(
				"%s/agent/tickets/%d",
				c.Zendesk.URL,
				ticket.ID,
			),
			Footer: ticket.CreatedAt.String(),
		}
		params.Attachments = append(params.Attachments, attachment)

	}
	message := fmt.Sprintf("The following tickets were received since the last loop:")

	if len(params.Attachments) != 0 {
		_, _, channelID, err := api.OpenIMChannel(Triager)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		api.PostMessage(channelID, message, params)
	}
}

// StatusMessage responds to @slab status with the version hash and current
// uptime for the Slab process
func StatusMessage() {
	attachment := slack.Attachment{
		Title: "Slab Status",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Version",
				Value: version,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Uptime",
				Value: time.Now().Sub(uptime).String(),
				Short: true,
			},
		},
	}
	SendMessage("...", attachment)
}

// HelpMessage responds to @slab help with a help message outlining all
// available commands
func HelpMessage() {

	params := slack.PostMessageParameters{}

	setCommand := slack.Attachment{
		Title: "@slab set",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Command Name",
				Value: "Set",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Command Description",
				Value: "Used to set the active Triager, returns a dropdown of users",
				Short: true,
			},
		},
		Footer:     fmt.Sprintf("Current triager: %v", Triager),
		FooterIcon: "https://slack-files2.s3-us-west-2.amazonaws.com/avatars/2018-01-05/294943756277_b467ce1bf3a88bdb8a6a_512.png",
	}
	whoisCommand := slack.Attachment{
		Title: "@slab whois",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Command Name",
				Value: "Whois",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Command Description",
				Value: "Returns the name of the user currently set as Triager",
				Short: true,
			},
		},
		Footer:     fmt.Sprintf("Current triager: %v", Triager),
		FooterIcon: "https://slack-files2.s3-us-west-2.amazonaws.com/avatars/2018-01-05/294943756277_b467ce1bf3a88bdb8a6a_512.png",
	}
	statusCommand := slack.Attachment{
		Title: "@slab status",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Command Name",
				Value: "Status",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Command Description",
				Value: "Returns metadata about the Slab instance currently running",
				Short: true,
			},
		},
		Footer:     fmt.Sprintf("Current uptime: %v", time.Now().Sub(uptime).String()),
		FooterIcon: "https://slack-files2.s3-us-west-2.amazonaws.com/avatars/2018-01-05/294943756277_b467ce1bf3a88bdb8a6a_512.png",
	}

	attachments := []slack.Attachment{setCommand, whoisCommand, statusCommand}
	params.Attachments = attachments
	message := "..."
	api.PostMessage(c.Slack.ChannelID, message, params)

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
	case "status":
		StatusMessage()
	case "help":
		HelpMessage()
	}

}

// VerifyUser takes a User ID string and runs the Slack GetUserInfo request. If
// the user exists, the function returns true.
func VerifyUser(user string) bool {
	_, err := api.GetUserInfo(user)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	return true
}

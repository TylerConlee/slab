package slack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slack"
)

var (
	c       = config.LoadConfig()
	uptime  time.Time
	version string
	// Sent represents the NotifySent from the zendesk package
	Sent interface{}
)

// Ticket represents an individual ticket to be used in SLAMessage and
// NewTicketMessage
type Ticket struct {
	ID          int
	Requester   int64
	Subject     string
	SLA         []interface{}
	Tags        []string
	Level       string
	Priority    interface{}
	CreatedAt   time.Time
	Description string
}

// NotifySent is represetative of an individual ticket, what kind of
// notification was last sent for that ticket, and when the SLA breach time is.
type NotifySent struct {
	ID     int
	Type   int64
	Expire time.Time
}

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
	log.Debug("Message sent successfully.", map[string]interface{}{
		"module":    "slack",
		"channel":   channelID,
		"timestamp": timestamp,
		"message":   message,
	})
}

// SendEphemeralMessage takes a message, attachment and a user ID and sends a
// message to that user ID.
func SendEphemeralMessage(message string, attachment slack.Attachment, user string) {
	params := slack.PostMessageParameters{}
	params.Attachments = []slack.Attachment{attachment}
	params.LinkNames = 1

	// Send a message to the given channel with pretext and the parameters
	timestamp, err := api.PostEphemeral(c.Slack.ChannelID, user, slack.MsgOptionText(message, params.EscapeText),
		slack.MsgOptionAttachments(params.Attachments...),
		slack.MsgOptionPostMessageParameters(params))
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// Log message if succesfully sent.
	log.Debug("Message sent successfully.", map[string]interface{}{
		"module":    "slack",
		"timestamp": timestamp,
		"message":   message,
	})
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
				Name:  "triage_select",
				Text:  ":white_check_mark: Set",
				Type:  "button",
				Value: "ack",
				Style: "primary",
				Confirm: &slack.ConfirmationField{
					Text:        "Are you sure?",
					OkText:      "Take it",
					DismissText: "Leave it",
				},
			},
		},
	}
	SendMessage("...", attachment)
}

// UnsetMessage resets the Triager role to the slab bot.
func UnsetMessage() {
	message := "Triager has been reset. Please use `@slab set` to set Triager."

	Triager = "None"
	t := fmt.Sprintf("Triager has been reset to %s", Triager)
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the triager here.",
		CallbackID: "triager_dropdown",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	SendMessage(message, attachment)
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

// SLAMessage sends off the SLA notification to Slack using the configured API key
func SLAMessage(n string, ticket Ticket, color string, user string, uid int64) {
	description := ticket.Description
	if len(ticket.Description) > 100 {
		description = description[0:100] + "..."
	}
	url := fmt.Sprintf("%s/agent/tickets/%d", c.Zendesk.URL, ticket.ID)
	link := fmt.Sprintf("%s/agent/users/%d", c.Zendesk.URL, uid)
	attachment := slack.Attachment{
		// Uncomment the following part to send a field too
		Title:      ticket.Subject,
		TitleLink:  url,
		AuthorName: user,
		AuthorLink: link,
		AuthorIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/google/119/bust-in-silhouette_1f464.png",
		CallbackID: "sla",
		Color:      color,
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
				Text:  ":white_check_mark: Acknowledge",
				Type:  "button",
				Value: "ack",
				Style: "primary",
				Confirm: &slack.ConfirmationField{
					Text:        "Are you sure?",
					OkText:      "Take it",
					DismissText: "Leave it",
				},
			},
			slack.AttachmentAction{
				Name:  "more_info_sla",
				Value: strconv.FormatInt(ticket.Requester, 10),
				Text:  ":mag: More Info",
				Type:  "button",
				Style: "default",
			},
		},
	}
	SendMessage(n, attachment)
}

// DiagMessage sends a DM to requestor with the current state of SLA
// notifications for tickets
func DiagMessage(user string) {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Title: "Slab Diagnostic Tool",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Notification Status",
				Value: fmt.Sprintf("%x", Sent),
			},
			slack.AttachmentField{
				Title: "Uptime",
				Value: time.Now().Sub(uptime).String(),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Current Triager",
				Value: Triager,
				Short: true,
			},
		},
		Footer:     fmt.Sprintf("Version: %s", version),
		FooterIcon: "https://slack-files2.s3-us-west-2.amazonaws.com/avatars/2018-01-05/294943756277_b467ce1bf3a88bdb8a6a_512.png",
	}
	params.Attachments = append(params.Attachments, attachment)
	message := ""
	if len(params.Attachments) != 0 {
		_, _, channelID, err := api.OpenIMChannel(user)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		api.PostMessage(channelID, message, params)
	}
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
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Ticket ID",
					Value: strconv.Itoa(ticket.ID),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Created At",
					Value: ticket.CreatedAt.String(),
					Short: true,
				},
			},
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
		Footer:     fmt.Sprintf("Current triager: <@%s>", Triager),
		FooterIcon: "https://slack-files2.s3-us-west-2.amazonaws.com/avatars/2018-01-05/294943756277_b467ce1bf3a88bdb8a6a_512.png",
	}
	unsetCommand := slack.Attachment{
		Title: "@slab unset",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Command Name",
				Value: "Unset",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Command Description",
				Value: "Used to unset the active Triager. Sets the Triager to 'None'",
				Short: true,
			},
		},
		Footer:     fmt.Sprintf("Current triager: <@%s>", Triager),
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
		Footer:     fmt.Sprintf("Current triager: <@%s>", Triager),
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
	diagCommand := slack.Attachment{
		Title: "@slab diag",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Command Name",
				Value: "Diag",
				Short: true,
			},
			slack.AttachmentField{
				Title: "Command Description",
				Value: "Sends a private message to the requestor with diagnostic information about Slab",
				Short: true,
			},
		},
		Footer:     fmt.Sprintf("Current uptime: %v", time.Now().Sub(uptime).String()),
		FooterIcon: "https://slack-files2.s3-us-west-2.amazonaws.com/avatars/2018-01-05/294943756277_b467ce1bf3a88bdb8a6a_512.png",
	}
	attachments := []slack.Attachment{
		setCommand,
		unsetCommand,
		whoisCommand,
		statusCommand,
		diagCommand,
	}
	params.Attachments = attachments
	params.LinkNames = 1
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
	log.Debug("Message updated.", map[string]interface{}{
		"module":    "slack",
		"channel":   channelID,
		"timestamp": timestamp,
		"message":   t,
		"error":     err,
	})
}

// VerifyUser takes a User ID string and runs the Slack GetUserInfo request. If
// the user exists, the function returns true.
func VerifyUser(user string) bool {
	_, err := api.GetUserInfo(user)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	return true
}

// PrepSLANotification takes a given ticket and what notification level and returns a string to be sent to Slack.
func PrepSLANotification(ticket Ticket, notify int64) (notification string, color string) {
	log.Info("Preparing SLA notification message.", map[string]interface{}{
		"module": "slack",
		"ticket": ticket.ID,
	})
	var t, p string
	var r bool

	switch ticket.Level {
	case "LevelOne":
		p = c.SLA.LevelOne.Tag
		r = c.SLA.LevelOne.Notify
	case "LevelTwo":
		p = c.SLA.LevelTwo.Tag
		r = c.SLA.LevelTwo.Notify

	case "LevelThree":
		p = c.SLA.LevelThree.Tag
		r = c.SLA.LevelThree.Notify

	case "LevelFour":
		p = c.SLA.LevelFour.Tag
		r = c.SLA.LevelFour.Notify
	}

	var n, c string

	switch notify {
	case 1:
		t = "15 minutes"
		c = "danger"
	case 2:
		t = "30 minutes"
		c = "warning"
	case 3:
		t = "1 hour"
		c = "#ffec1e"
	case 4:
		t = "2 hours"
		c = "#439fe0"
	case 5:
		t = "3 hours"
		c = "#43e0d3"
	}
	if r {
		n = fmt.Sprintf("@here SLA for *%s* ticket #%d has less than %s until expiration.", p, ticket.ID, t)
		if notify == 9 {
			n = fmt.Sprintf("@here Expired *%s* SLA! Ticket #%d has an expired SLA.", p, ticket.ID)
			c = "danger"
		}
	} else {
		n = fmt.Sprintf("SLA for *%s* ticket #%d has less than %s until expiration.", p, ticket.ID, t)
		if notify == 9 {
			n = fmt.Sprintf("Expired *%s* SLA! Ticket #%d has an expired SLA.", p, ticket.ID)
			c = "danger"
		}
	}

	return n, c

}

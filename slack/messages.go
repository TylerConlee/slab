package slack

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/datastore"
	"github.com/tylerconlee/slab/zendesk"
)

var (
	c       = config.LoadConfig()
	uptime  time.Time
	version string
	// Sent represents the NotifySent from the zendesk package
	Sent interface{}
	// NumTickets is the number of tickets processed on the last loop
	NumTickets int
	// LastProcessed is a timestamp of when the last loop was ran
	LastProcessed time.Time
	// deleteTag is a tag ID that will be deleted
	deleteTag string
	// update is a tag ID that will be updated
	update string
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
	UpdatedAt   time.Time
	Description string
}

// NotifySent is represetative of an individual ticket, what kind of
// notification was last sent for that ticket, and when the SLA breach time is.
type NotifySent struct {
	ID     int
	Type   int64
	Expire time.Time
}

// SetMessage creates and sends a message to Slack with a menu attachment,
// allowing users to set the triager staff member.
func SetMessage() (attachment slack.Attachment) {
	attachment = slack.Attachment{
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
	return attachment
}

// UnsetMessage resets the Triager role to the slab bot.
func UnsetMessage(user *slack.User) (attachment slack.Attachment) {
	Triager = "None"
	if err := datastore.SaveActivity(user.ID, user.Name, "unset"); err != nil {
		log.Error("Unable to save activity", map[string]interface{}{
			"module":   "slack",
			"activity": "unset",
			"error":    err,
		})
	}
	t := fmt.Sprintf("Triager has been reset to %s", Triager)
	attachment = slack.Attachment{
		Fallback:   "You would be able to select the triager here.",
		CallbackID: "triager_dropdown",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	return attachment
}

// WhoIsMessage creates and sends a Slack message that sends out the value of
// Triager.
func WhoIsMessage(user *slack.User) (attachment slack.Attachment) {
	if err := datastore.SaveActivity(user.ID, user.Name, "whois"); err != nil {
		log.Error("Unable to save activity", map[string]interface{}{
			"module":   "slack",
			"activity": "whois",
			"error":    err,
		})
	}
	attachment = slack.Attachment{
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
	return attachment
}

// SLAMessage sends off the SLA notification to Slack using the configured API key
func SLAMessage(ticket Ticket, color string, user string, uid int64, org string) (attachment slack.Attachment) {
	description := ticket.Description
	if len(ticket.Description) > 100 {
		description = description[0:100] + "..."
	}
	url := fmt.Sprintf("%s/agent/tickets/%d", c.Zendesk.URL, ticket.ID)
	link := fmt.Sprintf("%s/agent/users/%d", c.Zendesk.URL, uid)
	attachment = slack.Attachment{
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
				Title: "Organization",
				Value: org,
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
	return attachment
}

// DiagMessage sends a DM to requestor with the current state of SLA
// notifications for tickets
func DiagMessage(user *slack.User) {
	s := Sent.([]zendesk.NotifySent)
	attachment := slack.Attachment{

		Title: "Slab Diagnostic Tool",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Number of Ticket Notifications",
				Value: fmt.Sprintf("%v", len(s)),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Number of Tickets Processed",
				Value: fmt.Sprintf("%v", len(s)),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Last Process Loop Ran",
				Value: fmt.Sprintf("%v", LastProcessed.Format("Mon Jan 2 15:04:05 MST")),
			},

			slack.AttachmentField{
				Title: "Current Notification Status",
				Value: fmt.Sprintf("%v", s),
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
	attachments := []slack.Attachment{attachment}
	message := ""
	if len(attachments) != 0 {
		_, _, channelID, err := api.OpenIMChannel(user.ID)
		if err != nil {
			fmt.Printf("%s\n", err)
		}
		api.PostMessage(channelID, slack.MsgOptionText(message, false), slack.MsgOptionAttachments(attachments...))
	}
}

// NewTicketMessage takes a slice of tickets that have been created in the last
// loop interval and sends the IDs and links to the tickets to the user
// currently set as triager.
func NewTicketMessage(tickets []Ticket, tag string) (newTickets []slack.Attachment, message string) {
	attachments := []slack.Attachment{}
	for _, ticket := range tickets {
		description := ticket.Description
		if len(ticket.Description) > 100 {
			description = description[0:100] + "..."
		}
		attachment := slack.Attachment{
			Title: ticket.Subject,
			TitleLink: fmt.Sprintf(
				"%s/agent/tickets/%d",
				c.Zendesk.URL,
				ticket.ID,
			),
			ID:         ticket.ID,
			CallbackID: "newticket",
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Description",
					Value: description,
				},
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
				slack.AttachmentField{
					Title: "Tag",
					Value: tag,
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
			},
		}
		attachments = []slack.Attachment{attachment}

	}
	message = ""

	return attachments, message

}

// HistoryMessage responds to @slab history with the last 10 commands that were run
func HistoryMessage(user *slack.User) (attachments []slack.Attachment) {
	opts := datastore.ActivityOptions{
		Quantity: 10,
	}
	activities, err := datastore.LoadActivity(opts)
	if err != nil {
		log.Error("Unable to run command", map[string]interface{}{
			"module":   "slack",
			"activity": "history",
			"error":    err,
		})
	}
	for _, activity := range activities {
		attachment := slack.Attachment{
			Title:      activity.ActivityType,
			AuthorName: activity.SlackName,
			AuthorLink: activity.SlackID,
			AuthorIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/google/119/bust-in-silhouette_1f464.png",
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "Started At",
					Value: activity.StartedAt.String(),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Ended At",
					Value: activity.EndedAt.String(),
					Short: true,
				},
			},
		}
		attachments = append(attachments, attachment)
	}

	if err := datastore.SaveActivity(user.ID, user.Name, "history"); err != nil {
		log.Error("Unable to save activity", map[string]interface{}{
			"module":   "slack",
			"activity": "history",
			"error":    err,
		})
	}
	return attachments
}

// StatusMessage responds to @slab status with the version hash and current
// uptime for the Slab process
func StatusMessage(user *slack.User) (attachments []slack.Attachment) {
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
	if err := datastore.SaveActivity(user.ID, user.Name, "status"); err != nil {
		log.Error("Unable to save activity", map[string]interface{}{
			"module":   "slack",
			"activity": "status",
			"error":    err,
		})
	}
	attachments = []slack.Attachment{attachment}
	return attachments
}

// HelpMessage responds to @slab help with a help message outlining all
// available commands
func HelpMessage(user *slack.User) (message string) {
	if err := datastore.SaveActivity(user.ID, user.Name, "help"); err != nil {
		log.Error("Unable to save activity", map[string]interface{}{
			"module":   "slack",
			"activity": "help",
			"error":    err,
		})
	}
	params := slack.PostMessageParameters{}

	params.LinkNames = 1
	message = "Help for Slab can be found at <https://github.com/TylerConlee/slab/wiki|the Slab wiki>"
	return message

}

// ShowConfigMessage takes a user string and sends that user the value of the
// config.toml configuration file. Used for identifying configuration issues.
func ShowConfigMessage(user string) {
	attachment := []slack.Attachment{
		slack.Attachment{
			Title: "Config",
		},
	}
	message := "Test direct message for config."
	SendDirectMessage(message, attachment, user)

}

// UnknownCommandMessage sends a direct message to the user provided indicating
// that the command that they attempted is not a valid command.
func UnknownCommandMessage(text string, user string) {
	message := fmt.Sprintf("Sorry, the command `%s` is an invalid command. Please type `help` for a list of all available commands", text)
	attachment := []slack.Attachment{}
	SendDirectMessage(message, attachment, user)
}

// VerifyUser takes a User ID string and runs the Slack GetUserInfo request. If
// the user exists, the function returns true.
func VerifyUser(user string) bool {
	_, err := api.GetUserInfo(user)
	if err != nil {
		log.Error("Error verifying user", map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	return true
}

// PrepSLANotification takes a given ticket and what notification level and returns a string to be sent to Slack.
func PrepSLANotification(ticket Ticket, notify int64, tag string) (notification string, color string) {
	log.Info("Preparing SLA notification message.", map[string]interface{}{
		"module": "slack",
		"ticket": ticket.ID,
	})
	var t, n, c string

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
	n = fmt.Sprintf("<!here> SLA for *%s* ticket #%d has less than %s until expiration.", tag, ticket.ID, t)
	if notify == 9 {
		n = fmt.Sprintf("<!here> Expired *%s* SLA! Ticket #%d has an expired SLA.", tag, ticket.ID)
		c = "danger"
	}

	return n, c

}

// UpdateMessage sends a message to the channel indicating a ticket with a
// premium SLA tag associated with it has received an update. This functionality
// is a mirror of the official Zendesk > Slack integration.
func UpdateMessage(ticket Ticket, user string, uid int64) (attachment slack.Attachment) {
	description := ticket.Description
	if len(ticket.Description) > 100 {
		description = description[0:100] + "..."
	}
	url := fmt.Sprintf("%s/agent/tickets/%d", c.Zendesk.URL, ticket.ID)
	link := fmt.Sprintf("%s/agent/users/%d", c.Zendesk.URL, uid)
	attachment = slack.Attachment{
		// Uncomment the following part to send a field too
		Title:      ticket.Subject,
		TitleLink:  url,
		AuthorName: user,
		AuthorLink: link,
		AuthorIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/google/119/bust-in-silhouette_1f464.png",
		CallbackID: "sla",
		Color:      "primary",
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
				Title: "Updated At",
				Value: ticket.UpdatedAt.String(),
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
		},
	}
	return attachment
}

// CreateTagMessage responds to @slab tag create, taking the tag name provided
// and responds with the first step in the create tag config wizard
func CreateTagMessage(user *slack.User) {
	attachment := slack.Attachment{
		Title:      "Create Tag",
		CallbackID: "createtag",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "launch",
				Text:  "Launch Tag Creator",
				Type:  "button",
				Value: "createtag",
			},
		},
	}
	api.PostMessage(c.Slack.ChannelID, slack.MsgOptionAttachments(attachment))
}

// ListTagMessage grabs the tags stored in the database and outputs them to
// Slack in a DM
func ListTagMessage(user *slack.User) {
	tags := datastore.LoadTags()
	log.Info("Tags received from database", map[string]interface{}{
		"module": "slack",
		"tags":   tags,
	})
	var attachments []slack.Attachment
	for _, tag := range tags {
		i := int(tag["id"].(int64))
		id := strconv.Itoa(i)
		attachment := slack.Attachment{
			Title:      strings.ToTitle(tag["tag"].(string)),
			AuthorName: fmt.Sprintf("<@%s>", tag["user"]),
			Footer:     fmt.Sprintf("Last updated at %s", tag["updated_at"].(string)),
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "ID",
					Value: id,
					Short: true,
				},
				slack.AttachmentField{
					Title: "Channel",
					Value: fmt.Sprintf("<#%s>", tag["channel"]),
					Short: true,
				},
				slack.AttachmentField{
					Title: "Notification Type",
					Value: getNotificationType(tag["notify_type"].(string)),
					Short: true,
				},
			},
		}
		attachments = append(attachments, attachment)
	}
	message := "These are the current active tags for Slab"
	SendDirectMessage(message, attachments, user.ID)
}

// UpdateTagMessage takes the id for a given tag and opens a dialog for
// updating the entry
func UpdateTagMessage(user *slack.User, id string) {
	attachment := slack.Attachment{
		Title:      "Update Tag",
		CallbackID: "updatetag",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "launch",
				Text:  "Launch Tag Creator",
				Type:  "button",
				Value: id,
			},
		},
	}
	update = id
	api.PostMessage(c.Slack.ChannelID, slack.MsgOptionAttachments(attachment))
}

// DeleteTagMessage takes an id and deletes the corresponding tag in the
// database
func DeleteTagMessage(user *slack.User, id string) {
	tagID, err := strconv.Atoi(id)
	if err != nil {
		log.Error("Error converting ID to integer", map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	tag := datastore.LoadTag(tagID)
	attachment := slack.Attachment{
		Title:      strings.ToTitle(tag["tag"].(string)),
		CallbackID: "tagdelete",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Channel",
				Value: fmt.Sprintf("<#%s>", tag["channel"]),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Created By",
				Value: fmt.Sprintf("<@%s>", tag["user"]),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Notification Type",
				Value: getNotificationType(tag["notify_type"].(string)),
				Short: true,
			},
		},
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "deletetag",
				Text:  ":x: Delete Tag",
				Type:  "button",
				Value: "deletetag",
				Style: "Danger",
				Confirm: &slack.ConfirmationField{
					Text:        "Are you sure?",
					OkText:      "Delete",
					DismissText: "Cancel",
				},
			},
		},
	}
	deleteTag = id
	api.PostMessage(c.Slack.ChannelID, slack.MsgOptionAttachments(attachment))

}

func getNotificationType(notify string) (full string) {
	switch notify {
	case "new":
		return "New Tickets"
	case "sla":
		return "SLA Breach Alerts"
	case "updates":
		return "Ticket Updates"
	}
	return
}

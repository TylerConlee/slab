package slack

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nlopes/slack"
	"github.com/tylerconlee/slab/datastore"
	"github.com/tylerconlee/slab/zendesk"
)

// newTS stores the timestamp of the new tag dialog message
// so that the original message can be updated upon success
var newTS string

// updateTS stores the timestamp of the update tag dialog message
// so that the original message can be updated upon success
var updateTS string

var newAttachmentID string

var updateAttachmentID string

// SetTriager generates a new Slack attachment to update the
// original message and set the Triager role
func SetTriager(payload *slack.InteractionCallback) {
	if len(payload.Actions) == 0 {
		return
	}

	if VerifyUser(payload.User.ID) {
		Triager = payload.User.ID
		datastore.RSave("triager", payload.User.ID)
		if err := datastore.SaveActivity(payload.User.ID, payload.User.Name, "set"); err != nil {
			log.Error("Unable to save activity", map[string]interface{}{
				"module":   "slack",
				"activity": "set",
				"triager":  Triager,
				"error":    err,
			})
		}
		t := fmt.Sprintf("<@%s> is now set as Triager", Triager)
		log.Info("Triager set.", map[string]interface{}{
			"module":  "slack",
			"triager": Triager,
		})
		attachment := slack.Attachment{
			Fallback:   "You would be able to select the triager here.",
			CallbackID: "triager_dropdown",
			Footer:     t,
			FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
		}
		ChatUpdate(payload, attachment)
	}
}

// AcknowledgeSLA generates a new Slack attachment to state that a user has
// acknowledged a ticket.
func AcknowledgeSLA(payload *slack.InteractionCallback) {
	f, _ := strconv.ParseFloat(payload.ActionTs, 10)
	i := int64(f)
	ts := time.Unix(i, 0)
	log.Info("Ticket acknowledged.", map[string]interface{}{
		"module":   "slack",
		"actionts": payload.ActionTs,
		"ts":       ts.String(),
		"i":        i,
	})
	t := fmt.Sprintf("<@%s> acknowledged this ticket at %s", payload.User.Name, ts.String())
	log.Info("SLA ticket acknowledged.", map[string]interface{}{
		"module": "slack",
		"ack":    payload.User.Name,
	})
	attachment := slack.Attachment{}
	for i := range payload.OriginalMessage.Attachments {
		id := strconv.Itoa(payload.OriginalMessage.Attachments[i].ID)
		if id == payload.AttachmentID {
			attachment = slack.Attachment{
				Title:      payload.OriginalMessage.Attachments[i].Title,
				TitleLink:  payload.OriginalMessage.Attachments[i].TitleLink,
				Fallback:   "User acknowledged a ticket.",
				CallbackID: "newticket",
				Footer:     t,
				FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
			}
		}
	}

	ChatUpdate(payload, attachment)
}

// AcknowledgeNewTicket generates a new Slack attachment to state that a user
//has acknowledged a ticket.
func AcknowledgeNewTicket(payload *slack.InteractionCallback) {
	f, _ := strconv.ParseFloat(payload.ActionTs, 10)
	i := int64(f)
	ts := time.Unix(i, 0)
	log.Info("Ticket acknowledged.", map[string]interface{}{
		"module":   "slack",
		"actionts": payload.ActionTs,
		"ts":       ts.String(),
		"i":        i,
	})
	t := fmt.Sprintf("<@%s> acknowledged this ticket at %s", payload.User.Name, ts.String())
	log.Info("New ticket acknowledged.", map[string]interface{}{
		"module": "slack",
		"ack":    payload.User.Name,
	})
	attachment := slack.Attachment{}
	for i := range payload.OriginalMessage.Attachments {
		id := strconv.Itoa(payload.OriginalMessage.Attachments[i].ID)
		if id == payload.AttachmentID {
			attachment = slack.Attachment{
				Title:      payload.OriginalMessage.Attachments[i].Title,
				TitleLink:  payload.OriginalMessage.Attachments[i].TitleLink,
				Fallback:   "User acknowledged a ticket.",
				CallbackID: "newticket",
				Footer:     t,
				FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
			}
		}
	}

	ChatUpdate(payload, attachment)
}

// MoreInfoSLA grabs additional information from Zendesk using the information
// from the More Info button. It then sends  an ephemeral message to the
// requester with additional Zendesk information.
func MoreInfoSLA(payload *slack.InteractionCallback) {
	log.Info("More info SLA button clicked.", map[string]interface{}{
		"module": "slack",
		"ticket": payload.Actions[0].Value,
	})
	if err := datastore.SaveActivity(payload.User.ID, payload.User.Name, "moreinfo"); err != nil {
		log.Error("Unable to save activity", map[string]interface{}{
			"module":   "slack",
			"activity": "moreinfo",
			"triager":  Triager,
			"error":    err,
		})
	}
	id, _ := strconv.Atoi(payload.Actions[0].Value)
	// ORG Details
	org := zendesk.GetOrganization(id)
	// REQUESTED tickets
	requested := zendesk.GetRequestedTickets(id)
	name := "Not Set"
	if requested.Tickets[0].AssigneeID != nil {
		assignee := zendesk.GetTicketRequester(int(requested.Tickets[0].AssigneeID.(float64)))
		name = assignee.Name
	}
	var t string
	var s string
	for _, ticket := range requested.Tickets {
		i := strconv.Itoa(ticket.ID)
		status := statusDecode(ticket.Status)
		sat := satisfactionDecode(ticket.SatisfactionRating.Score)
		link := fmt.Sprintf("%s/agent/tickets/%d", c.Zendesk.URL, ticket.ID)
		t = t + "<" + link + "| #" + i + " (" + status + ")> "
		s = s + sat + ", "
	}
	o := ""
	if len(org) > 0 {
		orglink := fmt.Sprintf("%s/agent/organizations/%d", c.Zendesk.URL, org[0].ID)

		o = "<" + orglink + "| " + org[0].Name + "> "
	}
	len := len(requested.Tickets) - 1
	attachment := slack.Attachment{
		Fallback: "User acknowledged a ticket.",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Organization",
				Value: o,
			},
			slack.AttachmentField{
				Title: "Ticket Last Updated",
				Value: requested.Tickets[len].UpdatedAt.String(),
				Short: true,
			},
			slack.AttachmentField{
				Title: "Ticket Assigned To",
				Value: name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Tickets from this user",
				Value: t,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Satisfaction history",
				Value: s,
				Short: true,
			},
		},
	}
	SendEphemeralMessage("More information on ticket", attachment, payload.User.ID)
}

// CreateTagDialog receives the payload from the incoming callback
// and opens a dialog box allowing the user to create a tag they want to be
// notified on.
func CreateTagDialog(payload *slack.InteractionCallback) {
	newTS = payload.OriginalMessage.Timestamp
	newAttachmentID = payload.AttachmentID
	log.Info("Create tag dialog launching", map[string]interface{}{
		"module":    "slack",
		"timestamp": newTS,
	})
	dialog := slack.Dialog{
		TriggerID:      payload.TriggerID,
		CallbackID:     "process_create_tag",
		Title:          "Create Tag",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			slack.DialogInput{
				Type:        "text",
				Label:       "Tag",
				Name:        "tag",
				Placeholder: "Insert the tag to be notified on",
				Optional:    false,
			},
			slack.DialogInputSelect{
				DialogInput: slack.DialogInput{
					Type:     "select",
					Label:    "Channel",
					Name:     "channel",
					Optional: false,
				},
				DataSource: "conversations",
			},
			slack.DialogInputSelect{
				DialogInput: slack.DialogInput{
					Type:     "select",
					Label:    "Notification Type",
					Name:     "notify_type",
					Optional: false,
				},
				Options: []slack.DialogSelectOption{
					slack.DialogSelectOption{
						Label: "New Tickets",
						Value: "new",
					},
					slack.DialogSelectOption{
						Label: "SLA Breaches",
						Value: "sla",
					},
					slack.DialogSelectOption{
						Label: "Ticket Updates",
						Value: "updates",
					},
				},
			},
		},
	}
	api.OpenDialog(payload.TriggerID, dialog)
}

// UpdateTagDialog receives the payload from the incoming callback
// and opens a dialog box allowing the user to update a tag they want to be
// notified on.
func UpdateTagDialog(payload *slack.InteractionCallback) {
	updateTS = payload.OriginalMessage.Timestamp
	updateAttachmentID = payload.AttachmentID
	log.Info("Update tag dialog launching", map[string]interface{}{
		"module":    "slack",
		"timestamp": updateTS,
	})
	id, err := strconv.Atoi(payload.Actions[0].Value)
	if err != nil {
		log.Error("Error converting ID to integer", map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	tag := datastore.LoadTag(id)
	dialog := slack.Dialog{
		TriggerID:      payload.TriggerID,
		CallbackID:     "process_update_tag",
		Title:          "Update Tag",
		NotifyOnCancel: false,
		Elements: []slack.DialogElement{
			slack.DialogInput{
				Type:        "text",
				Label:       "Tag",
				Name:        "tag",
				Placeholder: tag["tag"].(string),
				Optional:    false,
			},
			slack.DialogInputSelect{
				DialogInput: slack.DialogInput{
					Type:        "select",
					Label:       "Channel",
					Name:        "channel",
					Optional:    false,
					Placeholder: tag["channel"].(string),
				},
				DataSource: "conversations",
			},
			slack.DialogInputSelect{
				DialogInput: slack.DialogInput{
					Type:        "select",
					Label:       "Notification Type",
					Name:        "notify_type",
					Optional:    false,
					Placeholder: tag["notify_type"].(string),
				},

				Options: []slack.DialogSelectOption{
					slack.DialogSelectOption{
						Label: "New Tickets",
						Value: "new",
					},
					slack.DialogSelectOption{
						Label: "SLA Breaches",
						Value: "sla",
					},
					slack.DialogSelectOption{
						Label: "Ticket Updates",
						Value: "updates",
					},
				},
			},
		},
	}
	api.OpenDialog(payload.TriggerID, dialog)
}

// UpdateTag takes the input collected from the user and updates a tag based
// on the tag ID provided
func UpdateTag(payload *slack.InteractionCallback) {
	var data map[string]string
	log.Info("Dialog saved", map[string]interface{}{
		"module": "slack",
		"user":   payload.User.ID,
		"data":   payload.DialogSubmissionCallback.Submission,
	})
	data = payload.DialogSubmissionCallback.Submission
	data["user"] = payload.User.ID
	data["id"] = update
	datastore.SaveTagUpdate(data)
	update = ""
	t := fmt.Sprintf("Tag '%s' created by <@%s>", data["tag"], data["user"])
	attachment := slack.Attachment{
		Fallback:   t,
		CallbackID: "triager_dropdown",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	payload.AttachmentID = updateAttachmentID
	payload.OriginalMessage.Timestamp = updateTS
	log.Info("Updating tag", map[string]interface{}{
		"module":    "slack",
		"timestamp": payload.OriginalMessage.Timestamp,
	})
	updateTS = ""
	ChatUpdate(payload, attachment)
}

// SaveDialog takes the input collected from the Create Tag Dialog and
// sends the data to Postgres to be saved
func SaveDialog(payload *slack.InteractionCallback) {
	var data map[string]string
	log.Info("Dialog saved", map[string]interface{}{
		"module": "slack",
		"user":   payload.User.ID,
		"data":   payload.DialogSubmissionCallback.Submission,
	})
	data = payload.DialogSubmissionCallback.Submission
	data["user"] = payload.User.ID
	datastore.SaveNewTag(data)
	t := fmt.Sprintf("Tag '%s' created by <@%s>", data["tag"], data["user"])
	attachment := slack.Attachment{
		Fallback:   t,
		CallbackID: "triager_dropdown",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	payload.AttachmentID = newAttachmentID
	payload.OriginalMessage.Timestamp = newTS
	log.Info("Saving new tag", map[string]interface{}{
		"module":    "slack",
		"timestamp": payload.OriginalMessage.Timestamp,
	})
	newTS = ""
	ChatUpdate(payload, attachment)
}

// DeleteTag deletes a tag based on the tag ID provided
func DeleteTag(payload *slack.InteractionCallback) {
	datastore.DeleteTag(deleteTag)
	deleteTag = ""
	attachment := slack.Attachment{
		Fallback:   "Tag deleted",
		CallbackID: "triager_dropdown",
		Footer:     "Tag deleted successfully",
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	ChatUpdate(payload, attachment)
}

func statusDecode(status string) (img string) {
	switch status {
	case "solved":
		img = ":white_check_mark:"
	case "new":
		img = ":new:"
	case "open":
		img = ":o2:"
	case "pending":
		img = ":parking:"
	case "closed":
		img = ":lock:"
	}
	return
}

func satisfactionDecode(sat string) (s string) {
	switch sat {
	case "good":
		s = ":white_check_mark:"
	case "bad":
		s = ":x:"
	case "unoffered":
		s = ":heavy_minus_sign:"
	}
	return
}

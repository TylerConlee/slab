package slack

import (
	"fmt"
	"strconv"

	"github.com/tylerconlee/slab/zendesk"
	"github.com/tylerconlee/slack"
)

// SetTriager generates a new Slack attachment to update the
// original message and set the Triager role
func SetTriager(payload *slack.AttachmentActionCallback) {
	if len(payload.Actions) == 0 {
		return
	}

	if VerifyUser(payload.Actions[0].SelectedOptions[0].Value) {
		Triager = payload.Actions[0].SelectedOptions[0].Value
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
func AcknowledgeSLA(payload *slack.AttachmentActionCallback) {
	t := fmt.Sprintf("<@%s> acknowledged this ticket", payload.User.Name)
	log.Info("SLA ticket acknowledged.", map[string]interface{}{
		"module": "slack",
		"ack":    payload.User.Name,
	})
	attachment := slack.Attachment{
		Title:      payload.OriginalMessage.Attachments[0].Title,
		TitleLink:  payload.OriginalMessage.Attachments[0].TitleLink,
		Fallback:   "User acknowledged a ticket.",
		CallbackID: "sla",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	ChatUpdate(payload, attachment)
}

func MoreInfoSLA(payload *slack.AttachmentActionCallback) {
	log.Info("More info SLA button clicked.", map[string]interface{}{
		"module": "slack",
		"ticket": payload.Actions[0].Value,
	})
	id, _ := strconv.Atoi(payload.Actions[0].Value)
	// ORG Details
	org := zendesk.GetOrganization(id)
	attachment := slack.Attachment{
		Fallback:   "User acknowledged a ticket.",
		Footer:     payload.Actions[0].Value,
		FooterIcon: ":mag:",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Organization",
				Value: org[0].Name,
				Short: true,
			},
		},
	}
	SendEphemeralMessage("More information on ticket", attachment, payload.User.ID)
}

package slack

import (
	"fmt"
	"reflect"
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
	// REQUESTED tickets
	requested := zendesk.GetRequestedTickets(id)
	name := "Not Set"
	if !reflect.ValueOf(requested.Tickets[0].AssigneeID.(int)).IsNil() {
		assignee := zendesk.GetTicketRequester(requested.Tickets[0].AssigneeID.(int))
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
	orglink := fmt.Sprintf("%s/agent/organizations/%d", c.Zendesk.URL, org[0].ID)
	o := "<" + orglink + "| " + org[0].Name + "> "
	attachment := slack.Attachment{
		Fallback: "User acknowledged a ticket.",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Organization",
				Value: o,
			},
			slack.AttachmentField{
				Title: "Ticket Last Updated",
				Value: requested.Tickets[0].UpdatedAt.String(),
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

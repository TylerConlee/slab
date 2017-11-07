package slack

import (
	"fmt"
	"os"
	"strings"

	"github.com/tylerconlee/slab/zendesk"

	"github.com/nlopes/slack"
	"github.com/tylerconlee/slab/config"
)

var c = config.LoadConfig()

func Send(n string, ticket zendesk.ActiveTicket) {
	api := slack.New(c.Slack.APIKey)

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

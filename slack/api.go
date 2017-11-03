package slack

import (
	"fmt"

	"github.com/nlopes/slack"
	"github.com/tylerconlee/slab/config"
)

var c = config.LoadConfig()

func Send(n string) {
	api := slack.New(c.Slack.APIKey)

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text: n,
		// Uncomment the following part to send a field too
		/*
			Fields: []slack.AttachmentField{
				slack.AttachmentField{
					Title: "a",
					Value: "no",
				},
			},
		*/
	}
	params.Attachments = []slack.Attachment{attachment}
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, "SLA Warning:", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

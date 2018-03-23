package slack

import (
	"fmt"

	"github.com/tylerconlee/slack"
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

// SendDirectMessage takes a message, an attachment and a user and sends a
// direct message to the user.
func SendDirectMessage(message string, attachment slack.Attachment, user string) {

}

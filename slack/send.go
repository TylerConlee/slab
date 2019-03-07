package slack

import (
	"fmt"
	"strconv"

	"github.com/nlopes/slack"
)

// SendMessage takes an attachment and message and composes a message to be
// sent to the configured Slack channel ID
func SendMessage(message string, channel string, attachments []slack.Attachment) {
	params := slack.PostMessageParameters{}
	params.LinkNames = 1
	// Send a message to the given channel with pretext and the parameters
	channelID, timestamp, err := api.PostMessage(channel, slack.MsgOptionText(message, false), slack.MsgOptionAttachments(attachments...))
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
	attachments := []slack.Attachment{attachment}
	params.LinkNames = 1

	// Send a message to the given channel with pretext and the parameters
	timestamp, err := api.PostEphemeral(c.Slack.ChannelID, user, slack.MsgOptionText(message, params.EscapeText),
		slack.MsgOptionAttachments(attachments...),
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
func SendDirectMessage(message string, attachments []slack.Attachment, user string) {
	params := slack.PostMessageParameters{}
	params.LinkNames = 1
	_, _, channelID, err := api.OpenIMChannel(user)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	api.PostMessage(channelID, slack.MsgOptionText(message, false), slack.MsgOptionAttachments(attachments...))
}

// ChatUpdate takes a channel ID, a timestamp and message text
// and updated the message in the given Slack channel at the given
// timestamp with the given message text. Currently, it also updates the
// attachment specifically for the Set message output.
func ChatUpdate(
	payload *slack.InteractionCallback,
	attachment slack.Attachment,
) {

	for i := range payload.OriginalMessage.Attachments {
		id := strconv.Itoa(payload.OriginalMessage.Attachments[i].ID)
		if id == payload.AttachmentID {
			payload.OriginalMessage.Attachments[i] = attachment
		}
	}

	attachments := payload.OriginalMessage.Attachments
	// Send an update to the given channel with pretext and the parameters
	channelID, timestamp, t, err := api.UpdateMessage(
		payload.Channel.ID,
		payload.OriginalMessage.Timestamp,
		slack.MsgOptionText(payload.OriginalMessage.Text, false),
		slack.MsgOptionAttachments(attachments...),
	)
	log.Info("Message updated.", map[string]interface{}{
		"module":    "slack",
		"channel":   channelID,
		"timestamp": timestamp,
		"message":   t,
		"error":     err,
	})
}

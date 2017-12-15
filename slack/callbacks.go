package slack

import (
	"fmt"

	"github.com/tylerconlee/slack"
)

func SetTriager(payload *slack.AttachmentActionCallback) {
	if len(payload.Actions) == 0 {
		return
	}
	log.Debug("Parsing action for callback")
	if VerifyUser(payload.Actions[0].SelectedOptions[0].Value) {
		OnCall = payload.Actions[0].SelectedOptions[0].Value
		t := fmt.Sprintf("<@%s> is now set as Triager", OnCall)
		attachment := slack.Attachment{
			Fallback:   "You would be able to select the triager here.",
			CallbackID: "triager_dropdown",
			Footer:     t,
			FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
		}
		ChatUpdate(payload.Channel.ID, payload.MessageTs, attachment)
	}
}

func AcknowledgeSLA(payload *slack.AttachmentActionCallback) {
	t := fmt.Sprintf("<@%s> acknowledged this ticket", payload.User.Name)
	attachment := slack.Attachment{
		Fallback:   "User acknowledged a ticket.",
		CallbackID: "ack_sla",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	ChatUpdate(payload.Channel.ID, payload.MessageTs, attachment)
}

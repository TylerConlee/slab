package slack

import (
	"strconv"

	"github.com/tylerconlee/slack"
)

var (
	activeWizard      bool
	activeUser        configUser
	ChannelsRemaining int
	ChannelSelect     bool
)

type configUser struct {
	user string
	step int
}

// StartWizard takes the user ID, sets the configUser and starts the
// ConfigSetupMessage function
func StartWizard(user string) {
	activeUser.user = user
	activeUser.step = 0
	ConfirmWizard()

}

// ConfigInProgressMessage takes a user ID string and sends a message to that
// user letting them know that there's already a configuration wizard in
// progress to avoid overlap.
func ConfigInProgressMessage(user string) {
	message := "Oops! The configuration wizard is currently being used by another user. Please try again later."
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, user)
}

// ConfirmWizard sends a confirmation message to the user letting them know
// that the changes made will overwrite the current configuration.
func ConfirmWizard() {
	message := "Hi! Let's get started using Slab!"
	attachment := slack.Attachment{
		Title:      "Warning! Using the Slab configuration wizard will overwrite the current configuration. Please select an option below.",
		CallbackID: "cfgwiz",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:  "confirm",
				Text:  "Start Wizard",
				Type:  "button",
				Value: "start",
			},
			slack.AttachmentAction{
				Name:  "view",
				Text:  "View Current Configuration",
				Type:  "button",
				Value: "view",
			},
			// TODO: Add in button for diagnostic info
		},
	}
	SendDirectMessage(message, attachment, activeUser.user)
}

// ConfigSetupMessage sends the first message to the specified user to start
// the configuration setup wizard process.
func ConfigSetupMessage() {
	activeUser.step = 1
	message := "To start, how many channels need Slab alerts?"
	attachment := slack.Attachment{
		Title: "Number of Channels",

		CallbackID: "cfgwiz",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name: "channels_list",
				Text: "Channel for Slab",
				Type: "select",
				Options: []slack.AttachmentActionOption{
					slack.AttachmentActionOption{
						Text:  "1",
						Value: "1",
					},
					slack.AttachmentActionOption{
						Text:  "2",
						Value: "2",
					},
					slack.AttachmentActionOption{
						Text:  "3",
						Value: "3",
					},
				},
			},
		},
	}
	SendDirectMessage(message, attachment, activeUser.user)
}

// ChannelSelectMessage takes a user string and sends that user a direct message
// asking for a channel to be selected that Slab will monitor/send alerts to.
func ChannelSelectMessage() {
	if ChannelsRemaining > 0 {
		attachment := slack.Attachment{
			Title:      "Channels",
			CallbackID: "cfgwiz",
			Actions: []slack.AttachmentAction{
				slack.AttachmentAction{
					Name:       "channels_list",
					Text:       "Channel for Slab",
					Type:       "select",
					DataSource: "channels",
				},
			},
		}
		ChannelSelect = true
		SendDirectMessage("Select a channel", attachment, activeUser.user)
		ChannelsRemaining--
	} else {
		ChannelSelect = false
		activeUser.step = 2
	}
}

func GetZendeskURL() {
	activeUser.step = 3
	message := "Please enter your Zendesk URL"
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, activeUser.user)
}

func GetZendeskUser() {
	activeUser.step = 4
	message := "Please enter your Zendesk username"
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, activeUser.user)
}

func GetZendeskAPIKey() {
	activeUser.step = 5
	message := "Please enter your Zendesk API Key"
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, activeUser.user)
}

func NextStep(msg string) {
	log.Info("Processing next step", map[string]interface{}{
		"module":  "slack",
		"message": msg,
		"step":    activeUser.step,
	})
	switch activeUser.step {
	case 0:
		ConfigSetupMessage()
		// set ChannelsRemaining, have the callback check ChannelsRemaining and subtract one until all channels are taken care of
	case 1:
		if ChannelsRemaining == 0 {
			var err error
			ChannelsRemaining, err = strconv.Atoi(msg)
			if err != nil {
				log.Error("error parsing channels remaining value", map[string]interface{}{
					"error": err,
				})
			}
		}
		ChannelSelect = true
		ChannelSelectMessage()
	case 2:
		GetZendeskURL()
	case 3:
		log.Info("Zendesk URL received.", map[string]interface{}{
			"module":  "slack",
			"message": msg,
		})
		GetZendeskUser()
	case 4:
		log.Info("Zendesk username received.", map[string]interface{}{
			"module":  "slack",
			"message": msg,
		})
		GetZendeskAPIKey()
	case 5:
		log.Info("Zendesk API key received.", map[string]interface{}{
			"module":  "slack",
			"message": msg,
		})
	}

}

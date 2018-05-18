package slack

import (
	"strings"
	"time"

	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slack"
)

var (
	activeWizard      bool
	activeUser        configUser
	ChannelsRemaining int
	ChannelSelect     bool
	ZenAPI            string
	ZenUser           string
	ZenURL            string
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
	activeWizard = true
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
	activeUser.step = 1
	SendDirectMessage(message, attachment, activeUser.user)
}

func ViewConfig() {
	message := "Here's the current configuration for Slab:"
	attachment := prepConfigLoad()
	SendDirectMessage(message, attachment, activeUser.user)
}

// ChannelSelectMessage takes a user string and sends that user a direct message
// asking for a channel to be selected that Slab will monitor/send alerts to.
func ChannelSelectMessage() {
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
	activeUser.step = 2
	SendDirectMessage("Select a channel for Slab to report in.", attachment, activeUser.user)

}

// GetZendeskURL is a step in the configuration wizard that asks the user to
// provide a URL to their Zendesk instance.
func GetZendeskURL() {
	activeUser.step = 3
	ChannelSelect = false
	message := "Please enter your Zendesk URL"
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, activeUser.user)
}

// GetZendeskUser is a step in the configuration wizard that asks the user to
// provide the username to access Zendesk with.
func GetZendeskUser() {
	activeUser.step = 4
	message := "Please enter your Zendesk username"
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, activeUser.user)
}

// GetZendeskAPIKey is a step in the configuration wizard that asks the user to
// provide the API key to access Zendesk with.
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
	case 1:
		ChannelSelectMessage()
	case 2:
		GetZendeskURL()
	case 3:
		log.Info("Zendesk URL received.", map[string]interface{}{
			"module":  "slack",
			"message": msg,
		})
		ZenURL = parseInput(msg)
		GetZendeskUser()
	case 4:
		log.Info("Zendesk username received.", map[string]interface{}{
			"module":  "slack",
			"message": msg,
		})
		ZenUser = parseInput(msg)
		GetZendeskAPIKey()
	case 5:
		log.Info("Zendesk API key received.", map[string]interface{}{
			"module":  "slack",
			"message": msg,
		})
		ZenAPI = msg

		prepConfigSave()

	}

}

func parseInput(input string) (output string) {
	r := strings.NewReplacer("<", "", ">", "")
	output = r.Replace(input)
	output = strings.TrimLeft(output, "|")
	return
}
func prepConfigLoad() (attachment slack.Attachment) {
	log.Info("Loading current configuration", map[string]interface{}{
		"module": "slack",
	})
	con := config.LoadConfig()

	attachment = slack.Attachment{
		Title: "Current Configuration",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Zendesk URL",
				Value: con.Zendesk.URL,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Zendesk User",
				Value: con.Zendesk.User,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Zendesk API Key",
				Value: con.Zendesk.APIKey[len(con.Zendesk.APIKey)-5:],
			},
			slack.AttachmentField{
				Title: "Update Frequency",
				Value: con.UpdateFreq.Duration.String(),
			},
		},
	}
	return attachment
}

func prepConfigSave() {
	con := config.Config{
		Zendesk: config.Zendesk{
			APIKey: ZenAPI,
			User:   ZenUser,
			URL:    ZenURL,
		},
		SLA: config.SLA{
			LevelOne: config.Level{
				Tag: "platinum",
				Low: config.Duration{
					time.Duration(180 * time.Minute),
				},
				Normal: config.Duration{
					time.Duration(120 * time.Minute),
				},
				High: config.Duration{
					time.Duration(60 * time.Minute),
				},
				Urgent: config.Duration{
					time.Duration(59 * time.Minute),
				},
				Notify: true,
			},
			LevelTwo: config.Level{
				Tag: "gold",
				Low: config.Duration{
					time.Duration(6 * time.Hour),
				},
				Normal: config.Duration{
					time.Duration(3 * time.Hour),
				},
				High: config.Duration{
					time.Duration(2 * time.Hour),
				},
				Urgent: config.Duration{
					time.Duration(1 * time.Hour),
				},
				Notify: true,
			},
			LevelThree: config.Level{
				Tag: "silver",
				Low: config.Duration{
					time.Duration(12 * time.Hour),
				},
				Normal: config.Duration{
					time.Duration(6 * time.Hour),
				},
				High: config.Duration{
					time.Duration(3 * time.Hour),
				},
				Urgent: config.Duration{
					time.Duration(2 * time.Hour),
				},
				Notify: true,
			},
			LevelFour: config.Level{
				Tag: "bronze",
				Low: config.Duration{
					time.Duration(48 * time.Hour),
				},
				Normal: config.Duration{
					time.Duration(24 * time.Hour),
				},
				High: config.Duration{
					time.Duration(8 * time.Hour),
				},
				Urgent: config.Duration{
					time.Duration(3 * time.Hour),
				},
				Notify: true,
			},
		},
		Metadata:      config.Metadata{},
		TriageEnabled: true,
		Slack: config.Slack{
			ChannelID: ChannelList[0].ID,
		},
	}
	log.Info("Saving current configuration", map[string]interface{}{
		"module": "slack",
		"config": con,
	})
	success := config.SaveConfig(con)
	attachment := slack.Attachment{}
	if success {

		SendDirectMessage("Configuration successfully saved", attachment, activeUser.user)
	} else {
		SendDirectMessage("An error occurred when attempting to save the configuration.", attachment, activeUser.user)
	}
	activeUser.step = 0
	activeUser.user = ""
	activeWizard = false
}

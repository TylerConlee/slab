package slack

import "github.com/tylerconlee/slack"

var (
	activeWizard bool
	activeUser   configUser
)

type configUser struct {
	user string
	step string
}

// ConfigInProgressMessage takes a user ID string and sends a message to that
// user letting them know that there's already a configuration wizard in
// progress to avoid overlap.
func ConfigInProgressMessage(user string) {
	message := "Oops! The configuration wizard is currently being used by another user. Please try again later."
	attachment := slack.Attachment{}
	SendDirectMessage(message, attachment, user)
}

// ConfigSetupMessage sends the first message to the specified user to start
// the configuration setup wizard process.
func ConfigSetupMessage(user string) {
	activeUser.step = "1"
	activeUser.user = user
	message := "Hi! Let's get Slab set up! First, how many channels need access?"
	attachment := slack.Attachment{
		Title: "Number of Channels",
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
	SendDirectMessage(message, attachment, user)
}

// ChannelSelectMessage takes a user string and sends that user a direct message
// asking for a channel to be selected that Slab will monitor/send alerts to.
func ChannelSelectMessage(user string) {
	attachment := slack.Attachment{
		Title: "Channels",
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:       "channels_list",
				Text:       "Channel for Slab",
				Type:       "select",
				DataSource: "channels",
			},
		},
	}
	SendDirectMessage("Select a channel", attachment, user)
}

func nextStep() {

}

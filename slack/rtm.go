package slack

import (
	"fmt"
	"os"
	"strings"

	"github.com/tylerconlee/slack"
)

// api is an instance of the tylerconlee/slack Client
var api *slack.Client

// OnCall holds the User ID of the current person set as "OnCall"
var OnCall string

// StartSlack initializes a connection with the given slack instance, gets
// team information, and starts a Go channel with the Real Time Messaging
// API watcher.
func StartSlack() {
	log.Info("Starting connection to Slack")
	// start a connection to Slack using the Slack Bot token

	api = slack.New(c.Slack.APIKey)

	// retrieve the team info for the newly connected Slack team
	d, err := api.GetTeamInfo()
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	log.Info("Connected to Slack:", d.Domain)

	// Set the initial value of OnCall
	OnCall = "None"

	// Start monitoring Slack
	startRTM()

}

// SetMessage creates and sends a message to Slack with a menu attachment,
// allowing users to set the OnCall staff member.
func SetMessage() {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "oncall_dropdown",
		// Show the current OnCall member
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Currently OnCall",
				Value: fmt.Sprintf("<@%s>", OnCall),
			},
		},
		// Show a dropdown of all users to select new OnCall target
		Actions: []slack.AttachmentAction{
			slack.AttachmentAction{
				Name:       "oncall_select",
				Text:       "Select Team Member",
				Type:       "select",
				Style:      "primary",
				DataSource: "users",
			},
		},
	}

	// Add the attachment to the parameters of a new message
	params.Attachments = []slack.Attachment{attachment}

	// Send a message to the given channel with pretext and the parameters
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, "...", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// Log message if succesfully sent.
	log.Infof("Set message successfully sent to channel %s at %s", channelID, timestamp)
}

// WhoIsMessage creates and sends a Slack message that sends out the value of
// OnCall.
func WhoIsMessage() {
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "oncall_dropdown",
		// Show the current OnCall member
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Currently OnCall",
				Value: fmt.Sprintf("<@%s>", OnCall),
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}
	// Send a message to the given channel with pretext and the parameters
	channelID, timestamp, err := api.PostMessage(c.Slack.ChannelID, "...", params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	// Send a message to the given channel with pretext and the parameters
	log.Infof("WhoIs message successfully sent to channel %s at %s", channelID, timestamp)
}

// startRTM creates a separate Go channel which monitors the Slack instance.
// The RTM tracks each and every event within Slack and allows the bot to act
// accordingly.
func startRTM() {
	log.Debug(api)
	rtm := api.NewRTM()
	chk := 0
	var user *slack.User
	var err error
	go rtm.ManageConnection()

	// When a new event occurs in Slack, track it here
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {

		// When a user connects to Slack for the first time. Logged message
		// shows number of already connected users.
		case *slack.ConnectedEvent:
			log.Debug("Connection counter:", ev.ConnectionCount)

		// If a new message is sent, check to see if the bot user is mentioned.
		case *slack.MessageEvent:
			if chk == 1 {
				if strings.Contains(ev.Msg.Text, user.ID) {
					parseCommand(ev.Msg.Text)
				}
			}

		// On bot startup, the bot goes from Offline to Online, and is likely
		// the first presence change for a bot that RTM will detect. Once
		// detected, grab the ID for the bot user
		case *slack.PresenceChangeEvent:
			log.Debugf("Presence Change: %v\n", ev)
			if chk == 0 {
				user, err = api.GetUserInfo(ev.User)
				if err != nil {
					log.Critical(err)
					os.Exit(1)
				}
				if user.Name == "oncall" && user.IsBot == true {
					chk = 1
				}

			}
		case *slack.RTMError:
			log.Debugf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Debugf("Invalid credentials")
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

// ChatUpdate takes a channel ID, a timestamp and message text
// and updated the message in the given Slack channel at the given
// timestamp with the given message text. Currently, it also updates the
// attachment specifically for the Set message output.
func ChatUpdate(channel string, ts string, text string) {
	t := fmt.Sprintf("Updated OnCall person to <@%s>", text)
	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Fallback:   "You would be able to select the oncall person here.",
		CallbackID: "oncall_dropdown",
		Footer:     t,
		FooterIcon: "https://emojipedia-us.s3.amazonaws.com/thumbs/120/apple/114/white-heavy-check-mark_2705.png",
	}
	params.Attachments = []slack.Attachment{attachment}
	text = "..."
	// Send an update to the given channel with pretext and the parameters
	channelID, timestamp, t, err := api.UpdateMessageWithParams(channel, ts, text, params)
	log.Debug(channelID, timestamp, t, err)
}

// parseCommand takes the message that mentions the bot user and identifies
// what the user is asking for.
func parseCommand(text string) {
	t := strings.Fields(text)
	switch t[1] {
	case "set":
		SetMessage()
	case "whois":
		WhoIsMessage()
	}

}

// verifyUser takes a User ID string and runs the Slack GetUserInfo request. If
// the user exists, the function returns true.
func VerifyUser(user string) bool {
	_, err := api.GetUserInfo(user)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	return true
}

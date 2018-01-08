package slack

import (
	"os"
	"strings"
	"time"

	"github.com/tylerconlee/slack"
)

var (
	// api is an instance of the tylerconlee/slack Client
	api *slack.Client
	// Triager holds the User ID of the current person set as "Triager"
	Triager string
	// SlabUser is the user ID of the slab Slack bot
	SlabUser string
)

// StartSlack initializes a connection with the given slack instance, gets
// team information, and starts a Go channel with the Real Time Messaging
// API watcher.
func StartSlack(v string) {
	log.Info("Starting connection to Slack")
	version = v
	uptime = time.Now()
	// start a connection to Slack using the Slack Bot token

	api = slack.New(c.Slack.APIKey)

	// retrieve the team info for the newly connected Slack team
	d, err := api.GetTeamInfo()
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	log.Info("Connected to Slack:", d.Domain)

	// Set the initial value of Triager
	Triager = "None"

	// Start monitoring Slack
	startRTM()

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
		log.Debug(msg.Data)
		log.Debug(msg.Type)
		log.Debug(rtm)
		switch ev := msg.Data.(type) {

		// When a user connects to Slack for the first time. Logged message
		// shows number of already connected users.
		case *slack.ConnectedEvent:
			log.Debug("Connection counter:", ev.ConnectionCount)
			if chk == 0 {
				user, err = api.GetUserInfo(ev.Info.User.ID)
				log.Debug(user.Name)
				if err != nil {
					log.Critical(err)
					os.Exit(1)
				}
				if user.Name == "slab" && user.IsBot == true {
					rtm.SendMessage(rtm.NewOutgoingMessage("Hello world. Slab connected.", c.Slack.ChannelID))

					log.Debug("Slab user identified")
					chk = 1
					Triager = SlabUser
				}
			}

		// If a new message is sent, check to see if the bot user is mentioned.
		case *slack.MessageEvent:
			log.Debug("Parsing Slack message")
			if chk == 1 {
				if strings.Contains(ev.Msg.Text, user.ID) && c.TriageEnabled {
					parseCommand(ev.Msg.Text)
				}
			}

		// On bot startup, the bot goes from Offline to Online, and is likely
		// the first presence change for a bot that RTM will detect. Once
		// detected, grab the ID for the bot user
		case *slack.PresenceChangeEvent:
			if chk == 0 {
				user, err = api.GetUserInfo(ev.User)
				log.Debug(user.Name)
				if err != nil {
					log.Critical(err)
					os.Exit(1)
				}
				if user.Name == "slab" && user.IsBot == true {
					chk = 1
					Triager = SlabUser
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

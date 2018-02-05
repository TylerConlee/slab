package slack

import (
	"strings"
	"time"

	l "github.com/tylerconlee/slab/log"
	"github.com/tylerconlee/slack"
)

var (
	// api is an instance of the tylerconlee/slack Client
	api *slack.Client
	// Triager holds the User ID of the current person set as "Triager"
	Triager string
	log     = l.Log
)

// StartSlack initializes a connection with the given slack instance, gets
// team information, and starts a Go channel with the Real Time Messaging
// API watcher.
func StartSlack(v string) {
	log.Info("Starting connection to Slack", map[string]interface{}{
		"module": "slack",
	})
	version = v
	uptime = time.Now()
	// start a connection to Slack using the Slack Bot token

	api = slack.New(c.Slack.APIKey)

	// retrieve the team info for the newly connected Slack team
	d, err := api.GetTeamInfo()
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "slack",
			"error":  err,
		})
	}
	log.Info("Connected to Slack", map[string]interface{}{
		"module": "slack",
		"team":   d.Domain,
	})

	// Set the initial value of Triager
	Triager = "None"

	// Start monitoring Slack
	startRTM()

}

// startRTM creates a separate Go channel which monitors the Slack instance.
// The RTM tracks each and every event within Slack and allows the bot to act
// accordingly.
func startRTM() {

	options := slack.RTMOptions{
		UseRTMStart: false,
	}
	rtm := api.NewRTMWithOptions(&options)
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

			if chk == 0 {
				user, err = api.GetUserInfo(ev.Info.User.ID)
				log.Debug("New user connected", map[string]interface{}{
					"module":   "slack",
					"count":    ev.ConnectionCount,
					"username": user.Name,
				})

				if err != nil {
					log.Fatal(map[string]interface{}{
						"module": "slack",
						"error":  err,
					})
				}
				if user.Name == "slab" && user.IsBot == true {
					log.Debug("Slab user identified", map[string]interface{}{
						"module": "slack",
						"id":     user.ID,
					})
					chk = 1
				}
			}

		// If a new message is sent, check to see if the bot user is mentioned.
		case *slack.MessageEvent:
			if chk == 1 {
				if strings.Contains(ev.Msg.Text, user.ID) && c.TriageEnabled {
					parseCommand(ev.Msg.Text, ev.User)
				}
			}

		// On bot startup, the bot goes from Offline to Online, and is likely
		// the first presence change for a bot that RTM will detect. Once
		// detected, grab the ID for the bot user
		case *slack.PresenceChangeEvent:
		case *slack.RTMError:
			log.Error("RTMError Encountered", map[string]interface{}{
				"module": "slack",
				"error":  ev.Error(),
			})

		case *slack.ConnectionErrorEvent:
			log.Error("Connection Error Encountered", map[string]interface{}{
				"module": "slack",
				"error":  ev.Error(),
			})

		case *slack.InvalidAuthEvent:
			log.Error("Authentication Error Encountered. Invalid Credentials", map[string]interface{}{
				"module": "slack",
			})
			return

		default:

			// Ignore other events..
			// fmt.Printf("Unexpected: %v\n", msg.Data)
		}
	}
}

package plugins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	l "github.com/tylerconlee/slab/log"
	"github.com/tylerconlee/slack"
)

var log = l.Log

// TwilioPhone is the "to" phone number that's set through Slack (@slab twilio
//	set)
var TwilioPhone string

// TwilioFrom is the "from" phone number that's set through Slack (@slab twilio
// configure)
var TwilioFrom string

// EnableTwilio changes the Enabled Twilio option to true.
func (p *Plugins) EnableTwilio() {
	p.Twilio.Enabled = true
}

// DisableTwilio changes the Enabled Twilio option to false.
func (p *Plugins) DisableTwilio() {
	p.Twilio.Enabled = false
}

// TwilioSet changes the TwilioPhone to the value of the number passed to
// it.
func TwilioSet(n string) {
	TwilioPhone = n
	log.Info("Phone number set.", map[string]interface{}{
		"module": "plugin",
		"plugin": "Twilio",
		"phone":  TwilioPhone,
	})
}

// TwilioUnset sets the TwilioPhone to `none`
func TwilioUnset() {
	TwilioPhone = ""
}

// TwilioStatus returns the current setting
func (p *Plugins) TwilioStatus() (attachment slack.Attachment) {
	s := ":x:"
	if p.Twilio.Enabled {
		s = ":white_check_mark:"
	}
	attachment = slack.Attachment{
		Title: "Twilio",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Enabled",
				Value: s,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Current Phone Number",
				Value: TwilioPhone,
				Short: true,
			},
		},
	}
	return attachment
}

// SendTwilio sends a message to the phone number currently set
// as TwilioPhone using the connection data found in the config
func (p *Plugins) SendTwilio(message string) {

	// Prep text message
	msgData := url.Values{}
	msgData.Set("To", TwilioPhone)
	msgData.Set("From", p.Twilio.Phone)
	msgData.Set("Body", message)
	msgDataReader := *strings.NewReader(msgData.Encode())

	// Connect to Twilio
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + p.Twilio.AccountID + "/Messages.json"
	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(p.Twilio.AccountID, p.Twilio.Auth)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, _ := client.Do(req)

	// Parse response
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
	}
}

// TwilioConfigure sets the "from" phone number to allow for international
// numbers to be set properly
func TwilioConfigure(n string) {
	TwilioFrom = n
	log.Info("From phone number set.", map[string]interface{}{
		"module": "plugin",
		"plugin": "Twilio",
		"from":   TwilioFrom,
	})
}

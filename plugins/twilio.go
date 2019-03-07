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

var (
	log = l.Log

	// TwilioPhone is the "to" phone number that's set through Slack (@slab twilio
	//	set)
	TwilioPhone string

	// TwilioFrom is the "from" phone number that's set through Slack (@slab twilio
	// configure)
	TwilioFrom string

	// TwilioEnabled holds whether the Twilio plugin is enabled or disabled.
	TwilioEnabled bool
)

// EnableTwilio changes the Enabled Twilio option to true.
func (p *Plugins) EnableTwilio() (attachment slack.Attachment) {
	TwilioEnabled = true
	return p.checkStatus()
}

// DisableTwilio changes the Enabled Twilio option to false.
func (p *Plugins) DisableTwilio() (attachment slack.Attachment) {
	TwilioEnabled = false
	return p.checkStatus()
}

// TwilioSet changes the TwilioPhone to the value of the number passed to
// it.
func TwilioSet(n string) (attachment slack.Attachment) {
	TwilioPhone = n
	log.Info("Phone number set.", map[string]interface{}{
		"module": "plugin",
		"plugin": "Twilio",
		"phone":  TwilioPhone,
	})
	attachment = slack.Attachment{
		Title: "Twilio 'To' Phone Number Set",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Phone Number",
				Value: TwilioPhone,
				Short: true,
			},
		},
	}
	return attachment
}

// TwilioUnset sets the TwilioPhone to `none`
func TwilioUnset() (attachment slack.Attachment) {
	TwilioPhone = ""
	attachment = slack.Attachment{
		Title: "Twilio 'To' Phone Number Unset",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Phone Number",
				Value: TwilioPhone,
				Short: true,
			},
		},
	}
	return attachment
}

// TwilioStatus returns the current setting
func (p *Plugins) TwilioStatus() (attachment slack.Attachment) {
	return p.checkStatus()
}

// SendTwilio sends a message to the phone number currently set
// as TwilioPhone using the connection data found in the config
func (p *Plugins) SendTwilio(message string) {

	// Prep text message
	msgData := url.Values{}
	msgData.Set("To", TwilioPhone)
	msgData.Set("From", TwilioFrom)
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
func TwilioConfigure(n string) (attachment slack.Attachment) {
	TwilioFrom = n
	log.Info("From phone number set.", map[string]interface{}{
		"module": "plugin",
		"plugin": "Twilio",
		"from":   TwilioFrom,
	})
	attachment = slack.Attachment{
		Title: "Twilio 'From' Phone Number Set",
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Current Phone Number",
				Value: TwilioFrom,
				Short: true,
			},
		},
	}
	return attachment
}

func (p *Plugins) checkStatus() (attachment slack.Attachment) {
	s := ":x:"
	if p.Twilio.Enabled {
		s = ":white_check_mark:"
	}
	attachment = slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Enabled",
				Value: s,
			},
			slack.AttachmentField{
				Title: "Current 'From' Phone Number",
				Value: TwilioFrom,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Current 'To' Phone Number",
				Value: TwilioPhone,
				Short: true,
			},
		},
	}
	return attachment
}

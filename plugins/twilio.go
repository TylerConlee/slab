package plugins

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// TwilioPhone is the "to" phone number that's set through Slack (@slab twilio
//	set)
var TwilioPhone string

// EnableTwilio changes the Enabled Twilio option to the setting passed to it.
func EnableTwilio(e bool) {
	p.Twilio.Enabled = e
}

// TwilioSet changes the TwilioPhone to the value of the number passed to
// it.
func TwilioSet(n string) {
	TwilioPhone = n
}

// TwilioUnset sets the TwilioPhone to `none`
func TwilioUnset() {

}

// TwilioStatus returns the current setting
func TwilioStatus() {

}

// SendTwilio sends a message to the phone number currently set
// as TwilioPhone using the connection data found in the config
func SendTwilio(message string) {

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

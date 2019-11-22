package plugins

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nlopes/slack"
)

var (

	// PagerDutyEnabled holds whether the PagerDuty plugin is enabled or disabled.
	PagerDutyEnabled bool
)

func init() {
	Commands["pagerduty"] = pagerdutyCommands
	Send["pagerduty"] = (*Plugins).SendPagerDuty
}

func pagerdutyCommands(t []string) (attachments []slack.Attachment, message string) {
	switch t[1] {
	case "pagerduty":
		if len(t) > 1 {
			p := LoadPlugins()
			switch t[2] {
			case "set":
				if len(t) > 3 {
					s := PagerDutySet(t[3])
					attachments = []slack.Attachment{s}
					message = "Plugin message"
				}
			case "unset":
				s := PagerDutyUnset()
				attachments = []slack.Attachment{s}
				message = "Plugin message"
			case "configure":
				if len(t) > 3 {
					s := PagerDutyConfigure(t[3])
					attachments = []slack.Attachment{s}
					message = "Plugin message"
				}
			case "status":
				s := p.PagerDutyStatus()
				attachments = []slack.Attachment{s}
				message = "Plugin status"
			case "enable":
				p.EnablePagerDuty()
				a := slack.Attachment{
					Title: "PagerDuty Plugin",
					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Title: "Enabled",
							Value: ":white_check_mark:",
						},
					},
				}
				attachments = []slack.Attachment{a}
				message = "Plugin PagerDuty has been updated"

			case "disable":
				p.DisablePagerDuty()
				a := slack.Attachment{
					Title: "PagerDuty Plugin",
					Fields: []slack.AttachmentField{
						slack.AttachmentField{
							Title: "Enabled",
							Value: ":x:",
						},
					},
				}
				attachments = []slack.Attachment{a}
				message = "Plugin PagerDuty has been updated"
			}
		}
	}
	return attachments, message
}

// PagerDuty contains the connection details for the PagerDuty API:
type PagerDuty struct {
	AccountID string
	Auth      string
	Enabled   bool
}

// EnablePagerDuty changes the Enabled PagerDuty option to true.
func (p *Plugins) EnablePagerDuty() (attachment slack.Attachment) {
	PagerDutyEnabled = true
	return p.checkPagerdutyStatus()
}

// DisablePagerDuty changes the Enabled PagerDuty option to false.
func (p *Plugins) DisablePagerDuty() (attachment slack.Attachment) {
	PagerDutyEnabled = false
	return p.checkPagerdutyStatus()
}

// PagerDutyStatus returns the current setting
func (p *Plugins) PagerDutyStatus() (attachment slack.Attachment) {
	return p.checkPagerdutyStatus()
}

// SendPagerDuty sends a message to PagerDuty and returns a list of incidents
func (p *Plugins) SendPagerDuty(message string) {

	// Connect to PagerDuty
	urlStr := "https://api.pagerduty.com/incidents"
	req, _ := http.NewRequest("GET", urlStr, nil)
	req.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
	req.Header.Set("Authorization", "Token token=y_NbAkKc66ryYTWUXYEu")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Info("Error receiving response from PagerDuty", map[string]interface{}{
			"module": "plugin",
			"plugin": "pagerduty",
			"error":  err,
		})
	}

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

func (p *Plugins) checkPagerdutyStatus() (attachment slack.Attachment) {
	s := ":x:"
	if p.PagerDuty.Enabled {
		s = ":white_check_mark:"
	}
	attachment = slack.Attachment{
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Enabled",
				Value: s,
			},
		},
	}
	return attachment
}

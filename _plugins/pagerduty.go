package plugins

import (
	"net/http"

	"github.com/PagerDuty/go-pagerduty"
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
			case "status":
				s := p.PagerDutyStatus()
				attachments = []slack.Attachment{s}
				message = "Plugin status"
			case "enable":
				a := p.EnablePagerDuty()
				attachments = []slack.Attachment{a}
				message = "Plugin PagerDuty has been updated"

			case "demo":
				p.DemoPagerDuty()

			case "disable":
				a := p.DisablePagerDuty()
				attachments = []slack.Attachment{a}
				message = "Plugin PagerDuty has been updated"
			}
		}
	}
	return attachments, message
}

// PagerDuty contains the connection details for the PagerDuty API:
type PagerDuty struct {
	APIKey    string
	ServiceID string
	Enabled   bool
}

// EnablePagerDuty changes the Enabled PagerDuty option to true.
func (p *Plugins) EnablePagerDuty() (attachment slack.Attachment) {
	p.PagerDuty.Enabled = true
	PagerDutyEnabled = true
	return p.checkPagerdutyStatus()
}

// DisablePagerDuty changes the Enabled PagerDuty option to false.
func (p *Plugins) DisablePagerDuty() (attachment slack.Attachment) {
	p.PagerDuty.Enabled = false
	PagerDutyEnabled = false
	return p.checkPagerdutyStatus()
}

// PagerDutyStatus returns the current setting
func (p *Plugins) PagerDutyStatus() (attachment slack.Attachment) {
	return p.checkPagerdutyStatus()
}

// SendPagerDuty sends a message to PagerDuty and returns a list of incidents
func (p *Plugins) SendPagerDuty(message string) {
	if PagerDutyEnabled {

		event := pagerduty.Event{
			Type:        "trigger",
			ServiceKey:  p.PagerDuty.ServiceID,
			Description: message,
		}
		log.Debug("Pagerduty request created", map[string]interface{}{
			"module":  "plugins",
			"plugin":  "pagerduty",
			"request": event,
		})
		resp, err := pagerduty.CreateEvent(event)

		if err != nil {
			log.Error("Error sending PagerDuty Event", map[string]interface{}{
				"module": "plugins",
				"plugin": "pagerduty",
				"error":  err,
			})
		}
		log.Info("Response from PagerDuty", map[string]interface{}{
			"module":   "plugins",
			"plugin":   "pagerduty",
			"response": resp,
		})
	}
}

// DemoPagerDuty grabs a demo list of incidents from PagerDuty
func (p *Plugins) DemoPagerDuty() {

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

	log.Info("Response received from PagerDuty", map[string]interface{}{
		"module":   "plugin",
		"plugin":   "pagerduty",
		"response": resp.Body,
	})
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

package plugins

import "github.com/tylerconlee/slab/config"

// plugins contains a list of all available plugins
type plugins struct {
	Twilio    Twilio
	PagerDuty PagerDuty
}

// Twilio contains the connection details for the Twilio API:
// https://www.twilio.com/docs/api
type Twilio struct {
	AccountID string
	Auth      string
	Phone     string
	Enabled   bool
}

// PagerDuty contains the connection details for the PagerDuty API:
// https://v2.developer.pagerduty.com/docs/rest-api
type PagerDuty struct {
	Email   string
	Key     string
	Enabled bool
}

var p plugins

// LoadPlugins is sent a map of the plugin configuration. It parses the
// configuration and determines which plugins are available.
func LoadPlugins(c config.Config) {
	p = plugins{
		Twilio{
			c.Twilio.AccountID,
			c.Twilio.Auth,
			c.Twilio.Phone,
			true,
		},
		PagerDuty{
			c.PagerDuty.Email,
			c.PagerDuty.Key,
			true,
		},
	}
}

// SendDispatcher receives the message from the process loop and checks which
// plugins are enabled and sends the appropriate notifications through them.
func SendDispatcher(message string) {
	log.Info("Plugins reached.", map[string]interface{}{
		"module": "plugin",
		"plugin": p,
	})
	if (p.Twilio.Enabled) && (p.Twilio.Phone != "") {
		log.Info("Plugin loaded. Sending Twilio message.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
		SendTwilio(message)
	}
}

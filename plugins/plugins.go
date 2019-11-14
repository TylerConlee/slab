package plugins

import "github.com/tylerconlee/slab/config"

// Plugins contains a list of all available plugins
type Plugins struct {
	Twilio    Twilio
	PagerDuty PagerDuty
}



// PagerDuty contains the connection details for the PagerDuty API:
// https://v2.developer.pagerduty.com/docs/rest-api
type PagerDuty struct {
	Email   string
	Key     string
	Enabled bool
}

// LoadPlugins is sent a map of the plugin configuration. It parses the
// configuration and determines which plugins are available.
func LoadPlugins(c config.Config) (p Plugins) {
	return Plugins{
		Twilio{
			c.Plugins.Twilio.AccountID,
			c.Plugins.Twilio.Auth,
			true,
		},
		PagerDuty{
			c.Plugins.PagerDuty.Email,
			c.Plugins.PagerDuty.Key,
			true,
		},
	}
}

// SendDispatcher receives the message from the process loop and checks which
// plugins are enabled and sends the appropriate notifications through them.
func (p *Plugins) SendDispatcher(message string) {
	log.Info("Plugins reached.", map[string]interface{}{
		"module": "plugin",
		"plugin": p,
	})
	if TwilioPhone == "" {
		log.Info("To phone number for Twilio not set.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
	}
	if TwilioFrom == "" {
		log.Info("From phone number for Twilio not set.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
	}
	if (TwilioEnabled) && (TwilioPhone != "") {
		log.Info("Plugin loaded. Sending Twilio message.", map[string]interface{}{
			"module": "plugin",
			"plugin": "twilio",
		})
		p.SendTwilio(message)
	}
}

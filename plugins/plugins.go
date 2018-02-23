package plugins

import "github.com/tylerconlee/slab/config"

type Twilio struct {
	Email string
	Key   string
}

type PagerDuty struct {
	Email string
	Key   string
}

type plugins struct {
	Twilio    Twilio
	PagerDuty PagerDuty
}

var p plugins

// LoadPlugins is sent a map of the plugin configuration. It parses the
// configuration and determines which plugins are available.
func LoadPlugins(c config.Config) {
	p = plugins{
		Twilio{
			c.Twilio.Email,
			c.Twilio.Key,
		},
		PagerDuty{
			c.PagerDuty.Email,
			c.PagerDuty.Key,
		},
	}
}

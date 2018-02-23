package config

type Plugins struct {
	Twilio
	PagerDuty
}

type Twilio struct {
	Email string
	Key   string
}

type PagerDuty struct {
	Email string
	Key   string
}

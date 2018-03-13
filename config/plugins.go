package config

type Plugins struct {
	Twilio    Twilio
	PagerDuty PagerDuty
}

type Twilio struct {
	AccountID string
	Auth      string
	Phone     string
}

type PagerDuty struct {
	Email string
	Key   string
}

package config

type Plugins struct {
	Twilio
	PagerDuty
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

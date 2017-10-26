package main

import "github.com/BurntSushi/toml"

// Config maps the values of the configuration file to a struct usable by the
// rest of the app
type Config struct {
	Zendesk  Zendesk
	SlackAPI string
}

type Zendesk struct {
	User   string
	APIKey string
	URL    string
}

func loadConfig(path string) (config Config) {
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Critical(err)
		return
	}
	return config
}

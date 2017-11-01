package config

import (
	"os"
	"time"

	logging "github.com/op/go-logging"

	"github.com/BurntSushi/toml"
)

var log = logging.MustGetLogger("config")

// Config maps the values of the configuration file to a struct usable by the
// rest of the app
type Config struct {
	Zendesk  Zendesk
	SlackAPI string
	SLA      SLA
}

type Zendesk struct {
	User   string
	APIKey string
	URL    string
}

type SLA struct {
	LevelOne   Level
	LevelTwo   Level
	LevelThree Level
	LevelFour  Level
}

type Level struct {
	Low    duration
	Normal duration
	High   duration
	Urgent duration
}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func LoadConfig() (config Config) {
	if len(os.Args) > 1 {
		if _, err := toml.DecodeFile(os.Args[1], &config); err != nil {
			log.Critical(err)
			return
		}
		log.Info("Configuration file", os.Args[1], "loaded successfully.")
		return config
	} else {
		log.Critical("Error. Configuration file must be specified when launching SLAB")
		os.Exit(1)
		return
	}
}

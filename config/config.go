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
	Zendesk    Zendesk
	SlackAPI   string
	SLA        SLA
	UpdateFreq Duration
}

// Zendesk contains configuration values specific to the Zendesk interactions
type Zendesk struct {
	User   string
	APIKey string
	URL    string
}

// SLA supports up to 4 levels of SLA in the configuration
type SLA struct {
	LevelOne   Level
	LevelTwo   Level
	LevelThree Level
	LevelFour  Level
}

// Level reflects the 4 priority levels Zendesk uses for SLA.
type Level struct {
	Low    Duration
	Normal Duration
	High   Duration
	Urgent Duration
}

// Duration allows for configurations to contain "3h", "8m", etc.
type Duration struct {
	time.Duration
}

// UnmarshalText takes the Duration and returns a time.Duration in place of the
// string.
func (d *Duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

// LoadConfig grabs the command line argument for where the configuration file
// is located and loads that into memory.
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

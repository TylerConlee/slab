package config

import (
	"bytes"
	"os"
	"time"

	l "github.com/tylerconlee/slab/log"

	"github.com/BurntSushi/toml"
)

var log = l.Log

// Config maps the values of the configuration file to a struct usable by the
// rest of the app
type Config struct {
	Zendesk       Zendesk
	Slack         Slack
	LogLevel      string
	SLA           SLA
	UpdateFreq    Duration
	TriageEnabled bool
	Metadata      Metadata
	Port          int
	Plugins       Plugins
}

// Metadata holds configuration related to the server metadata used in status calls
type Metadata struct {
	Server string
}

// Slack API key and Channel ID tell SLAB where to post notifications
type Slack struct {
	APIKey    string
	ChannelID string
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
	Tag    string
	Low    Duration
	Normal Duration
	High   Duration
	Urgent Duration
	Notify bool
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
	if _, err := os.Stat("config.toml"); err == nil {
		if _, err := toml.DecodeFile("config.toml", &config); err != nil {
			log.Error("Configuration file not found.", map[string]interface{}{
				"module": "config",
				"error":  err,
			})
			config = defaultConfig()
			return
		}
		log.Info("Configuration loaded successfully", map[string]interface{}{
			"module": "config",
			"file":   "config.toml",
		})
		return
	}
	config = defaultConfig()
	return
}

// SaveConfig takes a config and saves it to the local file, config.toml.
func SaveConfig(config Config) bool {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		log.Error("Error creating new buffer for config", map[string]interface{}{
			"module": "config",
			"config": config,
			"error":  err,
		})
		return false
	}

	f, err := os.Create("config.toml")
	if nil != err {
		log.Error("error saving file", map[string]interface{}{
			"module": "config",
			"error":  err,
		})
	}
	defer f.Close()
	n, err := f.WriteString(buf.String())
	if nil != err {
		log.Error("error saving file", map[string]interface{}{
			"module": "config",
		})
	}
	f.Sync()
	log.Debug("Saved configuration file", map[string]interface{}{
		"module": "config",
		"output": n,
	})
	return true

}

func defaultConfig() (config Config) {

	config = Config{
		Zendesk: Zendesk{
			APIKey: "",
			User:   "",
			URL:    "",
		},
		SLA: SLA{
			LevelOne: Level{
				Tag: "platinum",
			},
			LevelTwo: Level{
				Tag: "gold",
			},
			LevelThree: Level{
				Tag: "silver",
			},
			LevelFour: Level{
				Tag: "bronze",
			},
		},
		Metadata:      Metadata{},
		TriageEnabled: true,
	}
	return
}

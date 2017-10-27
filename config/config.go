package config

import (
	"os"

	logging "github.com/op/go-logging"

	"github.com/BurntSushi/toml"
)

var log = logging.MustGetLogger("config")

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

func LoadConfig() (config Config) {
	if len(os.Args) > 1 {
		if _, err := toml.DecodeFile(os.Args[1], &config); err != nil {
			log.Critical(err)
			return
		}
		return config
	} else {
		log.Critical("Error. Configuration file must be specified when launching SLAB")
		os.Exit(1)
		return
	}
}

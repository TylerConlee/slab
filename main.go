package main

import (
	"os"

	logging "github.com/op/go-logging"
	Zen "github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("slab")

const VERSION = "0.0.1"

func main() {
	// Get the path to the configuration file and use it to load the config
	path := os.Args[1]
	config := loadConfig(path)

	// Start up the logging system
	initLog()
	log.Notice("SLABot by Tyler Conlee")
	log.Noticef("Version: %s", VERSION)

	// Get all tickets from Zendesk using the configuration values
	Zen.GetAllTickets(config.Zendesk.User, config.Zendesk.APIKey, config.Zendesk.URL)
}

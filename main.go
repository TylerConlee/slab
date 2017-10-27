package main

import (
	logging "github.com/op/go-logging"
	c "github.com/tylerconlee/slab/config"
	Zen "github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("slab")

const VERSION = "0.0.1"

func main() {
	// Start up the logging system
	initLog()
	log.Notice("SLABot by Tyler Conlee")
	log.Noticef("Version: %s", VERSION)

	config := c.LoadConfig()

	// Get all tickets from Zendesk using the configuration values
	Zen.GetAllTickets(config.Zendesk.User, config.Zendesk.APIKey, config.Zendesk.URL)
}

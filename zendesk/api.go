package zendesk

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("zendesk")

func verifyAPIKey() (key string) {
	key = os.Getenv("SLAB_ZENDESK_API")
	if "" == key {
		log.Critical("No key provided for Zendesk API")
		os.Exit(1)
	}
	return key
}

// GetAllTickets grabs the latest tickets from Zendesk and returns the JSON
func GetAllTickets() {
	key := verifyAPIKey()
	log.Debugf("Zendesk API Key Found: %s", key)
}

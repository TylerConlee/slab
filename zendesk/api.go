package zendesk

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("zemdesk")

func GetAllTickets() {
	key := os.Getenv("SLAB_ZENDESK_API")
	if "" == key {
		log.Critical("No key provided for Zendesk API")
		os.Exit(1)
	}
}

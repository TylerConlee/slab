package zendesk

import (
	"net/http"
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("zendesk")

// GetAllTickets grabs the latest tickets from Zendesk and returns the JSON
func GetAllTickets(user string, key string, url string) {
	log.Info("Starting request to Zendesk for tickets")

	log.Debugf("Zendesk API User Found: %s", user)
	log.Debugf("Zendesk API Key Found: %s", key)

	zenURL := url + "/api/v2/tickets.json?include=slas"

	req, err := http.NewRequest("GET", zenURL, nil)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	req.SetBasicAuth(user, key)

	log.Debug("Request to Zendesk made for all tickets")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Output raw response
	log.Debug(resp)
}

package zendesk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	l "github.com/tylerconlee/slab/log"
)

var log = l.Log

// GetAllTickets grabs the latest tickets from Zendesk and returns the JSON
// Zendesk Endpoint: /incremental/tickets.json?include=slas
func GetAllTickets() (tickets ZenOutput) {
	log.Info("Starting request to Zendesk for tickets", map[string]interface{}{
		"module": "zendesk",
	})

	t := time.Now().AddDate(0, 0, -5).Unix()
	zen := c.Zendesk.URL + "/api/v2/incremental/tickets.json?include=slas&start_time=" + strconv.FormatInt(t, 10)
	resp := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)
	tickets = parseJSON(resp)
	log.Info("Request Complete. Parsing Ticket Data", map[string]interface{}{
		"module":      "zendesk",
		"num_tickets": len(tickets.Tickets),
	})
	return tickets
}

// GetTicketRequester takes the requester ID from the tickets grabbed in
// GetAllTickets and sends a request to Zendesk for the user info
// Zendesk Endpoint /users/{USER-ID}.json
func GetTicketRequester(user int) (output User) {
	log.Info("Starting request to Zendesk for user info", map[string]interface{}{
		"module": "zendesk",
		"user":   user,
	})

	zen := c.Zendesk.URL + "/api/v2/users/" + strconv.Itoa(user) + ".json"
	data := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)
	log.Info("Request Complete. Parsing Ticket Data", map[string]interface{}{
		"module": "zendesk",
		"user":   user,
	})
	resp := json.RawMessage(data)
	users := Users{}
	err := json.Unmarshal(resp, &users)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return users.User
}

// GetOrganization takes the org ID from the tickets grabbed in
// GetAllTickets and sends a request to Zendesk for the Org information
// Zendesk Endpoint /users/{USER-ID}/organizations.json
func GetOrganization(user int) (org Orgs) {
	log.Info("Starting request to Zendesk for organization info", map[string]interface{}{
		"module": "zendesk",
		"user":   user,
	})
	zen := c.Zendesk.URL + "/api/v2/users/" + strconv.Itoa(user) + "/organizations.json"
	data := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)
	log.Info("Request Complete. Parsing Organization Data", map[string]interface{}{
		"module": "zendesk",
		"user":   user,
	})
	resp := json.RawMessage(data)
	orgs := Organizations{}
	err := json.Unmarshal(resp, &orgs)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return orgs.Orgs

}

// GetRequestedTickets takes a user ID and sends a request to Zendesk to grab
// the IDs of tickets requested by that user
// Zendesk Endpoint /users/{USER-ID}/tickets/requested.json
func GetRequestedTickets(user int) (output ZenOutput) {
	log.Info("Starting request to Zendesk for requested ticket info", map[string]interface{}{
		"module": "zendesk",
		"user":   user,
	})
	zen := c.Zendesk.URL + "/api/v2/users/" + strconv.Itoa(user) + "/tickets/requested.json"
	data := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)
	log.Info("Request Complete. Parsing Organization Data", map[string]interface{}{
		"module": "zendesk",
		"user":   user,
	})
	resp := json.RawMessage(data)
	err := json.Unmarshal(resp, &output)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return
}

// makeRequests takes the Zendesk auth information and sends the curl request
// to Zendesk and returns a JSON blob
func makeRequest(user string, key string, url string) (responseData []byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	req.SetBasicAuth(user, key)

	// create custom http.Client to manually set timeout and disable keepalive
	// in an attempt to avoid EOF errors
	var netClient = &http.Client{
		Timeout: time.Second * 60,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	resp, err := netClient.Do(req)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	defer resp.Body.Close()
	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return responseData
}

// parseJSON takes the JSON from makeRequest and unmarshals it into the
// ZenOutput struct, allowing the data to be accessed
func parseJSON(data []byte) (output ZenOutput) {
	// Read response from HTTP client
	bytes := json.RawMessage(data)
	err := json.Unmarshal(bytes, &output)
	if err != nil {
		log.Fatal(map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return output
}

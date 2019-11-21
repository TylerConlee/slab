package zendesk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	l "github.com/tylerconlee/slab/log"
)

var log = l.Log

// NumTickets is the number of tickets processed on the last loop
var NumTickets int

// LastProcessed is a timestamp of when the last loop was ran
var LastProcessed time.Time

// GetAllTickets grabs the latest tickets from Zendesk and returns the JSON
// Zendesk Endpoint: /incremental/tickets.json?include=slas
// TODO: Handle paging from the Incremental API
func GetAllTickets() (tickets ZenOutput) {
	log.Info("Requesting all tickets from Zendesk for SLA", map[string]interface{}{
		"module": "zendesk",
	})

	t := time.Now().AddDate(0, 0, -3).Unix()
	zen := c.Zendesk.URL + "/api/v2/incremental/tickets.json?include=slas,metric_events&start_time=" + strconv.FormatInt(t, 10)
	resp := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)

	tickets = parseTicketJSON(resp)
	nextPage := ""
	if tickets.NextPage != nil {
		if tickets.NextPage.(string) != zen && strings.Contains(tickets.NextPage.(string), "page") {
			nextPage = tickets.NextPage.(string)
		}

		// While NextPage is blank
		for nextPage != "" && strings.Contains(nextPage, "page") {

			// Run GetNextPage and store the results in output
			output := getNextPage(nextPage)

			// Append the results to the main tickets list
			tickets.Tickets = append(tickets.Tickets, output.Tickets...)

			// Reset next page to blank to avoid unnecessary calls
			nextPage = ""

			// Check to see if NextPage from output is blank
			if output.NextPage != nil {
				l.Log.Info("Ticket output contains a non-nil next page", map[string]interface{}{
					"module": "zendesk",
					"ticket": tickets.NextPage,
					"output": output.NextPage,
				})

				// If the output's next page is different from the ticket's next page
				if output.NextPage.(string) != tickets.NextPage.(string) {

					// Set the next page as the output's next page to run the loop again
					nextPage = output.NextPage.(string)
				}
			}
		}
	} else {
		l.Log.Info("Ticket output contains a nil next page", map[string]interface{}{
			"module": "zendesk",
			"ticket": tickets.NextPage,
		})
	}
	NumTickets = len(tickets.Tickets)
	LastProcessed = time.Now()
	log.Info("Request Complete. Parsing Ticket Data", map[string]interface{}{
		"module":      "zendesk",
		"num_tickets": NumTickets,
	})
	return tickets
}

func getNextPage(nextURL string) (output ZenOutput) {
	nextResp := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, nextURL)
	l.Log.Info("Requesting next page of Zendesk tickets", map[string]interface{}{
		"module":   "zendesk",
		"nextpage": nextURL,
	})
	return parseTicketJSON(nextResp)
}

// GetTicketEvents grabs the latest ticket events from Zendesk and returns the
// JSON
// Zendesk Endpoint: /api/v2/incremental/ticket_events.json
func GetTicketEvents() (tickets EventOutput) {
	log.Info("Requesting latest ticket events for updates", map[string]interface{}{
		"module": "zendesk",
	})
	hour := time.Duration(1 * time.Hour)
	t := time.Now().Add(-hour)
	zen := c.Zendesk.URL + "/api/v2/incremental/ticket_events.json?start_time=" + strconv.FormatInt(t.Unix(), 10)
	resp := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)

	tickets = parseEventJSON(resp)
	log.Info("Request Complete. Parsing Ticket Data", map[string]interface{}{
		"module":      "zendesk",
		"num_tickets": len(tickets.Event),
	})
	return tickets
}

// GetTicket gets the details about an individual ticket from Zendesk
func GetTicket(id int) (ticket Ticket) {
	log.Info("Requesting data on individual ticket", map[string]interface{}{
		"module": "zendesk",
		"ticket": id,
	})
	zen := c.Zendesk.URL + "/api/v2/tickets/" + strconv.Itoa(id) + ".json"
	resp := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)
	bytes := json.RawMessage(resp)
	tick := TicketGroup{}
	err := json.Unmarshal(bytes, &tick)
	if err != nil {
		log.Error("Error parsing Zendesk JSON", map[string]interface{}{
			"module": "zendesk",
			"error":  err,
			"resp":   bytes,
		})
	}
	log.Info("Request Complete. Parsing Ticket Data", map[string]interface{}{
		"module": "zendesk",
	})
	return tick.Ticket
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
		log.Error("Error with Zendesk requestor parsing", map[string]interface{}{
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
		log.Error(
			"Error with Zendesk org parsing.", map[string]interface{}{
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
		log.Error("Error with Zendesk ticket parsing", map[string]interface{}{
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
		log.Error("Error forming HTTP request", map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	req.SetBasicAuth(user, key)

	// create custom http.Client to manually set timeout and disable keepalive
	// in an attempt to avoid EOF errors
	var netClient = &http.Client{
		Timeout: time.Second * 240,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}

	resp, err := netClient.Do(req)

	if err != nil {
		log.Error("Error reaching Zendesk", map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
		return
	}
	log.Info("Zendesk request received", map[string]interface{}{
		"module": "zendesk",
		"status": resp.StatusCode,
	})
	defer resp.Body.Close()
	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading Zendesk response", map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
		return
	}
	return responseData
}

// parseJSON takes the JSON from makeRequest and unmarshals it into the
// ZenOutput struct, allowing the data to be accessed
func parseTicketJSON(data []byte) (output ZenOutput) {
	// Read response from HTTP client
	bytes := json.RawMessage(data)
	err := json.Unmarshal(bytes, &output)
	if err != nil {
		log.Error("Error parsing Zendesk JSON", map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return output
}

// parseJSON takes the JSON from makeRequest and unmarshals it into the
// ZenOutput struct, allowing the data to be accessed
func parseEventJSON(data []byte) (output EventOutput) {
	// Read response from HTTP client
	bytes := json.RawMessage(data)
	err := json.Unmarshal(bytes, &output)
	if err != nil {
		log.Error("Error parsing Zendesk JSON", map[string]interface{}{
			"module": "zendesk",
			"error":  err,
		})
	}
	return output
}

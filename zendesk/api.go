package zendesk

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("zendesk")

// ZenOutput is the top level JSON-based struct that whatever is
// returned by Zendesk goes into
type ZenOutput struct {
	Tickets      `json:"tickets"`
	NextPage     interface{} `json:"next_page"`
	PreviousPage interface{} `json:"previous_page"`
	Count        int         `json:"count"`
}

// Tickets is a subset of ZenOutput that contains the details of the tickets
// outputted from the request to Zendesk
// TODO: use the OrgID to make a request for Org name using a different API call
type Tickets []struct {
	URL        string      `json:"url"`
	ID         int         `json:"id"`
	ExternalID interface{} `json:"external_id"`
	Via        struct {
		Channel string `json:"channel"`
		Source  struct {
			From struct {
			} `json:"from"`
			To struct {
			} `json:"to"`
			Rel string `json:"rel"`
		} `json:"source"`
	} `json:"via"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Type            interface{}   `json:"type"`
	Subject         string        `json:"subject"`
	RawSubject      string        `json:"raw_subject"`
	Description     string        `json:"description"`
	Priority        interface{}   `json:"priority"`
	Status          string        `json:"status"`
	Recipient       interface{}   `json:"recipient"`
	RequesterID     int64         `json:"requester_id"`
	SubmitterID     int64         `json:"submitter_id"`
	AssigneeID      interface{}   `json:"assignee_id"`
	OrganizationID  int64         `json:"organization_id"`
	GroupID         int           `json:"group_id"`
	CollaboratorIds []interface{} `json:"collaborator_ids"`
	FollowerIds     []interface{} `json:"follower_ids"`
	ForumTopicID    interface{}   `json:"forum_topic_id"`
	ProblemID       interface{}   `json:"problem_id"`
	HasIncidents    bool          `json:"has_incidents"`
	IsPublic        bool          `json:"is_public"`
	DueAt           interface{}   `json:"due_at"`
	Tags            []string      `json:"tags"`
	CustomFields    []struct {
		ID    int         `json:"id"`
		Value interface{} `json:"value"`
	} `json:"custom_fields"`
	SatisfactionRating struct {
		Score string `json:"score"`
	} `json:"satisfaction_rating"`
	SharingAgreementIds []interface{} `json:"sharing_agreement_ids"`
	Fields              []struct {
		ID    int         `json:"id"`
		Value interface{} `json:"value"`
	} `json:"fields"`
	TicketFormID            int         `json:"ticket_form_id"`
	BrandID                 int         `json:"brand_id"`
	SatisfactionProbability interface{} `json:"satisfaction_probability"`
	Slas                    struct {
		PolicyMetrics []interface{} `json:"policy_metrics"`
	} `json:"slas"`
	AllowChannelback bool `json:"allow_channelback"`
}

// GetAllTickets grabs the latest tickets from Zendesk and returns the JSON
func GetAllTickets(user string, key string, url string) (tickets ZenOutput) {
	log.Info("Starting request to Zendesk for tickets")

	log.Debugf("Zendesk API User Found: %s", user)
	log.Debugf("Zendesk API Key Found: %s", key)

	t := time.Now().AddDate(0, 0, -3).Unix()
	log.Debugf("Time: %d", t)
	zenURL := url + "/api/v2/incremental/tickets.json?include=slas&start_time=" + strconv.FormatInt(t, 10)
	log.Debugf("URL: %s", zenURL)
	resp := makeRequest(user, key, zenURL)
	tickets = parseJSON(resp)
	log.Info("Request Complete. Parsing Ticket Data for", len(tickets.Tickets), "tickets")

	return tickets
}

// makeRequests takes the Zendesk auth information and sends the curl request
// to Zendesk and returns a JSON blob
func makeRequest(user string, key string, url string) (responseData []byte) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
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
		log.Critical(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	responseData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Critical(err)
		os.Exit(1)
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
		log.Error("error:", err)
		os.Exit(1)
	}
	return output
}

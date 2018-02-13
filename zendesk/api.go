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

// ZenOutput is the top level JSON-based struct that whatever is
// returned by Zendesk goes into
// TODO: Change Tickets to Tickets []Ticket
type ZenOutput struct {
	Tickets      `json:"tickets"`
	NextPage     interface{} `json:"next_page"`
	PreviousPage interface{} `json:"previous_page"`
	Count        int         `json:"count"`
}

// Tickets is a subset of ZenOutput that contains the details of the tickets
// outputted from the request to Zendesk
// TODO: use the OrgID to make a request for Org name using a different API call
// TODO: rename this Ticket, as it represents a singular entity
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
		Score   string `json:"score"`
		Comment string `json:"comment"`
		ID      int    `json:"id"`
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

type User struct {
	ID                   int64         `json:"id"`
	URL                  string        `json:"url"`
	Name                 string        `json:"name"`
	Email                string        `json:"email"`
	CreatedAt            time.Time     `json:"created_at"`
	UpdatedAt            time.Time     `json:"updated_at"`
	TimeZone             string        `json:"time_zone"`
	Phone                interface{}   `json:"phone"`
	SharedPhoneNumber    interface{}   `json:"shared_phone_number"`
	Photo                interface{}   `json:"photo"`
	LocaleID             int           `json:"locale_id"`
	Locale               string        `json:"locale"`
	OrganizationID       int64         `json:"organization_id"`
	Role                 string        `json:"role"`
	Verified             bool          `json:"verified"`
	ExternalID           interface{}   `json:"external_id"`
	Tags                 []interface{} `json:"tags"`
	Alias                interface{}   `json:"alias"`
	Active               bool          `json:"active"`
	Shared               bool          `json:"shared"`
	SharedAgent          bool          `json:"shared_agent"`
	LastLoginAt          time.Time     `json:"last_login_at"`
	TwoFactorAuthEnabled bool          `json:"two_factor_auth_enabled"`
	Signature            interface{}   `json:"signature"`
	Details              interface{}   `json:"details"`
	Notes                interface{}   `json:"notes"`
	RoleType             interface{}   `json:"role_type"`
	CustomRoleID         interface{}   `json:"custom_role_id"`
	Moderator            bool          `json:"moderator"`
	TicketRestriction    string        `json:"ticket_restriction"`
	OnlyPrivateComments  bool          `json:"only_private_comments"`
	RestrictedAgent      bool          `json:"restricted_agent"`
	Suspended            bool          `json:"suspended"`
	ChatOnly             bool          `json:"chat_only"`
	DefaultGroupID       interface{}   `json:"default_group_id"`
	UserFields           struct {
		Mrr                      int         `json:"mrr"`
		SystemEmbeddableLastSeen interface{} `json:"system::embeddable_last_seen"`
	} `json:"user_fields"`
}

// GetAllTickets grabs the latest tickets from Zendesk and returns the JSON
// Zendesk Endpoint: /incremental/tickets.json?include=slas
// TODO: update tickets.Tickets with new naimg scheme
func GetAllTickets(user string, key string, url string) (tickets ZenOutput) {
	log.Info("Starting request to Zendesk for tickets", map[string]interface{}{
		"module": "zendesk",
	})

	t := time.Now().AddDate(0, 0, -5).Unix()
	zenURL := url + "/api/v2/incremental/tickets.json?include=slas&start_time=" + strconv.FormatInt(t, 10)
	resp := makeRequest(user, key, zenURL)
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
func GetTicketRequester(user int) (output json.RawMessage) {
	log.Info("Starting request to Zendesk for user info", map[string]interface{}{
		"module": "zendesk",
	})

	zen := c.Zendesk.URL + "/api/v2/users/" + strconv.Itoa(user)
	data := makeRequest(c.Zendesk.User, c.Zendesk.APIKey, zen)
	log.Info("Request Complete. Parsing Ticket Data", map[string]interface{}{
		"module": "zendesk",
		"resp":   data,
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

// GetOrganization takes the org ID from the tickets grabbed in
// GetAllTickets and sends a request to Zendesk for the Org information
// Zendesk Endpoint /users/{USER-ID}/organizations.json
func GetOrganization(url string) {

}

// GetRequestedTickets takes a user ID and sends a request to Zendesk to grab
// the IDs of tickets requested by that user
// Zendesk Endpoint /users/{USER-ID}/tickets/requested.json
func GetRequestedTickets(url string) {

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

package zendesk

import "time"

// ZenOutput is the top level JSON-based struct that whatever is
// returned by Zendesk goes into
// TODO: Change Tickets to Tickets []Ticket
type ZenOutput struct {
	Tickets      `json:"tickets"`
	NextPage     interface{} `json:"next_page"`
	PreviousPage interface{} `json:"previous_page"`
	Count        int         `json:"count"`
}

// EventOutput is the top level JSON-based struct that whatever is
// returned by Zendesk goes into
type EventOutput struct {
	Event        `json:"ticket_events"`
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
	MetricEvents struct {
		PeriodicUpdateTime []struct {
			ID         int64     `json:"id"`
			TicketID   int       `json:"ticket_id"`
			Metric     string    `json:"metric"`
			InstanceID int       `json:"instance_id"`
			Type       string    `json:"type"`
			Time       time.Time `json:"time"`
			Status     struct {
				Calendar int `json:"calendar"`
				Business int `json:"business"`
			} `json:"status,omitempty"`
		} `json:"periodic_update_time"`
		RequesterWaitTime []struct {
			ID         int64     `json:"id"`
			TicketID   int       `json:"ticket_id"`
			Metric     string    `json:"metric"`
			InstanceID int       `json:"instance_id"`
			Type       string    `json:"type"`
			Time       time.Time `json:"time"`
		} `json:"requester_wait_time"`
		ResolutionTime []struct {
			ID         int64     `json:"id"`
			TicketID   int       `json:"ticket_id"`
			Metric     string    `json:"metric"`
			InstanceID int       `json:"instance_id"`
			Type       string    `json:"type"`
			Time       time.Time `json:"time"`
		} `json:"resolution_time"`
		PausableUpdateTime []struct {
			ID         int64     `json:"id"`
			TicketID   int       `json:"ticket_id"`
			Metric     string    `json:"metric"`
			InstanceID int       `json:"instance_id"`
			Type       string    `json:"type"`
			Time       time.Time `json:"time"`
			Status     struct {
				Calendar int `json:"calendar"`
				Business int `json:"business"`
			} `json:"status,omitempty"`
		} `json:"pausable_update_time"`
		AgentWorkTime []struct {
			ID         int64     `json:"id"`
			TicketID   int       `json:"ticket_id"`
			Metric     string    `json:"metric"`
			InstanceID int       `json:"instance_id"`
			Type       string    `json:"type"`
			Time       time.Time `json:"time"`
		} `json:"agent_work_time"`
		ReplyTime []struct {
			ID         int64     `json:"id"`
			TicketID   int       `json:"ticket_id"`
			Metric     string    `json:"metric"`
			InstanceID int       `json:"instance_id"`
			Type       string    `json:"type"`
			Time       time.Time `json:"time"`
			SLA        struct {
				Target        int  `json:"target"`
				BusinessHours bool `json:"business_hours"`
				Policy        struct {
					ID          int         `json:"id"`
					Title       string      `json:"title"`
					Description interface{} `json:"description"`
				} `json:"policy"`
			} `json:"sla,omitempty"`
			Deleted bool `json:"deleted,omitempty"`
			Status  struct {
				Calendar int `json:"calendar"`
				Business int `json:"business"`
			} `json:"status,omitempty"`
		} `json:"reply_time"`
	} `json:"metric_events"`
	AllowChannelback bool `json:"allow_channelback"`
}

type TicketGroup struct {
	Ticket `json:"ticket"`
}

type Ticket struct {
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
			Rel interface{} `json:"rel"`
		} `json:"source"`
	} `json:"via"`
	CreatedAt       time.Time     `json:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at"`
	Type            interface{}   `json:"type"`
	Subject         string        `json:"subject"`
	RawSubject      string        `json:"raw_subject"`
	Description     string        `json:"description"`
	Priority        string        `json:"priority"`
	Status          string        `json:"status"`
	Recipient       interface{}   `json:"recipient"`
	RequesterID     int64         `json:"requester_id"`
	SubmitterID     int64         `json:"submitter_id"`
	AssigneeID      int64         `json:"assignee_id"`
	OrganizationID  int64         `json:"organization_id"`
	GroupID         int           `json:"group_id"`
	CollaboratorIds []interface{} `json:"collaborator_ids"`
	FollowerIds     []interface{} `json:"follower_ids"`
	EmailCcIds      []interface{} `json:"email_cc_ids"`
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
	FollowupIds             []interface{} `json:"followup_ids"`
	TicketFormID            int           `json:"ticket_form_id"`
	BrandID                 int           `json:"brand_id"`
	SatisfactionProbability interface{}   `json:"satisfaction_probability"`
	AllowChannelback        bool          `json:"allow_channelback"`
}

type Event []struct {
	ChildEvents []struct {
		ID             int64       `json:"id"`
		Via            string      `json:"via"`
		ViaReferenceID interface{} `json:"via_reference_id"`
		Priority       string      `json:"priority"`
		EventType      string      `json:"event_type"`
		PreviousValue  string      `json:"previous_value"`
	} `json:"child_events"`
	ID        int64     `json:"id"`
	TicketID  int       `json:"ticket_id"`
	Timestamp int       `json:"timestamp"`
	CreatedAt time.Time `json:"created_at"`
	UpdaterID int64     `json:"updater_id"`
	Via       string    `json:"via"`
	System    struct {
		Client    string  `json:"client"`
		Location  string  `json:"location"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"system"`
	EventType string `json:"event_type"`
}

type Users struct {
	User `json:"user"`
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

type Organizations struct {
	Orgs         `json:"organizations"`
	NextPage     interface{} `json:"next_page"`
	PreviousPage interface{} `json:"previous_page"`
	Count        int         `json:"count"`
}

type Orgs []struct {
	URL                string      `json:"url"`
	ID                 int64       `json:"id"`
	Name               string      `json:"name"`
	SharedTickets      bool        `json:"shared_tickets"`
	SharedComments     bool        `json:"shared_comments"`
	ExternalID         interface{} `json:"external_id"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	DomainNames        []string    `json:"domain_names"`
	Details            string      `json:"details"`
	Notes              string      `json:"notes"`
	GroupID            interface{} `json:"group_id"`
	Tags               []string    `json:"tags"`
	OrganizationFields struct {
		SLALevel string `json:"sla_level"`
	} `json:"organization_fields"`
}

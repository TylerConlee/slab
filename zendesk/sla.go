package zendesk

import (
	c "github.com/tylerconlee/slab/config"
)

type ActiveTicket struct {
	ID       int
	Subject  string
	SLA      []interface{}
	Tags     []string
	Priority int
}

// CheckSLA will grab the tickets from GetAllTickets, parse the SLA fields and // compare them to the current time
func CheckSLA() (sla []ActiveTicket) {
	config := c.LoadConfig()
	zenResp := GetAllTickets(config.Zendesk.User, config.Zendesk.APIKey, config.Zendesk.URL)

	for _, ticket := range zenResp.Tickets {
		priority := getPriorityLevel(ticket.Tags)

		Log.Debug("ID:", ticket.ID, ", Title:", ticket.Subject, ", SLA:", ticket.Slas, ", Tags:", ticket.Tags, ", Priority:", priority)

		if priority > 0 {
			t := ActiveTicket{
				ID:       ticket.ID,
				Priority: priority,
				SLA:      ticket.Slas.PolicyMetrics,
				Tags:     ticket.Tags,
				Subject:  ticket.Subject,
			}
			sla = append(sla, t)
		}

	}
	return sla
}

// getPriorityLevel takes an individual ticket row from the Zendesk output and // returns a string of what priority level the ticket is tagged with
func getPriorityLevel(tags []string) (priLvl int) {
	for _, v := range tags {
		if v == "platinum" {
			return 1
		}
		if v == "gold" {
			return 2
		}
		if v == "silver" {
			return 3
		}
	}
	return
}

package zendesk

import (
	c "github.com/tylerconlee/slab/config"
)

// ActiveTicket is the individual ticket details for a ticket
// that's nearing SLA breach. This is passed to the main function so the
// breach time can be compared
type ActiveTicket struct {
	ID       int
	Subject  string
	SLA      []interface{}
	Tags     []string
	Level    string
	Priority interface{}
}

// CheckSLA will grab the tickets from GetAllTickets, parse the SLA fields and // compare them to the current time
func CheckSLA() (sla []ActiveTicket) {
	config := c.LoadConfig()
	zenResp := GetAllTickets(config.Zendesk.User, config.Zendesk.APIKey, config.Zendesk.URL)

	for _, ticket := range zenResp.Tickets {
		priority := getPriorityLevel(ticket.Tags)

		Log.Debug("ID:", ticket.ID, ", Title:", ticket.Subject, ", SLA:", ticket.Slas, ", Tags:", ticket.Tags, ", Priority:", priority)

		if priority != "" {
			t := ActiveTicket{
				ID:       ticket.ID,
				Level:    priority,
				SLA:      ticket.Slas.PolicyMetrics,
				Tags:     ticket.Tags,
				Subject:  ticket.Subject,
				Priority: ticket.Priority,
			}
			sla = append(sla, t)
		}

	}
	return sla
}

// getPriorityLevel takes an individual ticket row from the Zendesk output and // returns a string of what priority level the ticket is tagged with
// TODO: make agnostic to our levels
func getPriorityLevel(tags []string) (priLvl string) {
	for _, v := range tags {
		if v == "platinum" {
			return "LevelOne"
		}
		if v == "gold" {
			return "LevelTwo"
		}
		if v == "silver" {
			return "LevelThree"
		}
	}
	return
}

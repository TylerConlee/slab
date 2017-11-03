package zendesk

import (
	"time"

	c "github.com/tylerconlee/slab/config"
)

var config = c.LoadConfig()

// ActiveTicket is the individual ticket details for a ticket
// that's nearing SLA breach. This is passed to the main function so the
// breach time can be compared
type ActiveTicket struct {
	ID          int
	Subject     string
	SLA         []interface{}
	Tags        []string
	Level       string
	Priority    interface{}
	CreatedAt   time.Time
	Description string
}

// CheckSLA will grab the tickets from GetAllTickets, parse the SLA fields and // compare them to the current time
func CheckSLA() (sla []ActiveTicket) {

	zenResp := GetAllTickets(config.Zendesk.User, config.Zendesk.APIKey, config.Zendesk.URL)

	for _, ticket := range zenResp.Tickets {
		priority := getPriorityLevel(ticket.Tags)

		if priority != "" {
			t := ActiveTicket{
				ID:          ticket.ID,
				Level:       priority,
				SLA:         ticket.Slas.PolicyMetrics,
				Tags:        ticket.Tags,
				Subject:     ticket.Subject,
				Priority:    ticket.Priority,
				CreatedAt:   ticket.CreatedAt,
				Description: ticket.Description,
			}
			Log.Debug("Ticket", ticket.ID, "successfully parsed. SLA found:", ticket.Slas.PolicyMetrics)
			sla = append(sla, t)
		} else {
			Log.Debug("Ticket", ticket.ID, "successfully parsed.")
		}

	}
	return sla
}

// getPriorityLevel takes an individual ticket row from the Zendesk output and
// returns a string of what priority level the ticket is tagged with
func getPriorityLevel(tags []string) (priLvl string) {
	for _, v := range tags {
		if v == config.SLA.LevelOne.Tag {
			return "LevelOne"
		}
		if v == config.SLA.LevelTwo.Tag {
			return "LevelTwo"
		}
		if v == config.SLA.LevelThree.Tag {
			return "LevelThree"
		}
	}
	return
}

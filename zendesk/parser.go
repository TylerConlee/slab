// Package zendesk provides functions to grab and manipulate ticket data from a
// Zendesk instance
package zendesk

import (
	"time"

	"github.com/tylerconlee/slab/config"
)

// config loads the configuration
var c = config.LoadConfig()

// ActiveTicket is the individual ticket details for a ticket
// that's nearing SLA breach. This is passed to the main function so the
// breach time can be compared
type ActiveTicket struct {
	ID          int
	Requester   int64
	Subject     string
	SLA         []interface{}
	Tags        []string
	Level       string
	Priority    interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Description string
}

// CheckSLA will grab the tickets from GetAllTickets, parse the SLA fields and // compare them to the current time
func CheckSLA(tick ZenOutput) (sla []ActiveTicket) {

	for _, ticket := range tick.Tickets {
		priority := getPriorityLevel(ticket.Tags)

		if priority != "" {
			t := ActiveTicket{
				ID:          ticket.ID,
				Level:       priority,
				SLA:         ticket.Slas.PolicyMetrics,
				Tags:        ticket.Tags,
				Requester:   ticket.RequesterID,
				Subject:     ticket.Subject,
				Priority:    ticket.Priority,
				CreatedAt:   ticket.CreatedAt,
				Description: ticket.Description,
				UpdatedAt:   ticket.UpdatedAt,
			}
			sla = append(sla, t)
		}

	}
	return sla
}

// CheckNewTicket loops over the Zendesk output from GetAllTickets and
// determines if there are tickets that have been created since the last loop
func CheckNewTicket(tick ZenOutput, interval time.Duration) (new []ActiveTicket) {
	previousLoop := time.Now().Add(-interval * 3)
	nowLoop := time.Now()
	log.Info("Checking for new tickets", map[string]interface{}{
		"module":      "zendesk",
		"currentloop": nowLoop,
		"prevloop":    previousLoop,
	})
	for _, ticket := range tick.Tickets {
		if ticket.CreatedAt.After(previousLoop) && ticket.CreatedAt.Before(nowLoop) {
			t := ActiveTicket{
				ID:          ticket.ID,
				SLA:         ticket.Slas.PolicyMetrics,
				Tags:        ticket.Tags,
				Subject:     ticket.Subject,
				Priority:    ticket.Priority,
				CreatedAt:   ticket.CreatedAt,
				Description: ticket.Description,
			}
			new = append(new, t)
		}
	}
	return new
}

// CheckUpdatedTicket loops over the Zendesk output from GetAllTickets
func CheckUpdatedTicket(interval time.Duration) (new []ActiveTicket) {
	previousLoop := time.Now().Add(-interval)
	nowLoop := time.Now()
	tick := GetTicketEvents()
	var ids []int64
	for _, event := range tick.Event {

		if event.CreatedAt.After(previousLoop) && event.CreatedAt.Before(nowLoop) && !eventExists(event.ID, ids) {

			user := GetTicketRequester(int(event.UpdaterID))
			if user.Role == "end-user" {
				ticket := GetTicket(event.TicketID)
				priority := getPriorityLevel(ticket.Tags)
				log.Info("Parsing updated ticket", map[string]interface{}{
					"module":   "zendesk",
					"priority": priority,
					"ticketID": ticket.ID,
				})
				if priority != "" && priority != "LevelFour" {
					t := ActiveTicket{
						ID:          ticket.ID,
						Tags:        ticket.Tags,
						Subject:     ticket.Subject,
						Priority:    ticket.Priority,
						CreatedAt:   ticket.CreatedAt,
						UpdatedAt:   ticket.UpdatedAt,
						Description: ticket.Description,
					}
					new = append(new, t)
				}
			}

			ids = append(ids, event.ID)

		}
	}
	return new
}

func eventExists(a int64, list []int64) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// getPriorityLevel takes an individual ticket row from the Zendesk output and
// returns a string of what priority level the ticket is tagged with
func getPriorityLevel(tags []string) (priLvl string) {
	for _, v := range tags {
		switch v {
		case c.SLA.LevelOne.Tag:
			return "LevelOne"
		case c.SLA.LevelTwo.Tag:
			return "LevelTwo"
		case c.SLA.LevelThree.Tag:
			return "LevelThree"
		case c.SLA.LevelFour.Tag:
			return "LevelFour"
		}
	}
	return
}

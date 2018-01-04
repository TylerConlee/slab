package zendesk

import (
	"fmt"
	"time"

	"github.com/tylerconlee/slab/slack"
)

// RunTimer takes the interval from the config, and at each loop iteration,
// grabs the latest tickets, checks for upcoming SLAs and send notifications if
// appropriate
func RunTimer(interval time.Duration) {
	log.Info("Starting timer with ", interval, " intervals")
	t := time.NewTicker(interval)
	for {
		tick := GetAllTickets(
			c.Zendesk.User,
			c.Zendesk.APIKey,
			c.Zendesk.URL,
		)

		log.Info("Successfully grabbed and parsed tickets from Zendesk")
		log.Info("Checking ticket notifications...")

		// Returns a list of all upcoming SLA breaches
		active := CheckSLA(tick)

		// Loop through all active SLA tickets and prepare SLA notification
		// for each.
		for _, ticket := range active {

			if ticket.Priority != nil {
				send, notify := UpdateCache(ticket)
				if send {
					n := PrepSLANotification(ticket, notify)
					m := slack.Ticket(ticket)
					slack.SLAMessage(n, m)
				}
			}
		}

		// Returns a list of all new tickets within the last loop
		new := CheckNewTicket(tick, interval)
		var newTickets []slack.Ticket
		// Loop through all tickets and check
		for _, ticket := range new {
			m := slack.Ticket(ticket)
			log.Debug("Adding new ticket to notification: %x", m)
			newTickets = append(newTickets, m)
		}
		slack.NewTicketMessage(newTickets)

		log.Info("Ticket notifications sent. Returning to idle state.")
		<-t.C
	}
}

// PrepSLANotification takes a given ticket and what notification level and returns a string to be sent to Slack.
func PrepSLANotification(ticket ActiveTicket, notify int64) (notification string) {
	log.Debug("Preparing notification for", ticket.ID)
	var t, p string
	var r bool

	switch ticket.Level {
	case "LevelOne":
		p = c.SLA.LevelOne.Tag
		r = c.SLA.LevelOne.Notify
	case "LevelTwo":
		p = c.SLA.LevelTwo.Tag
		r = c.SLA.LevelTwo.Notify

	case "LevelThree":
		p = c.SLA.LevelThree.Tag
		r = c.SLA.LevelThree.Notify

	case "LevelFour":
		p = c.SLA.LevelFour.Tag
		r = c.SLA.LevelFour.Notify
	}

	var n string

	switch notify {
	case 1:
		t = "15 minutes"
	case 2:
		t = "30 minutes"
	case 3:
		t = "1 hour"
	case 4:
		t = "2 hours"
	case 5:
		t = "3 hours"
	}
	if r {
		n = fmt.Sprintf("@here SLA for *%s* ticket #%d has less than %s until expiration.", p, ticket.ID, t)
	} else {
		n = fmt.Sprintf("SLA for *%s* ticket #%d has less than %s until expiration.", p, ticket.ID, t)
	}

	return n

}

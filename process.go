package main

import (
	"time"

	"github.com/tylerconlee/slab/slack"
	"github.com/tylerconlee/slab/zendesk"
)

// RunTimer takes the interval from the config, and at each loop iteration,
// grabs the latest tickets, checks for upcoming SLAs and send notifications if
// appropriate
func RunTimer(interval time.Duration) {
	log.Info("Starting timer", map[string]interface{}{
		"module":   "main",
		"interval": interval,
	})
	t := time.NewTicker(interval)
	for {
		// TODO: Add func to connect to Zendesk and pass single config
		tick := zendesk.GetAllTickets(
			c.Zendesk.User,
			c.Zendesk.APIKey,
			c.Zendesk.URL,
		)

		log.Info("Successfully grabbed and parsed tickets from Zendesk", map[string]interface{}{
			"module": "main",
		})
		log.Info("Checking ticket notifications...", map[string]interface{}{
			"module": "main",
		})
		// Returns a list of all upcoming SLA breaches
		active := zendesk.CheckSLA(tick)

		// Loop through all active SLA tickets and prepare SLA notification
		// for each.
		for _, ticket := range active {

			if ticket.Priority != nil {
				send, notify := zendesk.UpdateCache(ticket)
				if send {
					log.Info("Preparing SLA notification for ticket", map[string]interface{}{
						"module": "main",
						"ticket": ticket.ID,
					})
					m := slack.Ticket(ticket)
					n, c := slack.PrepSLANotification(m, notify)
					slack.SLAMessage(n, m, c)
				}
			}
		}

		slack.Sent = zendesk.Sent

		// Returns a list of all new tickets within the last loop
		new := zendesk.CheckNewTicket(tick, interval)
		var newTickets []slack.Ticket
		// Loop through all tickets and add to Slack package friendly slice
		for _, ticket := range new {
			m := slack.Ticket(ticket)
			log.Info("Adding new ticket to notification", map[string]interface{}{
				"module": "main",
				"ticket": m.ID,
			})
			log.Debug("Ticket information", map[string]interface{}{
				"module": "main",
				"ticket": m,
			})
			newTickets = append(newTickets, m)
		}
		slack.NewTicketMessage(newTickets)

		log.Info("Ticket notifications sent. Returning to idle state.", map[string]interface{}{
			"module": "main",
		})
		<-t.C
	}
}

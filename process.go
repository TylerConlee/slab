package main

import (
	"time"

	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/plugins"
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
		// reload the config on each pass to allow for changes to the config to
		// be recognized
		// TODO: Update this to only update when the file is modified, rather
		// than every pass
		c = config.LoadConfig()
		p := plugins.LoadPlugins(c)
		if c.Slack.ChannelID != "" {
			channel := slack.GetChannel(c.Slack.ChannelID)
			if channel == 0 {
				slack.AddChannel(c.Slack.ChannelID, 2)
			}
		}
		log.Info("Loaded plugins.", map[string]interface{}{
			"module":  "main",
			"plugins": p,
		})
		if c.Zendesk.URL != "" && c.Zendesk.APIKey != "" {
			tick := zendesk.GetAllTickets()

			log.Info("Grabbed and parsed tickets from Zendesk", map[string]interface{}{
				"module": "main",
			})
			log.Info("Checking ticket notifications...", map[string]interface{}{
				"module": "main",
			})
			// Returns a list of all upcoming SLA breaches
			active := zendesk.CheckSLA(tick)
			updated := zendesk.CheckUpdatedTicket(interval)

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
						p.SendDispatcher(n)
						user := zendesk.GetTicketRequester(int(ticket.Requester))
						slack.SLAMessage(n, m, c, user.Name, user.ID)
					}
				}
			}
			for _, ticket := range updated {
				log.Info("Preparing update notification for ticket", map[string]interface{}{
					"module": "main",
					"ticket": ticket.ID,
				})
				m := slack.Ticket(ticket)
				user := zendesk.GetTicketRequester(int(ticket.Requester))
				slack.UpdateMessage(m, user.Name, user.ID)
			}

			slack.Sent = zendesk.Sent

			// Returns a list of all new tickets within the last loop
			new := zendesk.CheckNewTicket(tick, interval)
			if new != nil {

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

			}
			log.Info("Ticket notifications sent. Returning to idle state.", map[string]interface{}{

				"module": "main",
			})
		} else {
			log.Info("Zendesk authorization required. Please run @slab start config to begin.", map[string]interface{}{

				"module": "main",
			})
		}
		<-t.C
	}
}

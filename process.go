package main

import (
	"fmt"
	"strings"
	"time"

	sl "github.com/nlopes/slack"
	plugins "github.com/tylerconlee/slab/_plugins"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/datastore"
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
		iteration(t, interval)
		<-t.C

	}
}

func iteration(t *time.Ticker, interval time.Duration) {
	// reload the config on each pass to allow for changes to the config to
	// be recognized
	// TODO: Update this to only update when the file is modified, rather
	// than every pass
	c := config.LoadConfig()

	p := plugins.LoadPlugins(c)

	if c.Slack.ChannelID != "" {
		channel := slack.GetChannel(c.Slack.ChannelID)
		if channel == 0 {
			slack.AddChannel(c.Slack.ChannelID, 2)

		}
		log.Info("Loaded plugins.", map[string]interface{}{
			"module":  "main",
			"plugins": p,
		})

		if c.Zendesk.URL != "" && c.Zendesk.APIKey != "" {
			tick := zendesk.GetAllTickets()

			// Grab tags
			tags := datastore.LoadTags()

			log.Info("Grabbed and parsed tickets from Zendesk", map[string]interface{}{
				"module": "main",
			})
			log.Info("Checking ticket notifications...", map[string]interface{}{
				"module": "main",
			})
			processSLAAlerts(tick, tags, p)
			processUpdateAlerts(tick, tags, p, interval)
			processNewAlerts(tick, tags, p, interval)
			slack.Sent = zendesk.Sent
			slack.NumTickets = zendesk.NumTickets
			slack.LastProcessed = zendesk.LastProcessed

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

func getOrgName(id int) (o string) {
	org := zendesk.GetOrganization(id)
	if len(org) > 0 {
		orglink := fmt.Sprintf("%s/agent/organizations/%d", c.Zendesk.URL, org[0].ID)

		o = "<" + orglink + "| " + org[0].Name + "> "
	}
	return o

}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func processSLAAlerts(tick zendesk.ZenOutput, tags []map[string]interface{}, p plugins.Plugins) {
	// Returns a list of all upcoming SLA breaches
	active := zendesk.CheckSLA(tick)

	// Loop through all active SLA tickets and prepare SLA notification
	// for each.
	for _, ticket := range active {
		// Loop through all available tags
		for _, tag := range tags {
			if contains(ticket.Tags, tag["tag"].(string)) && tag["notify_type"].(string) == "sla" {
				if ticket.Priority != nil {
					send, notify := zendesk.UpdateCache(ticket, tag["channel"].(string))
					if send {
						log.Info("Preparing SLA notification for ticket", map[string]interface{}{
							"module":  "main",
							"ticket":  ticket.ID,
							"channel": tag["channel"].(string),
						})
						m := slack.Ticket(ticket)
						n, c := slack.PrepSLANotification(m, notify, tag["tag"].(string))
						p.SendDispatcher(n)
						user := zendesk.GetTicketRequester(int(ticket.Requester))
						org := getOrgName(ticket.ID)
						attach := slack.SLAMessage(m, c, user.Name, user.ID, org)
						attachments := []sl.Attachment{attach}
						if strings.HasPrefix(tag["channel"].(string), "U") {

							slack.SendDirectMessage(n, attachments, tag["channel"].(string))
						} else {
							slack.SendMessage(n, tag["channel"].(string), attachments)
						}

					}
				}
			}
		}

	}
}

func processUpdateAlerts(tick zendesk.ZenOutput, tags []map[string]interface{}, p plugins.Plugins, interval time.Duration) {
	updated := zendesk.CheckUpdatedTicket(interval)
	for _, ticket := range updated {
		for _, tag := range tags {
			if contains(ticket.Tags, tag["tag"].(string)) && tag["notify_type"].(string) == "updates" {
				log.Info("Preparing update notification for ticket", map[string]interface{}{
					"module": "main",
					"ticket": ticket.ID,
				})
				n := fmt.Sprintf("Ticket #%d updated. Priority: %s, Tag: %s", ticket.ID, ticket.Priority, tag["tag"])
				m := slack.Ticket(ticket)
				user := zendesk.GetTicketRequester(int(ticket.Requester))
				p.SendDispatcher(n)
				attach := slack.UpdateMessage(m, user.Name, user.ID)
				attachments := []sl.Attachment{attach}
				if strings.HasPrefix(tag["channel"].(string), "U") {
					slack.SendDirectMessage(n, attachments, tag["channel"].(string))
				} else {
					slack.SendMessage(n, tag["channel"].(string), attachments)
				}
			}
		}

	}
}

func processNewAlerts(tick zendesk.ZenOutput, tags []map[string]interface{}, p plugins.Plugins, interval time.Duration) {
	// Returns a list of all new tickets within the last loop
	new := zendesk.CheckNewTicket(tick, interval)
	if new != nil {

		var newTickets []slack.Ticket
		// Loop through all tickets and add to Slack package friendly slice
		for _, ticket := range new {
			for _, tag := range tags {
				if contains(ticket.Tags, tag["tag"].(string)) && tag["notify_type"].(string) == "new" {
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
					attachments, message := slack.NewTicketMessage(newTickets, tag["tag"].(string))
					if strings.HasPrefix(tag["channel"].(string), "U") {
						slack.SendDirectMessage(message, attachments, tag["channel"].(string))
					} else {
						if tag["channel"].(string) == c.Slack.ChannelID {
							if slack.Triager != "None" {
								message = fmt.Sprintf("<@%s> The following tickets were received since the last loop:", slack.Triager)
							} else {
								message = fmt.Sprintf("The following tickets were received since the last loop:")
							}
						}
						slack.SendMessage(message, tag["channel"].(string), attachments)
					}
				}
			}
		}
	}
}

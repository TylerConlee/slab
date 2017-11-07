package main

import (
	"time"

	"github.com/tylerconlee/slab/sla"
	"github.com/tylerconlee/slab/slack"

	Zen "github.com/tylerconlee/slab/zendesk"
)

// RunTimer takes the interval from the config, and at each loop iteration,
// grabs the latest tickets, checks for upcoming SLAs and send notifications if
// appropriate
func RunTimer(interval time.Duration) {
	log.Info("Starting timer with ", interval, " intervals")
	t := time.NewTicker(interval)
	for {
		active := Zen.CheckSLA()
		log.Info("Successfully grabbed and parsed tickets from Zendesk")
		log.Info("Checking ticket notifications...")
		for _, ticket := range active {

			if ticket.Priority != nil {
				send, notify := sla.UpdateCache(ticket)
				if send {
					n := slack.PrepNotification(ticket, notify)
					slack.Send(n, ticket)
				}
			}
		}
		log.Info("Ticket notifications sent. Returning to idle state.")
		<-t.C
	}
}

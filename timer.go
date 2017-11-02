package main

import (
	"time"

	"github.com/tylerconlee/slab/sla"

	Zen "github.com/tylerconlee/slab/zendesk"
)

func RunTimer(interval time.Duration) {
	log.Debug(interval)
	t := time.NewTicker(interval)
	for {
		active := Zen.CheckSLA()
		log.Debug(active)

		for _, ticket := range active {
			if ticket.Priority != nil {
				send := sla.UpdateCache(ticket)
				log.Debug(send)
			}
		}
		<-t.C
	}
}

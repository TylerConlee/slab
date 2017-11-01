package main

import (
	"time"

	"github.com/tylerconlee/slab/sla"

	Zen "github.com/tylerconlee/slab/zendesk"
)

func RunTimer(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Minute)
	for {
		active := Zen.CheckSLA()
		log.Debug(active)

		sla.SetSLA()
		<-t.C
	}
}

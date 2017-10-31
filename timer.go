package main

import (
	"time"

	Zen "github.com/tylerconlee/slab/zendesk"
)

func RunTimer(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Minute)
	for {
		active := Zen.CheckSLA()
		log.Debug(active)
		<-t.C
	}
}

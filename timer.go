package main

import (
	"time"

	Zen "github.com/tylerconlee/slab/zendesk"
)

func StartTimer(interval int) {
	t := time.NewTicker(time.Duration(interval) * time.Minute)
	for {
		Zen.CheckSLA()
		<-t.C
	}
}

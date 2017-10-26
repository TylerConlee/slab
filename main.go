package main

import (
	Zen "github.com/tylerconlee/slab/zendesk"
)

func main() {
	initLog()
	log.Debugf("debug %s")
	log.Info("info")
	log.Notice("notice")
	log.Warning("warning")
	log.Error("err")
	log.Critical("crit")

	Zen.GetAllTickets()
}

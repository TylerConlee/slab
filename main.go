package main

import (
	logging "github.com/op/go-logging"
	Zen "github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("slab")

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

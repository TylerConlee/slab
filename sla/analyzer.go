package sla

import (
	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
)

var c config.Config
var log = logging.MustGetLogger("sla")

func SetSLA() {
	c = config.LoadConfig()
	log.Debug(c.SLA)
}

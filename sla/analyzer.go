package sla

import (
	"reflect"
	"strings"
	"time"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/zendesk"
)

var c config.Config
var log = logging.MustGetLogger("sla")

// InitSLA loads the configuration and any other package-wide variables
// needed for processing SLAs
func InitSLA() {
	c = config.LoadConfig()
}

// GetTimer takes an instance of a ticket and determines the SLA timer that it // uses. For example, if a ticket has a tag that matches the tag for LevelOne,
// and it has a priority of 'high', it will use config.SLA.LevelOne.High as the
// timer.
func GetTimer(ticket zendesk.ActiveTicket) (breach time.Duration) {
	r := reflect.ValueOf(c.SLA)
	f := reflect.Indirect(r).FieldByName(ticket.Level)
	priority := strings.Title(ticket.Priority.(string))
	p := reflect.Indirect(f).FieldByName(priority)
	str := p.Interface().(config.Duration)
	return time.Duration(str.Duration)
}

// GetBreach takes an instance of a ticket and returns the value of the next SLA
// breach.
func GetBreach(ticket zendesk.ActiveTicket) {
	p := ticket.SLA[0].(map[string]interface{})
	log.Debug(p["breach_at"])
}

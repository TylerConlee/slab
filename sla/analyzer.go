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

func InitSLA() {
	c = config.LoadConfig()
}
func GetTimer(ticket zendesk.ActiveTicket) (breach time.Duration) {
	r := reflect.ValueOf(c.SLA)
	f := reflect.Indirect(r).FieldByName(ticket.Level)
	priority := strings.Title(ticket.Priority.(string))
	p := reflect.Indirect(f).FieldByName(priority)
	str := p.Interface().(config.Duration)
	return time.Duration(str.Duration)
}

func GetBreach(ticket zendesk.ActiveTicket) {
	p := ticket.SLA[0].(map[string]interface{})
	log.Debug(p["breach_at"])
}

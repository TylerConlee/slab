package sla

import (
	"os"
	"time"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("sla")

// GetTimeRemaining takes an instance of a ticket and returns the value of the next SLA
// breach.
func GetTimeRemaining(ticket zendesk.ActiveTicket) {
	p := ticket.SLA[0].(map[string]interface{})
	breach, err := time.Parse(time.RFC3339, p["breach_at"].(string))
	if nil != err {
		log.Critical(err)
		os.Exit(1)
	}

	remain := time.Until(breach)
	log.Debug(remain)
}

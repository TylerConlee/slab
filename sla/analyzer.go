package sla

import (
	"fmt"
	"os"
	"time"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("sla")

// GetTimeRemaining takes an instance of a ticket and returns the value of the next SLA
// breach.
func GetTimeRemaining(ticket zendesk.ActiveTicket) (remain time.Duration) {
	p := ticket.SLA[0].(map[string]interface{})
	breach, err := time.Parse(time.RFC3339, p["breach_at"].(string))
	if nil != err {
		log.Critical(err)
		os.Exit(1)
	}

	remain = time.Until(breach)
	return remain
}

func GetNotifyTime(remain time.Duration) (notifyType int) {
	p, _ := time.ParseDuration("3h")
	q, _ := time.ParseDuration("2h")
	r, _ := time.ParseDuration("1h")
	s, _ := time.ParseDuration("30m")
	t, _ := time.ParseDuration("15m")

	switch {
	case remain < t:
		log.Debug("Send 15 minute notification")
		return 1
	case remain < s:
		log.Debug("Send 30 minute notification")
		return 2
	case remain < r:
		log.Debug("Send 1 hour notification")
		return 3
	case remain < q:
		log.Debug("Send 2 hour notification")
		return 4
	case remain < p:
		log.Debug("Send 3 hour notification")
		return 5
	default:
		fmt.Println("Too far away.")
		return 0
	}
}

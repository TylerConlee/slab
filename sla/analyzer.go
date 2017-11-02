package sla

import (
	"os"
	"reflect"
	"time"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("sla")

// Sent is a collection of all NotifySent tickets that is checked before each // notification is sent.
var Sent = []NotifySent{}

// NotifySent is represetative of an individual ticket, what kind of
// notification was last sent for that ticket, and when the SLA breach time is.
type NotifySent struct {
	ID     int
	Type   int64
	Expire time.Time
}

// GetTimeRemaining takes an instance of a ticket and returns the value of the next SLA
// breach.
func GetTimeRemaining(ticket zendesk.ActiveTicket) (remain time.Time) {
	p := ticket.SLA[0].(map[string]interface{})
	breach, err := time.Parse(time.RFC3339, p["breach_at"].(string))
	if nil != err {
		log.Critical(err)
		os.Exit(1)
	}

	return breach
}

// GetNotifyType - Based off of the time remaining on the ticket, return a
// integer representing the closest time marker to a notification time.
func GetNotifyType(remain time.Duration) (notifyType int64) {
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
		log.Debug("SLA is longer than 3 hours away")
		return 0
	}
}

// UpdateCache checks the time remaining on a ticket, what the closest marker
// for notifications is, and then checks to see if that ticket ID and
// notification type have been sent already. If yes, it returns True,
// indicating a notifcation needs to be sent.
func UpdateCache(ticket zendesk.ActiveTicket) bool {
	cleanCache()
	expire := GetTimeRemaining(ticket)
	notify := GetNotifyType(time.Until(expire))
	t := expire.Add(15 * time.Minute)
	if t.After(time.Now()) && notify != 0 {
		log.Debug(Sent, notify, ticket.ID)
		rangeOnMe := reflect.ValueOf(Sent)
		for i := 0; i < rangeOnMe.Len(); i++ {
			s := rangeOnMe.Index(i)
			f := s.FieldByName("ID")
			if f.IsValid() {
				if f.Interface() == ticket.ID && s.FieldByName("Type").Int() == notify {
					return false
				}

			}

		}
		Sent = append(Sent, NotifySent{ticket.ID, notify, expire})

		return true
	}
	return false

}

// cleanCache checks the Sent slice and loops through the tickets listed. If
// any have gone 15 minutes past the expiration time, they are removed from the
// slice and the length of the slice is shortened.
func cleanCache() {
	for i := 0; i < len(Sent); i++ {
		item := Sent[i]
		t := item.Expire.Add(15 * time.Minute)
		if t.Before(time.Now()) {
			Sent = append(Sent[:i], Sent[i+1:]...)
			i--
		}
	}

}

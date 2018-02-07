package zendesk

import (
	"reflect"
	"time"
)

// Sent is a collection of all NotifySent tickets that is checked before each // notification is sent.
var Sent = []NotifySent{}

// NotifySent is represetative of an individual ticket, what kind of
// notification was last sent for that ticket, and when the SLA breach time is.
type NotifySent struct {
	ID     int
	Type   int64
	Expire time.Time
}

// GetTimeRemaining takes an instance of a ticket and returns the value of the
// next SLA breach.
func GetTimeRemaining(ticket ActiveTicket) (remain time.Time) {
	if len(ticket.SLA) >= 1 {
		p := ticket.SLA[0].(map[string]interface{})
		if p["breach_at"] != nil {
			breach, err := time.Parse(time.RFC3339, p["breach_at"].(string))
			if nil != err {
				log.Fatal(map[string]interface{}{
					"module": "zendesk",
					"error":  err,
				})
			}

			return breach
		}
	}
	return
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
		return 1
	case remain < s:
		return 2
	case remain < r:
		return 3
	case remain < q:
		return 4
	case remain < p:
		return 5
	default:
		return 0
	}
}

// UpdateCache checks the time remaining on a ticket, what the closest marker
// for notifications is, and then checks to see if that ticket ID and
// notification type have been sent already. If yes, it returns True,
// indicating a notifcation needs to be sent.
func UpdateCache(ticket ActiveTicket) (bool, int64) {
	//cleanCache()

	// get the expiration timestamp
	expire := GetTimeRemaining(ticket)
	notify := GetNotifyType(time.Until(expire))

	// take the ticket expiration time and add 15 minutes
	t := expire.Add(15 * time.Minute)

	// if the ticket expiration time is after 15 minutes from now and there's a
	// valid notification type
	if t.After(time.Now()) && notify != 0 {
		rangeOnMe := reflect.ValueOf(Sent)
		for i := 0; i < rangeOnMe.Len(); i++ {
			s := rangeOnMe.Index(i)
			f := s.FieldByName("ID")
			if f.IsValid() {
				if f.Interface() == ticket.ID && s.FieldByName("Type").Int() == notify {
					log.Info("Ticket has already received a notification", map[string]interface{}{
						"module":      "zendesk",
						"ticket":      ticket.ID,
						"notify_type": notify,
						"expires":     expire,
					})
					return false, 0
				}

			}

		}
		Sent = append(Sent, NotifySent{ticket.ID, notify, expire})
		log.Info("Ticket should receive a notification", map[string]interface{}{
			"module":      "zendesk",
			"ticket":      ticket.ID,
			"notify_type": notify,
			"expires":     expire,
		})
		return true, notify
	}

	return false, 0

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

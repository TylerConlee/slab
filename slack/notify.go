package slack

import (
	"fmt"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("slack")

// PrepNotification takes a given ticket and what notification level and returns a string to be sent to Slack.
func PrepNotification(ticket zendesk.ActiveTicket, notify int64) (notification string) {
	log.Debug("Preparing notification for", ticket.ID)
	c := config.LoadConfig()
	var t, p string
	var r bool

	switch ticket.Level {
	case "LevelOne":
		p = c.SLA.LevelOne.Tag
		r = c.SLA.LevelOne.Notify
	case "LevelTwo":
		p = c.SLA.LevelTwo.Tag
		r = c.SLA.LevelTwo.Notify

	case "LevelThree":
		p = c.SLA.LevelThree.Tag
		r = c.SLA.LevelThree.Notify

	case "LevelFour":
		p = c.SLA.LevelFour.Tag
		r = c.SLA.LevelFour.Notify
	}

	var n string

	switch notify {
	case 1:
		t = "15 minutes"
		if r {
			n = fmt.Sprintf("@here SLA for *%s* ticket has less than %s until expiration.", p, t)
		} else {
			n = fmt.Sprintf("SLA for *%s* ticket has less than %s until expiration.", p, t)
		}

	case 2:
		t = "30 minutes"
		if r {
			n = fmt.Sprintf("@here SLA for *%s* ticket has less than %s until expiration.", p, t)
		} else {
			n = fmt.Sprintf("SLA for *%s* ticket has less than %s until expiration.", p, t)
		}
	case 3:
		t = "1 hour"
		n = fmt.Sprintf("@here SLA for *%s* ticket has less than %s until expiration.", p, t)
	case 4:
		t = "2 hours"
		n = fmt.Sprintf("SLA for *%s* ticket has less than %s until expiration.", p, t)
	case 5:
		t = "3 hours"
		n = fmt.Sprintf("SLA for *%s* ticket has less than %s until expiration.", p, t)
	}

	return n

}

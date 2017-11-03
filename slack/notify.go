package slack

import (
	"fmt"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("notification")

func PrepNotification(ticket zendesk.ActiveTicket, notify int64) (notification string) {
	log.Debug("Preparing notification for", ticket.ID)
	c := config.LoadConfig()
	var t string
	switch notify {
	case 1:
		t = "15 minutes"
	case 2:
		t = "30 minutes"
	case 3:
		t = "1 hour"
	case 4:
		t = "2 hours"
	case 5:
		t = "3 hours"
	}

	var p string
	switch ticket.Level {
	case "LevelOne":
		p = c.SLA.LevelOne.Tag

	case "LevelTwo":
		p = c.SLA.LevelTwo.Tag

	case "LevelThree":
		p = c.SLA.LevelThree.Tag

	case "LevelFour":
		p = c.SLA.LevelFour.Tag
	}
	var n string
	n = fmt.Sprintf("SLA for %s ticket, %d - %s, has less than %s until expiration. %s/%d", p, ticket.ID, ticket.Subject, t, c.Zendesk.URL, ticket.ID)
	return n

}

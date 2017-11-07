package slack

import (
	"fmt"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/zendesk"
)

var log = logging.MustGetLogger("slack")

func PrepNotification(ticket zendesk.ActiveTicket, notify int64) (notification string) {
	log.Debug("Preparing notification for", ticket.ID)
	c := config.LoadConfig()
	var t string

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

	switch notify {
	case 1:
		t = "15 minutes"
		n = fmt.Sprintf("@sup SLA for %s ticket has less than %s until expiration.", p, t)
	case 2:
		t = "30 minutes"
		n = fmt.Sprintf("@sup SLA for %s ticket has less than %s until expiration.", p, t)
	case 3:
		t = "1 hour"
		n = fmt.Sprintf("@sup SLA for %s ticket has less than %s until expiration.", p, t)
	case 4:
		t = "2 hours"
		n = fmt.Sprintf("SLA for %s ticket has less than %s until expiration.", p, t)
	case 5:
		t = "3 hours"
		n = fmt.Sprintf("SLA for %s ticket has less than %s until expiration.", p, t)
	}

	return n

}

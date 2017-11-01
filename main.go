package main

import (
	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
)

var log = logging.MustGetLogger("slab")

// VERSION lists the version number. Attempts to follow SemVer
// (http://semver.org/)
const VERSION = "0.0.1"

func main() {
	// Start up the logging system
	initLog()
	log.Notice("SLABot by Tyler Conlee")
	log.Noticef("Version: %s", VERSION)

	c := config.LoadConfig()
	// Start timer process. Takes an int as the number of minutes to loop
	RunTimer(c.UpdateFreq.Duration)
}

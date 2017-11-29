package main

import (
	"flag"
	"fmt"
	"os"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
)

var log *logging.Logger

// VERSION lists the version number. Attempts to follow SemVer
// (http://semver.org/)
var (
	Version = "undefined"
)

func main() {
	flagCheck()
	// Start up the logging system
	log = logging.MustGetLogger("slab")
	initLog()
	log.Notice("SLABot by Tyler Conlee")
	log.Noticef("Version: %s", Version)

	c := config.LoadConfig()
	// Start timer process. Takes an int as the number of minutes to loop
	RunTimer(c.UpdateFreq.Duration)
}

func flagCheck() {
	var help bool

	var helpText = "SLAB is a utility designed to integrate Zendesk SLAs with Slack notifications.\nUsage: ./slab [configuration-file-path]"

	flag.BoolVar(&help, "help", false, helpText)

	var version *bool
	version = flag.Bool("version", false, Version)

	flag.Parse()

	if *version {
		fmt.Printf("Version %s\n", (flag.Lookup("version")).Usage)
		os.Exit(0)
	}
	if help {
		fmt.Printf("%s\n", flag.Lookup("help").Usage)
		os.Exit(0)
	}
}

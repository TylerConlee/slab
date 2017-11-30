package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/slack"
)

var log *logging.Logger
var c config.Config

// VERSION lists the version number. On build, uses the git hash as a version ID
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

	c = config.LoadConfig()
	// Start timer process. Takes an int as the number of minutes to loop

	termChan := make(chan os.Signal, 1)
	s := startServer()
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:

		case <-termChan:
			shutdown(ticker, s)
		}
	}

}

func startServer() *slack.Server {
	s := &slack.Server{
		Info: &slack.ServerInfo{
			Server:  "CircleCI-Support",
			Version: Version,
		},
		Uptime: time.Now(),
	}
	go func() {
		s.StartServer()
		RunTimer(c.UpdateFreq.Duration)
	}()
	return s
}

func shutdown(ticker *time.Ticker, s *slack.Server) {

	if ticker != nil {
		ticker.Stop()
	}

	log.Info("Shutdown complete.")
	os.Exit(0)
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

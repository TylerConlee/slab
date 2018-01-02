// SLAB is a full support bot for integrating a Zendesk workflow and Slack.
// A Slack app must be set up for this to run properly.
// (https://api.slack.com/apps)
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	logging "github.com/op/go-logging"
	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/server"
	"github.com/tylerconlee/slab/zendesk"
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
	initLog()
	log = logging.MustGetLogger("slab")

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

// startServer initializes the metadata for the status page, starts the timer
// for the SLA breach monitor loop, and starts an HTTP server for running Slacks
// real time messaging monitoring API.
func startServer() *server.Server {
	s := &server.Server{
		Info: &server.ServerInfo{
			Server:  c.Metadata.Server,
			Version: Version,
			Port:    c.Port,
		},
		Uptime: time.Now(),
	}
	go func() {
		zendesk.RunTimer(c.UpdateFreq.Duration)
	}()
	go func() {
		s.StartServer()

	}()
	return s
}

// shutdown stops the ticker and gracefully shuts down the server.
func shutdown(ticker *time.Ticker, s *server.Server) {

	if ticker != nil {
		ticker.Stop()
	}

	log.Info("Shutdown complete.")
	os.Exit(0)
}

// flagCheck parses any flags that are passed when calling SLAB on the
// command line.
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

// SLAB is a full support bot for integrating a Zendesk workflow and Slack.
// A Slack app must be set up for this to run properly.
// (https://api.slack.com/apps)
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tylerconlee/slab/config"
	l "github.com/tylerconlee/slab/log"
)

// VERSION lists the version number. On build, uses the git hash as a version ID
var (
	Version  = "undefined"
	log      = l.Log
	c        config.Config
	slackKey string
)

func main() {
	flagCheck()
	keyCheck()
	// Start up the logging system
	c = config.LoadConfig()
	log.SetLogLevel(c.LogLevel)
	log.Info("SLABot by Tyler Conlee", map[string]interface{}{
		"module": "main",
	})

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
func startServer() *Server {
	s := &Server{
		Info: &ServerInfo{
			Server:  c.Metadata.Server,
			Version: Version,
			Port:    c.Port,
		},
		Uptime: time.Now(),
	}
	go func() {
		RunTimer(c.UpdateFreq.Duration)
	}()
	go func() {
		s.StartServer()

	}()
	return s
}

// shutdown stops the ticker and gracefully shuts down the server.
func shutdown(ticker *time.Ticker, s *Server) {

	if ticker != nil {
		ticker.Stop()
	}

	log.Info("Shutdown complete.", map[string]interface{}{
		"module": "main",
	})
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

// keyCheck looks for any passed arguments. If there are none, an error is
// displayed and the app exits.
func keyCheck() bool {
	// Check to see if there are a sufficient number of arguments passed
	if len(os.Args) > 1 {
		// Check to see if key starts with `xoxb`. All Slack keys start with
		// `xoxb`, so it's a simple validation test
		if strings.HasPrefix(os.Args[1], "xoxb") {
			slackKey = os.Args[1]
			return true
		}

		log.Fatal(map[string]interface{}{
			"module": "main",
			"error":  "Key provided does not appear to be a valid Slack API key",
		})
		return false
	}
	log.Fatal(map[string]interface{}{
		"module": "main",
		"error":  "Slack key not provided. Slack key must be present to run Slab",
	})
	return false
}

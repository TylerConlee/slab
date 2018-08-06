// SLAB is a full support bot for integrating a Zendesk workflow and Slack.
// A Slack app must be set up for this to run properly.
// (https://api.slack.com/apps)
package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/tylerconlee/slab/config"
	"github.com/tylerconlee/slab/datastore"
	l "github.com/tylerconlee/slab/log"
)

// VERSION lists the version number. On build, uses the git hash as a version ID
var (
	Version    = "undefined"
	log        = l.Log
	c          config.Config
	slackKey   string
	serverPort int
	freq       time.Duration
)

func main() {
	keyCheck()
	log.Info("SLABot by Tyler Conlee", map[string]interface{}{
		"module": "main",
	})

	// Start up the logging system
	c = config.LoadConfig()

	// Start timer process. Takes an int as the number of minutes to loop

	datastore.RedisConnect(serverPort)
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
			Port:    serverPort,
		},
		Uptime: time.Now(),
	}
	go func() {
		RunTimer(freq)
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

// keyCheck looks for any passed arguments. If there are none, an error is
// displayed and the app exits.
func keyCheck() bool {
	// use an int to check if both port and key are valid
	valid := 0
	// Check to see that the proper flags have been passed - port and Slack key
	k := flag.String("key", "APIKey", "a valid Slack API key")
	t := flag.String("time", "2m", "the amount of time between Zendesk checks")
	p := flag.Int("port", 8090, "the port Slab will listen on")

	flag.Parse()

	key := *k
	port := *p
	// Check to see if key starts with `xoxb`. All Slack keys start with
	// `xoxb`, so it's a simple validation test
	if strings.HasPrefix(key, "xoxb") {
		slackKey = key
		valid++
	} else {
		log.Fatal(map[string]interface{}{
			"module": "main",
			"error":  "Key provided does not appear to be a valid Slack API key",
			"key":    key,
		})
	}
	loop, err := time.ParseDuration(*t)
	if err != nil {
		log.Error("Invalid loop time provided", map[string]interface{}{
			"module": "main",
			"error":  err,
			"time":   *t,
		})
		freq = time.Duration(10 * time.Minute)
	} else {
		freq = loop
	}

	if port < 65534 && port > 1 {
		serverPort = port
		valid++
	}
	if valid == 2 {
		return true
	}
	return false
}

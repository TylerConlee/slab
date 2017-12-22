// SLAB is a full support bot for integrating a Zendesk workflow and Slack.
// A Slack app must be set up for this to run properly.
// (https://api.slack.com/apps)
package main

import (
	"os"

	"github.com/op/go-logging"
)

// format is the log format used for most log messages.
var format = logging.MustStringFormatter(
	`%{color}%{time:2006-01-02T15:04:05.000} %{module} ▶ %{level:.4s} %{color:reset} %{message}`,
)

// errorFormatStr is a special log format reserved for logging any major errors.
var errorFormatStr = logging.MustStringFormatter(
	`%{color} %{longpkg} %{shortfunc} ▶ %{shortfile}`,
)

// initLog starts an instance of go-logging and formats it to show the
// line numbers for ERROR and CRITICAL level log messages
func initLog() {
	// Create a new backend for logs
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	// Create a second one for errors and critical log messages
	errorBackend := logging.NewLogBackend(os.Stderr, "", 0)
	// Format the first backend using the overall formatter
	formatter := logging.NewBackendFormatter(backend, format)
	// Format the errors using the error formatter
	errorFormat := logging.NewBackendFormatter(errorBackend, errorFormatStr)
	// Make the error backend leveled
	errorLogBackend := logging.AddModuleLevel(errorFormat)
	// Only show error and above messages
	errorLogBackend.SetLevel(logging.ERROR, "")
	// Apply the backend changes
	logging.SetBackend(formatter, errorLogBackend)

}

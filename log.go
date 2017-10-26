package main

import (
	"os"

	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("slab")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset} %{message} %{color}`,
)
var errorFormatStr = logging.MustStringFormatter(
	`%{color} %{longpkg} %{shortfunc} ▶ %{shortfile}`,
)

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

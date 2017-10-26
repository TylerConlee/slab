package main

import (
	"github.com/op/go-logging"
	"os"
  )

var log = logging.MustGetLogger("slab")

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{color:reset} %{message} %{color}`,
)
var errorFormatStr = logging.MustStringFormatter(
	`%{color} %{longpkg} %{shortfunc} ▶ %{shortfile}`,
)


func initLog(){
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	errorBackend := logging.NewLogBackend(os.Stderr, "", 0)
	formatter := logging.NewBackendFormatter(backend, format)
  
  
	errorFormat := logging.NewBackendFormatter(errorBackend, errorFormatStr)
	errorLogBackend := logging.AddModuleLevel(errorFormat)
	errorLogBackend.SetLevel(logging.ERROR, "")
  
	logging.SetBackend(formatter, errorLogBackend)
  
  }
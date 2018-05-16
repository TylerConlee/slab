package log

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/doublefree/sumorus"
)

type Logger struct{ e *logrus.Logger }

var Log *Logger
var endpoint string = os.Getenv("SLAB_SUMO_URL")
var host = os.Getenv("SLAB_SUMO_HOST")

func init() {
	Log = &Logger{logrus.New()}
	if endpoint != "" {
		sumoLogicHook := sumorus.NewSumoLogicHook(endpoint, host, logrus.InfoLevel, "tag1", "tag2")
		Log.e.Hooks.Add(sumoLogicHook)
	}
}

func (s *Logger) Fatal(c map[string]interface{}) {
	s.e.WithFields(c).Fatal("Error")
}
func (s *Logger) Error(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Error(msg)
}
func (s *Logger) Warn(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Warn(msg)
}
func (s *Logger) Info(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Info(msg)
}
func (s *Logger) Debug(msg string, c map[string]interface{}) {
	s.e.WithFields(c).Debug(msg)
}
func (s *Logger) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		os.Exit(1)
	}
	logrus.SetLevel(lvl)

}

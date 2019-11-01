package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct{ e *logrus.Logger }

var Log *Logger

func init() {
	Log = &Logger{logrus.New()}
	Log.e.Formatter = &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
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
	Log.e.SetLevel(lvl)

}

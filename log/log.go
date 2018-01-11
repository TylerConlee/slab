package log

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct{ e *logrus.Entry }

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
	if s.DebugEnabled() {
		s.e.WithFields(c).Debug(msg)
	}
}
func (s *Logger) DebugEnabled() bool {
	lvl := logrus.GetLevel()
	return lvl >= logrus.DebugLevel
}
func (s *Logger) SetLogLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		os.Exit(1)
	}
	logrus.SetLevel(lvl)
}

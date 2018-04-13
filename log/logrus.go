package log

import (
	"github.com/sirupsen/logrus"
)

// GetLogLevel return Logrus Level
func GetLogLevel(level string) logrus.Level {
	switch level {
	case "error", "err":
		return logrus.ErrorLevel
	case "debug":
		return logrus.DebugLevel
	case "warning", "war":
		return logrus.WarnLevel
	case "info":
		return logrus.InfoLevel
	}
	return logrus.ErrorLevel
}

package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Retrieve the log level based on
func GetLogLevel() logrus.Level {
	level := os.Getenv(EnvLogLevel)
	if len(level) == 0 {
		return logrus.InfoLevel
	}
	for _, a_level := range logrus.AllLevels {
		if level == a_level.String() {
			return a_level
		}
	}
	return logrus.InfoLevel
}
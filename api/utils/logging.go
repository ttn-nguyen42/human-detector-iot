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
	for _, aLevel := range logrus.AllLevels {
		if level == aLevel.String() {
			return aLevel
		}
	}
	return logrus.InfoLevel
}

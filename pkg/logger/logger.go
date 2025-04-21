package logger

import (
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogger(level string) *logrus.Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	log.SetOutput(os.Stdout)

	level = strings.TrimSpace(strings.ToLower(level))
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	log.SetLevel(logLevel)
	return log
}

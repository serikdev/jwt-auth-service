package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger = *logrus.Entry

func NewLogger() Logger {
	log := logrus.New()

	log.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyTime:  "timestamp",
		},

		TimestampFormat: time.RFC3339Nano,
	})

	log.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)
	logger := log.WithFields(logrus.Fields{
		"service": os.Getenv("SERVICE_NAME"),
		"version": os.Getenv("APP_VERSION"),
	})

	return logger

}

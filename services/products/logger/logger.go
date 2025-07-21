package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

// SetupLogger configures the logger with structured format
func SetupLogger() *logrus.Logger {
	logger := logrus.New()

	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Warn("Invalid log level, using info")
		level = logrus.InfoLevel
	}

	logger.SetLevel(level)

	logger.SetOutput(os.Stdout)

	logger.WithFields(logrus.Fields{
		"service": "products-service",
		"version": "1.0.0",
	}).Info("Logger initialized")

	return logger
}

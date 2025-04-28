package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"lowerthirdsapi/internal/helpers"
	"os"
)

// defaultLogLevel defines the default log level
const defaultLogLevel = logrus.DebugLevel // TODO: switch to for deploymnet: logrus.InfoLevel

// New creates a new Logrus logger
func New() *logrus.Entry {
	// Create the logger
	logger := logrus.New()

	// Determine the API logger level, using `logrus.InfoLevel` as the default
	logLevel := defaultLogLevel
	logLevelString := os.Getenv("GPSI_API_LOG_LEVEL")
	if logLevelString != "" {
		reqLogLevel, err := logrus.ParseLevel(logLevelString)
		if err == nil {
			logLevel = reqLogLevel
		}
	}
	logger.SetLevel(logLevel)

	// Set the logger format
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Get the hostname and environment to add as default fields
	hostname, _ := os.Hostname()
	environment := os.Getenv("ENVIRONMENT")

	// Return the logger with the default fields
	return logger.WithFields(logrus.Fields{
		"host":        hostname,
		"environment": environment,
	})
}

func WithContext(ctx context.Context, log *logrus.Entry) *logrus.Entry {
	if userID, ok := ctx.Value(helpers.UserIDKey).(string); ok {
		log = log.WithField("userID", userID)
	}

	return log
}

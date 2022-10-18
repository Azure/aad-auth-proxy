package telemetry

import (
	"aad-auth-proxy/contracts"
	"os"

	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
)

var (
	_, isDebugMode = os.LookupEnv("TRACE_LOGGING")
)

type Logger struct {
	logger *logrus.Entry
}

// The formatter for message templates: "<timestamp> <level> <message>"
type guestAgentFormatter struct {
	logrus.Formatter
}

//
// Get Logger singleton
//
func NewLogger() contracts.ILogger {
	logger := setupLogrusLogger(isDebugMode)

	loggerEntry := logrus.NewEntry(logger)
	return &Logger{
		logger: loggerEntry,
	}
}

//
// Helper function to set up the configurations of logrus.
//
func setupLogrusLogger(isDebugModeParam bool) *logrus.Logger {
	logrusLogger := logrus.New()
	format := "%time% [%lvl%] %msg% %" + "%\n"

	//set up the logger
	logrusLogger.SetFormatter(&guestAgentFormatter{
		&easy.Formatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
			LogFormat:       format,
		},
	})

	if isDebugModeParam {
		logrusLogger.SetLevel(logrus.TraceLevel)
	} else {
		logrusLogger.SetLevel(logrus.InfoLevel)
	}

	return logrusLogger
}

//
// Logger print error level message.
//
func (l *Logger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

//
// Logger print info level message.
//
func (l *Logger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

//
// Logger print warn level message.
//
func (l *Logger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

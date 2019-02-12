package utils

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	LogFormatText = "text"
	LogFormatJSON = "json"

	ContextKeyLogger = "logger"
)

var (
	logLevel                   = logrus.DebugLevel
	logFormat logrus.Formatter = &logrus.TextFormatter{}
	logOut    io.Writer
)

func InitLogger(ll, lf string) {
	logLevel = parseLogrusLevel(ll)
	logrus.SetLevel(logLevel)

	logFormat = parseLogrusFormat(lf)
	logrus.SetFormatter(logFormat)

	logOut = os.Stdout
	logrus.SetOutput(logOut)
}

func GetLoggerFromCtx(c *gin.Context) *logrus.Entry {
	if logger, ok := c.Get(ContextKeyLogger); ok {
		logEntry, assertionOk := logger.(*logrus.Entry)
		if assertionOk {
			return logEntry
		}
	}
	return logrus.NewEntry(GetLogger())
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()
	logger.Formatter = logFormat
	logger.Level = logLevel
	logger.Out = logOut
	return logger
}

func parseLogrusLevel(logLevelStr string) logrus.Level {
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logrus.WithError(err).Errorf("error while parsing log level. %v is set as default.", logLevel)
		logLevel = logrus.DebugLevel
	}
	return logLevel
}

func parseLogrusFormat(logFormatStr string) logrus.Formatter {
	var formatter logrus.Formatter
	switch logFormatStr {
	case LogFormatText:
		formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	case LogFormatJSON:
		formatter = &logrus.JSONFormatter{}
	default:
		logrus.Errorf("error while parsing log format. %v is set as default.", formatter)
		formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	}
	return formatter
}

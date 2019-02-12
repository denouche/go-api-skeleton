package utils

import "github.com/sirupsen/logrus"

const (
	LogFormatText = "text"
	LogFormatJSON = "json"
)

func ParseLogrusLevel(logLevelStr string) logrus.Level {
	logLevel, err := logrus.ParseLevel(logLevelStr)
	if err != nil {
		logrus.WithError(err).Errorf("error while parsing log level. %v is set as default.", logLevel)
		logLevel = logrus.DebugLevel
	}
	return logLevel
}

func ParseLogrusFormat(logFormatStr string) logrus.Formatter {
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

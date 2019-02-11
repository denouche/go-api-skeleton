package utils

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	ContextKeyLogger = "logger"
)

func GetLoggerForEnvironment(environment string) *zap.Logger {
	var logger *zap.Logger

	switch environment {
	case "development":
		fallthrough
	default:
		logger, _ = zap.NewDevelopment()
	case "production":
		logger, _ = zap.NewProduction()
	}
	return logger
}

func GetLogger(c *gin.Context) *zap.SugaredLogger {
	if c == nil {
		logger, _ := zap.NewProduction()
		return logger.Sugar()
	}

	if l, ok := c.Get(ContextKeyLogger); ok {
		sugar, assertionOk := l.(*zap.SugaredLogger)
		if assertionOk {
			return sugar
		}
	}

	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

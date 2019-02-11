package middlewares

import (
	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"time"
)

func GetHTTPLoggerMiddleware(environment string) gin.HandlerFunc {
	logger := utils.GetLoggerForEnvironment(environment)
	return ginzap.Ginzap(logger, time.RFC3339, true)
}

func GetLoggerMiddleware(environment string) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := utils.GetLoggerForEnvironment(environment)
		defer logger.Sync()

		sugar := logger.Sugar()
		c.Set(utils.ContextKeyLogger, sugar)

		c.Next()
	}
}

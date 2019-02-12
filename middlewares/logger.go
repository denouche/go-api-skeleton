package middlewares

import (
	"math/rand"
	"time"

	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func GetLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationID := c.Request.Header.Get(utils.HeaderNameCorrelationID)
		if correlationID == "" {
			correlationID = randStringBytesMaskImprSrc(30)
			c.Writer.Header().Set(utils.HeaderNameCorrelationID, correlationID)
		}

		logger := utils.GetLogger()
		logEntry := logger.WithField(utils.HeaderNameCorrelationID, correlationID)

		c.Set(utils.ContextKeyLogger, logEntry)
	}
}

func GetHTTPLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		utils.GetLoggerFromCtx(c).
			WithField("method", c.Request.Method).
			WithField("url", c.Request.RequestURI).
			WithField("from", c.ClientIP()).
			Info("start handling HTTP request")

		c.Next()
		d := time.Since(start)

		utils.GetLoggerFromCtx(c).
			WithField("status", c.Writer.Status()).
			WithField("duration", d.String()).
			Info("end handling HTTP request")
	}
}

package middlewares

import (
	"net/http"

	"github.com/denouche/go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

func GetCORSMiddlewareForOthersHTTPMethods() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if v := c.Request.Header.Get(httputils.HeaderNameOrigin); v != "" {
			c.Writer.Header().Set(httputils.HeaderNameAccessControlAllowOrigin, v)
		}
		c.Next()
	}
}

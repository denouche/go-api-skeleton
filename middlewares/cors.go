package middlewares

import (
	"net/http"

	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func CORSMiddlewareForOthersHTTPMethods() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		v := c.Request.Header.Get(utils.HeaderNameOrigin)
		if v != "" {
			c.Writer.Header().Set(utils.HeaderNameAccessControlAllowOrigin, v)
		}
		c.Next()
	}
}

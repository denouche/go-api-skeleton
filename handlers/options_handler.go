package handlers

import (
	"net/http"
	"strings"

	"github.com/denouche/go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

func (hc *Context) GetOptionsHandler(allowedHeaders []string, allowedMethods ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		headerOrigin := c.Request.Header.Get(httputils.HeaderNameOrigin)
		if headerOrigin != "" {
			optionsInArray := false
			for _, method := range allowedMethods {
				if method == http.MethodOptions {
					optionsInArray = true
					break
				}
			}
			if !optionsInArray {
				allowedMethods = append(allowedMethods, http.MethodOptions)
			}

			c.Writer.Header().Set(httputils.HeaderNameAccessControlAllowOrigin, headerOrigin)
			if headerOrigin == "*" {
				c.Writer.Header().Set(httputils.HeaderNameAccessControlAllowCredentials, "false")
			} else {
				c.Writer.Header().Set(httputils.HeaderNameAccessControlAllowCredentials, "true")
			}

			c.Writer.Header().Set(httputils.HeaderNameAccessControlAllowMethods, strings.Join(allowedMethods, ","))
			c.Writer.Header().Set(httputils.HeaderNameAccessControlAllowHeaders, strings.Join(allowedHeaders, ","))

			c.Status(http.StatusOK)
			return
		}
		c.Status(http.StatusMethodNotAllowed)
	}
}

package handlers

import (
	"net/http"
	"strings"

	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func (hc *handlersContext) GetOptionsHandler(allowedHeaders []string, allowedMethods ...string) func(*gin.Context) {
	return func(c *gin.Context) {
		headerOrigin := c.Request.Header.Get(utils.HeaderNameOrigin)
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

			c.Writer.Header().Set(utils.HeaderNameAccessControlAllowOrigin, headerOrigin)
			if headerOrigin == "*" {
				c.Writer.Header().Set(utils.HeaderNameAccessControlAllowCredentials, "false")
			} else {
				c.Writer.Header().Set(utils.HeaderNameAccessControlAllowCredentials, "true")
			}

			c.Writer.Header().Set(utils.HeaderNameAccessControlAllowMethods, strings.Join(allowedMethods, ","))
			c.Writer.Header().Set(utils.HeaderNameAccessControlAllowHeaders, strings.Join(allowedHeaders, ","))

			c.Status(http.StatusOK)
			return
		}
		c.Status(http.StatusMethodNotAllowed)
	}
}

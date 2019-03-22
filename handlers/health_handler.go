package handlers

import (
	"net/http"

	"github.com/denouche/go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

func (hc *Context) GetHealth(c *gin.Context) {
	conf := map[string]string{
		"ApplicationName":      ApplicationName,
		"ApplicationVersion":   ApplicationVersion,
		"ApplicationGitHash":   ApplicationGitHash,
		"ApplicationBuildDate": ApplicationBuildDate,
	}
	httputils.JSON(c.Writer, http.StatusOK, conf)
}

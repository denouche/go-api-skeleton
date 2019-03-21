package handlers

import (
	"net/http"

	"github.com/denouche/go-api-skeleton/utils/httputils"
	"github.com/gin-gonic/gin"
)

func (hc *Context) GetHealth(c *gin.Context) {
	httputils.JSON(c.Writer, http.StatusNoContent, nil)
}

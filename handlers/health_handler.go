package handlers

import (
	"net/http"

	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func (hc *handlersContext) GetHealth(c *gin.Context) {
	utils.JSON(c.Writer, http.StatusNoContent, nil)
}

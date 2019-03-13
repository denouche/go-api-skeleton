package handlers

import (
	"net/http"

	"github.com/denouche/go-api-skeleton/utils"
	"github.com/gin-gonic/gin"
)

func (hc *handlersContext) GetOpenAPISchema(c *gin.Context) {
	utils.YAML(c.Writer, http.StatusOK, utils.OpenAPISchema)
}

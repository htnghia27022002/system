package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/common/response"
)

// AppHandler serves public, non-sensitive facts about this deployment so a
// separately deployed FE can adapt (e.g. hide UI for modules that are off).
type AppHandler struct {
	modules []string
}

func NewAppHandler(modules []string) *AppHandler {
	if modules == nil {
		modules = []string{}
	}
	return &AppHandler{modules: modules}
}

// Info returns the enabled optional modules (MODULES env / registry defaults).
func (h *AppHandler) Info(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{"modules": h.modules})
}

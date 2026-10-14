package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/common/response"
	"be/internal/settings"
)

// AppHandler serves public, non-sensitive facts about this deployment so a
// separately deployed FE can adapt (e.g. hide UI for modules that are off, or
// the sign-up link when registration is closed).
type AppHandler struct {
	modules  []string
	settings *settings.Service
}

func NewAppHandler(modules []string, settingsSvc *settings.Service) *AppHandler {
	if modules == nil {
		modules = []string{}
	}
	return &AppHandler{modules: modules, settings: settingsSvc}
}

// Info returns the enabled optional modules and the public settings.
func (h *AppHandler) Info(c *gin.Context) {
	public := map[string]any{}
	if h.settings != nil {
		values, err := h.settings.Public(c.Request.Context())
		if err != nil {
			response.HandleError(c, err)
			return
		}
		public = values
	}
	response.JSON(c, http.StatusOK, gin.H{"modules": h.modules, "settings": public})
}

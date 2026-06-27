package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/internal/common/response"
	permissionsvc "be/internal/services/permission"
)

type PermissionHandler struct {
	svc *permissionsvc.Service
}

func NewPermissionHandler(svc *permissionsvc.Service) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, permissions)
}

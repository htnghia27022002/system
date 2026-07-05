package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/internal/common/response"
	permissiondto "be/internal/dto/permission"
	permissionsvc "be/internal/services/permission"
)

type PermissionHandler struct {
	svc *permissionsvc.Service
}

func NewPermissionHandler(svc *permissionsvc.Service) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

func (h *PermissionHandler) ListAll(c *gin.Context) {
	permissions, err := h.svc.ListAll(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, permissions)
}

func (h *PermissionHandler) List(c *gin.Context) {
	var query permissiondto.ListPermissionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	permissions, total, page, pageSize, err := h.svc.List(c.Request.Context(), query)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, permissiondto.PaginatedPermissionsResponse{
		Items:    permissions,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

package handlers

import (
	"github.com/gin-gonic/gin"

	"be/common/httpx"
	permissiondto "be/internal/dto/permission"
	permissionsvc "be/internal/services/permission"
	"be/pkg/query"
)

type PermissionHandler struct {
	svc *permissionsvc.Service
}

func NewPermissionHandler(svc *permissionsvc.Service) *PermissionHandler {
	return &PermissionHandler{svc: svc}
}

func (h *PermissionHandler) ListAll(c *gin.Context) {
	permissions, err := h.svc.ListAll(c.Request.Context())
	httpx.OK(c, permissions, err)
}

func (h *PermissionHandler) List(c *gin.Context) {
	var form permissiondto.ListPermissionsQuery
	if !httpx.BindQuery(c, &form) {
		return
	}

	permissions, total, page, pageSize, err := h.svc.List(c.Request.Context(), form)
	httpx.OK(c, query.NewPage(permissions, total, page, pageSize), err)
}

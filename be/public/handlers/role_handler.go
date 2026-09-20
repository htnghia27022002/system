package handlers

import (
	"github.com/gin-gonic/gin"

	"be/common/httpx"
	roledto "be/internal/dto/role"
	rolesvc "be/internal/services/role"
	"be/pkg/query"
)

type RoleHandler struct {
	svc *rolesvc.Service
}

func NewRoleHandler(svc *rolesvc.Service) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req roledto.CreateRoleRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	role, err := h.svc.Create(c.Request.Context(), req)
	httpx.Created(c, role, err)
}

func (h *RoleHandler) ListAll(c *gin.Context) {
	roles, err := h.svc.ListAll(c.Request.Context())
	httpx.OK(c, roles, err)
}

func (h *RoleHandler) List(c *gin.Context) {
	var form roledto.ListRolesQuery
	if !httpx.BindQuery(c, &form) {
		return
	}

	roles, total, page, pageSize, err := h.svc.List(c.Request.Context(), form)
	httpx.OK(c, query.NewPage(roles, total, page, pageSize), err)
}

func (h *RoleHandler) Get(c *gin.Context) {
	role, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	httpx.OK(c, role, err)
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req roledto.UpdateRoleRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	role, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	httpx.OK(c, role, err)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	httpx.NoContent(c, h.svc.Delete(c.Request.Context(), c.Param("id")))
}

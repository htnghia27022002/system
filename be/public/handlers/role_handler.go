package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/internal/common/response"
	roledto "be/internal/dto/role"
	rolesvc "be/internal/services/role"
)

type RoleHandler struct {
	svc *rolesvc.Service
}

func NewRoleHandler(svc *rolesvc.Service) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req roledto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, role)
}

func (h *RoleHandler) List(c *gin.Context) {
	roles, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, roles)
}

func (h *RoleHandler) Get(c *gin.Context) {
	role, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	var req roledto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	role, err := h.svc.Update(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id")); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

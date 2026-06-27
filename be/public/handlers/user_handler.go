package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/internal/common/response"
	userdto "be/internal/dto/user"
	"be/internal/middleware"
	usersvc "be/internal/services/user"
)

type UserHandler struct {
	svc *usersvc.Service
}

func NewUserHandler(svc *usersvc.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req userdto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, usersvc.ToResponse(user))
}

func (h *UserHandler) Get(c *gin.Context) {
	user, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, usersvc.ToResponse(user))
}

func (h *UserHandler) List(c *gin.Context) {
	var query userdto.ListUsersQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	users, total, page, pageSize, err := h.svc.List(c.Request.Context(), query)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	items := make([]userdto.UserResponse, 0, len(users))
	for i := range users {
		items = append(items, usersvc.ToResponse(&users[i]))
	}
	response.JSON(c, http.StatusOK, userdto.PaginatedUsersResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	var req userdto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	user, err := h.svc.Update(c.Request.Context(), c.Param("id"), req, middleware.GetUserID(c))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, usersvc.ToResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Request.Context(), c.Param("id"), middleware.GetUserID(c)); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

package handlers

import (
	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/response"
	userdto "be/internal/dto/user"
	"be/internal/middleware"
	usersvc "be/internal/services/user"
	"be/pkg/query"
)

type UserHandler struct {
	svc *usersvc.Service
}

func NewUserHandler(svc *usersvc.Service) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Create(c *gin.Context) {
	var req userdto.CreateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	user, err := h.svc.Create(c.Request.Context(), req, actorFrom(c))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	item, err := h.svc.ResponseForUser(c.Request.Context(), user)
	httpx.Created(c, item, err)
}

func (h *UserHandler) Get(c *gin.Context) {
	user, err := h.svc.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	item, err := h.svc.ResponseForUser(c.Request.Context(), user)
	httpx.OK(c, item, err)
}

func (h *UserHandler) List(c *gin.Context) {
	var form userdto.ListUsersQuery
	if !httpx.BindQuery(c, &form) {
		return
	}

	users, total, page, pageSize, err := h.svc.List(c.Request.Context(), form)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	items, err := h.svc.ResponsesForUsers(c.Request.Context(), users)
	httpx.OK(c, query.NewPage(items, total, page, pageSize), err)
}

func (h *UserHandler) Update(c *gin.Context) {
	var req userdto.UpdateUserRequest
	if !httpx.BindJSON(c, &req) {
		return
	}

	user, err := h.svc.Update(c.Request.Context(), c.Param("id"), req, actorFrom(c))
	if err != nil {
		response.HandleError(c, err)
		return
	}
	item, err := h.svc.ResponseForUser(c.Request.Context(), user)
	httpx.OK(c, item, err)
}

func (h *UserHandler) UploadAvatar(c *gin.Context) {
	file, header, ok := httpx.FormFile(c, "avatar file is required (field: file or avatar)", "file", "avatar")
	if !ok {
		return
	}
	defer file.Close()

	user, err := h.svc.UploadAvatar(c.Request.Context(), c.Param("id"), file, header)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	item, err := h.svc.ResponseForUser(c.Request.Context(), user)
	httpx.OK(c, item, err)
}

func (h *UserHandler) Delete(c *gin.Context) {
	httpx.NoContent(c, h.svc.Delete(c.Request.Context(), c.Param("id"), actorFrom(c)))
}

func actorFrom(c *gin.Context) usersvc.Actor {
	return usersvc.Actor{
		ID:         middleware.GetUserID(c),
		SuperAdmin: middleware.IsSuperAdmin(c),
	}
}

package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/common/httpx"
	"be/common/response"
	searchdto "be/internal/dto/search"
	"be/internal/middleware"
	searchsvc "be/internal/services/search"
)

type SearchHandler struct {
	search *searchsvc.Service
}

func NewSearchHandler(search *searchsvc.Service) *SearchHandler {
	return &SearchHandler{search: search}
}

func (h *SearchHandler) Search(c *gin.Context) {
	var form searchdto.SearchQuery
	if !httpx.BindQuery(c, &form) {
		return
	}

	permissions, ok := middleware.GetPermissions(c)
	if !ok {
		response.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	result, err := h.search.Search(c.Request.Context(), form, permissions)
	httpx.OK(c, result, err)
}

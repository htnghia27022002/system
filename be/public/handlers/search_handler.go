package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"be/internal/common/response"
	searchdto "be/internal/dto/search"
	"be/internal/middleware"
	searchsvc "be/internal/services/search"
)

type SearchHandler struct {
	search    *searchsvc.Service
	processor *searchsvc.IndexProcessor
	outbox    *searchsvc.OutboxService
}

func NewSearchHandler(
	search *searchsvc.Service,
	processor *searchsvc.IndexProcessor,
	outbox *searchsvc.OutboxService,
) *SearchHandler {
	return &SearchHandler{search: search, processor: processor, outbox: outbox}
}

func (h *SearchHandler) Search(c *gin.Context) {
	var query searchdto.SearchQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	permissions, ok := middleware.GetPermissions(c)
	if !ok {
		response.Error(c, http.StatusForbidden, "forbidden")
		return
	}

	result, err := h.search.Search(c.Request.Context(), query, permissions)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, result)
}

func (h *SearchHandler) Reindex(c *gin.Context) {
	var req searchdto.ReindexRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	entityType := req.EntityType
	if entityType == "" {
		entityType = "all"
	}

	count, err := h.processor.Reindex(c.Request.Context(), entityType)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, searchdto.ReindexResponse{
		EntityType: entityType,
		Indexed:    count,
		Status:     "completed",
	})
}

func (h *SearchHandler) OutboxStats(c *gin.Context) {
	stats, err := h.processor.OutboxStats(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, searchdto.OutboxStatsResponse{
		PendingCount:            stats.PendingCount,
		FailedCount:             stats.FailedCount,
		OldestPendingAgeSeconds: stats.OldestPendingAgeSeconds,
	})
}

func (h *SearchHandler) ReplayOutbox(c *gin.Context) {
	var req searchdto.ReplayOutboxRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if req.ID == "" {
		response.Error(c, http.StatusBadRequest, "id is required")
		return
	}

	if err := h.outbox.Replay(c.Request.Context(), req.ID); err != nil {
		response.HandleError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

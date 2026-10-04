package searchdto

import (
	"time"

	"be/pkg/query"
)

type SearchQuery struct {
	query.PageParams
	Q     string `form:"q" binding:"omitempty,max=200"`
	Types string `form:"types" binding:"omitempty,max=100"`
}

type SearchHitResponse struct {
	EntityType string            `json:"entityType"`
	EntityID   string            `json:"entityId"`
	Title      string            `json:"title"`
	Snippet    string            `json:"snippet,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	UpdatedAt  time.Time         `json:"updatedAt"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type SearchResponse struct {
	Hits       []SearchHitResponse `json:"hits"`
	Pagination PaginationResponse  `json:"pagination"`
}

package searchsvc

import (
	"context"
	"strings"

	"be/common/rbac"
	searchdto "be/internal/dto/search"
	"be/internal/repository/interfaces"
	listquery "be/pkg/query"
)

// Service is admin global search over users, roles, and permissions, backed by Postgres.
type Service struct {
	repo interfaces.SearchRepository
}

func NewService(repo interfaces.SearchRepository) *Service {
	return &Service{repo: repo}
}

// entityPermission maps each searchable entity type to the RBAC resource that gates it.
var entityPermission = []struct {
	entityType string
	resource   string
}{
	{interfaces.SearchEntityUser, "users"},
	{interfaces.SearchEntityRole, "roles"},
	{interfaces.SearchEntityPermission, "permissions"},
}

func (s *Service) Search(ctx context.Context, form searchdto.SearchQuery, permissions []string) (*searchdto.SearchResponse, error) {
	q := strings.TrimSpace(form.Q)
	page, pageSize := listquery.NormalizePage(form.Page, form.PageSize)
	if q == "" {
		return emptyResponse(page, pageSize), nil
	}

	types := filterTypesByPermission(parseTypes(form.Types), permissions)
	if len(types) == 0 {
		return emptyResponse(page, pageSize), nil
	}

	hits, total, err := s.repo.Search(ctx, q, types, listquery.Offset(page, pageSize), pageSize)
	if err != nil {
		return nil, err
	}

	out := make([]searchdto.SearchHitResponse, 0, len(hits))
	for _, hit := range hits {
		out = append(out, searchdto.SearchHitResponse{
			EntityType: hit.EntityType,
			EntityID:   hit.EntityID,
			Title:      hit.Title,
			Snippet:    hit.Snippet,
			Metadata:   hit.Metadata,
			UpdatedAt:  hit.UpdatedAt,
		})
	}

	return &searchdto.SearchResponse{
		Hits: out,
		Pagination: searchdto.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: listquery.TotalPages(total, pageSize),
		},
	}, nil
}

// filterTypesByPermission keeps only entity types the caller may view.
// An empty request means "all types the caller may view".
func filterTypesByPermission(requested []string, permissions []string) []string {
	out := make([]string, 0, len(entityPermission))
	for _, item := range entityPermission {
		if len(requested) > 0 && !contains(requested, item.entityType) {
			continue
		}
		if rbac.Allowed(permissions, rbac.Key(item.resource, rbac.ActionView)) {
			out = append(out, item.entityType)
		}
	}
	return out
}

func parseTypes(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func emptyResponse(page, pageSize int) *searchdto.SearchResponse {
	return &searchdto.SearchResponse{
		Hits: []searchdto.SearchHitResponse{},
		Pagination: searchdto.PaginationResponse{
			Page:     page,
			PageSize: pageSize,
		},
	}
}

package searchsvc

import (
	"context"
	"fmt"
	"strings"

	listquery "be/pkg/query"
	"be/internal/common/rbac"
	searchdto "be/internal/dto/search"
	searchpkg "be/internal/search"
	"be/internal/repository/interfaces"
)

type Service struct {
	client  *searchpkg.Client
	users   interfaces.UserRepository
	builder *DocumentBuilder
}

func NewService(client *searchpkg.Client, users interfaces.UserRepository, builder *DocumentBuilder) *Service {
	return &Service{client: client, users: users, builder: builder}
}

func (s *Service) Search(ctx context.Context, form searchdto.SearchQuery, permissions []string) (*searchdto.SearchResponse, error) {
	q := strings.TrimSpace(form.Q)
	page, pageSize := listquery.NormalizePage(form.Page, form.PageSize)
	if q == "" {
		return emptyResponse(page, pageSize), nil
	}

	allowed := viewPermissions(permissions)
	types := parseTypes(form.Types)
	types = filterTypesByPermission(types, permissions)
	if len(types) == 0 && len(allowed) == 0 {
		return emptyResponse(page, pageSize), nil
	}

	if s.client != nil && s.client.Enabled() {
		if err := s.client.Ping(ctx); err == nil {
			result, err := s.client.Search(ctx, searchpkg.SearchRequest{
				Query:              q,
				EntityTypes:        types,
				AllowedPermissions: allowed,
				From:               listquery.Offset(page, pageSize),
				Size:               pageSize,
			})
			if err == nil {
				return toResponse(result, page, pageSize, false), nil
			}
		}
	}

	return s.fallbackUsers(ctx, q, types, permissions, page, pageSize)
}

func (s *Service) fallbackUsers(
	ctx context.Context,
	q string,
	types []string,
	permissions []string,
	page, pageSize int,
) (*searchdto.SearchResponse, error) {
	resp := emptyResponse(page, pageSize)
	resp.Degraded = true

	if !rbac.Allowed(permissions, rbac.Key("users", rbac.ActionView)) {
		return resp, nil
	}
	if len(types) > 0 && !containsType(types, searchpkg.EntityUser) {
		return resp, nil
	}

	listQ := listquery.New(page, pageSize).
		OrderBy("created_at DESC").
		WhereLikeAny([]string{"email", "full_name"}, q)
	users, total, err := s.users.List(ctx, listQ)
	if err != nil {
		return nil, err
	}

	hits := make([]searchdto.SearchHitResponse, 0, len(users))
	for _, user := range users {
		doc, err := s.builder.Build(ctx, searchpkg.EntityUser, user.ID)
		if err != nil || doc == nil {
			continue
		}
		hits = append(hits, searchdto.SearchHitResponse{
			EntityType: doc.EntityType,
			EntityID:   doc.EntityID,
			Title:      doc.Title,
			Snippet:    doc.SearchableText,
			Metadata:   doc.Metadata,
			UpdatedAt:  doc.UpdatedAt,
		})
	}

	return &searchdto.SearchResponse{
		Hits: hits,
		Pagination: searchdto.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      total,
			TotalPages: listquery.TotalPages(total, pageSize),
		},
		Degraded: true,
	}, nil
}

func viewPermissions(permissions []string) []string {
	keys := []string{
		rbac.Key("users", rbac.ActionView),
		rbac.Key("roles", rbac.ActionView),
		rbac.Key("permissions", rbac.ActionView),
	}
	out := make([]string, 0, len(keys))
	for _, key := range keys {
		if rbac.Allowed(permissions, key) {
			out = append(out, key)
		}
	}
	return out
}

func filterTypesByPermission(types []string, permissions []string) []string {
	if len(types) == 0 {
		out := make([]string, 0, 3)
		if rbac.Allowed(permissions, rbac.Key("users", rbac.ActionView)) {
			out = append(out, searchpkg.EntityUser)
		}
		if rbac.Allowed(permissions, rbac.Key("roles", rbac.ActionView)) {
			out = append(out, searchpkg.EntityRole)
		}
		if rbac.Allowed(permissions, rbac.Key("permissions", rbac.ActionView)) {
			out = append(out, searchpkg.EntityPermission)
		}
		return out
	}

	out := make([]string, 0, len(types))
	for _, t := range types {
		switch t {
		case searchpkg.EntityUser:
			if rbac.Allowed(permissions, rbac.Key("users", rbac.ActionView)) {
				out = append(out, t)
			}
		case searchpkg.EntityRole:
			if rbac.Allowed(permissions, rbac.Key("roles", rbac.ActionView)) {
				out = append(out, t)
			}
		case searchpkg.EntityPermission:
			if rbac.Allowed(permissions, rbac.Key("permissions", rbac.ActionView)) {
				out = append(out, t)
			}
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
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func containsType(types []string, target string) bool {
	for _, t := range types {
		if t == target {
			return true
		}
	}
	return false
}

func emptyResponse(page, pageSize int) *searchdto.SearchResponse {
	return &searchdto.SearchResponse{
		Hits: []searchdto.SearchHitResponse{},
		Pagination: searchdto.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      0,
			TotalPages: 0,
		},
	}
}

func toResponse(result *searchpkg.SearchResult, page, pageSize int, degraded bool) *searchdto.SearchResponse {
	hits := make([]searchdto.SearchHitResponse, 0, len(result.Hits))
	for _, hit := range result.Hits {
		hits = append(hits, searchdto.SearchHitResponse{
			EntityType: hit.EntityType,
			EntityID:   hit.EntityID,
			Title:      hit.Title,
			Snippet:    hit.Snippet,
			Metadata:   hit.Metadata,
			UpdatedAt:  hit.UpdatedAt,
		})
	}

	totalPages := listquery.TotalPages(result.Total, pageSize)

	return &searchdto.SearchResponse{
		Hits: hits,
		Pagination: searchdto.PaginationResponse{
			Page:       page,
			PageSize:   pageSize,
			Total:      result.Total,
			TotalPages: totalPages,
		},
		Degraded: degraded,
	}
}

func (s *Service) Enabled() bool {
	return s.client != nil && s.client.Enabled()
}

func (s *Service) EnsureIndex(ctx context.Context) error {
	if s.client == nil || !s.client.Enabled() {
		return fmt.Errorf("search is disabled")
	}
	return s.client.EnsureIndex(ctx)
}

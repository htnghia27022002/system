package search_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"be/common/rbac"
	searchdto "be/internal/dto/search"
	"be/internal/repository/interfaces"
	searchsvc "be/internal/services/search"
	"be/pkg/query"
)

type searchCall struct {
	term   string
	types  []string
	offset int
	limit  int
}

type stubSearchRepo struct {
	hits  []interfaces.SearchHit
	total int64
	err   error
	calls []searchCall
}

func (s *stubSearchRepo) Search(_ context.Context, term string, types []string, offset, limit int) ([]interfaces.SearchHit, int64, error) {
	s.calls = append(s.calls, searchCall{term: term, types: types, offset: offset, limit: limit})
	return s.hits, s.total, s.err
}

func viewKeys(resources ...string) []string {
	keys := make([]string, 0, len(resources))
	for _, r := range resources {
		keys = append(keys, rbac.Key(r, rbac.ActionView))
	}
	return keys
}

func TestSearchEmptyQueryReturnsEmptyWithoutQuerying(t *testing.T) {
	t.Parallel()
	repo := &stubSearchRepo{}
	svc := searchsvc.NewService(repo)

	resp, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "   "}, viewKeys("users"))
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(resp.Hits) != 0 || len(repo.calls) != 0 {
		t.Fatalf("expected no hits and no repo call, got hits=%d calls=%d", len(resp.Hits), len(repo.calls))
	}
}

func TestSearchWithoutViewPermissionsSkipsRepo(t *testing.T) {
	t.Parallel()
	repo := &stubSearchRepo{}
	svc := searchsvc.NewService(repo)

	resp, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "admin"}, []string{})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(resp.Hits) != 0 || len(repo.calls) != 0 {
		t.Fatalf("expected empty result without *:view, got hits=%d calls=%d", len(resp.Hits), len(repo.calls))
	}
}

func TestSearchRestrictsTypesToPermittedEntities(t *testing.T) {
	t.Parallel()
	repo := &stubSearchRepo{}
	svc := searchsvc.NewService(repo)

	// Caller asks for users+roles but may only view roles.
	_, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "admin", Types: "user, role"}, viewKeys("roles"))
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(repo.calls) != 1 {
		t.Fatalf("expected 1 repo call, got %d", len(repo.calls))
	}
	got := repo.calls[0].types
	if len(got) != 1 || got[0] != interfaces.SearchEntityRole {
		t.Fatalf("expected types [role], got %v", got)
	}
}

func TestSearchDefaultsToAllPermittedTypes(t *testing.T) {
	t.Parallel()
	repo := &stubSearchRepo{}
	svc := searchsvc.NewService(repo)

	_, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "x"}, viewKeys("users", "roles", "permissions"))
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if got := repo.calls[0].types; len(got) != 3 {
		t.Fatalf("expected 3 types, got %v", got)
	}
}

func TestSearchMapsHitsAndPagination(t *testing.T) {
	t.Parallel()
	now := time.Now().UTC()
	repo := &stubSearchRepo{
		hits: []interfaces.SearchHit{{
			EntityType: interfaces.SearchEntityUser,
			EntityID:   "u1",
			Title:      "Admin User",
			Snippet:    "admin@example.com",
			Metadata:   map[string]string{"email": "admin@example.com"},
			UpdatedAt:  now,
		}},
		total: 41,
	}
	svc := searchsvc.NewService(repo)

	resp, err := svc.Search(context.Background(), searchdto.SearchQuery{
		PageParams: query.PageParams{Page: 3, PageSize: 20},
		Q:          " admin ",
	}, viewKeys("users"))
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if call := repo.calls[0]; call.term != "admin" || call.offset != 40 || call.limit != 20 {
		t.Fatalf("unexpected repo call %+v", call)
	}
	if len(resp.Hits) != 1 || resp.Hits[0].EntityID != "u1" || resp.Hits[0].Title != "Admin User" {
		t.Fatalf("unexpected hits %+v", resp.Hits)
	}
	if resp.Pagination.Total != 41 || resp.Pagination.TotalPages != 3 || resp.Pagination.Page != 3 {
		t.Fatalf("unexpected pagination %+v", resp.Pagination)
	}
}

func TestSearchPropagatesRepoError(t *testing.T) {
	t.Parallel()
	repo := &stubSearchRepo{err: errors.New("db down")}
	svc := searchsvc.NewService(repo)

	if _, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "a"}, viewKeys("users")); err == nil {
		t.Fatal("expected error")
	}
}

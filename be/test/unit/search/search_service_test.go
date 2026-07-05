package search_test

import (
	"context"
	"strings"
	"testing"

	"be/internal/common/query"
	"be/internal/common/rbac"
	searchdto "be/internal/dto/search"
	usermodel "be/internal/models/user"
	searchpkg "be/internal/search"
	searchsvc "be/internal/services/search"
)

type stubUserRepo struct {
	users []usermodel.User
}

func (s *stubUserRepo) Create(context.Context, *usermodel.User) error { return nil }
func (s *stubUserRepo) GetByID(_ context.Context, id string) (*usermodel.User, error) {
	for _, user := range s.users {
		if user.ID == id {
			copy := user
			return &copy, nil
		}
	}
	return nil, nil
}
func (s *stubUserRepo) GetByEmail(context.Context, string) (*usermodel.User, error) {
	return nil, nil
}
func (s *stubUserRepo) List(_ context.Context, q *query.Query) ([]usermodel.User, int64, error) {
	term := ""
	offset, limit := 0, len(s.users)
	if q != nil {
		term = strings.ToLower(strings.TrimSpace(q.LikeSearchTerm()))
		offset = q.Offset
		limit = q.Limit
	}

	items := make([]usermodel.User, 0, len(s.users))
	for _, user := range s.users {
		if term != "" {
			email := strings.ToLower(user.Email)
			name := strings.ToLower(user.FullName)
			if !strings.Contains(email, term) && !strings.Contains(name, term) {
				continue
			}
		}
		items = append(items, user)
	}
	total := int64(len(items))
	if offset >= len(items) {
		return []usermodel.User{}, total, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], total, nil
}
func (s *stubUserRepo) ListAll(context.Context) ([]usermodel.User, error) { return s.users, nil }
func (s *stubUserRepo) Update(context.Context, *usermodel.User) error       { return nil }
func (s *stubUserRepo) Delete(context.Context, string) error                 { return nil }

func TestSearchEmptyQueryReturnsEmpty(t *testing.T) {
	t.Parallel()

	client, _ := searchpkg.NewClient("http://localhost:9200", false)
	builder := searchsvc.NewDocumentBuilder(&stubUserRepo{}, nil, nil)
	svc := searchsvc.NewService(client, &stubUserRepo{}, builder)

	resp, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "   "}, []string{
		rbac.Key("users", rbac.ActionView),
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(resp.Hits) != 0 {
		t.Fatalf("expected empty hits, got %d", len(resp.Hits))
	}
}

func TestSearchFallbackFiltersByPermission(t *testing.T) {
	t.Parallel()

	repo := &stubUserRepo{
		users: []usermodel.User{
			{ID: "u1", Email: "admin@example.com", FullName: "Admin User", Status: usermodel.StatusActive},
		},
	}
	client, _ := searchpkg.NewClient("http://localhost:9200", false)
	builder := searchsvc.NewDocumentBuilder(repo, nil, nil)
	svc := searchsvc.NewService(client, repo, builder)

	resp, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "admin"}, []string{})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(resp.Hits) != 0 {
		t.Fatalf("expected no hits without users.view, got %d", len(resp.Hits))
	}
	if resp.Degraded {
		t.Fatal("expected non-degraded empty result when entity view permission is missing")
	}
}

func TestSearchFallbackReturnsUsersWhenSearchDisabled(t *testing.T) {
	t.Parallel()

	repo := &stubUserRepo{
		users: []usermodel.User{
			{ID: "u1", Email: "admin@example.com", FullName: "Admin User", Status: usermodel.StatusActive},
		},
	}
	client, _ := searchpkg.NewClient("http://localhost:9200", false)
	builder := searchsvc.NewDocumentBuilder(repo, nil, nil)
	svc := searchsvc.NewService(client, repo, builder)

	resp, err := svc.Search(context.Background(), searchdto.SearchQuery{Q: "admin"}, []string{
		rbac.Key("users", rbac.ActionView),
	})
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(resp.Hits) != 1 {
		t.Fatalf("expected 1 hit, got %d", len(resp.Hits))
	}
	if !resp.Degraded {
		t.Fatal("expected degraded fallback response")
	}
}

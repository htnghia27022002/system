package query_test

import (
	"testing"

	"be/pkg/query"
)

func TestNormalizePage(t *testing.T) {
	page, pageSize := query.NormalizePage(0, 200)
	if page != 1 || pageSize != query.MaxPageSize {
		t.Fatalf("got page=%d pageSize=%d", page, pageSize)
	}
}

func TestWhereLikeAnyCreatesOrGroup(t *testing.T) {
	q := query.New(2, 10).
		WhereEqual("role_id", "role-1").
		WhereLikeAny([]string{"email", "full_name"}, "admin")

	if len(q.Conditions.Predicates) != 1 {
		t.Fatalf("predicates = %d, want 1", len(q.Conditions.Predicates))
	}
	if len(q.Conditions.Groups) != 1 || q.Conditions.Groups[0].Logic != query.Or {
		t.Fatalf("expected one OR group")
	}
	if len(q.Conditions.Groups[0].Predicates) != 2 {
		t.Fatalf("or predicates = %d, want 2", len(q.Conditions.Groups[0].Predicates))
	}
}

func TestNestedOrGroup(t *testing.T) {
	q := query.New(1, 10).Or(func(g *query.Group) {
		g.Predicates = append(g.Predicates,
			query.Predicate{Column: "slug", Operator: query.OpEqual, Value: "admin"},
			query.Predicate{Column: "name", Operator: query.OpEqual, Value: "Admin"},
		)
	})

	if len(q.Conditions.Groups) != 1 {
		t.Fatalf("groups = %d, want 1", len(q.Conditions.Groups))
	}
}

func TestPageParamsPagination(t *testing.T) {
	page, pageSize, offset, limit := query.PageParams{Page: 3, PageSize: 5}.Pagination()
	if page != 3 || pageSize != 5 || offset != 10 || limit != 5 {
		t.Fatalf("unexpected pagination page=%d pageSize=%d offset=%d limit=%d", page, pageSize, offset, limit)
	}
}

func TestLikeSearchTerm(t *testing.T) {
	q := query.New(1, 10).WhereLikeAny([]string{"email"}, "admin")
	if q.LikeSearchTerm() != "admin" {
		t.Fatalf("LikeSearchTerm() = %q, want admin", q.LikeSearchTerm())
	}
}

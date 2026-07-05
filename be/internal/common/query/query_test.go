package query

import "testing"

func TestNormalizePage(t *testing.T) {
	page, pageSize := NormalizePage(0, 200)
	if page != 1 || pageSize != MaxPageSize {
		t.Fatalf("got page=%d pageSize=%d", page, pageSize)
	}
}

func TestWhereLikeAnyCreatesOrGroup(t *testing.T) {
	q := New(2, 10).
		WhereEqual("role_id", "role-1").
		WhereLikeAny([]string{"email", "full_name"}, "admin")

	if len(q.Conditions.Predicates) != 1 {
		t.Fatalf("predicates = %d, want 1", len(q.Conditions.Predicates))
	}
	if len(q.Conditions.Groups) != 1 || q.Conditions.Groups[0].Logic != Or {
		t.Fatalf("expected one OR group")
	}
	if len(q.Conditions.Groups[0].Predicates) != 2 {
		t.Fatalf("or predicates = %d, want 2", len(q.Conditions.Groups[0].Predicates))
	}
}

func TestNestedOrGroup(t *testing.T) {
	q := New(1, 10).Or(func(g *Group) {
		g.Predicates = append(g.Predicates,
			Predicate{Column: "slug", Operator: OpEqual, Value: "admin"},
			Predicate{Column: "name", Operator: OpEqual, Value: "Admin"},
		)
	})

	if len(q.Conditions.Groups) != 1 {
		t.Fatalf("groups = %d, want 1", len(q.Conditions.Groups))
	}
}

func TestPageParamsPagination(t *testing.T) {
	page, pageSize, offset, limit := PageParams{Page: 3, PageSize: 5}.Pagination()
	if page != 3 || pageSize != 5 || offset != 10 || limit != 5 {
		t.Fatalf("unexpected pagination page=%d pageSize=%d offset=%d limit=%d", page, pageSize, offset, limit)
	}
}

func TestLikeSearchTerm(t *testing.T) {
	q := New(1, 10).WhereLikeAny([]string{"email"}, "admin")
	if q.LikeSearchTerm() != "admin" {
		t.Fatalf("LikeSearchTerm() = %q, want admin", q.LikeSearchTerm())
	}
}

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

func TestCompileWhereEqualAndLikeOr(t *testing.T) {
	q := query.New(1, 10).
		WhereEqual("role_id", "role-1").
		WhereLikeAny([]string{"email", "full_name"}, "admin")

	sql, args, err := query.CompileWhere(q)
	if err != nil {
		t.Fatal(err)
	}
	if sql == "" {
		t.Fatal("expected SQL")
	}
	if len(args) != 3 {
		t.Fatalf("args = %d, want 3; sql=%s args=%v", len(args), sql, args)
	}
}

func TestCompileWhereEmpty(t *testing.T) {
	sql, args, err := query.CompileWhere(query.New(1, 10))
	if err != nil {
		t.Fatal(err)
	}
	if sql != "" || args != nil {
		t.Fatalf("expected empty where, got sql=%q args=%v", sql, args)
	}
}

func TestNewPageNormalizesNilItems(t *testing.T) {
	page := query.NewPage[string](nil, 3, 0, 200)
	if page.Items == nil {
		t.Fatal("items must be empty slice, not nil")
	}
	if len(page.Items) != 0 {
		t.Fatalf("items len = %d", len(page.Items))
	}
	if page.Total != 3 || page.Page != 1 || page.PageSize != query.MaxPageSize {
		t.Fatalf("unexpected page meta: %+v", page)
	}
}

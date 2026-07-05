package query

import "strings"

// PageParams is the shared HTTP pagination shape for list endpoints.
type PageParams struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"pageSize" binding:"omitempty,min=1,max=100"`
}

// Query describes pagination, ordering, and nested filter conditions.
type Query struct {
	Page         int
	PageSize     int
	Offset       int
	Limit        int
	Conditions   Group
	OrderClause  string
}

// New creates a query with normalized pagination and an AND root group.
func New(page, pageSize int) *Query {
	page, pageSize = NormalizePage(page, pageSize)
	return &Query{
		Page:     page,
		PageSize: pageSize,
		Offset:   Offset(page, pageSize),
		Limit:    pageSize,
		Conditions: Group{Logic: And},
	}
}

// Unbounded returns a large page query for internal bulk reads.
func Unbounded() *Query {
	return New(1, MaxPageSize*100)
}

// Pagination normalizes embedded page params.
func (p PageParams) Pagination() (page, pageSize, offset, limit int) {
	page, pageSize = NormalizePage(p.Page, p.PageSize)
	return page, pageSize, Offset(page, pageSize), pageSize
}

// Where adds an AND predicate.
func (q *Query) Where(column string, op Operator, value any) *Query {
	q.Conditions.Predicates = append(q.Conditions.Predicates, Predicate{
		Column:   column,
		Operator: op,
		Value:    value,
	})
	return q
}

// WhereEqual adds column = value when value is non-empty.
func (q *Query) WhereEqual(column, value string) *Query {
	value = strings.TrimSpace(value)
	if value == "" || strings.TrimSpace(column) == "" {
		return q
	}
	return q.Where(column, OpEqual, value)
}

// WhereRaw adds a custom SQL predicate to the AND root group.
func (q *Query) WhereRaw(sql string, args ...any) *Query {
	if strings.TrimSpace(sql) == "" {
		return q
	}
	q.Conditions.Predicates = append(q.Conditions.Predicates, Predicate{
		Operator: OpRaw,
		SQL:      sql,
		Args:     args,
	})
	return q
}

// WhereLikeAny adds a nested OR group: (col1 LIKE ? OR col2 LIKE ? ...).
func (q *Query) WhereLikeAny(columns []string, term string) *Query {
	term = strings.TrimSpace(term)
	if term == "" || len(columns) == 0 {
		return q
	}

	likeValue := "%" + strings.ToLower(term) + "%"
	orGroup := Group{Logic: Or}
	for _, column := range columns {
		if strings.TrimSpace(column) == "" {
			continue
		}
		orGroup.Predicates = append(orGroup.Predicates, Predicate{
			Column:   column,
			Operator: OpLike,
			Value:    likeValue,
		})
	}
	if len(orGroup.Predicates) == 0 {
		return q
	}

	q.Conditions.Groups = append(q.Conditions.Groups, orGroup)
	return q
}

// Or appends a nested OR group built by fn.
func (q *Query) Or(fn func(*Group)) *Query {
	group := Group{Logic: Or}
	if fn != nil {
		fn(&group)
	}
	if len(group.Predicates) > 0 || len(group.Groups) > 0 {
		q.Conditions.Groups = append(q.Conditions.Groups, group)
	}
	return q
}

// And appends a nested AND group built by fn.
func (q *Query) And(fn func(*Group)) *Query {
	group := Group{Logic: And}
	if fn != nil {
		fn(&group)
	}
	if len(group.Predicates) > 0 || len(group.Groups) > 0 {
		q.Conditions.Groups = append(q.Conditions.Groups, group)
	}
	return q
}

// OrderBy sets the SQL ORDER BY clause.
func (q *Query) OrderBy(clause string) *Query {
	q.OrderClause = strings.TrimSpace(clause)
	return q
}

// RoleByPermissionKeySQL filters roles that grant a permission key.
const RoleByPermissionKeySQL = `id IN (
	SELECT role_permissions.role_id FROM role_permissions
	JOIN permissions ON permissions.id = role_permissions.permission_id
	WHERE permissions.key = ?
)`

// LikeSearchTerm returns the first LIKE pattern term without wildcards, if any.
func (q *Query) LikeSearchTerm() string {
	if q == nil {
		return ""
	}
	return likeSearchTermFromGroup(q.Conditions)
}

func likeSearchTermFromGroup(group Group) string {
	for _, predicate := range group.Predicates {
		if predicate.Operator != OpLike {
			continue
		}
		if value, ok := predicate.Value.(string); ok {
			return strings.Trim(value, "%")
		}
	}
	for _, nested := range group.Groups {
		if term := likeSearchTermFromGroup(nested); term != "" {
			return term
		}
	}
	return ""
}

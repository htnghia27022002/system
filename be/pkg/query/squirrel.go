package query

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Masterminds/squirrel"

	"be/pkg/postgres"
)

// ApplySelect attaches filter conditions from q onto a squirrel select.
func ApplySelect(sb squirrel.SelectBuilder, q *Query) squirrel.SelectBuilder {
	if q == nil {
		return sb
	}
	if sqlizer := GroupSqlizer(q.Conditions); sqlizer != nil {
		sb = sb.Where(sqlizer)
	}
	return sb
}

// ApplyOrderLimit adds ORDER BY / OFFSET / LIMIT from q.
func ApplyOrderLimit(sb squirrel.SelectBuilder, q *Query) squirrel.SelectBuilder {
	if q == nil {
		q = New(DefaultPage, DefaultPageSize)
	}
	orderBy := strings.TrimSpace(q.OrderClause)
	if orderBy == "" {
		orderBy = "created_at DESC"
	}
	return sb.OrderBy(orderBy).Offset(uint64(q.Offset)).Limit(uint64(q.Limit))
}

// GroupSqlizer compiles a nested AND/OR group into a squirrel predicate.
func GroupSqlizer(group Group) squirrel.Sqlizer {
	parts := make([]squirrel.Sqlizer, 0, len(group.Predicates)+len(group.Groups))
	for _, predicate := range group.Predicates {
		if sqlizer := predicateSqlizer(predicate); sqlizer != nil {
			parts = append(parts, sqlizer)
		}
	}
	for _, nested := range group.Groups {
		if sqlizer := GroupSqlizer(nested); sqlizer != nil {
			parts = append(parts, sqlizer)
		}
	}
	if len(parts) == 0 {
		return nil
	}
	if group.Logic == Or {
		return squirrel.Or(parts)
	}
	return squirrel.And(parts)
}

func predicateSqlizer(predicate Predicate) squirrel.Sqlizer {
	switch predicate.Operator {
	case OpEqual:
		if strings.TrimSpace(predicate.Column) == "" {
			return nil
		}
		return squirrel.Expr(postgres.QuoteIdent(predicate.Column)+" = ?", predicate.Value)
	case OpNotEqual:
		if strings.TrimSpace(predicate.Column) == "" {
			return nil
		}
		return squirrel.Expr(postgres.QuoteIdent(predicate.Column)+" != ?", predicate.Value)
	case OpLike:
		if strings.TrimSpace(predicate.Column) == "" {
			return nil
		}
		return squirrel.Expr("LOWER("+postgres.QuoteIdent(predicate.Column)+") LIKE ?", predicate.Value)
	case OpIn:
		if strings.TrimSpace(predicate.Column) == "" {
			return nil
		}
		if isEmptySlice(predicate.Value) {
			return squirrel.Expr("FALSE")
		}
		return squirrel.Eq{postgres.QuoteIdent(predicate.Column): predicate.Value}
	case OpIsNull:
		if strings.TrimSpace(predicate.Column) == "" {
			return nil
		}
		return squirrel.Expr(postgres.QuoteIdent(predicate.Column) + " IS NULL")
	case OpNotNull:
		if strings.TrimSpace(predicate.Column) == "" {
			return nil
		}
		return squirrel.Expr(postgres.QuoteIdent(predicate.Column) + " IS NOT NULL")
	case OpRaw:
		if strings.TrimSpace(predicate.SQL) == "" {
			return nil
		}
		return squirrel.Expr(predicate.SQL, predicate.Args...)
	default:
		return nil
	}
}

func isEmptySlice(value any) bool {
	if value == nil {
		return true
	}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false
	}
	return rv.Len() == 0
}

// CompileWhere returns SQL and args for tests and ad-hoc debugging.
func CompileWhere(q *Query) (string, []any, error) {
	sqlizer := GroupSqlizer(q.Conditions)
	if sqlizer == nil {
		return "", nil, nil
	}
	sql, args, err := sqlizer.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("compile where: %w", err)
	}
	return sql, args, nil
}

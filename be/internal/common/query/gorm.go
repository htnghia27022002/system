package query

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

// Apply attaches query conditions to a GORM query.
func Apply(db *gorm.DB, q *Query) *gorm.DB {
	if q == nil {
		return db
	}
	return applyGroup(db, q.Conditions)
}

// Paginate counts and loads a paginated result set for model T.
func Paginate[T any](ctx context.Context, db *gorm.DB, q *Query, dest *[]T) (int64, error) {
	if q == nil {
		q = New(DefaultPage, DefaultPageSize)
	}

	base := Apply(db.WithContext(ctx).Model(new(T)), q)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return 0, err
	}

	orderBy := strings.TrimSpace(q.OrderClause)
	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	if err := base.Order(orderBy).Offset(q.Offset).Limit(q.Limit).Find(dest).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func applyGroup(db *gorm.DB, group Group) *gorm.DB {
	if len(group.Predicates) == 0 && len(group.Groups) == 0 {
		return db
	}

	if group.Logic == Or {
		orDB := db.Session(&gorm.Session{NewDB: true})
		started := false

		for _, predicate := range group.Predicates {
			clause := applyPredicate(orDB, predicate)
			if !started {
				orDB = clause
				started = true
				continue
			}
			orDB = orDB.Or(clause)
		}

		for _, nested := range group.Groups {
			clause := applyGroup(orDB.Session(&gorm.Session{NewDB: true}), nested)
			if !started {
				orDB = clause
				started = true
				continue
			}
			orDB = orDB.Or(clause)
		}

		return db.Where(orDB)
	}

	for _, predicate := range group.Predicates {
		db = applyPredicate(db, predicate)
	}
	for _, nested := range group.Groups {
		db = db.Where(applyGroup(db.Session(&gorm.Session{NewDB: true}), nested))
	}

	return db
}

func applyPredicate(db *gorm.DB, predicate Predicate) *gorm.DB {
	switch predicate.Operator {
	case OpEqual:
		return db.Where(predicate.Column+" = ?", predicate.Value)
	case OpNotEqual:
		return db.Where(predicate.Column+" != ?", predicate.Value)
	case OpLike:
		return db.Where("LOWER("+predicate.Column+") LIKE ?", predicate.Value)
	case OpIn:
		return db.Where(predicate.Column+" IN ?", predicate.Value)
	case OpIsNull:
		return db.Where(predicate.Column + " IS NULL")
	case OpNotNull:
		return db.Where(predicate.Column + " IS NOT NULL")
	case OpRaw:
		return db.Where(predicate.SQL, predicate.Args...)
	default:
		return db
	}
}

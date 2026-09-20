package repository

import (
	"context"
	"strings"

	categorymodel "be/internal/models/category"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type CategoryRepository struct {
	*repo.Repository[categorymodel.Category]
}

var _ interfaces.CategoryRepository = (*CategoryRepository)(nil)

func NewCategoryRepository(db *postgres.Postgres) *CategoryRepository {
	return &CategoryRepository{
		Repository: repo.New[categorymodel.Category](db, repo.Opts{Table: "categories", PK: "id"}),
	}
}

func (r *CategoryRepository) ListActive(ctx context.Context) ([]categorymodel.Category, error) {
	q := query.New(1, query.MaxPageSize).
		WhereRaw("is_active = TRUE").
		OrderBy("sort_order ASC, name ASC")
	return r.Find(ctx, q)
}

func (r *CategoryRepository) GetByKey(ctx context.Context, key string) (*categorymodel.Category, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("key", key))
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*categorymodel.Category, error) {
	return r.FindByID(ctx, id)
}

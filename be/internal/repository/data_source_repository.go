package repository

import (
	"context"

	dsmodel "be/internal/models/datasource"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type DataSourceRepository struct {
	*repo.Repository[dsmodel.Source]
}

var _ interfaces.DataSourceRepository = (*DataSourceRepository)(nil)

func NewDataSourceRepository(db *postgres.Postgres) *DataSourceRepository {
	return &DataSourceRepository{
		Repository: repo.New[dsmodel.Source](db, repo.Opts{Table: "data_sources", PK: "id"}),
	}
}

func (r *DataSourceRepository) GetByID(ctx context.Context, id string) (*dsmodel.Source, error) {
	return r.FindByID(ctx, id)
}

func (r *DataSourceRepository) List(ctx context.Context, page, limit int) ([]dsmodel.Source, int64, error) {
	q := query.New(page, limit).OrderBy("updated_at DESC")
	q.Page = page
	q.PageSize = limit
	q.Offset = (page - 1) * limit
	q.Limit = limit
	return r.Paginate(ctx, q)
}

func (r *DataSourceRepository) ListEnabled(ctx context.Context) ([]dsmodel.Source, error) {
	return r.Find(ctx, query.New(1, query.MaxPageSize*10).
		WhereRaw("enabled = TRUE").
		OrderBy("updated_at DESC"))
}

func (r *DataSourceRepository) Create(ctx context.Context, src *dsmodel.Source) error {
	if src.Headers == nil {
		src.Headers = map[string]any{}
	}
	if src.QueryParams == nil {
		src.QueryParams = map[string]any{}
	}
	if src.FieldMapping == nil {
		src.FieldMapping = map[string]any{}
	}
	return r.Insert(ctx, src)
}

func (r *DataSourceRepository) Update(ctx context.Context, src *dsmodel.Source) error {
	if src.Headers == nil {
		src.Headers = map[string]any{}
	}
	if src.QueryParams == nil {
		src.QueryParams = map[string]any{}
	}
	if src.FieldMapping == nil {
		src.FieldMapping = map[string]any{}
	}
	return r.Repository.Update(ctx, src)
}

func (r *DataSourceRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteByID(ctx, id)
}

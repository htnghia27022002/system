package repository

import (
	"context"
	"strings"
	"time"

	dsmodel "be/internal/models/datasource"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type IngestRunRepository struct {
	*repo.Repository[dsmodel.IngestRun]
	sources *repo.Repository[dsmodel.IngestRunSource]
}

var _ interfaces.IngestRunRepository = (*IngestRunRepository)(nil)

func NewIngestRunRepository(db *postgres.Postgres) *IngestRunRepository {
	return &IngestRunRepository{
		Repository: repo.New[dsmodel.IngestRun](db, repo.Opts{Table: "data_ingest_runs", PK: "id"}),
		sources:    repo.New[dsmodel.IngestRunSource](db, repo.Opts{Table: "data_ingest_run_sources", PK: "id"}),
	}
}

func (r *IngestRunRepository) HasActive(ctx context.Context) (bool, error) {
	count, err := r.Count(ctx, query.New(1, 1).WhereRaw("status IN (?, ?)", dsmodel.RunQueued, dsmodel.RunRunning))
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *IngestRunRepository) Insert(ctx context.Context, run *dsmodel.IngestRun) error {
	return r.Repository.Insert(ctx, run)
}

func (r *IngestRunRepository) GetByID(ctx context.Context, id string) (*dsmodel.IngestRun, error) {
	return r.FindByID(ctx, id)
}

func (r *IngestRunRepository) Latest(ctx context.Context) (*dsmodel.IngestRun, error) {
	return r.FindOne(ctx, query.New(1, 1).OrderBy("created_at DESC"))
}

func (r *IngestRunRepository) ClaimQueued(ctx context.Context, id string) (*dsmodel.IngestRun, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil
	}
	now := time.Now().UTC()
	n, err := r.UpdateWhere(ctx,
		query.New(1, 1).WhereEqual("id", id).WhereRaw("status = ?", dsmodel.RunQueued),
		map[string]any{
			"status":     dsmodel.RunRunning,
			"started_at": now,
			"updated_at": now,
		},
	)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, nil
	}
	return r.FindByID(ctx, id)
}

func (r *IngestRunRepository) ListQueued(ctx context.Context, limit int) ([]dsmodel.IngestRun, error) {
	if limit < 1 {
		limit = 10
	}
	return r.Find(ctx, query.New(1, limit).
		WhereRaw("status = ?", dsmodel.RunQueued).
		OrderBy("created_at ASC"))
}

func (r *IngestRunRepository) Update(ctx context.Context, run *dsmodel.IngestRun) error {
	return r.Repository.Update(ctx, run)
}

func (r *IngestRunRepository) InsertSource(ctx context.Context, row *dsmodel.IngestRunSource) error {
	return r.sources.Insert(ctx, row)
}

func (r *IngestRunRepository) UpdateSource(ctx context.Context, row *dsmodel.IngestRunSource) error {
	return r.sources.Update(ctx, row)
}

func (r *IngestRunRepository) ListSources(ctx context.Context, runID string) ([]dsmodel.IngestRunSource, error) {
	return r.sources.Find(ctx, query.New(1, query.MaxPageSize).
		WhereEqual("run_id", runID).
		OrderBy("created_at ASC"))
}

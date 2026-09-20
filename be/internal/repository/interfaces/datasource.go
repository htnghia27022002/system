package interfaces

import (
	"context"

	dsmodel "be/internal/models/datasource"
)

// DataSourceRepository persists HTTP ingest Sources.
type DataSourceRepository interface {
	GetByID(ctx context.Context, id string) (*dsmodel.Source, error)
	List(ctx context.Context, page, limit int) ([]dsmodel.Source, int64, error)
	ListEnabled(ctx context.Context) ([]dsmodel.Source, error)
	Create(ctx context.Context, src *dsmodel.Source) error
	Update(ctx context.Context, src *dsmodel.Source) error
	Delete(ctx context.Context, id string) error
}

// IngestRunRepository persists ingest jobs and per-source outcomes.
type IngestRunRepository interface {
	HasActive(ctx context.Context) (bool, error)
	Insert(ctx context.Context, run *dsmodel.IngestRun) error
	GetByID(ctx context.Context, id string) (*dsmodel.IngestRun, error)
	Latest(ctx context.Context) (*dsmodel.IngestRun, error)
	ClaimQueued(ctx context.Context, id string) (*dsmodel.IngestRun, error)
	ListQueued(ctx context.Context, limit int) ([]dsmodel.IngestRun, error)
	Update(ctx context.Context, run *dsmodel.IngestRun) error
	InsertSource(ctx context.Context, row *dsmodel.IngestRunSource) error
	UpdateSource(ctx context.Context, row *dsmodel.IngestRunSource) error
	ListSources(ctx context.Context, runID string) ([]dsmodel.IngestRunSource, error)
}

package interfaces

import (
	"context"

	searchmodel "be/internal/models/search"
)

type OutboxStats struct {
	PendingCount            int64
	FailedCount             int64
	OldestPendingAgeSeconds int64
}

type SearchOutboxRepository interface {
	Enqueue(ctx context.Context, entityType, entityID, operation string) (string, error)
	GetByID(ctx context.Context, id string) (*searchmodel.OutboxEntry, error)
	ClaimByID(ctx context.Context, id string) (*searchmodel.OutboxEntry, error)
	ClaimPending(ctx context.Context, limit int) ([]searchmodel.OutboxEntry, error)
	MarkCompleted(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, errMsg string) error
	ResetFailed(ctx context.Context, id string) error
	Stats(ctx context.Context) (*OutboxStats, error)
}

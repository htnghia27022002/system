package interfaces

import (
	"context"
	"time"

	auditmodel "be/internal/models/audit"
	"be/pkg/query"
)

// AuditLogRepository persists the audit trail.
type AuditLogRepository interface {
	Insert(ctx context.Context, log *auditmodel.Log) error
	List(ctx context.Context, q *query.Query) ([]auditmodel.Log, int64, error)
	// CountByActionPerDay counts entries of action per UTC day since `since`.
	CountByActionPerDay(ctx context.Context, action string, since time.Time) (map[string]int64, error)
	// CountByAction counts entries of action in [from, to).
	CountByAction(ctx context.Context, action string, from, to time.Time) (int64, error)
	// DeleteBefore removes entries older than t and returns how many.
	DeleteBefore(ctx context.Context, t time.Time) (int64, error)
}

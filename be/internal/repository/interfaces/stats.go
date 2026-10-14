package interfaces

import (
	"context"
	"time"
)

// StatsRepository answers the counts behind the admin dashboard.
type StatsRepository interface {
	// CountUsers counts non-deleted users created before t.
	CountUsers(ctx context.Context, before time.Time) (int64, error)
	// CountUsersCreated counts non-deleted users created in [from, to).
	CountUsersCreated(ctx context.Context, from, to time.Time) (int64, error)
	// CountActiveSessions counts unrevoked, unexpired refresh tokens at t.
	CountActiveSessions(ctx context.Context, at time.Time) (int64, error)
	CountRoles(ctx context.Context) (int64, error)
}

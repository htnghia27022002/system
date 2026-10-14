package repository

import (
	"context"
	"time"

	"be/internal/repository/interfaces"
	"be/pkg/postgres"
)

type StatsRepository struct {
	db *postgres.Postgres
}

var _ interfaces.StatsRepository = (*StatsRepository)(nil)

func NewStatsRepository(db *postgres.Postgres) *StatsRepository {
	return &StatsRepository{db: db}
}

func (r *StatsRepository) count(ctx context.Context, sql string, args ...any) (int64, error) {
	var n int64
	err := r.db.Querier(ctx).QueryRow(ctx, sql, args...).Scan(&n)
	return n, err
}

func (r *StatsRepository) CountUsers(ctx context.Context, before time.Time) (int64, error) {
	return r.count(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at < $1`, before)
}

func (r *StatsRepository) CountUsersCreated(ctx context.Context, from, to time.Time) (int64, error) {
	return r.count(ctx, `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND created_at >= $1 AND created_at < $2`, from, to)
}

func (r *StatsRepository) CountActiveSessions(ctx context.Context, at time.Time) (int64, error) {
	return r.count(ctx, `SELECT COUNT(*) FROM refresh_tokens WHERE revoked_at IS NULL AND expires_at > $1`, at)
}

func (r *StatsRepository) CountRoles(ctx context.Context) (int64, error) {
	return r.count(ctx, `SELECT COUNT(*) FROM roles`)
}

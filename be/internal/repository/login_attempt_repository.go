package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	"be/internal/repository/interfaces"
	"be/pkg/postgres"
)

// LoginAttemptRepository tracks consecutive failed logins per user (table login_attempts).
type LoginAttemptRepository struct {
	db *postgres.Postgres
}

var _ interfaces.LoginAttemptRepository = (*LoginAttemptRepository)(nil)

func NewLoginAttemptRepository(db *postgres.Postgres) *LoginAttemptRepository {
	return &LoginAttemptRepository{db: db}
}

func (r *LoginAttemptRepository) LockedUntil(ctx context.Context, userID string) (*time.Time, error) {
	var lockedUntil *time.Time
	err := r.db.Querier(ctx).QueryRow(ctx,
		`SELECT locked_until FROM login_attempts WHERE user_id = $1`, userID,
	).Scan(&lockedUntil)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return lockedUntil, err
}

// RecordFailure counts one failure; on reaching maxFailures it locks the account
// for lockFor and restarts the count. It returns the lock expiry when it locked.
func (r *LoginAttemptRepository) RecordFailure(ctx context.Context, userID string, maxFailures int, lockFor time.Duration) (*time.Time, error) {
	var count int
	err := r.db.Querier(ctx).QueryRow(ctx, `
INSERT INTO login_attempts (user_id, failed_count, updated_at) VALUES ($1, 1, NOW())
ON CONFLICT (user_id) DO UPDATE SET failed_count = login_attempts.failed_count + 1, updated_at = NOW()
RETURNING failed_count`, userID).Scan(&count)
	if err != nil || count < maxFailures {
		return nil, err
	}

	lockedUntil := time.Now().Add(lockFor)
	_, err = r.db.Querier(ctx).Exec(ctx,
		`UPDATE login_attempts SET failed_count = 0, locked_until = $2, updated_at = NOW() WHERE user_id = $1`,
		userID, lockedUntil,
	)
	if err != nil {
		return nil, err
	}
	return &lockedUntil, nil
}

func (r *LoginAttemptRepository) Reset(ctx context.Context, userID string) error {
	_, err := r.db.Querier(ctx).Exec(ctx, `DELETE FROM login_attempts WHERE user_id = $1`, userID)
	return err
}

package interfaces

import (
	"context"
	"time"
)

// LoginAttemptRepository backs temporary account lockout after repeated failed logins.
type LoginAttemptRepository interface {
	// LockedUntil returns the lock expiry, or nil when the user has no lock record.
	LockedUntil(ctx context.Context, userID string) (*time.Time, error)
	// RecordFailure counts a failed login and locks the account once maxFailures is
	// reached; it returns the new lock expiry when it locked, else nil.
	RecordFailure(ctx context.Context, userID string, maxFailures int, lockFor time.Duration) (*time.Time, error)
	// Reset clears the failure count after a successful login.
	Reset(ctx context.Context, userID string) error
}

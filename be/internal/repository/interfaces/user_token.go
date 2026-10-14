package interfaces

import (
	"context"
	"time"

	accountmodel "be/internal/models/account"
)

// UserTokenRepository stores single-use email tokens (reset, verify, invite).
type UserTokenRepository interface {
	Create(ctx context.Context, token *accountmodel.UserToken) error
	// FindValid returns the unused, unexpired token with this hash and purpose, or nil.
	FindValid(ctx context.Context, purpose accountmodel.TokenPurpose, tokenHash string) (*accountmodel.UserToken, error)
	// Consume marks the token used. It returns false when it was already used,
	// so two concurrent requests cannot both redeem one token.
	Consume(ctx context.Context, id string) (bool, error)
	// InvalidateOpen marks every open token of this purpose for the user as used.
	InvalidateOpen(ctx context.Context, userID string, purpose accountmodel.TokenPurpose) error
	// LatestCreatedAt returns when the newest open token of this purpose was issued, or nil.
	LatestCreatedAt(ctx context.Context, userID string, purpose accountmodel.TokenPurpose) (*time.Time, error)
}

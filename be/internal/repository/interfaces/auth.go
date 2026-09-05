package interfaces

import (
	"context"
	"time"

	authmodel "be/internal/models/auth"
)

type AuthRepository interface {
	CreateRefreshToken(ctx context.Context, token *authmodel.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*authmodel.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string, revokedAt time.Time) error
	// ListActiveByUserID returns unrevoked, unexpired refresh tokens for the user (no pagination).
	ListActiveByUserID(ctx context.Context, userID string) ([]authmodel.RefreshToken, error)
	// RevokeByID sets revoked_at on an active row owned by userID. Returns rows affected (0 = not found).
	RevokeByID(ctx context.Context, userID, id string) (int64, error)
	// RevokeAllExcept revokes every active token for userID except keepID.
	RevokeAllExcept(ctx context.Context, userID, keepID string) (int64, error)
	FindOAuthAccount(ctx context.Context, provider, providerUserID string) (*authmodel.OAuthAccount, error)
	CreateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error
	UpdateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error
	// ListProvidersByUserIDs returns distinct OAuth provider ids keyed by user id.
	ListProvidersByUserIDs(ctx context.Context, userIDs []string) (map[string][]string, error)
}

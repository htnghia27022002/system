package interfaces

import (
	"context"
	"time"

	authmodel "be/internal/models/auth"
	usermodel "be/internal/models/user"
)

type AuthRepository interface {
	FindUserByEmail(ctx context.Context, email string) (*usermodel.User, error)
	FindUserByID(ctx context.Context, id string) (*usermodel.User, error)
	CreateUser(ctx context.Context, user *usermodel.User) error
	CreateRefreshToken(ctx context.Context, token *authmodel.RefreshToken) error
	FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*authmodel.RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, tokenHash string, revokedAt time.Time) error
	FindOAuthAccount(ctx context.Context, provider, providerUserID string) (*authmodel.OAuthAccount, error)
	CreateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error
	UpdateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error
}

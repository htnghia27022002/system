package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"

	authmodel "be/internal/models/auth"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
)

type AuthRepository struct {
	db *gorm.DB
}

var _ interfaces.AuthRepository = (*AuthRepository)(nil)

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*usermodel.User, error) {
	var user usermodel.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) FindUserByID(ctx context.Context, id string) (*usermodel.User, error) {
	var user usermodel.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *usermodel.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, token *authmodel.RefreshToken) error {
	if token.TokenHash == "" && token.RawToken != "" {
		token.TokenHash = hashToken(token.RawToken)
	}
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *AuthRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*authmodel.RefreshToken, error) {
	var token authmodel.RefreshToken
	if err := r.db.WithContext(ctx).
		Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", tokenHash, time.Now()).
		First(&token).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	return r.db.WithContext(ctx).
		Model(&authmodel.RefreshToken{}).
		Where("token_hash = ?", tokenHash).
		Update("revoked_at", revokedAt).Error
}

func (r *AuthRepository) FindOAuthAccount(ctx context.Context, provider, providerUserID string) (*authmodel.OAuthAccount, error) {
	var account authmodel.OAuthAccount
	if err := r.db.WithContext(ctx).
		Where("provider = ? AND provider_user_id = ?", provider, providerUserID).
		First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &account, nil
}

func (r *AuthRepository) CreateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *AuthRepository) UpdateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func HashRefreshToken(value string) string {
	return hashToken(value)
}

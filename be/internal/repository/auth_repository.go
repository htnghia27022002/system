package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"

	authmodel "be/internal/models/auth"
	"be/internal/repository/interfaces"
	"be/pkg/hash"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type AuthRepository struct {
	tokens   *repo.Repository[authmodel.RefreshToken]
	accounts *repo.Repository[authmodel.OAuthAccount]
}

var _ interfaces.AuthRepository = (*AuthRepository)(nil)

func NewAuthRepository(db *postgres.Postgres) *AuthRepository {
	return &AuthRepository{
		tokens:   repo.New[authmodel.RefreshToken](db, repo.Opts{Table: "refresh_tokens", PK: "id"}),
		accounts: repo.New[authmodel.OAuthAccount](db, repo.Opts{Table: "oauth_accounts", PK: "id"}),
	}
}

func (r *AuthRepository) CreateRefreshToken(ctx context.Context, token *authmodel.RefreshToken) error {
	if token.TokenHash == "" && token.RawToken != "" {
		token.TokenHash = hash.SHA256Hex(token.RawToken)
	}
	return r.tokens.Insert(ctx, token)
}

func (r *AuthRepository) FindRefreshTokenByHash(ctx context.Context, tokenHash string) (*authmodel.RefreshToken, error) {
	return r.tokens.FindOne(ctx, query.New(1, 1).
		WhereEqual("token_hash", tokenHash).
		Where("revoked_at", query.OpIsNull, nil).
		WhereRaw("expires_at > ?", time.Now()))
}

func (r *AuthRepository) RevokeRefreshToken(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	_, err := r.tokens.UpdateWhere(ctx,
		query.New(1, 1).WhereEqual("token_hash", tokenHash),
		map[string]any{"revoked_at": revokedAt},
	)
	return err
}

func (r *AuthRepository) ListActiveByUserID(ctx context.Context, userID string) ([]authmodel.RefreshToken, error) {
	if strings.TrimSpace(userID) == "" {
		return []authmodel.RefreshToken{}, nil
	}
	sb := r.tokens.Select().
		Where(squirrel.Eq{postgres.QuoteIdent("user_id"): userID}).
		Where(squirrel.Expr(postgres.QuoteIdent("revoked_at") + " IS NULL")).
		Where(squirrel.Expr(postgres.QuoteIdent("expires_at")+" > ?", time.Now())).
		OrderBy(postgres.QuoteIdent("last_used_at") + " DESC")
	return r.tokens.Query(ctx, sb)
}

func (r *AuthRepository) RevokeByID(ctx context.Context, userID, id string) (int64, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(id) == "" {
		return 0, nil
	}
	now := time.Now()
	return r.tokens.UpdateWhere(ctx,
		query.New(1, 1).
			WhereEqual("id", id).
			WhereEqual("user_id", userID).
			Where("revoked_at", query.OpIsNull, nil).
			WhereRaw("expires_at > ?", now),
		map[string]any{"revoked_at": now},
	)
}

func (r *AuthRepository) RevokeAllExcept(ctx context.Context, userID, keepID string) (int64, error) {
	if strings.TrimSpace(userID) == "" || strings.TrimSpace(keepID) == "" {
		return 0, nil
	}
	now := time.Now()
	return r.tokens.UpdateWhere(ctx,
		query.New(1, 1).
			WhereEqual("user_id", userID).
			Where("id", query.OpNotEqual, keepID).
			Where("revoked_at", query.OpIsNull, nil).
			WhereRaw("expires_at > ?", now),
		map[string]any{"revoked_at": now},
	)
}

func (r *AuthRepository) FindOAuthAccount(ctx context.Context, provider, providerUserID string) (*authmodel.OAuthAccount, error) {
	return r.accounts.FindOne(ctx, query.New(1, 1).
		WhereEqual("provider", provider).
		WhereEqual("provider_user_id", providerUserID))
}

func (r *AuthRepository) CreateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error {
	return r.accounts.Insert(ctx, account)
}

func (r *AuthRepository) UpdateOAuthAccount(ctx context.Context, account *authmodel.OAuthAccount) error {
	return r.accounts.Update(ctx, account)
}

func (r *AuthRepository) ListProvidersByUserIDs(ctx context.Context, userIDs []string) (map[string][]string, error) {
	out := make(map[string][]string)
	if len(userIDs) == 0 {
		return out, nil
	}

	sql, args, err := r.accounts.DB().Builder.
		Select("user_id", "provider").
		From(postgres.QuoteIdent("oauth_accounts")).
		Where(squirrel.Eq{postgres.QuoteIdent("user_id"): userIDs}).
		OrderBy("provider ASC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.accounts.DB().Querier(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	seen := make(map[string]map[string]struct{})
	for rows.Next() {
		var userID, provider string
		if err := rows.Scan(&userID, &provider); err != nil {
			return nil, err
		}
		if userID == "" || provider == "" {
			continue
		}
		if _, ok := seen[userID]; !ok {
			seen[userID] = map[string]struct{}{}
		}
		if _, ok := seen[userID][provider]; ok {
			continue
		}
		seen[userID][provider] = struct{}{}
		out[userID] = append(out[userID], provider)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

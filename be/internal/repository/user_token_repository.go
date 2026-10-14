package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"

	accountmodel "be/internal/models/account"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type UserTokenRepository struct {
	*repo.Repository[accountmodel.UserToken]
	db *postgres.Postgres
}

var _ interfaces.UserTokenRepository = (*UserTokenRepository)(nil)

func NewUserTokenRepository(db *postgres.Postgres) *UserTokenRepository {
	return &UserTokenRepository{
		Repository: repo.New[accountmodel.UserToken](db, repo.Opts{Table: "user_tokens", PK: "id"}),
		db:         db,
	}
}

func (r *UserTokenRepository) Create(ctx context.Context, token *accountmodel.UserToken) error {
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}
	return r.Insert(ctx, token)
}

func (r *UserTokenRepository) FindValid(ctx context.Context, purpose accountmodel.TokenPurpose, tokenHash string) (*accountmodel.UserToken, error) {
	return r.FindOne(ctx, query.New(1, 1).
		WhereEqual("token_hash", tokenHash).
		WhereEqual("purpose", string(purpose)).
		Where("used_at", query.OpIsNull, nil).
		WhereRaw("expires_at > ?", time.Now()))
}

func (r *UserTokenRepository) Consume(ctx context.Context, id string) (bool, error) {
	tag, err := r.db.Querier(ctx).Exec(ctx,
		`UPDATE user_tokens SET used_at = NOW() WHERE id = $1 AND used_at IS NULL`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() == 1, nil
}

func (r *UserTokenRepository) InvalidateOpen(ctx context.Context, userID string, purpose accountmodel.TokenPurpose) error {
	_, err := r.db.Querier(ctx).Exec(ctx,
		`UPDATE user_tokens SET used_at = NOW() WHERE user_id = $1 AND purpose = $2 AND used_at IS NULL`,
		userID, string(purpose))
	return err
}

func (r *UserTokenRepository) LatestCreatedAt(ctx context.Context, userID string, purpose accountmodel.TokenPurpose) (*time.Time, error) {
	var created time.Time
	err := r.db.Querier(ctx).QueryRow(ctx,
		`SELECT created_at FROM user_tokens WHERE user_id = $1 AND purpose = $2 AND used_at IS NULL
		 ORDER BY created_at DESC LIMIT 1`, userID, string(purpose)).Scan(&created)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &created, nil
}

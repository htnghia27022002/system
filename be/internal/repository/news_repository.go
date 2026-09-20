package repository

import (
	"context"
	"strings"

	newsmodel "be/internal/models/news"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type NewsRepository struct {
	*repo.Repository[newsmodel.News]
}

var _ interfaces.NewsRepository = (*NewsRepository)(nil)

func NewNewsRepository(db *postgres.Postgres) *NewsRepository {
	return &NewsRepository{
		Repository: repo.New[newsmodel.News](db, repo.Opts{Table: "news", PK: "id"}),
	}
}

func (r *NewsRepository) GetByOriginalURL(ctx context.Context, originalURL string) (*newsmodel.News, error) {
	originalURL = strings.TrimSpace(originalURL)
	if originalURL == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("original_url", originalURL))
}

func (r *NewsRepository) ListByPlaceID(ctx context.Context, placeID string) ([]newsmodel.News, error) {
	return r.Find(ctx, query.New(1, query.MaxPageSize).
		WhereEqual("place_id", placeID).
		OrderBy("created_at DESC"))
}

func (r *NewsRepository) Create(ctx context.Context, item *newsmodel.News) error {
	return r.Insert(ctx, item)
}

func (r *NewsRepository) Update(ctx context.Context, item *newsmodel.News) error {
	return r.Repository.Update(ctx, item)
}

package interfaces

import (
	"context"

	newsmodel "be/internal/models/news"
)

// NewsRepository upserts articles by original URL.
type NewsRepository interface {
	GetByOriginalURL(ctx context.Context, originalURL string) (*newsmodel.News, error)
	ListByPlaceID(ctx context.Context, placeID string) ([]newsmodel.News, error)
	Create(ctx context.Context, item *newsmodel.News) error
	Update(ctx context.Context, item *newsmodel.News) error
}

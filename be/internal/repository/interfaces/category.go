package interfaces

import (
	"context"

	categorymodel "be/internal/models/category"
)

// CategoryRepository reads the shared category catalog.
// GetByKey never inserts; unknown keys are resolved to uncategorized by the service.
type CategoryRepository interface {
	ListActive(ctx context.Context) ([]categorymodel.Category, error)
	GetByKey(ctx context.Context, key string) (*categorymodel.Category, error)
	GetByID(ctx context.Context, id string) (*categorymodel.Category, error)
}

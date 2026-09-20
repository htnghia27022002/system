package interfaces

import (
	"context"

	newsmodel "be/internal/models/news"
	placemodel "be/internal/models/place"
)

// PlaceListFilter is the pinnable Places query.
type PlaceListFilter struct {
	CategoryKey     string
	Q               string
	AdminPathPrefix string
	CountryCode     string
	Page            int
	Limit           int
}

// PlaceAdminListFilter lists Places for operator CRUD (any status).
type PlaceAdminListFilter struct {
	LocationID string
	Q          string
	Status     *int16
	Page       int
	Limit      int
}

// PlaceListRow is one joined pin row.
type PlaceListRow struct {
	Place        placemodel.Place
	CategoryKey  string
	LocationID   string
	LocationName string
}

// PlaceDetailRow is Place + Location + category + geo for detail.
type PlaceDetailRow struct {
	Place           placemodel.Place
	CategoryKey     string
	Location        PlaceLocationRow
	News            []newsmodel.News
	NewsCategoryKey map[string]string
}

// PlaceLocationRow is Location plus optional country/division labels.
type PlaceLocationRow struct {
	ID              string
	Name            string
	LocationKey     *string
	CountryID       *string
	CountryCode     *string
	AdminDivisionID *string
	AdminPath       *string
	Street          string
	PostalCode      string
	Formatted       string
	Lat             *float64
	Lng             *float64
}

// PlaceRepository persists Places and list/detail joins.
type PlaceRepository interface {
	GetByID(ctx context.Context, id string) (*placemodel.Place, error)
	FindByLocationAndKey(ctx context.Context, locationID, placeKey string) (*placemodel.Place, error)
	Create(ctx context.Context, p *placemodel.Place) error
	Update(ctx context.Context, p *placemodel.Place) error
	ListPinnable(ctx context.Context, filter PlaceListFilter) ([]PlaceListRow, int64, error)
	ListAdmin(ctx context.Context, filter PlaceAdminListFilter) ([]PlaceListRow, int64, error)
	CountByLocationID(ctx context.Context, locationID string) (int64, error)
	Delete(ctx context.Context, id string) error
	GetDetail(ctx context.Context, id string) (*PlaceDetailRow, error)
}

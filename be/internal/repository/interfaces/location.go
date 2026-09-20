package interfaces

import (
	"context"

	locationmodel "be/internal/models/location"
)

// LocationMatch describes reuse keys for find-or-create.
type LocationMatch struct {
	LocationKey     string
	CountryID       *string
	AdminDivisionID *string
	Street          string
	Formatted       string
	Lat             *float64
	Lng             *float64
}

// LocationListFilter pages operator Location lists.
type LocationListFilter struct {
	Q               string
	CountryCode     string
	AdminDivisionID string
	Page            int
	Limit           int
}

// LocationListRow is a Location plus optional geo labels.
type LocationListRow struct {
	Location    locationmodel.Location
	CountryCode *string
	AdminPath   *string
}

// LocationRepository persists shared Location sites.
type LocationRepository interface {
	GetByID(ctx context.Context, id string) (*locationmodel.Location, error)
	FindByLocationKey(ctx context.Context, key string) (*locationmodel.Location, error)
	FindByCountryAdminStreet(ctx context.Context, countryID, adminDivisionID, street string) (*locationmodel.Location, error)
	FindByFormatted(ctx context.Context, formatted string) (*locationmodel.Location, error)
	FindByRoundedCoordsAndFormatted(ctx context.Context, lat, lng float64, formatted string) (*locationmodel.Location, error)
	List(ctx context.Context, filter LocationListFilter) ([]LocationListRow, int64, error)
	Create(ctx context.Context, loc *locationmodel.Location) error
	Update(ctx context.Context, loc *locationmodel.Location) error
	Delete(ctx context.Context, id string) error
}

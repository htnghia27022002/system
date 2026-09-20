package interfaces

import (
	"context"

	geomodel "be/internal/models/geo"
)

// CountryRepository reads the shared ISO country catalog.
type CountryRepository interface {
	ListActive(ctx context.Context) ([]geomodel.Country, error)
	GetByCode(ctx context.Context, code string) (*geomodel.Country, error)
	GetByID(ctx context.Context, id string) (*geomodel.Country, error)
}

// DivisionListFilter filters administrative divisions.
type DivisionListFilter struct {
	CountryID string
	ParentID  *string // nil = no parent filter; pointer to "" = top-level only
	Q         string
}

// DivisionRepository reads the administrative division tree.
type DivisionRepository interface {
	ListActive(ctx context.Context, filter DivisionListFilter) ([]geomodel.Division, error)
	GetByID(ctx context.Context, id string) (*geomodel.Division, error)
	GetByCountryAndCode(ctx context.Context, countryID, code string) (*geomodel.Division, error)
}

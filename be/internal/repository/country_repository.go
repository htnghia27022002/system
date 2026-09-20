package repository

import (
	"context"
	"strings"

	geomodel "be/internal/models/geo"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type CountryRepository struct {
	*repo.Repository[geomodel.Country]
}

var _ interfaces.CountryRepository = (*CountryRepository)(nil)

func NewCountryRepository(db *postgres.Postgres) *CountryRepository {
	return &CountryRepository{
		Repository: repo.New[geomodel.Country](db, repo.Opts{Table: "countries", PK: "id"}),
	}
}

func (r *CountryRepository) ListActive(ctx context.Context) ([]geomodel.Country, error) {
	q := query.New(1, query.MaxPageSize).
		WhereRaw("is_active = TRUE").
		OrderBy("name ASC")
	return r.Find(ctx, q)
}

func (r *CountryRepository) GetByCode(ctx context.Context, code string) (*geomodel.Country, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("code", code))
}

func (r *CountryRepository) GetByID(ctx context.Context, id string) (*geomodel.Country, error) {
	return r.FindByID(ctx, id)
}

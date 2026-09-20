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

type DivisionRepository struct {
	*repo.Repository[geomodel.Division]
}

var _ interfaces.DivisionRepository = (*DivisionRepository)(nil)

func NewDivisionRepository(db *postgres.Postgres) *DivisionRepository {
	return &DivisionRepository{
		Repository: repo.New[geomodel.Division](db, repo.Opts{Table: "administrative_divisions", PK: "id"}),
	}
}

func (r *DivisionRepository) ListActive(ctx context.Context, filter interfaces.DivisionListFilter) ([]geomodel.Division, error) {
	q := query.New(1, query.MaxPageSize*10).WhereRaw("is_active = TRUE")
	if filter.CountryID != "" {
		q.WhereEqual("country_id", filter.CountryID)
	}
	if filter.ParentID != nil {
		if strings.TrimSpace(*filter.ParentID) == "" {
			q.WhereRaw("parent_id IS NULL")
		} else {
			q.WhereEqual("parent_id", strings.TrimSpace(*filter.ParentID))
		}
	}
	if strings.TrimSpace(filter.Q) != "" {
		q.WhereLikeAny([]string{"name", "name_en", "full_name", "code"}, filter.Q)
	}
	q.OrderBy("name ASC")
	return r.Find(ctx, q)
}

func (r *DivisionRepository) GetByID(ctx context.Context, id string) (*geomodel.Division, error) {
	return r.FindByID(ctx, id)
}

func (r *DivisionRepository) GetByCountryAndCode(ctx context.Context, countryID, code string) (*geomodel.Division, error) {
	countryID = strings.TrimSpace(countryID)
	code = strings.TrimSpace(code)
	if countryID == "" || code == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).
		WhereEqual("country_id", countryID).
		WhereEqual("code", code))
}

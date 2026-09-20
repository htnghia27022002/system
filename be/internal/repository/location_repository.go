package repository

import (
	"context"
	"math"
	"strings"

	"github.com/Masterminds/squirrel"

	locationmodel "be/internal/models/location"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type LocationRepository struct {
	*repo.Repository[locationmodel.Location]
}

var _ interfaces.LocationRepository = (*LocationRepository)(nil)

func NewLocationRepository(db *postgres.Postgres) *LocationRepository {
	return &LocationRepository{
		Repository: repo.New[locationmodel.Location](db, repo.Opts{Table: "locations", PK: "id"}),
	}
}

func (r *LocationRepository) GetByID(ctx context.Context, id string) (*locationmodel.Location, error) {
	return r.FindByID(ctx, id)
}

func (r *LocationRepository) FindByLocationKey(ctx context.Context, key string) (*locationmodel.Location, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("location_key", key))
}

func (r *LocationRepository) FindByCountryAdminStreet(ctx context.Context, countryID, adminDivisionID, street string) (*locationmodel.Location, error) {
	countryID = strings.TrimSpace(countryID)
	adminDivisionID = strings.TrimSpace(adminDivisionID)
	street = strings.TrimSpace(street)
	if countryID == "" || adminDivisionID == "" || street == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).WhereRaw(
		"country_id = ? AND admin_division_id = ? AND lower(btrim(street)) = lower(btrim(?))",
		countryID, adminDivisionID, street,
	))
}

func (r *LocationRepository) FindByFormatted(ctx context.Context, formatted string) (*locationmodel.Location, error) {
	formatted = strings.TrimSpace(formatted)
	if formatted == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).WhereRaw(
		"formatted <> '' AND lower(btrim(formatted)) = lower(btrim(?))",
		formatted,
	))
}

func (r *LocationRepository) FindByRoundedCoordsAndFormatted(ctx context.Context, lat, lng float64, formatted string) (*locationmodel.Location, error) {
	roundedLat := math.Round(lat*1e6) / 1e6
	roundedLng := math.Round(lng*1e6) / 1e6
	return r.FindOne(ctx, query.New(1, 1).WhereRaw(
		"lat IS NOT NULL AND lng IS NOT NULL AND ROUND(lat::numeric, 6) = ROUND(?::numeric, 6) AND ROUND(lng::numeric, 6) = ROUND(?::numeric, 6) AND lower(btrim(formatted)) = lower(btrim(?))",
		roundedLat, roundedLng, strings.TrimSpace(formatted),
	))
}

func (r *LocationRepository) Create(ctx context.Context, loc *locationmodel.Location) error {
	return r.Insert(ctx, loc)
}

func (r *LocationRepository) Update(ctx context.Context, loc *locationmodel.Location) error {
	return r.Repository.Update(ctx, loc)
}

func (r *LocationRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteByID(ctx, id)
}

func (r *LocationRepository) List(ctx context.Context, filter interfaces.LocationListFilter) ([]interfaces.LocationListRow, int64, error) {
	page, limit := filter.Page, filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	base := r.DB().Builder.
		Select(
			"l.id", "l.name", "l.location_key", "l.country_id", "l.admin_division_id",
			"l.street", "l.postal_code", "l.formatted", "l.lat", "l.lng",
			"l.created_by", "l.updated_by", "l.created_at", "l.updated_at",
			"co.code AS country_code",
			"d.path AS admin_path",
		).
		From("locations l").
		LeftJoin("countries co ON co.id = l.country_id").
		LeftJoin("administrative_divisions d ON d.id = l.admin_division_id")

	if q := strings.TrimSpace(filter.Q); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		base = base.Where(squirrel.Expr(
			"(lower(l.name) LIKE ? OR lower(l.street) LIKE ? OR lower(l.formatted) LIKE ?)",
			like, like, like,
		))
	}
	if code := strings.ToUpper(strings.TrimSpace(filter.CountryCode)); code != "" {
		base = base.Where(squirrel.Eq{"co.code": code})
	}
	if id := strings.TrimSpace(filter.AdminDivisionID); id != "" {
		base = base.Where(squirrel.Eq{"l.admin_division_id": id})
	}

	countSB := r.DB().Builder.Select("COUNT(*)").FromSelect(base, "loc_list")
	total, err := r.QueryCount(ctx, countSB)
	if err != nil {
		return nil, 0, err
	}

	listSB := base.OrderBy("l.name ASC").Limit(uint64(limit)).Offset(uint64((page - 1) * limit))
	sql, args, err := listSB.ToSql()
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.DB().Querier(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]interfaces.LocationListRow, 0)
	for rows.Next() {
		var row interfaces.LocationListRow
		if err := rows.Scan(
			&row.Location.ID, &row.Location.Name, &row.Location.LocationKey,
			&row.Location.CountryID, &row.Location.AdminDivisionID,
			&row.Location.Street, &row.Location.PostalCode, &row.Location.Formatted,
			&row.Location.Lat, &row.Location.Lng,
			&row.Location.CreatedBy, &row.Location.UpdatedBy,
			&row.Location.CreatedAt, &row.Location.UpdatedAt,
			&row.CountryCode, &row.AdminPath,
		); err != nil {
			return nil, 0, err
		}
		out = append(out, row)
	}
	return out, total, rows.Err()
}

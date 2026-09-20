package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	newsmodel "be/internal/models/news"
	placemodel "be/internal/models/place"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type PlaceRepository struct {
	*repo.Repository[placemodel.Place]
}

var _ interfaces.PlaceRepository = (*PlaceRepository)(nil)

func NewPlaceRepository(db *postgres.Postgres) *PlaceRepository {
	return &PlaceRepository{
		Repository: repo.New[placemodel.Place](db, repo.Opts{Table: "places", PK: "id"}),
	}
}

func (r *PlaceRepository) GetByID(ctx context.Context, id string) (*placemodel.Place, error) {
	return r.FindByID(ctx, id)
}

func (r *PlaceRepository) FindByLocationAndKey(ctx context.Context, locationID, placeKey string) (*placemodel.Place, error) {
	locationID = strings.TrimSpace(locationID)
	placeKey = strings.TrimSpace(placeKey)
	if locationID == "" || placeKey == "" {
		return nil, nil
	}
	return r.FindOne(ctx, query.New(1, 1).
		WhereEqual("location_id", locationID).
		WhereEqual("place_key", placeKey))
}

func (r *PlaceRepository) Create(ctx context.Context, p *placemodel.Place) error {
	if p.Details == nil {
		p.Details = map[string]any{}
	}
	return r.Insert(ctx, p)
}

func (r *PlaceRepository) Update(ctx context.Context, p *placemodel.Place) error {
	if p.Details == nil {
		p.Details = map[string]any{}
	}
	return r.Repository.Update(ctx, p)
}

func (r *PlaceRepository) ListPinnable(ctx context.Context, filter interfaces.PlaceListFilter) ([]interfaces.PlaceListRow, int64, error) {
	page, limit := filter.Page, filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	base := r.DB().Builder.
		Select(
			"p.id", "p.location_id", "p.category_id", "p.place_key", "p.name", "p.unit",
			"p.status", "p.lat", "p.lng", "p.details", "p.created_by", "p.updated_by",
			"p.created_at", "p.updated_at",
			"c.key AS category_key",
			"l.id AS location_id_join",
			"l.name AS location_name",
		).
		From("places p").
		Join("categories c ON c.id = p.category_id").
		Join("locations l ON l.id = p.location_id").
		LeftJoin("administrative_divisions d ON d.id = l.admin_division_id").
		LeftJoin("countries co ON co.id = l.country_id").
		Where(squirrel.Expr(
			"p.status = ? AND p.lat IS NOT NULL AND p.lng IS NOT NULL AND p.lat BETWEEN -90 AND 90 AND p.lng BETWEEN -180 AND 180",
			placemodel.StatusActive,
		))

	if key := strings.TrimSpace(filter.CategoryKey); key != "" {
		base = base.Where(squirrel.Eq{"c.key": key})
	}
	if q := strings.TrimSpace(filter.Q); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		base = base.Where(squirrel.Expr(
			"(lower(p.name) LIKE ? OR lower(l.street) LIKE ? OR lower(l.formatted) LIKE ?)",
			like, like, like,
		))
	}
	if prefix := strings.TrimSpace(filter.AdminPathPrefix); prefix != "" {
		base = base.Where(squirrel.Expr("d.path LIKE ?", prefix+"%"))
	}
	if code := strings.ToUpper(strings.TrimSpace(filter.CountryCode)); code != "" {
		base = base.Where(squirrel.Eq{"co.code": code})
	}

	countSB := r.DB().Builder.Select("COUNT(*)").FromSelect(base, "pinnable")
	total, err := r.QueryCount(ctx, countSB)
	if err != nil {
		return nil, 0, err
	}

	listSB := base.OrderBy("p.name ASC").Limit(uint64(limit)).Offset(uint64((page - 1) * limit))
	sql, args, err := listSB.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("list pinnable: %w", err)
	}
	rows, err := r.DB().Querier(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]interfaces.PlaceListRow, 0)
	for rows.Next() {
		var row interfaces.PlaceListRow
		var locJoin string
		if err := rows.Scan(
			&row.Place.ID, &row.Place.LocationID, &row.Place.CategoryID, &row.Place.PlaceKey,
			&row.Place.Name, &row.Place.Unit, &row.Place.Status, &row.Place.Lat, &row.Place.Lng,
			&row.Place.Details, &row.Place.CreatedBy, &row.Place.UpdatedBy,
			&row.Place.CreatedAt, &row.Place.UpdatedAt,
			&row.CategoryKey, &locJoin, &row.LocationName,
		); err != nil {
			return nil, 0, err
		}
		row.LocationID = row.Place.LocationID
		out = append(out, row)
	}
	return out, total, rows.Err()
}

func (r *PlaceRepository) ListAdmin(ctx context.Context, filter interfaces.PlaceAdminListFilter) ([]interfaces.PlaceListRow, int64, error) {
	page, limit := filter.Page, filter.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}

	base := r.DB().Builder.
		Select(
			"p.id", "p.location_id", "p.category_id", "p.place_key", "p.name", "p.unit",
			"p.status", "p.lat", "p.lng", "p.details", "p.created_by", "p.updated_by",
			"p.created_at", "p.updated_at",
			"c.key AS category_key",
			"l.id AS location_id_join",
			"l.name AS location_name",
		).
		From("places p").
		Join("categories c ON c.id = p.category_id").
		Join("locations l ON l.id = p.location_id")

	if id := strings.TrimSpace(filter.LocationID); id != "" {
		base = base.Where(squirrel.Eq{"p.location_id": id})
	}
	if q := strings.TrimSpace(filter.Q); q != "" {
		like := "%" + strings.ToLower(q) + "%"
		base = base.Where(squirrel.Expr(
			"(lower(p.name) LIKE ? OR lower(l.name) LIKE ?)",
			like, like,
		))
	}
	if filter.Status != nil {
		base = base.Where(squirrel.Eq{"p.status": *filter.Status})
	}

	countSB := r.DB().Builder.Select("COUNT(*)").FromSelect(base, "place_admin")
	total, err := r.QueryCount(ctx, countSB)
	if err != nil {
		return nil, 0, err
	}

	listSB := base.OrderBy("p.name ASC").Limit(uint64(limit)).Offset(uint64((page - 1) * limit))
	sql, args, err := listSB.ToSql()
	if err != nil {
		return nil, 0, err
	}
	rows, err := r.DB().Querier(ctx).Query(ctx, sql, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := make([]interfaces.PlaceListRow, 0)
	for rows.Next() {
		var row interfaces.PlaceListRow
		var locJoin string
		if err := rows.Scan(
			&row.Place.ID, &row.Place.LocationID, &row.Place.CategoryID, &row.Place.PlaceKey,
			&row.Place.Name, &row.Place.Unit, &row.Place.Status, &row.Place.Lat, &row.Place.Lng,
			&row.Place.Details, &row.Place.CreatedBy, &row.Place.UpdatedBy,
			&row.Place.CreatedAt, &row.Place.UpdatedAt,
			&row.CategoryKey, &locJoin, &row.LocationName,
		); err != nil {
			return nil, 0, err
		}
		row.LocationID = row.Place.LocationID
		out = append(out, row)
	}
	return out, total, rows.Err()
}

func (r *PlaceRepository) CountByLocationID(ctx context.Context, locationID string) (int64, error) {
	locationID = strings.TrimSpace(locationID)
	if locationID == "" {
		return 0, nil
	}
	return r.Count(ctx, query.New(1, 1).WhereEqual("location_id", locationID))
}

func (r *PlaceRepository) Delete(ctx context.Context, id string) error {
	return r.DeleteByID(ctx, id)
}

func (r *PlaceRepository) GetDetail(ctx context.Context, id string) (*interfaces.PlaceDetailRow, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil
	}

	sql := `
SELECT
	p.id, p.location_id, p.category_id, p.place_key, p.name, p.unit, p.status,
	p.lat, p.lng, p.details, p.created_by, p.updated_by, p.created_at, p.updated_at,
	c.key,
	l.id, l.name, l.location_key, l.country_id, l.admin_division_id,
	l.street, l.postal_code, l.formatted, l.lat, l.lng,
	co.code, d.path
FROM places p
JOIN categories c ON c.id = p.category_id
JOIN locations l ON l.id = p.location_id
LEFT JOIN countries co ON co.id = l.country_id
LEFT JOIN administrative_divisions d ON d.id = l.admin_division_id
WHERE p.id = $1`

	var detail interfaces.PlaceDetailRow
	err := r.DB().Querier(ctx).QueryRow(ctx, sql, id).Scan(
		&detail.Place.ID, &detail.Place.LocationID, &detail.Place.CategoryID, &detail.Place.PlaceKey,
		&detail.Place.Name, &detail.Place.Unit, &detail.Place.Status, &detail.Place.Lat, &detail.Place.Lng,
		&detail.Place.Details, &detail.Place.CreatedBy, &detail.Place.UpdatedBy,
		&detail.Place.CreatedAt, &detail.Place.UpdatedAt,
		&detail.CategoryKey,
		&detail.Location.ID, &detail.Location.Name, &detail.Location.LocationKey,
		&detail.Location.CountryID, &detail.Location.AdminDivisionID,
		&detail.Location.Street, &detail.Location.PostalCode, &detail.Location.Formatted,
		&detail.Location.Lat, &detail.Location.Lng,
		&detail.Location.CountryCode, &detail.Location.AdminPath,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	newsRepo := NewNewsRepository(r.DB())
	items, err := newsRepo.ListByPlaceID(ctx, id)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []newsmodel.News{}
	}
	detail.News = items
	detail.NewsCategoryKey = map[string]string{}
	for _, item := range items {
		if item.CategoryID == detail.Place.CategoryID {
			detail.NewsCategoryKey[item.ID] = detail.CategoryKey
			continue
		}
		catRepo := NewCategoryRepository(r.DB())
		cat, catErr := catRepo.GetByID(ctx, item.CategoryID)
		if catErr != nil {
			return nil, catErr
		}
		if cat != nil {
			detail.NewsCategoryKey[item.ID] = cat.Key
		} else {
			detail.NewsCategoryKey[item.ID] = detail.CategoryKey
		}
	}
	return &detail, nil
}

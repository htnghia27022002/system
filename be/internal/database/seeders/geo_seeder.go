package seeders

import (
	"context"
	"encoding/json"
	"fmt"

	"be/internal/database/seeders/fixtures"
	"be/pkg/postgres"
)

// GeoSeeder upserts VN and the 2025 two-level administrative fixture.
type GeoSeeder struct{}

func NewGeoSeeder() *GeoSeeder {
	return &GeoSeeder{}
}

func (s *GeoSeeder) Name() string {
	return "GeoSeeder"
}

type vnFixture struct {
	Country struct {
		Code      string `json:"code"`
		Code3     string `json:"code3"`
		Name      string `json:"name"`
		NameLocal string `json:"nameLocal"`
		PhoneCode string `json:"phoneCode"`
		Currency  string `json:"currency"`
	} `json:"country"`
	Divisions []vnDivision `json:"divisions"`
}

type vnDivision struct {
	Code       string   `json:"code"`
	ParentCode string   `json:"parentCode"`
	Name       string   `json:"name"`
	NameEn     string   `json:"nameEn"`
	FullName   string   `json:"fullName"`
	Type       string   `json:"type"`
	Level      int16    `json:"level"`
	Lat        *float64 `json:"lat"`
	Lng        *float64 `json:"lng"`
}

func (s *GeoSeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	var fixture vnFixture
	if err := json.Unmarshal(fixtures.VietnamDivisionsJSON, &fixture); err != nil {
		return fmt.Errorf("decode vn divisions fixture: %w", err)
	}

	countrySQL, countryArgs, err := db.Builder.
		Insert(postgres.QuoteIdent("countries")).
		Columns("code", "code3", "name", "name_local", "phone_code", "currency", "is_active").
		Values(fixture.Country.Code, fixture.Country.Code3, fixture.Country.Name, fixture.Country.NameLocal, fixture.Country.PhoneCode, fixture.Country.Currency, true).
		Suffix("ON CONFLICT (code) DO UPDATE SET name = EXCLUDED.name, name_local = EXCLUDED.name_local, code3 = EXCLUDED.code3, updated_at = NOW()").
		ToSql()
	if err != nil {
		return err
	}
	if _, err := db.Querier(ctx).Exec(ctx, countrySQL, countryArgs...); err != nil {
		return err
	}

	var countryID string
	if err := db.Querier(ctx).QueryRow(ctx, `SELECT id FROM countries WHERE code = $1`, fixture.Country.Code).Scan(&countryID); err != nil {
		return err
	}

	byCode := map[string]string{}
	for _, div := range fixture.Divisions {
		if div.Level == 1 {
			path := fixture.Country.Code + "/" + div.Code + "/"
			id, err := upsertDivision(ctx, db, countryID, nil, div, path)
			if err != nil {
				return err
			}
			byCode[div.Code] = id
		}
	}
	for _, div := range fixture.Divisions {
		if div.Level != 2 {
			continue
		}
		parentID, ok := byCode[div.ParentCode]
		if !ok {
			return fmt.Errorf("missing parent province %q for commune %q", div.ParentCode, div.Code)
		}
		path := fixture.Country.Code + "/" + div.ParentCode + "/" + div.Code + "/"
		if _, err := upsertDivision(ctx, db, countryID, &parentID, div, path); err != nil {
			return err
		}
	}
	return nil
}

func upsertDivision(ctx context.Context, db *postgres.Postgres, countryID string, parentID *string, div vnDivision, path string) (string, error) {
	sql, args, err := db.Builder.
		Insert(postgres.QuoteIdent("administrative_divisions")).
		Columns("country_id", "parent_id", "level", "code", "name", "name_en", "full_name", "type", "path", "lat", "lng", "is_active").
		Values(countryID, parentID, div.Level, div.Code, div.Name, div.NameEn, div.FullName, div.Type, path, div.Lat, div.Lng, true).
		Suffix(`ON CONFLICT (country_id, code) DO UPDATE SET
			name = EXCLUDED.name,
			name_en = EXCLUDED.name_en,
			full_name = EXCLUDED.full_name,
			type = EXCLUDED.type,
			path = EXCLUDED.path,
			parent_id = EXCLUDED.parent_id,
			level = EXCLUDED.level,
			lat = EXCLUDED.lat,
			lng = EXCLUDED.lng,
			updated_at = NOW()
		RETURNING id`).
		ToSql()
	if err != nil {
		return "", err
	}
	var id string
	if err := db.Querier(ctx).QueryRow(ctx, sql, args...).Scan(&id); err != nil {
		return "", err
	}
	return id, nil
}

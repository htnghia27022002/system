package seeders

import (
	"context"

	"be/pkg/postgres"
)

// CategorySeeder upserts P1 Maps categories by unique key.
type CategorySeeder struct{}

func NewCategorySeeder() *CategorySeeder {
	return &CategorySeeder{}
}

func (s *CategorySeeder) Name() string {
	return "CategorySeeder"
}

func (s *CategorySeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	rows := []struct {
		Key       string
		Name      string
		NameLocal string
		Sort      int
	}{
		{"room_rental", "Room rental", "Cho thuê phòng", 1},
		{"restaurant", "Restaurant", "Nhà hàng", 2},
		{"hotel", "Hotel", "Khách sạn", 3},
		{"eatery", "Eatery", "Quán ăn", 4},
		{"uncategorized", "Uncategorized", "Chưa phân loại", 99},
	}
	for _, row := range rows {
		sql, args, err := db.Builder.
			Insert(postgres.QuoteIdent("categories")).
			Columns("key", "name", "name_local", "sort_order", "is_active").
			Values(row.Key, row.Name, row.NameLocal, row.Sort, true).
			Suffix("ON CONFLICT (key) DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}
		if _, err := db.Querier(ctx).Exec(ctx, sql, args...); err != nil {
			return err
		}
	}
	return nil
}

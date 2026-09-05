package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// QueryStrings runs sql and scans a single text column into []string.
func QueryStrings(ctx context.Context, db DBTX, sql string, args ...any) ([]string, error) {
	rows, err := db.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []string{}
	}
	return items, nil
}

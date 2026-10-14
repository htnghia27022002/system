package repository

import (
	"context"
	"encoding/json"

	"be/internal/settings"
	"be/pkg/postgres"
)

// SettingRepository stores settings overrides in app_settings.
type SettingRepository struct {
	db *postgres.Postgres
}

var _ settings.Store = (*SettingRepository)(nil)

func NewSettingRepository(db *postgres.Postgres) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) All(ctx context.Context) (map[string]json.RawMessage, error) {
	rows, err := r.db.Querier(ctx).Query(ctx, `SELECT key, value FROM app_settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]json.RawMessage{}
	for rows.Next() {
		var key string
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		out[key] = value
	}
	return out, rows.Err()
}

func (r *SettingRepository) Upsert(ctx context.Context, values map[string]json.RawMessage, updatedBy string) error {
	var actor any
	if updatedBy != "" {
		actor = updatedBy
	}
	return r.db.WithTx(ctx, func(ctx context.Context) error {
		for key, value := range values {
			if _, err := r.db.Querier(ctx).Exec(ctx, `
INSERT INTO app_settings (key, value, updated_by, updated_at) VALUES ($1, $2, $3, NOW())
ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_by = EXCLUDED.updated_by, updated_at = NOW()`,
				key, []byte(value), actor); err != nil {
				return err
			}
		}
		return nil
	})
}

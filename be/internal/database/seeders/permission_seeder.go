package seeders

import (
	"context"

	"be/pkg/postgres"
)

// PermissionSeeder seeds the permissions catalog (core + module permissions).
type PermissionSeeder struct {
	permissions []PermissionDef
}

// NewPermissionSeeder seeds the core catalog plus any extra module permissions.
func NewPermissionSeeder(extra ...PermissionDef) *PermissionSeeder {
	return &PermissionSeeder{permissions: append(DefaultPermissions(), extra...)}
}

func (s *PermissionSeeder) Name() string {
	return "PermissionSeeder"
}

func (s *PermissionSeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	for _, item := range s.permissions {
		sql, args, err := db.Builder.
			Insert(postgres.QuoteIdent("permissions")).
			Columns("id", "key", "name", `"group"`, "description").
			Values(item.ID, item.Key, item.Name, item.Group, item.Description).
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

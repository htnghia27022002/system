package seeders

import (
	"context"

	"be/pkg/postgres"
)

// RoleSeeder seeds default roles.
type RoleSeeder struct{}

func NewRoleSeeder() *RoleSeeder {
	return &RoleSeeder{}
}

func (s *RoleSeeder) Name() string {
	return "RoleSeeder"
}

func (s *RoleSeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	roles := []struct {
		ID          string
		Name        string
		Slug        string
		Description string
	}{
		{RoleAdminID, "Administrator", "admin", "Full system access"},
		{RoleUserID, "Member", "user", "Standard member access"},
	}

	for _, role := range roles {
		sql, args, err := db.Builder.
			Insert(postgres.QuoteIdent("roles")).
			Columns("id", "name", "slug", "description").
			Values(role.ID, role.Name, role.Slug, role.Description).
			Suffix("ON CONFLICT (slug) DO NOTHING").
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

package seeders

import (
	"context"

	"be/internal/common/rbac"
	"be/pkg/postgres"
)

// RolePermissionSeeder assigns permissions to roles.
type RolePermissionSeeder struct{}

func NewRolePermissionSeeder() *RolePermissionSeeder {
	return &RolePermissionSeeder{}
}

func (s *RolePermissionSeeder) Name() string {
	return "RolePermissionSeeder"
}

func (s *RolePermissionSeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	insert := func(roleID, permissionID string) error {
		sql, args, err := db.Builder.
			Insert(postgres.QuoteIdent("role_permissions")).
			Columns("role_id", "permission_id").
			Values(roleID, permissionID).
			Suffix("ON CONFLICT DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}
		_, err = db.Querier(ctx).Exec(ctx, sql, args...)
		return err
	}

	for _, item := range DefaultPermissions() {
		if err := insert(RoleAdminID, item.ID); err != nil {
			return err
		}
	}

	return insert(RoleUserID, PermissionIDByKey(rbac.Key("dashboard", rbac.ActionView)))
}


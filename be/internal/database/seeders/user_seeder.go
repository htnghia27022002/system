package seeders

import (
	"context"

	"be/pkg/hash"
	usermodel "be/internal/models/user"
	"be/pkg/postgres"
)

// UserSeeder seeds demo admin and member accounts.
type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

func (s *UserSeeder) Run(ctx context.Context, db *postgres.Postgres) error {
	adminHash, err := hash.HashPassword("admin1234")
	if err != nil {
		return err
	}
	demoHash, err := hash.HashPassword("password123")
	if err != nil {
		return err
	}

	users := []struct {
		ID           string
		Email        string
		PasswordHash string
		FullName     string
		RoleID       string
		Status       usermodel.Status
		SuperAdmin   bool
	}{
		{AdminUserID, "admin@example.com", adminHash, "Admin User", RoleAdminID, usermodel.StatusActive, true},
		{DemoUserID, "demo@example.com", demoHash, "Demo User", RoleUserID, usermodel.StatusActive, false},
	}

	for _, user := range users {
		sql, args, err := db.Builder.
			Insert(postgres.QuoteIdent("users")).
			Columns("id", "email", "password_hash", "full_name", "role_id", "status", "is_super_admin").
			Values(user.ID, user.Email, user.PasswordHash, user.FullName, user.RoleID, string(user.Status), user.SuperAdmin).
			Suffix("ON CONFLICT (email) DO NOTHING").
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

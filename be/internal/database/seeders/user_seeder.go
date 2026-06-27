package seeders

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"be/internal/common/hash"
	usermodel "be/internal/models/user"
)

// UserSeeder seeds demo admin and member accounts.
type UserSeeder struct{}

func NewUserSeeder() *UserSeeder {
	return &UserSeeder{}
}

func (s *UserSeeder) Name() string {
	return "UserSeeder"
}

func (s *UserSeeder) Run(ctx context.Context, db *gorm.DB) error {
	adminHash, err := hash.HashPassword("admin1234")
	if err != nil {
		return err
	}
	demoHash, err := hash.HashPassword("password123")
	if err != nil {
		return err
	}

	users := []usermodel.User{
		{
			ID:           AdminUserID,
			Email:        "admin@example.com",
			PasswordHash: adminHash,
			FullName:     "Admin User",
			RoleID:       RoleAdminID,
			Status:       usermodel.StatusActive,
		},
		{
			ID:           DemoUserID,
			Email:        "demo@example.com",
			PasswordHash: demoHash,
			FullName:     "Demo User",
			RoleID:       RoleUserID,
			Status:       usermodel.StatusActive,
		},
	}

	for _, user := range users {
		if err := db.WithContext(ctx).Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoNothing: true,
		}).Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

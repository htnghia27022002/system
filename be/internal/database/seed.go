package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"be/internal/common/hash"
	authmodel "be/internal/models/auth"
	permissionmodel "be/internal/models/permission"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
)

const (
	RoleAdminID = "11111111-1111-4111-8111-111111111111"
	RoleUserID  = "22222222-2222-4222-8222-222222222222"
	AdminUserID = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	DemoUserID  = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
)

type seedPermission struct {
	ID          string
	Key         string
	Name        string
	Group       string
	Description string
}

var seedPermissions = []seedPermission{
	{ID: "10000001-0000-4000-8000-000000000001", Key: "dashboard:read", Name: "View dashboard", Group: "dashboard", Description: "Access the admin dashboard overview"},
	{ID: "10000001-0000-4000-8000-000000000002", Key: "users:read", Name: "View users", Group: "users", Description: "List and view user accounts"},
	{ID: "10000001-0000-4000-8000-000000000003", Key: "users:create", Name: "Create users", Group: "users", Description: "Create new user accounts"},
	{ID: "10000001-0000-4000-8000-000000000004", Key: "users:update", Name: "Update users", Group: "users", Description: "Edit existing user accounts"},
	{ID: "10000001-0000-4000-8000-000000000005", Key: "users:delete", Name: "Delete users", Group: "users", Description: "Remove user accounts"},
	{ID: "10000001-0000-4000-8000-000000000006", Key: "roles:read", Name: "View roles", Group: "roles", Description: "List and view roles"},
	{ID: "10000001-0000-4000-8000-000000000007", Key: "roles:create", Name: "Create roles", Group: "roles", Description: "Create new roles"},
	{ID: "10000001-0000-4000-8000-000000000008", Key: "roles:update", Name: "Update roles", Group: "roles", Description: "Edit existing roles"},
	{ID: "10000001-0000-4000-8000-000000000009", Key: "roles:delete", Name: "Delete roles", Group: "roles", Description: "Remove roles"},
	{ID: "10000001-0000-4000-8000-000000000010", Key: "permissions:read", Name: "View permissions", Group: "permissions", Description: "View the permissions catalog"},
}

// SeedRBAC inserts default roles, permissions, and demo users when the database is empty.
func SeedRBAC(ctx context.Context, db *gorm.DB) error {
	var count int64
	if err := db.WithContext(ctx).Model(&rolemodel.Role{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count roles: %w", err)
	}
	if count > 0 {
		return nil
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range seedPermissions {
			perm := permissionmodel.Permission{
				ID:          item.ID,
				Key:         item.Key,
				Name:        item.Name,
				Group:       item.Group,
				Description: item.Description,
			}
			if err := tx.Create(&perm).Error; err != nil {
				return fmt.Errorf("seed permission %s: %w", item.Key, err)
			}
		}

		adminRole := rolemodel.Role{
			ID:          RoleAdminID,
			Name:        "Administrator",
			Slug:        "admin",
			Description: "Full system access",
		}
		if err := tx.Create(&adminRole).Error; err != nil {
			return fmt.Errorf("seed admin role: %w", err)
		}

		memberRole := rolemodel.Role{
			ID:          RoleUserID,
			Name:        "Member",
			Slug:        "user",
			Description: "Standard member access",
		}
		if err := tx.Create(&memberRole).Error; err != nil {
			return fmt.Errorf("seed member role: %w", err)
		}

		for _, item := range seedPermissions {
			if err := tx.Create(&rolemodel.RolePermission{
				RoleID:       RoleAdminID,
				PermissionID: item.ID,
			}).Error; err != nil {
				return fmt.Errorf("seed admin permission %s: %w", item.Key, err)
			}
		}

		if err := tx.Create(&rolemodel.RolePermission{
			RoleID:       RoleUserID,
			PermissionID: "10000001-0000-4000-8000-000000000001",
		}).Error; err != nil {
			return fmt.Errorf("seed member permission: %w", err)
		}

		adminHash, err := hash.HashPassword("admin1234")
		if err != nil {
			return fmt.Errorf("hash admin password: %w", err)
		}
		demoHash, err := hash.HashPassword("password123")
		if err != nil {
			return fmt.Errorf("hash demo password: %w", err)
		}

		adminUser := usermodel.User{
			ID:           AdminUserID,
			Email:        "admin@example.com",
			PasswordHash: adminHash,
			FullName:     "Admin User",
			RoleID:       RoleAdminID,
			Status:       usermodel.StatusActive,
		}
		if err := tx.Create(&adminUser).Error; err != nil {
			return fmt.Errorf("seed admin user: %w", err)
		}

		demoUser := usermodel.User{
			ID:           DemoUserID,
			Email:        "demo@example.com",
			PasswordHash: demoHash,
			FullName:     "Demo User",
			RoleID:       RoleUserID,
			Status:       usermodel.StatusActive,
		}
		if err := tx.Create(&demoUser).Error; err != nil {
			return fmt.Errorf("seed demo user: %w", err)
		}

		_ = authmodel.RefreshToken{}
		return nil
	})
}

//go:build integration

package integration_test

import (
	"context"
	"testing"

	"be/internal/common/hash"
	"be/internal/database"
	"be/internal/repository"
	usermodel "be/internal/models/user"
	"be/test/testutil"
)

func TestAuthRepositoryFindUserByEmail(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		_ = db.Exec("TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts RESTART IDENTITY CASCADE").Error
	})

	ctx := context.Background()
	if err := database.SeedRBAC(ctx, db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	authRepo := repository.NewAuthRepository(db)
	user, err := authRepo.FindUserByEmail(ctx, "admin@example.com")
	if err != nil {
		t.Fatalf("find user: %v", err)
	}
	if user == nil {
		t.Fatal("expected seeded admin user")
	}
	if user.Status != usermodel.StatusActive {
		t.Fatalf("expected active status, got %s", user.Status)
	}
}

func TestAuthServiceLoginWithSeededAdmin(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		_ = db.Exec("TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts RESTART IDENTITY CASCADE").Error
	})

	ctx := context.Background()
	if err := database.SeedRBAC(ctx, db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	container := testutil.NewTestContainer(t, db)
	result, err := container.AuthService.Login(ctx, "admin@example.com", "admin1234")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if result.AccessToken == "" {
		t.Fatal("expected access token")
	}
}

func TestUserRepositoryCreateAndFind(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		_ = db.Exec("TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts RESTART IDENTITY CASCADE").Error
	})

	ctx := context.Background()
	if err := database.SeedRBAC(ctx, db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	roleRepo := repository.NewRoleRepository(db)
	memberRole, err := roleRepo.GetBySlug(ctx, "user")
	if err != nil || memberRole == nil {
		t.Fatalf("member role: %v", err)
	}

	hashed, err := hash.HashPassword("password123")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	user := &usermodel.User{
		Email:        "integration@example.com",
		PasswordHash: hashed,
		FullName:     "Integration User",
		RoleID:       memberRole.ID,
		Status:       usermodel.StatusActive,
	}
	if err := userRepo.Create(ctx, user); err != nil {
		t.Fatalf("create: %v", err)
	}

	found, err := userRepo.GetByEmail(ctx, "integration@example.com")
	if err != nil {
		t.Fatalf("get by email: %v", err)
	}
	if found == nil || found.FullName != "Integration User" {
		t.Fatal("expected created user")
	}
}

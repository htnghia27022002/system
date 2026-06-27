package auth_test

import (
	"context"
	"testing"
	"time"

	"be/internal/common/hash"
	jwtmanager "be/internal/common/jwt"
	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	authsvc "be/internal/services/auth"
	"be/test/testutil"
)

func TestLoginIssuesTokensWithRoleClaims(t *testing.T) {
	hashed, err := hash.HashPassword("admin1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	roleID := "role-admin"
	authRepo := &testutil.MockAuthRepo{
		Users: map[string]*usermodel.User{
			"user-1": {
				ID:           "user-1",
				Email:        "admin@example.com",
				PasswordHash: hashed,
				FullName:     "Admin User",
				RoleID:       roleID,
				Status:       usermodel.StatusActive,
			},
		},
		RefreshTokens: map[string]*authmodel.RefreshToken{},
	}
	roleRepo := &testutil.MockRoleRepo{
		Roles: map[string]*rolemodel.Role{
			roleID: {ID: roleID, Slug: "admin", Name: "Administrator"},
		},
		Permissions: map[string][]string{
			roleID: {"users:view", "roles:view"},
		},
	}

	cfg := testutil.UnitConfig()
	jwtManager := jwtmanager.NewManager(cfg)
	svc := authsvc.NewService(authRepo, &testutil.MockUserRepo{}, roleRepo, jwtManager, 24*time.Hour)

	result, err := svc.Login(context.Background(), "admin@example.com", "admin1234")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatal("expected token pair")
	}
	if result.User.Role != "admin" {
		t.Fatalf("expected admin role, got %s", result.User.Role)
	}
	if len(result.User.Permissions) != 2 {
		t.Fatalf("expected permissions in auth user, got %v", result.User.Permissions)
	}

	claims, err := jwtManager.VerifyAccessToken(result.AccessToken)
	if err != nil {
		t.Fatalf("verify access token: %v", err)
	}
	if claims.RoleID != roleID {
		t.Fatalf("expected roleId claim %s, got %s", roleID, claims.RoleID)
	}
}

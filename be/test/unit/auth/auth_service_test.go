package auth_test

import (
	"context"
	"testing"
	"time"

	jwtmanager "be/common/jwt"
	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	authsvc "be/internal/services/auth"
	"be/pkg/hash"
	"be/test/testutil"
)

func TestLoginIssuesTokensWithRoleClaims(t *testing.T) {
	hashed, err := hash.HashPassword("admin1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}

	roleID := "role-admin"
	user := &usermodel.User{
		ID:           "user-1",
		Email:        "admin@example.com",
		PasswordHash: hashed,
		FullName:     "Admin User",
		RoleID:       roleID,
		Status:       usermodel.StatusActive,
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"user-1": user}}
	authRepo := &testutil.MockAuthRepo{
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
	svc := authsvc.NewService(authRepo, userRepo, roleRepo, jwtManager, 24*time.Hour, nil)

	result, err := svc.Login(context.Background(), "admin@example.com", "admin1234", "203.0.113.10", "Mozilla/5.0")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	if result.AccessToken == "" || result.RefreshToken == "" {
		t.Fatal("expected token pair")
	}
	if result.User.Role != "admin" {
		t.Fatalf("expected admin role, got %s", result.User.Role)
	}
	if result.SessionID == "" {
		t.Fatal("expected sessionId on login")
	}
	if result.RefreshToken == "" {
		t.Fatal("expected refresh token")
	}
	if len(result.User.Permissions) != 2 {
		t.Fatalf("expected permissions in auth user, got %v", result.User.Permissions)
	}
	if result.User.SuperAdmin {
		t.Fatal("expected superAdmin false when user flag is unset")
	}

	stored, ok := authRepo.RefreshTokens[hash.SHA256Hex(result.RefreshToken)]
	if !ok || stored == nil {
		t.Fatal("expected persisted refresh token")
	}
	if stored.ID != result.SessionID {
		t.Fatalf("expected sessionId %s, got %s", stored.ID, result.SessionID)
	}
	if stored.IPAddress != "203.0.113.10" || stored.UserAgent != "Mozilla/5.0" {
		t.Fatalf("expected ip/ua on token, got %q %q", stored.IPAddress, stored.UserAgent)
	}
	if stored.LastUsedAt.IsZero() {
		t.Fatal("expected last_used_at on create")
	}
	if stored.CreatedAt.IsZero() || stored.CreatedAt.Unix() <= 0 {
		t.Fatalf("expected real created_at, got %v", stored.CreatedAt)
	}
	if stored.DeviceID != "" {
		t.Fatalf("device_id must stay unused in P1, got %q", stored.DeviceID)
	}

	claims, err := jwtManager.VerifyAccessToken(result.AccessToken)
	if err != nil {
		t.Fatalf("verify access token: %v", err)
	}
	if claims.RoleID != roleID {
		t.Fatalf("expected roleId claim %s, got %s", roleID, claims.RoleID)
	}
}

func TestLogoutRevokesRefreshTokenSoItIsNoLongerActive(t *testing.T) {
	hashed, err := hash.HashPassword("admin1234")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	roleID := "role-admin"
	user := &usermodel.User{
		ID:           "user-1",
		Email:        "admin@example.com",
		PasswordHash: hashed,
		FullName:     "Admin User",
		RoleID:       roleID,
		Status:       usermodel.StatusActive,
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"user-1": user}}
	authRepo := &testutil.MockAuthRepo{RefreshTokens: map[string]*authmodel.RefreshToken{}}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{roleID: {ID: roleID, Slug: "admin", Name: "Administrator"}},
		Permissions: map[string][]string{roleID: {}},
	}
	svc := authsvc.NewService(authRepo, userRepo, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), 24*time.Hour, nil)

	result, err := svc.Login(context.Background(), "admin@example.com", "admin1234", "10.0.0.1", "UA")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	listed, err := svc.ListSessions(context.Background(), user.ID, result.SessionID)
	if err != nil {
		t.Fatalf("list before logout: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Fatalf("expected 1 active session, got %d", len(listed.Items))
	}

	if err := svc.Logout(context.Background(), result.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}

	listed, err = svc.ListSessions(context.Background(), user.ID, result.SessionID)
	if err != nil {
		t.Fatalf("list after logout: %v", err)
	}
	if len(listed.Items) != 0 {
		t.Fatalf("expected no active sessions after logout, got %d", len(listed.Items))
	}
}

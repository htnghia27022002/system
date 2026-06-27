package jwt_test

import (
	"testing"
	"time"

	jwtmanager "be/internal/common/jwt"
	"be/test/testutil"
)

func TestSignAndVerifyAccessToken(t *testing.T) {
	cfg := testutil.UnitConfig()
	cfg.JWTAccessTTL = time.Minute
	manager := jwtmanager.NewManager(cfg)

	token, err := manager.SignAccessToken(
		"user-1",
		"demo@example.com",
		"Demo User",
		"user",
		"role-user",
		[]string{"dashboard:view"},
	)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	claims, err := manager.VerifyAccessToken(token)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if claims.Subject != "user-1" {
		t.Fatalf("expected sub user-1, got %s", claims.Subject)
	}
	if claims.RoleID != "role-user" {
		t.Fatalf("expected roleId role-user, got %s", claims.RoleID)
	}
	if len(claims.Permissions) != 1 || claims.Permissions[0] != "dashboard:view" {
		t.Fatalf("unexpected permissions: %v", claims.Permissions)
	}
}

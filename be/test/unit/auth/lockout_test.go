package auth_test

import (
	"context"
	"errors"
	"testing"
	"time"

	apperrors "be/common/errors"
	jwtmanager "be/common/jwt"
	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	authsvc "be/internal/services/auth"
	"be/pkg/hash"
	"be/test/testutil"
)

func newLockoutService(t *testing.T, attempts *testutil.MockLoginAttemptRepo) *authsvc.Service {
	t.Helper()
	hashed, err := hash.HashPassword("correct-horse")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{"u1": {
		ID: "u1", Email: "u@example.com", PasswordHash: hashed, RoleID: "r1", Status: usermodel.StatusActive,
	}}}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{"r1": {ID: "r1", Slug: "user", Name: "User"}},
		Permissions: map[string][]string{},
	}
	authRepo := &testutil.MockAuthRepo{RefreshTokens: map[string]*authmodel.RefreshToken{}}
	svc := authsvc.NewService(authRepo, userRepo, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), time.Hour, nil)
	svc.SetLoginLockout(attempts, 3, 15*time.Minute)
	return svc
}

func TestLoginLocksAccountAfterRepeatedFailures(t *testing.T) {
	attempts := testutil.NewMockLoginAttemptRepo()
	svc := newLockoutService(t, attempts)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		_, err := svc.Login(ctx, "u@example.com", "wrong", "", "")
		if !errors.Is(err, apperrors.ErrUnauthorized) {
			t.Fatalf("attempt %d: expected unauthorized, got %v", i+1, err)
		}
	}

	// Locked: even the right password is refused, with a 429-class error.
	_, err := svc.Login(ctx, "u@example.com", "correct-horse", "", "")
	if !errors.Is(err, authsvc.ErrAccountLocked) || !errors.Is(err, apperrors.ErrTooManyRequests) {
		t.Fatalf("expected account locked, got %v", err)
	}

	// After the lock expires the right password works and resets the counter.
	attempts.Locked["u1"] = time.Now().Add(-time.Second)
	if _, err := svc.Login(ctx, "u@example.com", "correct-horse", "", ""); err != nil {
		t.Fatalf("login after lock expiry: %v", err)
	}
	if _, locked := attempts.Locked["u1"]; locked || attempts.Failed["u1"] != 0 {
		t.Fatalf("successful login must reset lockout state: %+v", attempts)
	}
}

func TestSuccessfulLoginResetsFailureCount(t *testing.T) {
	attempts := testutil.NewMockLoginAttemptRepo()
	svc := newLockoutService(t, attempts)
	ctx := context.Background()

	for round := 0; round < 3; round++ {
		for i := 0; i < 2; i++ {
			_, _ = svc.Login(ctx, "u@example.com", "wrong", "", "")
		}
		if _, err := svc.Login(ctx, "u@example.com", "correct-horse", "", ""); err != nil {
			t.Fatalf("round %d: failures below the limit must not lock: %v", round, err)
		}
	}
}

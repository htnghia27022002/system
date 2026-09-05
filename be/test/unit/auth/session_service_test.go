package auth_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	apperrors "be/internal/common/errors"
	jwtmanager "be/internal/common/jwt"
	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	authsvc "be/internal/services/auth"
	"be/pkg/hash"
	"be/test/testutil"
)

type sessionFixture struct {
	svc      *authsvc.Service
	authRepo *testutil.MockAuthRepo
	userA    *usermodel.User
	userB    *usermodel.User
}

func newSessionFixture(t *testing.T) sessionFixture {
	t.Helper()
	hashed, err := hash.HashPassword("password12")
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	roleID := "role-user"
	userA := &usermodel.User{
		ID:           "user-a",
		Email:        "a@example.com",
		PasswordHash: hashed,
		FullName:     "User A",
		RoleID:       roleID,
		Status:       usermodel.StatusActive,
	}
	userB := &usermodel.User{
		ID:           "user-b",
		Email:        "b@example.com",
		PasswordHash: hashed,
		FullName:     "User B",
		RoleID:       roleID,
		Status:       usermodel.StatusActive,
	}
	authRepo := &testutil.MockAuthRepo{RefreshTokens: map[string]*authmodel.RefreshToken{}}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{
		userA.ID: userA,
		userB.ID: userB,
	}}
	roleRepo := &testutil.MockRoleRepo{
		Roles:       map[string]*rolemodel.Role{roleID: {ID: roleID, Slug: "user", Name: "User"}},
		Permissions: map[string][]string{roleID: {}},
	}
	svc := authsvc.NewService(authRepo, userRepo, roleRepo, jwtmanager.NewManager(testutil.UnitConfig()), 24*time.Hour, nil)
	return sessionFixture{svc: svc, authRepo: authRepo, userA: userA, userB: userB}
}

func putToken(repo *testutil.MockAuthRepo, token *authmodel.RefreshToken) {
	if token.TokenHash == "" {
		token.TokenHash = token.ID + "-hash"
	}
	_ = repo.CreateRefreshToken(context.Background(), token)
	repo.RefreshTokens[token.TokenHash] = token
}

func TestListSessionsReturnsOnlyActiveOwnedRows(t *testing.T) {
	fx := newSessionFixture(t)
	now := time.Now()
	revokedAt := now.Add(-time.Hour)

	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID:         "active-a",
		UserID:     fx.userA.ID,
		ExpiresAt:  now.Add(time.Hour),
		CreatedAt:  now.Add(-2 * time.Hour),
		LastUsedAt: now.Add(-time.Minute),
		IPAddress:  "203.0.113.10",
		UserAgent:  "Mozilla/5.0",
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID:         "revoked-a",
		UserID:     fx.userA.ID,
		ExpiresAt:  now.Add(time.Hour),
		RevokedAt:  &revokedAt,
		CreatedAt:  now.Add(-3 * time.Hour),
		LastUsedAt: now.Add(-2 * time.Hour),
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID:         "expired-a",
		UserID:     fx.userA.ID,
		ExpiresAt:  now.Add(-time.Minute),
		CreatedAt:  now.Add(-2 * time.Hour),
		LastUsedAt: now.Add(-time.Hour),
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID:         "active-b",
		UserID:     fx.userB.ID,
		ExpiresAt:  now.Add(time.Hour),
		CreatedAt:  now,
		LastUsedAt: now,
	})

	listed, err := fx.svc.ListSessions(context.Background(), fx.userA.ID, "active-a")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Fatalf("expected only user A active session, got %d: %+v", len(listed.Items), listed.Items)
	}
	item := listed.Items[0]
	if item.ID != "active-a" {
		t.Fatalf("expected active-a, got %s", item.ID)
	}
	if !item.Current {
		t.Fatal("expected current when X-Session-Id matches")
	}
	if item.IPAddress == nil || *item.IPAddress != "203.0.113.10" {
		t.Fatalf("expected ip, got %v", item.IPAddress)
	}
	if item.UserAgent == nil || *item.UserAgent != "Mozilla/5.0" {
		t.Fatalf("expected ua, got %v", item.UserAgent)
	}

	other, err := fx.svc.ListSessions(context.Background(), fx.userB.ID, "")
	if err != nil {
		t.Fatalf("list B: %v", err)
	}
	if len(other.Items) != 1 || other.Items[0].ID != "active-b" {
		t.Fatalf("user B should only see own row, got %+v", other.Items)
	}
	if other.Items[0].Current {
		t.Fatal("current must be false when header is absent")
	}
}

func TestListSessionsReplacesEpochCreatedAt(t *testing.T) {
	fx := newSessionFixture(t)
	now := time.Now().UTC()
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID:         "epoch-row",
		UserID:     fx.userA.ID,
		ExpiresAt:  now.Add(time.Hour),
		CreatedAt:  time.Unix(0, 0).UTC(),
		LastUsedAt: now,
	})

	listed, err := fx.svc.ListSessions(context.Background(), fx.userA.ID, "epoch-row")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(listed.Items))
	}
	if listed.Items[0].CreatedAt.Unix() <= 0 {
		t.Fatalf("expected started time fallback, got %v", listed.Items[0].CreatedAt)
	}
}

func TestListSessionsNullIPAndUserAgentWhenEmpty(t *testing.T) {
	fx := newSessionFixture(t)
	result, err := fx.svc.Login(context.Background(), fx.userA.Email, "password12", "", "")
	if err != nil {
		t.Fatalf("login: %v", err)
	}

	listed, err := fx.svc.ListSessions(context.Background(), fx.userA.ID, result.SessionID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(listed.Items))
	}
	item := listed.Items[0]
	if item.IPAddress != nil {
		t.Fatalf("expected null ipAddress, got %v", *item.IPAddress)
	}
	if item.UserAgent != nil {
		t.Fatalf("expected null userAgent, got %v", *item.UserAgent)
	}

	raw, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if decoded["ipAddress"] != nil {
		t.Fatalf("json ipAddress should be null, got %#v", decoded["ipAddress"])
	}
	if decoded["userAgent"] != nil {
		t.Fatalf("json userAgent should be null, got %#v", decoded["userAgent"])
	}

	stored := fx.authRepo.RefreshTokens[hash.SHA256Hex(result.RefreshToken)]
	if stored == nil {
		t.Fatal("expected stored token")
	}
	if stored.DeviceID != "" {
		t.Fatalf("device_id must stay unused, got %q", stored.DeviceID)
	}
}

func TestRefreshCopiesCreatedAtAndReturnsNewSessionID(t *testing.T) {
	fx := newSessionFixture(t)
	started := time.Now().Add(-24 * time.Hour).UTC().Truncate(time.Second)
	login, err := fx.svc.Login(context.Background(), fx.userA.Email, "password12", "1.1.1.1", "old-ua")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	old := fx.authRepo.RefreshTokens[hash.SHA256Hex(login.RefreshToken)]
	if old == nil {
		t.Fatal("expected login token")
	}
	old.CreatedAt = started
	old.LastUsedAt = started
	oldID := old.ID

	pair, err := fx.svc.Refresh(context.Background(), login.RefreshToken, "9.9.9.9", "new-ua")
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if pair.SessionID == "" || pair.SessionID == oldID {
		t.Fatalf("expected new sessionId, got %q (old %q)", pair.SessionID, oldID)
	}

	rotated := fx.authRepo.RefreshTokens[hash.SHA256Hex(pair.RefreshToken)]
	if rotated == nil {
		t.Fatal("expected rotated token")
	}
	if !rotated.CreatedAt.Equal(started) {
		t.Fatalf("expected copied created_at %v, got %v", started, rotated.CreatedAt)
	}
	if !rotated.LastUsedAt.After(started) {
		t.Fatalf("expected new last_used_at, got %v", rotated.LastUsedAt)
	}
	if rotated.IPAddress != "9.9.9.9" || rotated.UserAgent != "new-ua" {
		t.Fatalf("expected renewed ip/ua, got %q %q", rotated.IPAddress, rotated.UserAgent)
	}
	if old.RevokedAt == nil {
		t.Fatal("predecessor must be revoked")
	}
}

func TestRevokeSessionOwnerScopeAndCurrentGuard(t *testing.T) {
	fx := newSessionFixture(t)
	now := time.Now()
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "keep-a", UserID: fx.userA.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "other-a", UserID: fx.userA.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "b-row", UserID: fx.userB.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})

	if err := fx.svc.RevokeSession(context.Background(), fx.userA.ID, "keep-a", "keep-a"); err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected 400 when revoking current, got %v", err)
	}

	if err := fx.svc.RevokeSession(context.Background(), fx.userA.ID, "other-a", "keep-a"); err != nil {
		t.Fatalf("revoke other: %v", err)
	}
	listed, err := fx.svc.ListSessions(context.Background(), fx.userA.ID, "keep-a")
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(listed.Items) != 1 || listed.Items[0].ID != "keep-a" {
		t.Fatalf("expected only keep-a, got %+v", listed.Items)
	}

	if err := fx.svc.RevokeSession(context.Background(), fx.userA.ID, "b-row", "keep-a"); err == nil || !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected 404 for other user's row, got %v", err)
	}
	if err := fx.svc.RevokeSession(context.Background(), fx.userA.ID, "other-a", "keep-a"); err == nil || !errors.Is(err, apperrors.ErrNotFound) {
		t.Fatalf("expected 404 for already revoked, got %v", err)
	}

	bList, err := fx.svc.ListSessions(context.Background(), fx.userB.ID, "b-row")
	if err != nil {
		t.Fatalf("list B: %v", err)
	}
	if len(bList.Items) != 1 {
		t.Fatalf("user B session must remain, got %d", len(bList.Items))
	}
}

func TestRevokeOtherSessionsKeepsCurrentAndNoOpsWhenAlone(t *testing.T) {
	fx := newSessionFixture(t)
	now := time.Now()
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "current", UserID: fx.userA.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "extra-1", UserID: fx.userA.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "extra-2", UserID: fx.userA.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})
	putToken(fx.authRepo, &authmodel.RefreshToken{
		ID: "b-row", UserID: fx.userB.ID, ExpiresAt: now.Add(time.Hour), CreatedAt: now, LastUsedAt: now,
	})

	remaining, err := fx.svc.RevokeOtherSessions(context.Background(), fx.userA.ID, "current")
	if err != nil {
		t.Fatalf("revoke others: %v", err)
	}
	if len(remaining.Items) != 1 || remaining.Items[0].ID != "current" || !remaining.Items[0].Current {
		t.Fatalf("expected only current remaining, got %+v", remaining.Items)
	}

	noop, err := fx.svc.RevokeOtherSessions(context.Background(), fx.userA.ID, "current")
	if err != nil {
		t.Fatalf("no-op revoke others: %v", err)
	}
	if len(noop.Items) != 1 || noop.Items[0].ID != "current" {
		t.Fatalf("no-op must keep current, got %+v", noop.Items)
	}

	if _, err := fx.svc.RevokeOtherSessions(context.Background(), fx.userA.ID, ""); err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected 400 for missing sessionId, got %v", err)
	}
	if _, err := fx.svc.RevokeOtherSessions(context.Background(), fx.userA.ID, "b-row"); err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected 400 for not owned sessionId, got %v", err)
	}
	if _, err := fx.svc.RevokeOtherSessions(context.Background(), fx.userA.ID, "extra-1"); err == nil || !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expected 400 for already inactive sessionId, got %v", err)
	}

	bList, err := fx.svc.ListSessions(context.Background(), fx.userB.ID, "b-row")
	if err != nil {
		t.Fatalf("list B: %v", err)
	}
	if len(bList.Items) != 1 {
		t.Fatalf("user B must be untouched, got %d", len(bList.Items))
	}
}

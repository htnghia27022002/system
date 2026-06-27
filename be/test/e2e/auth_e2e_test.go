//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"be/internal/database"
	"be/test/testutil"
)

func TestAuthLoginE2E(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		_ = db.Exec("TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts RESTART IDENTITY CASCADE").Error
	})

	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	router := testutil.NewTestRouter(t, testutil.NewTestContainer(t, db))

	body, _ := json.Marshal(map[string]string{
		"email":    "admin@example.com",
		"password": "admin1234",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}

	var payload struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload.AccessToken == "" || payload.RefreshToken == "" {
		t.Fatal("expected token pair in response")
	}
}

func TestAuthOAuthProvidersE2E(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)

	router := testutil.NewTestRouter(t, testutil.NewTestContainer(t, db))

	req := httptest.NewRequest(http.MethodGet, "/api/auth/oauth/providers", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var payload struct {
		Providers []string `json:"providers"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if payload.Providers == nil {
		t.Fatal("expected providers array")
	}
}

func TestAuthLoginInvalidCredentialsE2E(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		_ = db.Exec("TRUNCATE users, roles, role_permissions, permissions, refresh_tokens, oauth_accounts RESTART IDENTITY CASCADE").Error
	})

	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	router := testutil.NewTestRouter(t, testutil.NewTestContainer(t, db))

	body, _ := json.Marshal(map[string]string{
		"email":    "admin@example.com",
		"password": "wrong-password",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d body=%s", rec.Code, rec.Body.String())
	}
}

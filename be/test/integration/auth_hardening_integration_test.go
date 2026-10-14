//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"be/internal/config"
	"be/internal/database"
	"be/public/handlers"
	"be/test/testutil"
)

func hardeningRouter(t *testing.T, opts ...func(*config.Config)) *gin.Engine {
	t.Helper()
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() { testutil.TruncateAuthTables(t, db) })
	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}
	return testutil.NewTestRouter(t, testutil.NewTestContainer(t, db, opts...))
}

func postJSON(router *gin.Engine, path string, body any, mutate func(*http.Request)) *httptest.ResponseRecorder {
	raw, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")
	if mutate != nil {
		mutate(req)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

func refreshCookie(rec *httptest.ResponseRecorder) *http.Cookie {
	for _, c := range rec.Result().Cookies() {
		if c.Name == handlers.RefreshCookieName {
			return c
		}
	}
	return nil
}

func withCookie(c *http.Cookie, csrf bool) func(*http.Request) {
	return func(req *http.Request) {
		req.AddCookie(&http.Cookie{Name: c.Name, Value: c.Value})
		if csrf {
			req.Header.Set(handlers.CSRFHeader, "XMLHttpRequest")
		}
	}
}

func TestRefreshCookieRotationCSRFAndLogout(t *testing.T) {
	router := hardeningRouter(t)
	creds := map[string]string{"email": "admin@example.com", "password": "admin1234"}

	login := postJSON(router, "/api/auth/login", creds, nil)
	if login.Code != http.StatusOK {
		t.Fatalf("login: %d %s", login.Code, login.Body.String())
	}
	var body authPayload
	_ = json.Unmarshal(login.Body.Bytes(), &body)
	cookie := refreshCookie(login)
	if body.AccessToken == "" || body.RefreshToken != "" || cookie == nil || !cookie.HttpOnly ||
		cookie.Path != "/api/auth" || cookie.SameSite != http.SameSiteLaxMode || cookie.Secure {
		t.Fatalf("expected access token in JSON and httpOnly lax refresh cookie, body=%+v cookie=%+v", body, cookie)
	}

	// Cookie without the CSRF header is refused and does not consume the token.
	if rec := postJSON(router, "/api/auth/refresh", map[string]string{}, withCookie(cookie, false)); rec.Code != http.StatusForbidden {
		t.Fatalf("refresh without %s: expected 403, got %d", handlers.CSRFHeader, rec.Code)
	}

	refreshed := postJSON(router, "/api/auth/refresh", map[string]string{}, withCookie(cookie, true))
	rotated := refreshCookie(refreshed)
	var pair authPayload
	_ = json.Unmarshal(refreshed.Body.Bytes(), &pair)
	if refreshed.Code != http.StatusOK || pair.AccessToken == "" || pair.RefreshToken != "" ||
		rotated == nil || rotated.Value == "" || rotated.Value == cookie.Value {
		t.Fatalf("refresh via cookie: %d body=%s cookie=%+v", refreshed.Code, refreshed.Body.String(), rotated)
	}

	// The old cookie was rotated away; replaying it fails and clears the cookie.
	replay := postJSON(router, "/api/auth/refresh", map[string]string{}, withCookie(cookie, true))
	if replay.Code != http.StatusUnauthorized {
		t.Fatalf("replayed refresh cookie: expected 401, got %d", replay.Code)
	}
	if c := refreshCookie(replay); c == nil || c.MaxAge >= 0 {
		t.Fatalf("failed refresh must clear the cookie, got %+v", c)
	}

	logout := postJSON(router, "/api/auth/logout", map[string]string{}, withCookie(rotated, true))
	if logout.Code != http.StatusNoContent {
		t.Fatalf("logout: %d %s", logout.Code, logout.Body.String())
	}
	if c := refreshCookie(logout); c == nil || c.MaxAge >= 0 {
		t.Fatalf("logout must clear the cookie, got %+v", c)
	}
	if rec := postJSON(router, "/api/auth/refresh", map[string]string{}, withCookie(rotated, true)); rec.Code != http.StatusUnauthorized {
		t.Fatalf("refresh after logout: expected 401, got %d", rec.Code)
	}
}

func TestRefreshCookieSameSiteNoneForcesSecure(t *testing.T) {
	router := hardeningRouter(t, func(cfg *config.Config) { cfg.Auth.CookieSameSite = "none" })
	login := postJSON(router, "/api/auth/login", map[string]string{"email": "admin@example.com", "password": "admin1234"}, nil)
	if c := refreshCookie(login); c == nil || c.SameSite != http.SameSiteNoneMode || !c.Secure {
		t.Fatalf("SameSite=None cookie must be Secure, got %+v", c)
	}
}

func TestAccountLockoutReturns429(t *testing.T) {
	router := hardeningRouter(t, func(cfg *config.Config) {
		cfg.Auth.MaxFailedLogins = 3
		cfg.Auth.RateLimitPerMinute = 0 // isolate lockout from the per-IP limit
	})
	wrong := map[string]string{"email": "admin@example.com", "password": "wrong-password"}
	for i := 0; i < 3; i++ {
		if rec := postJSON(router, "/api/auth/login", wrong, nil); rec.Code != http.StatusUnauthorized {
			t.Fatalf("attempt %d: expected 401, got %d", i+1, rec.Code)
		}
	}
	right := map[string]string{"email": "admin@example.com", "password": "admin1234"}
	if rec := postJSON(router, "/api/auth/login", right, nil); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("locked account: expected 429, got %d %s", rec.Code, rec.Body.String())
	}
}

func TestLoginRateLimitPerIP(t *testing.T) {
	router := hardeningRouter(t, func(cfg *config.Config) {
		cfg.Auth.RateLimitPerMinute = 2
		cfg.Auth.MaxFailedLogins = 0
	})
	wrong := map[string]string{"email": "nobody@example.com", "password": "x"}
	for i := 0; i < 2; i++ {
		postJSON(router, "/api/auth/login", wrong, nil)
	}
	rec := postJSON(router, "/api/auth/login", wrong, nil)
	if rec.Code != http.StatusTooManyRequests || rec.Header().Get("Retry-After") == "" {
		t.Fatalf("expected 429 with Retry-After, got %d", rec.Code)
	}
}

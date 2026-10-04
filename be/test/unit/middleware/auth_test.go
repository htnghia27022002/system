package middleware_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	jwtmanager "be/common/jwt"
	"be/internal/middleware"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
	"be/test/testutil"
)

func TestRequirePermissionAllowsSuperAdminWithoutKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testutil.UnitConfig()
	cfg.JWTAccessTTL = time.Minute
	manager := jwtmanager.NewManager(cfg)

	token, err := manager.SignAccessToken("u1", "sa@example.com", "SA", "user", "role-user", []string{}, true)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	userRepo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"u1": {ID: "u1", Email: "sa@example.com", IsSuperAdmin: true},
		},
	}

	r := gin.New()
	r.GET("/admin/users", middleware.Auth(manager, nil, userRepo), middleware.RequireView("users"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for super admin, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestRequirePermissionForbidsNonSuperAdminWithoutKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := testutil.UnitConfig()
	cfg.JWTAccessTTL = time.Minute
	manager := jwtmanager.NewManager(cfg)

	token, err := manager.SignAccessToken("u2", "m@example.com", "M", "user", "role-user", []string{}, false)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	userRepo := &testutil.MemoryUserRepo{
		Users: map[string]*usermodel.User{
			"u2": {ID: "u2", Email: "m@example.com", IsSuperAdmin: false},
		},
	}

	r := gin.New()
	r.GET("/admin/users", middleware.Auth(manager, nil, userRepo), middleware.RequireView("users"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d body=%s", w.Code, w.Body.String())
	}
}

type failingUserRepo struct {
	testutil.MemoryUserRepo
}

func (r *failingUserRepo) GetByID(context.Context, string) (*usermodel.User, error) {
	return nil, errors.New("connection refused")
}

func serveWithAuth(t *testing.T, userRepo interfaces.UserRepository, subject string) *httptest.ResponseRecorder {
	t.Helper()
	gin.SetMode(gin.TestMode)
	cfg := testutil.UnitConfig()
	cfg.JWTAccessTTL = time.Minute
	manager := jwtmanager.NewManager(cfg)

	// Token claims super-admin; the middleware must not trust that when the DB lookup fails.
	token, err := manager.SignAccessToken(subject, "x@example.com", "X", "user", "role-user", []string{}, true)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}

	r := gin.New()
	r.GET("/admin/users", middleware.Auth(manager, nil, userRepo), middleware.RequireView("users"), func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	req := httptest.NewRequest(http.MethodGet, "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthRejectsTokenForDeletedUser(t *testing.T) {
	w := serveWithAuth(t, &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{}}, "gone")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing user, got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAuthFailsClosedWhenUserLookupErrors(t *testing.T) {
	w := serveWithAuth(t, &failingUserRepo{}, "u1")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when user lookup fails, got %d body=%s", w.Code, w.Body.String())
	}
	if strings.Contains(w.Body.String(), "connection refused") {
		t.Fatalf("internal error leaked to client: %s", w.Body.String())
	}
}

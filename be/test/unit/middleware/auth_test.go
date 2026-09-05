package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	jwtmanager "be/internal/common/jwt"
	"be/internal/middleware"
	usermodel "be/internal/models/user"
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

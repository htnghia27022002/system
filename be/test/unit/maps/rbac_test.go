package maps_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	jwtmanager "be/common/jwt"
	"be/internal/middleware"
	usermodel "be/internal/models/user"
	"be/common/rbac"
	"be/test/testutil"
)

func mapsRouter(t *testing.T, tokenUser *usermodel.User, perms []string, super bool) (*gin.Engine, string) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	cfg := testutil.UnitConfig()
	cfg.JWTAccessTTL = time.Minute
	manager := jwtmanager.NewManager(cfg)
	token, err := manager.SignAccessToken(tokenUser.ID, tokenUser.Email, tokenUser.FullName, "user", tokenUser.RoleID, perms, super)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	userRepo := &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{tokenUser.ID: tokenUser}}
	r := gin.New()
	admin := r.Group("/api/admin")
	admin.Use(middleware.Auth(manager, nil, userRepo))
	maps := admin.Group("/maps")
	maps.GET("/places", middleware.RequireView("maps"), func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"items": []any{}}) })
	maps.POST("/sources", middleware.RequireModify("maps"), func(c *gin.Context) { c.Status(http.StatusCreated) })
	maps.POST("/ingest", middleware.RequireModify("maps"), func(c *gin.Context) { c.Status(http.StatusAccepted) })
	maps.PATCH("/places/:id", middleware.RequireModify("maps"), func(c *gin.Context) { c.Status(http.StatusOK) })
	return r, token
}

func TestMapsViewForbiddenWithoutKey(t *testing.T) {
	t.Parallel()
	user := &usermodel.User{ID: "u-none", Email: "n@example.com", RoleID: "role-user"}
	r, token := mapsRouter(t, user, []string{}, false)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/maps/places", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestMapsViewOnlyForbiddenOnModify(t *testing.T) {
	t.Parallel()
	user := &usermodel.User{ID: "u-view", Email: "v@example.com", RoleID: "role-user"}
	r, token := mapsRouter(t, user, []string{rbac.Key("maps", rbac.ActionView)}, false)

	cases := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/admin/maps/ingest"},
		{http.MethodPost, "/api/admin/maps/sources"},
		{http.MethodPatch, "/api/admin/maps/places/p1"},
	}
	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusForbidden {
			t.Fatalf("%s %s expected 403, got %d", tc.method, tc.path, w.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/admin/maps/places", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("view list expected 200, got %d", w.Code)
	}
}

func TestMapsSuperAdminAllowedWithoutCatalogKeys(t *testing.T) {
	t.Parallel()
	user := &usermodel.User{ID: "u-sa", Email: "sa@example.com", RoleID: "role-user", IsSuperAdmin: true}
	r, token := mapsRouter(t, user, []string{}, true)
	paths := []struct {
		method string
		path   string
		want   int
	}{
		{http.MethodGet, "/api/admin/maps/places", http.StatusOK},
		{http.MethodPost, "/api/admin/maps/ingest", http.StatusAccepted},
		{http.MethodPost, "/api/admin/maps/sources", http.StatusCreated},
		{http.MethodPatch, "/api/admin/maps/places/p1", http.StatusOK},
	}
	for _, tc := range paths {
		req := httptest.NewRequest(tc.method, tc.path, strings.NewReader(`{}`))
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != tc.want {
			t.Fatalf("%s %s expected %d, got %d", tc.method, tc.path, tc.want, w.Code)
		}
	}
}

func TestMapsUnauthenticatedIs401(t *testing.T) {
	t.Parallel()
	user := &usermodel.User{ID: "u-x", Email: "x@example.com"}
	r, _ := mapsRouter(t, user, []string{rbac.Key("maps", rbac.ActionView)}, false)
	req := httptest.NewRequest(http.MethodGet, "/api/admin/maps/places", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

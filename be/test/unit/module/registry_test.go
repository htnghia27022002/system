package module_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"be/internal/config"
	"be/internal/database/seeders"
	"be/internal/module"
	"be/internal/modules"
	"be/public/handlers"
)

func build(t *testing.T, cfg config.Config) module.Set {
	t.Helper()
	set, err := modules.Build(&module.Kernel{Config: cfg, Auth: func(c *gin.Context) { c.Next() }})
	if err != nil {
		t.Fatalf("build modules: %v", err)
	}
	return set
}

func TestDefaultsEnableWebhooks(t *testing.T) {
	set := build(t, config.Config{})
	if got := strings.Join(set.Names(), ","); got != "webhooks" {
		t.Fatalf("default modules = %q, want webhooks", got)
	}
}

func TestModulesEnvSelectsExactList(t *testing.T) {
	none := build(t, config.Config{ModulesSet: true, Modules: nil})
	if len(none) != 0 {
		t.Fatalf("MODULES='' should enable nothing, got %v", none.Names())
	}

	one := build(t, config.Config{ModulesSet: true, Modules: []string{" Webhooks "}})
	if got := strings.Join(one.Names(), ","); got != "webhooks" {
		t.Fatalf("MODULES=webhooks gave %q", got)
	}
}

func TestUnknownModuleIsAnError(t *testing.T) {
	_, err := modules.Build(&module.Kernel{Config: config.Config{ModulesSet: true, Modules: []string{"maps"}}})
	if err == nil || !strings.Contains(err.Error(), `unknown module "maps"`) {
		t.Fatalf("expected unknown module error, got %v", err)
	}
}

func TestModulePermissionsMergeWithCoreWithoutDuplicates(t *testing.T) {
	set := build(t, config.Config{})
	all, err := seeders.MergePermissions(set.SeedContributions()...)
	if err != nil {
		t.Fatalf("merge: %v", err)
	}
	keys := map[string]string{}
	for _, p := range all {
		keys[p.Key] = p.ID
	}
	// IDs are stable across the move into the module.
	if keys["webhooks:view"] != "10000001-0000-4000-8000-000000000017" ||
		keys["webhooks:modify"] != "10000001-0000-4000-8000-000000000018" {
		t.Fatalf("webhooks permission IDs changed: %v", keys)
	}
}

func TestEachModuleHasItsOwnMigrationSet(t *testing.T) {
	sets := build(t, config.Config{}).MigrationSets()
	if len(sets) != 1 || sets[0].Name != "webhooks" || sets[0].FS == nil {
		t.Fatalf("unexpected migration sets: %+v", sets)
	}
}

func TestDisabledModuleRoutesAreNotMounted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cases := map[string]struct {
		set     module.Set
		mounted bool
	}{
		"enabled":  {build(t, config.Config{}), true},
		"disabled": {build(t, config.Config{ModulesSet: true}), false},
	}
	for name, tc := range cases {
		r := gin.New()
		api := r.Group("/api")
		tc.set.RegisterRoutes(module.Routes{API: api, Admin: api.Group("/admin")})

		mounted := false
		for _, info := range r.Routes() {
			if info.Method == http.MethodGet && info.Path == "/api/webhooks/inbox" {
				mounted = true
			}
		}
		if mounted != tc.mounted {
			t.Fatalf("%s: webhooks inbox route mounted=%v, want %v", name, mounted, tc.mounted)
		}
	}
}

func TestAppInfoListsEnabledModules(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for name, tc := range map[string]struct {
		set  module.Set
		want string
	}{
		"default":   {build(t, config.Config{}), `{"modules":["webhooks"]}`},
		"core only": {build(t, config.Config{ModulesSet: true}), `{"modules":[]}`},
	} {
		r := gin.New()
		r.GET("/api/app/info", handlers.NewAppHandler(tc.set.Names()).Info)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/api/app/info", nil))
		if w.Code != http.StatusOK || w.Body.String() != tc.want {
			t.Fatalf("%s: got %d %s, want %s", name, w.Code, w.Body.String(), tc.want)
		}
	}
}

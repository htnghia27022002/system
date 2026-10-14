package apidoc_test

import (
	"bytes"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/config"
	"be/pkg/postgres"
	"be/public"
	"be/public/handlers"
	"be/public/routes"
	"be/test/testutil"
)

// anyMethods are what gin's Any() registers; an "ANY" op covers them.
var anyMethods = []string{"GET", "POST", "PUT", "PATCH", "HEAD", "OPTIONS", "DELETE", "CONNECT", "TRACE"}

func TestEveryRouteIsDocumentedAndViceVersa(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := config.Config{JWTSecret: "x", CORSOrigins: []string{"http://localhost:3000"}, TrustedProxies: []string{"127.0.0.1"}, UploadDir: t.TempDir(), Storage: config.StorageConfig{Driver: "local", PublicURL: "/api/media"}}
	// Wiring only: nothing here touches the database.
	container, err := app.NewContainer(cfg, &postgres.Postgres{}, app.WithMailer(&testutil.RecordingMailer{}))
	if err != nil {
		t.Fatal(err)
	}
	router, err := public.NewRouter(cfg, container, handlers.NewHealthHandler())
	if err != nil {
		t.Fatal(err)
	}

	mounted := map[string]bool{}
	for _, r := range router.Routes() {
		mounted[r.Method+" "+r.Path] = true
	}
	documented := map[string]bool{}
	for _, op := range append(routes.CoreOperations(), container.Modules.APIDoc()...) {
		if op.Method == "ANY" {
			for _, m := range anyMethods {
				documented[m+" "+op.Path] = true
			}
			continue
		}
		key := op.Method + " " + op.Path
		if documented[key] {
			t.Errorf("documented twice: %s", key)
		}
		documented[key] = true
	}

	var missing, stale []string
	for key := range mounted {
		if !documented[key] {
			missing = append(missing, key)
		}
	}
	for key := range documented {
		if !mounted[key] {
			stale = append(stale, key)
		}
	}
	sort.Strings(missing)
	sort.Strings(stale)
	if len(missing) > 0 {
		t.Errorf("routes without docs (add them to routes.CoreOperations or the module's APIDoc):\n  %s", strings.Join(missing, "\n  "))
	}
	if len(stale) > 0 {
		t.Errorf("docs for routes that are not mounted:\n  %s", strings.Join(stale, "\n  "))
	}
}

func TestCommittedSpecIsUpToDate(t *testing.T) {
	container, err := app.NewContainer(config.Config{JWTSecret: "x", UploadDir: t.TempDir()}, &postgres.Postgres{}, app.WithMailer(&testutil.RecordingMailer{}))
	if err != nil {
		t.Fatal(err)
	}
	want, err := json.MarshalIndent(routes.OpenAPI(container.Modules), "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile("../../../api/openapi.json")
	if err != nil {
		t.Fatalf("read api/openapi.json: %v", err)
	}
	if !bytes.Equal(bytes.TrimSpace(got), bytes.TrimSpace(want)) {
		t.Fatal("api/openapi.json is stale: run `make openapi` (go run ./cmd/openapi -o api/openapi.json) and commit it")
	}
}

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

	"be/internal/database"
	"be/test/testutil"
)

func authed(router *gin.Engine, method, path, token string, body any) *httptest.ResponseRecorder {
	var buf bytes.Buffer
	if body != nil {
		_ = json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("User-Agent", "IntegrationTest/1")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

type auditPage struct {
	Items []struct {
		Action     string         `json:"action"`
		ActorEmail string         `json:"actorEmail"`
		ResourceID string         `json:"resourceId"`
		UserAgent  string         `json:"userAgent"`
		RequestID  string         `json:"requestId"`
		Metadata   map[string]any `json:"metadata"`
	} `json:"items"`
	Total int `json:"total"`
}

func TestAuditTrailRecordsSignInsAndAdminChanges(t *testing.T) {
	router, _ := accountRouter(t)
	_ = postJSON(router, "/api/auth/login", map[string]string{"email": "demo@example.com", "password": "wrong-password"}, nil)
	admin := loginHTTP(t, router, "admin@example.com", "admin1234", "", "")

	rec := authed(router, http.MethodPatch, "/api/admin/settings", admin.AccessToken,
		map[string]any{"values": map[string]any{"app.supportEmail": "help@example.com"}})
	if rec.Code != http.StatusOK {
		t.Fatalf("settings update: %d %s", rec.Code, rec.Body)
	}

	rec = authed(router, http.MethodGet, "/api/admin/audit-logs?pageSize=50", admin.AccessToken, nil)
	var page auditPage
	_ = json.Unmarshal(rec.Body.Bytes(), &page)
	if rec.Code != http.StatusOK {
		t.Fatalf("audit list: %d %s", rec.Code, rec.Body)
	}
	seen := map[string]bool{}
	for _, item := range page.Items {
		seen[item.Action] = true
		if item.Action == "settings.update" {
			if item.ActorEmail != "admin@example.com" || item.UserAgent != "IntegrationTest/1" || item.RequestID == "" {
				t.Fatalf("settings.update lacks actor/request facts: %+v", item)
			}
		}
	}
	for _, want := range []string{"auth.login_failed", "auth.login", "settings.update"} {
		if !seen[want] {
			t.Fatalf("missing %s in %+v", want, page.Items)
		}
	}

	rec = authed(router, http.MethodGet, "/api/admin/audit-logs?action=auth.login_failed", admin.AccessToken, nil)
	_ = json.Unmarshal(rec.Body.Bytes(), &page)
	if page.Total != 1 || page.Items[0].ActorEmail != "demo@example.com" {
		t.Fatalf("action filter: %+v", page)
	}

	demo := loginHTTP(t, router, "demo@example.com", "password123", "", "")
	if rec := authed(router, http.MethodGet, "/api/admin/audit-logs", demo.AccessToken, nil); rec.Code != http.StatusForbidden {
		t.Fatalf("audit log needs audit:view, got %d", rec.Code)
	}
}

func TestSettingsGateRegistrationAndAppInfo(t *testing.T) {
	router, _ := accountRouter(t)
	admin := loginHTTP(t, router, "admin@example.com", "admin1234", "", "")

	if rec := authed(router, http.MethodPatch, "/api/admin/settings", admin.AccessToken,
		map[string]any{"values": map[string]any{"auth.registrationEnabled": "nope"}}); rec.Code != http.StatusBadRequest {
		t.Fatalf("invalid value: %d", rec.Code)
	}
	if rec := authed(router, http.MethodPatch, "/api/admin/settings", admin.AccessToken,
		map[string]any{"values": map[string]any{"auth.registrationEnabled": false}}); rec.Code != http.StatusOK {
		t.Fatalf("disable registration: %d %s", rec.Code, rec.Body)
	}

	reg := postJSON(router, "/api/auth/register", map[string]string{"name": "Closed Door", "email": "x@example.com", "password": "password-123"}, nil)
	var body struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(reg.Body.Bytes(), &body)
	if reg.Code != http.StatusForbidden || body.Code != "registration_disabled" {
		t.Fatalf("register while closed: %d %s", reg.Code, reg.Body)
	}

	info := httptest.NewRecorder()
	router.ServeHTTP(info, httptest.NewRequest(http.MethodGet, "/api/app/info", nil))
	var app struct {
		Settings map[string]any `json:"settings"`
	}
	_ = json.Unmarshal(info.Body.Bytes(), &app)
	if app.Settings["auth.registrationEnabled"] != false {
		t.Fatalf("app info: %s", info.Body)
	}
	if _, leaked := app.Settings["audit.retentionDays"]; leaked {
		t.Fatal("private setting exposed on /app/info")
	}
}

func TestDashboardOverviewCountsRealData(t *testing.T) {
	router, _ := accountRouter(t)
	admin := loginHTTP(t, router, "admin@example.com", "admin1234", "", "")
	loginHTTP(t, router, "demo@example.com", "password123", "", "")

	rec := authed(router, http.MethodGet, "/api/admin/dashboard/overview", admin.AccessToken, nil)
	var overview struct {
		KPIs []struct {
			Key   string `json:"key"`
			Value int64  `json:"value"`
		} `json:"kpis"`
		WeeklyActivity []struct {
			Date  string `json:"date"`
			Value int64  `json:"value"`
		} `json:"weeklyActivity"`
	}
	_ = json.Unmarshal(rec.Body.Bytes(), &overview)
	if rec.Code != http.StatusOK || len(overview.KPIs) != 4 || len(overview.WeeklyActivity) != 7 {
		t.Fatalf("overview: %d %s", rec.Code, rec.Body)
	}
	kpi := map[string]int64{}
	for _, k := range overview.KPIs {
		kpi[k.Key] = k.Value
	}
	if kpi["totalUsers"] != 2 || kpi["activeSessions"] != 2 || kpi["roles"] != 2 || kpi["newUsers"] != 2 {
		t.Fatalf("kpis: %v", kpi)
	}
	if today := overview.WeeklyActivity[6]; today.Value != 2 {
		t.Fatalf("today's sign-ins: %+v", overview.WeeklyActivity)
	}
}

func TestJanitorPrunesByRetentionSetting(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() { testutil.TruncateAuthTables(t, db) })
	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}
	container := testutil.NewTestContainer(t, db)
	ctx := context.Background()
	if _, err := container.Settings.Update(ctx, map[string]any{"audit.retentionDays": float64(30)}, ""); err != nil {
		t.Fatal(err)
	}
	_, err := db.Pool.Exec(ctx, `INSERT INTO audit_logs (action, created_at) VALUES ('old.entry', NOW() - INTERVAL '40 days'), ('new.entry', NOW() - INTERVAL '10 days')`)
	if err != nil {
		t.Fatal(err)
	}
	container.Janitor.RunOnce(ctx)

	var actions []string
	rows, _ := db.Pool.Query(ctx, `SELECT action FROM audit_logs WHERE action LIKE '%.entry' ORDER BY action`)
	for rows.Next() {
		var a string
		_ = rows.Scan(&a)
		actions = append(actions, a)
	}
	rows.Close()
	if len(actions) != 1 || actions[0] != "new.entry" {
		t.Fatalf("expected only new.entry to survive, got %v", actions)
	}
}

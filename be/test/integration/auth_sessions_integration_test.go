//go:build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"be/internal/database"
	"be/test/testutil"
)

type authPayload struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	SessionID    string `json:"sessionId"`
}

type sessionListPayload struct {
	Items []struct {
		ID         string  `json:"id"`
		Current    bool    `json:"current"`
		CreatedAt  string  `json:"createdAt"`
		IPAddress  *string `json:"ipAddress"`
		UserAgent  *string `json:"userAgent"`
	} `json:"items"`
}

func loginHTTP(t *testing.T, router *gin.Engine, email, password, ip, ua string) authPayload {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"email": email, "password": password})
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if ip != "" {
		req.Header.Set("X-Forwarded-For", ip)
	}
	if ua != "" {
		req.Header.Set("User-Agent", ua)
	}
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("login: expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var payload authPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode login: %v", err)
	}
	if payload.AccessToken == "" || payload.SessionID == "" {
		t.Fatalf("expected accessToken and sessionId, got %+v", payload)
	}
	return payload
}

func TestAuthSessionsListMarksCurrentAndSurvivesRefresh(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		testutil.TruncateAuthTables(t, db)
	})
	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	router := testutil.NewTestRouter(t, testutil.NewTestContainer(t, db))
	first := loginHTTP(t, router, "admin@example.com", "admin1234", "203.0.113.10", "Browser/A")
	second := loginHTTP(t, router, "admin@example.com", "admin1234", "198.51.100.20", "Browser/B")

	req := httptest.NewRequest(http.MethodGet, "/api/auth/sessions", nil)
	req.Header.Set("Authorization", "Bearer "+first.AccessToken)
	req.Header.Set("X-Session-Id", first.SessionID)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list: expected 200, got %d body=%s", rec.Code, rec.Body.String())
	}
	var listed sessionListPayload
	if err := json.Unmarshal(rec.Body.Bytes(), &listed); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(listed.Items) != 2 {
		t.Fatalf("expected 2 sessions, got %d", len(listed.Items))
	}
	if !listed.Items[0].Current || listed.Items[0].ID != first.SessionID {
		t.Fatalf("expected current first, got %+v", listed.Items)
	}
	if listed.Items[0].CreatedAt == "" || strings.HasPrefix(listed.Items[0].CreatedAt, "1970-") {
		t.Fatalf("expected real createdAt, got %q", listed.Items[0].CreatedAt)
	}
	if listed.Items[1].Current || listed.Items[1].ID != second.SessionID {
		t.Fatalf("expected other session second, got %+v", listed.Items)
	}

	refreshBody, _ := json.Marshal(map[string]string{"refreshToken": first.RefreshToken})
	refreshReq := httptest.NewRequest(http.MethodPost, "/api/auth/refresh", bytes.NewReader(refreshBody))
	refreshReq.Header.Set("Content-Type", "application/json")
	refreshReq.Header.Set("X-Forwarded-For", "203.0.113.11")
	refreshReq.Header.Set("User-Agent", "Browser/A-renewed")
	refreshRec := httptest.NewRecorder()
	router.ServeHTTP(refreshRec, refreshReq)
	if refreshRec.Code != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d body=%s", refreshRec.Code, refreshRec.Body.String())
	}
	var renewed authPayload
	if err := json.Unmarshal(refreshRec.Body.Bytes(), &renewed); err != nil {
		t.Fatalf("decode refresh: %v", err)
	}
	if renewed.SessionID == "" || renewed.SessionID == first.SessionID {
		t.Fatalf("expected new sessionId after refresh, got %q", renewed.SessionID)
	}

	afterReq := httptest.NewRequest(http.MethodGet, "/api/auth/sessions", nil)
	afterReq.Header.Set("Authorization", "Bearer "+renewed.AccessToken)
	afterReq.Header.Set("X-Session-Id", renewed.SessionID)
	afterRec := httptest.NewRecorder()
	router.ServeHTTP(afterRec, afterReq)
	if afterRec.Code != http.StatusOK {
		t.Fatalf("list after refresh: expected 200, got %d body=%s", afterRec.Code, afterRec.Body.String())
	}
	var after sessionListPayload
	if err := json.Unmarshal(afterRec.Body.Bytes(), &after); err != nil {
		t.Fatalf("decode list after refresh: %v", err)
	}
	var currentCount int
	for _, item := range after.Items {
		if item.Current {
			currentCount++
			if item.ID != renewed.SessionID {
				t.Fatalf("current id %s != renewed sessionId %s", item.ID, renewed.SessionID)
			}
		}
	}
	if currentCount != 1 {
		t.Fatalf("expected exactly one current row after refresh, got %d items=%+v", currentCount, after.Items)
	}
}

func TestAuthSessionsRevokeOthersAndCurrentGuard(t *testing.T) {
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() {
		testutil.TruncateAuthTables(t, db)
	})
	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}

	router := testutil.NewTestRouter(t, testutil.NewTestContainer(t, db))
	keep := loginHTTP(t, router, "admin@example.com", "admin1234", "203.0.113.10", "Keep")
	other := loginHTTP(t, router, "admin@example.com", "admin1234", "198.51.100.20", "Other")

	delCurrent := httptest.NewRequest(http.MethodDelete, "/api/auth/sessions/"+keep.SessionID, nil)
	delCurrent.Header.Set("Authorization", "Bearer "+keep.AccessToken)
	delCurrent.Header.Set("X-Session-Id", keep.SessionID)
	delCurrentRec := httptest.NewRecorder()
	router.ServeHTTP(delCurrentRec, delCurrent)
	if delCurrentRec.Code != http.StatusBadRequest {
		t.Fatalf("delete current: expected 400, got %d body=%s", delCurrentRec.Code, delCurrentRec.Body.String())
	}

	delOther := httptest.NewRequest(http.MethodDelete, "/api/auth/sessions/"+other.SessionID, nil)
	delOther.Header.Set("Authorization", "Bearer "+keep.AccessToken)
	delOther.Header.Set("X-Session-Id", keep.SessionID)
	delOtherRec := httptest.NewRecorder()
	router.ServeHTTP(delOtherRec, delOther)
	if delOtherRec.Code != http.StatusNoContent {
		t.Fatalf("delete other: expected 204, got %d body=%s", delOtherRec.Code, delOtherRec.Body.String())
	}

	third := loginHTTP(t, router, "admin@example.com", "admin1234", "192.0.2.1", "Third")
	if third.SessionID == "" {
		t.Fatal("expected third session")
	}
	body, _ := json.Marshal(map[string]string{"sessionId": keep.SessionID})
	revokeReq := httptest.NewRequest(http.MethodPost, "/api/auth/sessions/revoke-others", bytes.NewReader(body))
	revokeReq.Header.Set("Authorization", "Bearer "+keep.AccessToken)
	revokeReq.Header.Set("Content-Type", "application/json")
	revokeRec := httptest.NewRecorder()
	router.ServeHTTP(revokeRec, revokeReq)
	if revokeRec.Code != http.StatusOK {
		t.Fatalf("revoke-others: expected 200, got %d body=%s", revokeRec.Code, revokeRec.Body.String())
	}
	var remaining sessionListPayload
	if err := json.Unmarshal(revokeRec.Body.Bytes(), &remaining); err != nil {
		t.Fatalf("decode revoke-others: %v", err)
	}
	if len(remaining.Items) != 1 || remaining.Items[0].ID != keep.SessionID {
		t.Fatalf("expected only keep session, got %+v", remaining.Items)
	}

	noopBody, _ := json.Marshal(map[string]string{"sessionId": keep.SessionID})
	noopReq := httptest.NewRequest(http.MethodPost, "/api/auth/sessions/revoke-others", bytes.NewReader(noopBody))
	noopReq.Header.Set("Authorization", "Bearer "+keep.AccessToken)
	noopReq.Header.Set("Content-Type", "application/json")
	noopRec := httptest.NewRecorder()
	router.ServeHTTP(noopRec, noopReq)
	if noopRec.Code != http.StatusOK {
		t.Fatalf("revoke-others no-op: expected 200, got %d body=%s", noopRec.Code, noopRec.Body.String())
	}

	missing := httptest.NewRequest(http.MethodPost, "/api/auth/sessions/revoke-others", bytes.NewReader([]byte(`{}`)))
	missing.Header.Set("Authorization", "Bearer "+keep.AccessToken)
	missing.Header.Set("Content-Type", "application/json")
	missingRec := httptest.NewRecorder()
	router.ServeHTTP(missingRec, missing)
	if missingRec.Code != http.StatusBadRequest {
		t.Fatalf("revoke-others missing id: expected 400, got %d body=%s", missingRec.Code, missingRec.Body.String())
	}
}

//go:build integration

package integration_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"

	"github.com/gin-gonic/gin"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/test/testutil"
)

var mailedToken = regexp.MustCompile(`\?token=([A-Za-z0-9_%-]+)`)

func accountRouter(t *testing.T, opts ...func(*config.Config)) (*gin.Engine, *testutil.RecordingMailer) {
	t.Helper()
	db := testutil.ConnectPostgres(t)
	testutil.MigrateTestSchema(t, db)
	t.Cleanup(func() { testutil.TruncateAuthTables(t, db) })
	if err := database.SeedRBAC(context.Background(), db); err != nil {
		t.Fatalf("seed rbac: %v", err)
	}
	mailer := &testutil.RecordingMailer{}
	opts = append([]func(*config.Config){func(cfg *config.Config) {
		cfg.Auth.RateLimitPerMinute = 0
		cfg.AppURL = "https://app.example.com"
	}}, opts...)
	container := testutil.NewTestContainerWith(t, db, []app.Option{app.WithMailer(mailer)}, opts...)
	return testutil.NewTestAPIRouter(t, container), mailer
}

func lastToken(t *testing.T, mailer *testutil.RecordingMailer) string {
	t.Helper()
	m := mailedToken.FindStringSubmatch(mailer.Last().Text)
	if m == nil {
		t.Fatalf("no token in email %q", mailer.Last().Text)
	}
	token, _ := url.QueryUnescape(m[1])
	return token
}

func TestForgotAndResetPasswordHTTP(t *testing.T) {
	router, mailer := accountRouter(t)

	unknown := postJSON(router, "/api/auth/password/forgot", map[string]string{"email": "nobody@example.com"}, nil)
	known := postJSON(router, "/api/auth/password/forgot", map[string]string{"email": "demo@example.com"}, nil)
	if unknown.Code != http.StatusAccepted || known.Code != http.StatusAccepted || unknown.Body.String() != known.Body.String() {
		t.Fatalf("forgot must answer identically: %d %q vs %d %q", unknown.Code, unknown.Body, known.Code, known.Body)
	}
	if mailer.Count() != 1 || mailer.Last().To != "demo@example.com" {
		t.Fatalf("expected exactly one reset email to demo, got %d", mailer.Count())
	}
	token := lastToken(t, mailer)

	bad := postJSON(router, "/api/auth/password/reset", map[string]string{"token": "nope", "password": "brand-new-pass"}, nil)
	if bad.Code != http.StatusBadRequest {
		t.Fatalf("bad token: %d", bad.Code)
	}
	ok := postJSON(router, "/api/auth/password/reset", map[string]string{"token": token, "password": "brand-new-pass"}, nil)
	if ok.Code != http.StatusNoContent {
		t.Fatalf("reset: %d %s", ok.Code, ok.Body)
	}
	if rec := postJSON(router, "/api/auth/login", map[string]string{"email": "demo@example.com", "password": "password123"}, nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("old password must stop working: %d", rec.Code)
	}
	if rec := postJSON(router, "/api/auth/login", map[string]string{"email": "demo@example.com", "password": "brand-new-pass"}, nil); rec.Code != http.StatusOK {
		t.Fatalf("new password login: %d %s", rec.Code, rec.Body)
	}
	if rec := postJSON(router, "/api/auth/password/reset", map[string]string{"token": token, "password": "third-password"}, nil); rec.Code != http.StatusBadRequest {
		t.Fatalf("reused token: %d", rec.Code)
	}
}

func TestRegisterSendsVerificationAndRequireFlagBlocksLogin(t *testing.T) {
	router, mailer := accountRouter(t, func(cfg *config.Config) { cfg.Auth.RequireEmailVerification = true })

	reg := postJSON(router, "/api/auth/register", map[string]string{"name": "New User", "email": "new@example.com", "password": "password-123"}, nil)
	var body struct {
		Code string `json:"code"`
	}
	_ = json.Unmarshal(reg.Body.Bytes(), &body)
	if reg.Code != http.StatusForbidden || body.Code != "email_not_verified" {
		t.Fatalf("register with required verification: %d %s", reg.Code, reg.Body)
	}
	if mailer.Count() != 1 || mailer.Last().To != "new@example.com" {
		t.Fatalf("expected a verification email, got %d", mailer.Count())
	}

	login := postJSON(router, "/api/auth/login", map[string]string{"email": "new@example.com", "password": "password-123"}, nil)
	_ = json.Unmarshal(login.Body.Bytes(), &body)
	if login.Code != http.StatusForbidden || body.Code != "email_not_verified" {
		t.Fatalf("unverified login: %d %s", login.Code, login.Body)
	}

	if rec := postJSON(router, "/api/auth/email/verify", map[string]string{"token": lastToken(t, mailer)}, nil); rec.Code != http.StatusNoContent {
		t.Fatalf("verify: %d %s", rec.Code, rec.Body)
	}
	login = postJSON(router, "/api/auth/login", map[string]string{"email": "new@example.com", "password": "password-123"}, nil)
	var payload struct {
		User struct {
			EmailVerified bool `json:"emailVerified"`
		} `json:"user"`
	}
	_ = json.Unmarshal(login.Body.Bytes(), &payload)
	if login.Code != http.StatusOK || !payload.User.EmailVerified {
		t.Fatalf("verified login: %d %s", login.Code, login.Body)
	}
}

func TestAdminInviteAcceptAndLogin(t *testing.T) {
	router, mailer := accountRouter(t)
	admin := loginHTTP(t, router, "admin@example.com", "admin1234", "", "")
	asAdmin := func(req *http.Request) { req.Header.Set("Authorization", "Bearer "+admin.AccessToken) }

	roles := httpGet(router, "/api/admin/roles/all", admin.AccessToken)
	var roleList []struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	}
	_ = json.Unmarshal(roles.Body.Bytes(), &roleList)
	var memberRole string
	for _, r := range roleList {
		if r.Slug == "user" {
			memberRole = r.ID
		}
	}
	if memberRole == "" {
		t.Fatalf("member role not found: %s", roles.Body)
	}

	invite := postJSON(router, "/api/admin/users/invite", map[string]string{"email": "invitee@example.com", "name": "Invitee", "roleId": memberRole}, asAdmin)
	var created struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		EmailVerified bool   `json:"emailVerified"`
	}
	_ = json.Unmarshal(invite.Body.Bytes(), &created)
	if invite.Code != http.StatusCreated || created.Status != "invited" || created.EmailVerified {
		t.Fatalf("invite: %d %s", invite.Code, invite.Body)
	}
	if rec := postJSON(router, "/api/auth/login", map[string]string{"email": "invitee@example.com", "password": "anything-at-all"}, nil); rec.Code != http.StatusUnauthorized {
		t.Fatalf("invited user cannot log in yet: %d", rec.Code)
	}

	resend := postJSON(router, "/api/admin/users/"+created.ID+"/invite", map[string]string{}, asAdmin)
	if resend.Code != http.StatusAccepted || mailer.Count() != 2 {
		t.Fatalf("resend: %d, sent %d", resend.Code, mailer.Count())
	}
	token := lastToken(t, mailer)

	inspect := postJSON(router, "/api/auth/invite/inspect", map[string]string{"token": token}, nil)
	if inspect.Code != http.StatusOK || !regexp.MustCompile(`"email":"invitee@example.com"`).Match(inspect.Body.Bytes()) {
		t.Fatalf("inspect: %d %s", inspect.Code, inspect.Body)
	}
	accept := postJSON(router, "/api/auth/invite/accept", map[string]string{"token": token, "password": "invitee-pass"}, nil)
	if accept.Code != http.StatusNoContent {
		t.Fatalf("accept: %d %s", accept.Code, accept.Body)
	}
	loginHTTP(t, router, "invitee@example.com", "invitee-pass", "", "")

	if rec := postJSON(router, "/api/admin/users/"+created.ID+"/invite", map[string]string{}, asAdmin); rec.Code != http.StatusConflict {
		t.Fatalf("resend after accept: %d", rec.Code)
	}
	demo := loginHTTP(t, router, "demo@example.com", "password123", "", "")
	asDemo := func(req *http.Request) { req.Header.Set("Authorization", "Bearer "+demo.AccessToken) }
	if rec := postJSON(router, "/api/admin/users/invite", map[string]string{"email": "x@example.com", "name": "X", "roleId": memberRole}, asDemo); rec.Code != http.StatusForbidden {
		t.Fatalf("invite needs users:modify: %d", rec.Code)
	}
}

func httpGet(router *gin.Engine, path, accessToken string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	return rec
}

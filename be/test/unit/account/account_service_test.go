package account_test

import (
	"context"
	"errors"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	apperrors "be/common/errors"
	authmodel "be/internal/models/auth"
	rolemodel "be/internal/models/role"
	usermodel "be/internal/models/user"
	accountsvc "be/internal/services/account"
	"be/pkg/hash"
	"be/test/testutil"
)

type fixture struct {
	svc      *accountsvc.Service
	users    *testutil.MemoryUserRepo
	tokens   *testutil.MockUserTokenRepo
	sessions *testutil.MockAuthRepo
	attempts *testutil.MockLoginAttemptRepo
	mailer   *testutil.RecordingMailer
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	pw, err := hash.HashPassword("old-password")
	if err != nil {
		t.Fatal(err)
	}
	f := &fixture{
		users: &testutil.MemoryUserRepo{Users: map[string]*usermodel.User{
			"u1": {ID: "u1", Email: "user@example.com", FullName: "User One", PasswordHash: pw, RoleID: "r1", Status: usermodel.StatusActive},
		}},
		tokens:   testutil.NewMockUserTokenRepo(),
		sessions: &testutil.MockAuthRepo{RefreshTokens: map[string]*authmodel.RefreshToken{}},
		attempts: testutil.NewMockLoginAttemptRepo(),
		mailer:   &testutil.RecordingMailer{},
	}
	roles := &testutil.MockRoleRepo{Roles: map[string]*rolemodel.Role{"r1": {ID: "r1", Slug: "user", Name: "User"}}}
	f.svc = accountsvc.NewService(f.users, roles, f.tokens, f.sessions, f.attempts, testutil.NoTx{}, f.mailer,
		accountsvc.DefaultConfig("Test App", "https://app.example.com/"))
	return f
}

var tokenInLink = regexp.MustCompile(`https://app\.example\.com(/[a-z-]+)\?token=([A-Za-z0-9_%-]+)`)

// linkFrom extracts the path and raw token from the last email.
func (f *fixture) linkFrom(t *testing.T) (string, string) {
	t.Helper()
	m := tokenInLink.FindStringSubmatch(f.mailer.Last().Text)
	if m == nil {
		t.Fatalf("no link in email: %q", f.mailer.Last().Text)
	}
	token, _ := url.QueryUnescape(m[2])
	return m[1], token
}

func TestPasswordResetFlow(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	_ = f.sessions.CreateRefreshToken(ctx, &authmodel.RefreshToken{UserID: "u1", RawToken: "s1", ExpiresAt: time.Now().Add(time.Hour)})
	f.attempts.Failed["u1"] = 3

	if err := f.svc.RequestPasswordReset(ctx, "user@example.com"); err != nil {
		t.Fatal(err)
	}
	path, token := f.linkFrom(t)
	if path != "/reset-password" || f.mailer.Last().To != "user@example.com" || !strings.Contains(f.mailer.Last().Subject, "Test App") {
		t.Fatalf("unexpected email: %+v", f.mailer.Last())
	}
	for _, stored := range f.tokens.Tokens {
		if stored.TokenHash == token || stored.TokenHash != hash.SHA256Hex(token) {
			t.Fatal("only the hash of the token may be stored")
		}
	}

	if err := f.svc.ResetPassword(ctx, token, "new-password-123"); err != nil {
		t.Fatalf("reset: %v", err)
	}
	user := f.users.Users["u1"]
	if hash.ComparePassword(user.PasswordHash, "new-password-123") != nil {
		t.Fatal("password not changed")
	}
	if user.EmailVerifiedAt == nil {
		t.Fatal("a reset link proves the mailbox; email should be verified")
	}
	if active, _ := f.sessions.ListActiveByUserID(ctx, "u1"); len(active) != 0 {
		t.Fatalf("reset must sign out every session, %d left", len(active))
	}
	if _, ok := f.attempts.Failed["u1"]; ok {
		t.Fatal("reset must clear the lockout counter")
	}

	err := f.svc.ResetPassword(ctx, token, "another-password")
	if !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("a used link must be rejected, got %v", err)
	}
}

func TestPasswordResetDoesNotRevealUnknownOrInactiveAccounts(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	if err := f.svc.RequestPasswordReset(ctx, "nobody@example.com"); err != nil {
		t.Fatalf("unknown email must not error: %v", err)
	}
	f.users.Users["u1"].Status = usermodel.StatusBlocked
	if err := f.svc.RequestPasswordReset(ctx, "user@example.com"); err != nil {
		t.Fatal(err)
	}
	if f.mailer.Count() != 0 {
		t.Fatalf("no email for unknown/blocked accounts, sent %d", f.mailer.Count())
	}
}

func TestResendCooldownAndNewestLinkWins(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	now := time.Now()
	f.svc.SetClock(func() time.Time { return now })

	_ = f.svc.RequestPasswordReset(ctx, "user@example.com")
	_, first := f.linkFrom(t)
	_ = f.svc.RequestPasswordReset(ctx, "user@example.com")
	if f.mailer.Count() != 1 {
		t.Fatalf("second request inside the cooldown must not send, sent %d", f.mailer.Count())
	}

	now = now.Add(2 * time.Minute)
	_ = f.svc.RequestPasswordReset(ctx, "user@example.com")
	_, second := f.linkFrom(t)
	if f.mailer.Count() != 2 || first == second {
		t.Fatal("after the cooldown a new link is sent")
	}
	if err := f.svc.ResetPassword(ctx, first, "new-password-123"); !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("older link must stop working, got %v", err)
	}
	if err := f.svc.ResetPassword(ctx, second, "new-password-123"); err != nil {
		t.Fatalf("newest link must work: %v", err)
	}
}

func TestExpiredTokenIsRejected(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	_ = f.svc.RequestPasswordReset(ctx, "user@example.com")
	_, token := f.linkFrom(t)
	for _, stored := range f.tokens.Tokens {
		stored.ExpiresAt = time.Now().Add(-time.Second)
	}
	if err := f.svc.ResetPassword(ctx, token, "new-password-123"); !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("expired link must be rejected, got %v", err)
	}
}

func TestEmailVerification(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()
	if err := f.svc.SendVerification(ctx, "u1"); err != nil {
		t.Fatal(err)
	}
	path, token := f.linkFrom(t)
	if path != "/verify-email" {
		t.Fatalf("path %q", path)
	}
	// A verify token cannot be used as a reset token.
	if err := f.svc.ResetPassword(ctx, token, "new-password-123"); !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("token purpose must be enforced, got %v", err)
	}
	if err := f.svc.VerifyEmail(ctx, token); err != nil {
		t.Fatal(err)
	}
	if f.users.Users["u1"].EmailVerifiedAt == nil {
		t.Fatal("email not marked verified")
	}
	if err := f.svc.SendVerification(ctx, "u1"); err != nil || f.mailer.Count() != 1 {
		t.Fatalf("verified users get no new email (err %v, sent %d)", err, f.mailer.Count())
	}
}

func TestInviteInspectAccept(t *testing.T) {
	f := newFixture(t)
	ctx := context.Background()

	if _, err := f.svc.Invite(ctx, accountsvc.InviteRequest{Email: "user@example.com", FullName: "Dup", RoleID: "r1"}, "admin@example.com"); !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("existing email must conflict, got %v", err)
	}
	if _, err := f.svc.Invite(ctx, accountsvc.InviteRequest{Email: "new@example.com", FullName: "New", RoleID: "missing"}, ""); !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("unknown role must be rejected, got %v", err)
	}

	invited, err := f.svc.Invite(ctx, accountsvc.InviteRequest{Email: "new@example.com", FullName: "New Person", RoleID: "r1"}, "admin@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if invited.Status != usermodel.StatusInvited || invited.PasswordHash != "" {
		t.Fatalf("invited user must have no password yet: %+v", invited)
	}
	path, token := f.linkFrom(t)
	if path != "/accept-invite" || !strings.Contains(f.mailer.Last().Text, "admin@example.com invited you") {
		t.Fatalf("unexpected invite email: %q", f.mailer.Last().Text)
	}

	info, err := f.svc.InspectInvite(ctx, token)
	if err != nil || info.Email != "new@example.com" || info.FullName != "New Person" {
		t.Fatalf("inspect: %+v %v", info, err)
	}

	if err := f.svc.AcceptInvite(ctx, token, "chosen-password", "Renamed"); err != nil {
		t.Fatal(err)
	}
	user := f.users.Users[invited.ID]
	if user.Status != usermodel.StatusActive || user.FullName != "Renamed" || user.EmailVerifiedAt == nil ||
		hash.ComparePassword(user.PasswordHash, "chosen-password") != nil {
		t.Fatalf("accept did not activate the account: %+v", user)
	}
	if _, err := f.svc.InspectInvite(ctx, token); !errors.Is(err, apperrors.ErrBadRequest) {
		t.Fatalf("used invite link must be invalid, got %v", err)
	}
	if err := f.svc.ResendInvite(ctx, invited.ID, ""); !errors.Is(err, apperrors.ErrConflict) {
		t.Fatalf("resend to an active user must conflict, got %v", err)
	}
}

// Package account owns the email-driven account lifecycle: password reset,
// email verification and admin invitations. Every flow uses a hashed,
// expiring, single-use token delivered by mail.
package account

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"

	apperrors "be/common/errors"
	"be/internal/mail"
	accountmodel "be/internal/models/account"
	usermodel "be/internal/models/user"
	"be/internal/repository/interfaces"
	"be/pkg/hash"
)

// Config holds link and lifetime settings.
type Config struct {
	AppName string
	// AppURL is the FE origin; links are AppURL + "/reset-password?token=…" etc.
	AppURL    string
	ResetTTL  time.Duration
	VerifyTTL time.Duration
	InviteTTL time.Duration
	// ResendCooldown skips sending a new email of the same kind to the same user
	// within this window (anti mail-bombing).
	ResendCooldown time.Duration
}

// DefaultConfig returns the standard lifetimes.
func DefaultConfig(appName, appURL string) Config {
	return Config{
		AppName:        appName,
		AppURL:         strings.TrimRight(appURL, "/"),
		ResetTTL:       time.Hour,
		VerifyTTL:      48 * time.Hour,
		InviteTTL:      7 * 24 * time.Hour,
		ResendCooldown: time.Minute,
	}
}

// Transactor runs fn in one database transaction (postgres.Postgres.WithTx).
type Transactor interface {
	WithTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type Service struct {
	users    interfaces.UserRepository
	roles    interfaces.RoleRepository
	tokens   interfaces.UserTokenRepository
	sessions interfaces.AuthRepository
	attempts interfaces.LoginAttemptRepository // optional
	tx       Transactor
	mailer   mail.Mailer
	cfg      Config
	now      func() time.Time
}

func NewService(
	users interfaces.UserRepository,
	roles interfaces.RoleRepository,
	tokens interfaces.UserTokenRepository,
	sessions interfaces.AuthRepository,
	attempts interfaces.LoginAttemptRepository,
	tx Transactor,
	mailer mail.Mailer,
	cfg Config,
) *Service {
	return &Service{
		users: users, roles: roles, tokens: tokens, sessions: sessions, attempts: attempts,
		tx: tx, mailer: mailer, cfg: cfg, now: time.Now,
	}
}

var errInvalidToken = fmt.Errorf("%w: this link is invalid or has expired", apperrors.ErrBadRequest)

// RequestPasswordReset emails a reset link. It never reveals whether the email
// exists: unknown, invited, blocked or OAuth-only addresses simply get nothing.
func (s *Service) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.users.GetByEmail(ctx, strings.TrimSpace(email))
	if err != nil || user == nil {
		return err
	}
	if user.Status != usermodel.StatusActive {
		return nil
	}
	return s.issueAndSend(ctx, user, accountmodel.PurposePasswordReset, s.cfg.ResetTTL, "/reset-password", mail.PasswordResetEmail, "")
}

// ResetPassword sets a new password from a reset link, signs out every session,
// clears any login lockout and marks the email verified (the link proved it).
func (s *Service) ResetPassword(ctx context.Context, rawToken, newPassword string) error {
	hashed, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}
	return s.redeem(ctx, accountmodel.PurposePasswordReset, rawToken, func(ctx context.Context, user *usermodel.User) error {
		user.PasswordHash = hashed
		s.markVerified(user)
		if err := s.users.Update(ctx, user); err != nil {
			return err
		}
		if err := s.revokeAllSessions(ctx, user.ID); err != nil {
			return err
		}
		if s.attempts != nil {
			return s.attempts.Reset(ctx, user.ID)
		}
		return nil
	})
}

// SendVerification emails a verification link to the user. No-op when the
// email is already verified or a link was sent within the cooldown.
func (s *Service) SendVerification(ctx context.Context, userID string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrNotFound
	}
	if user.EmailVerifiedAt != nil {
		return nil
	}
	return s.issueAndSend(ctx, user, accountmodel.PurposeEmailVerify, s.cfg.VerifyTTL, "/verify-email", mail.VerifyEmail, "")
}

// VerifyEmail marks the email verified from a verification link.
func (s *Service) VerifyEmail(ctx context.Context, rawToken string) error {
	return s.redeem(ctx, accountmodel.PurposeEmailVerify, rawToken, func(ctx context.Context, user *usermodel.User) error {
		if user.EmailVerifiedAt != nil {
			return nil
		}
		s.markVerified(user)
		return s.users.Update(ctx, user)
	})
}

// InviteRequest is an admin invitation.
type InviteRequest struct {
	Email    string
	FullName string
	RoleID   string
}

// Invite creates an account in status "invited" (no password) and emails a link
// to set one.
func (s *Service) Invite(ctx context.Context, req InviteRequest, invitedBy string) (*usermodel.User, error) {
	email := strings.TrimSpace(req.Email)
	existing, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: email already exists", apperrors.ErrConflict)
	}
	role, err := s.roles.GetByID(ctx, req.RoleID)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("%w: role not found", apperrors.ErrBadRequest)
	}

	user := &usermodel.User{
		Email:       email,
		FullName:    strings.TrimSpace(req.FullName),
		RoleID:      role.ID,
		Status:      usermodel.StatusInvited,
		SocialLinks: []usermodel.SocialLink{},
	}
	err = s.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := s.users.Create(ctx, user); err != nil {
			return err
		}
		return s.send(ctx, user, accountmodel.PurposeInvite, s.cfg.InviteTTL, "/accept-invite", mail.InviteEmail, invitedBy)
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ResendInvite issues a fresh invitation link; older links stop working.
func (s *Service) ResendInvite(ctx context.Context, userID, invitedBy string) error {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrNotFound
	}
	if user.Status != usermodel.StatusInvited {
		return fmt.Errorf("%w: user has already accepted the invitation", apperrors.ErrConflict)
	}
	// An explicit admin resend skips the cooldown.
	return s.send(ctx, user, accountmodel.PurposeInvite, s.cfg.InviteTTL, "/accept-invite", mail.InviteEmail, invitedBy)
}

// Invitation is what the accept page shows before the user sets a password.
type Invitation struct {
	Email    string `json:"email"`
	FullName string `json:"name"`
}

// InspectInvite returns the invited account for a valid link without using it.
func (s *Service) InspectInvite(ctx context.Context, rawToken string) (*Invitation, error) {
	token, err := s.tokens.FindValid(ctx, accountmodel.PurposeInvite, hash.SHA256Hex(strings.TrimSpace(rawToken)))
	if err != nil {
		return nil, err
	}
	if token == nil {
		return nil, errInvalidToken
	}
	user, err := s.users.GetByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Status != usermodel.StatusInvited {
		return nil, errInvalidToken
	}
	return &Invitation{Email: user.Email, FullName: user.FullName}, nil
}

// AcceptInvite sets the password (and optionally the name) and activates the account.
func (s *Service) AcceptInvite(ctx context.Context, rawToken, password, fullName string) error {
	hashed, err := hash.HashPassword(password)
	if err != nil {
		return err
	}
	return s.redeem(ctx, accountmodel.PurposeInvite, rawToken, func(ctx context.Context, user *usermodel.User) error {
		if user.Status != usermodel.StatusInvited {
			return errInvalidToken
		}
		user.PasswordHash = hashed
		user.Status = usermodel.StatusActive
		if name := strings.TrimSpace(fullName); name != "" {
			user.FullName = name
		}
		s.markVerified(user)
		return s.users.Update(ctx, user)
	})
}

type renderFunc func(to string, data mail.AccountEmail) (mail.Message, error)

// issueAndSend sends a new link unless one of the same kind went out within the cooldown.
func (s *Service) issueAndSend(ctx context.Context, user *usermodel.User, purpose accountmodel.TokenPurpose, ttl time.Duration, path string, render renderFunc, invitedBy string) error {
	if s.cfg.ResendCooldown > 0 {
		last, err := s.tokens.LatestCreatedAt(ctx, user.ID, purpose)
		if err != nil {
			return err
		}
		if last != nil && s.now().Sub(*last) < s.cfg.ResendCooldown {
			slog.InfoContext(ctx, "account email skipped (cooldown)", "purpose", purpose, "user_id", user.ID)
			return nil
		}
	}
	return s.send(ctx, user, purpose, ttl, path, render, invitedBy)
}

func (s *Service) send(ctx context.Context, user *usermodel.User, purpose accountmodel.TokenPurpose, ttl time.Duration, path string, render renderFunc, invitedBy string) error {
	raw, err := newRawToken()
	if err != nil {
		return err
	}
	// Only the newest link of a kind works.
	if err := s.tokens.InvalidateOpen(ctx, user.ID, purpose); err != nil {
		return err
	}
	now := s.now().UTC()
	if err := s.tokens.Create(ctx, &accountmodel.UserToken{
		UserID:    user.ID,
		Purpose:   purpose,
		TokenHash: hash.SHA256Hex(raw),
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}); err != nil {
		return err
	}

	msg, err := render(user.Email, mail.AccountEmail{
		AppName:   s.cfg.AppName,
		Name:      firstNonEmpty(user.FullName, user.Email),
		Link:      s.cfg.AppURL + path + "?token=" + url.QueryEscape(raw),
		ExpiresIn: humanDuration(ttl),
		InvitedBy: invitedBy,
	})
	if err != nil {
		return err
	}
	return s.mailer.Send(ctx, msg)
}

// redeem validates and consumes a token, then applies fn to its user, all in one
// transaction so a failed update does not burn the link.
func (s *Service) redeem(ctx context.Context, purpose accountmodel.TokenPurpose, rawToken string, fn func(ctx context.Context, user *usermodel.User) error) error {
	tokenHash := hash.SHA256Hex(strings.TrimSpace(rawToken))
	return s.tx.WithTx(ctx, func(ctx context.Context) error {
		token, err := s.tokens.FindValid(ctx, purpose, tokenHash)
		if err != nil {
			return err
		}
		if token == nil {
			return errInvalidToken
		}
		consumed, err := s.tokens.Consume(ctx, token.ID)
		if err != nil {
			return err
		}
		if !consumed {
			return errInvalidToken
		}
		user, err := s.users.GetByID(ctx, token.UserID)
		if err != nil {
			return err
		}
		if user == nil {
			return errInvalidToken
		}
		return fn(ctx, user)
	})
}

// revokeAllSessions signs the user out everywhere (after a password reset).
func (s *Service) revokeAllSessions(ctx context.Context, userID string) error {
	active, err := s.sessions.ListActiveByUserID(ctx, userID)
	if err != nil {
		return err
	}
	for _, session := range active {
		if _, err := s.sessions.RevokeByID(ctx, userID, session.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) markVerified(user *usermodel.User) {
	if user.EmailVerifiedAt == nil {
		now := s.now().UTC()
		user.EmailVerifiedAt = &now
	}
}

// SetClock replaces the time source (tests).
func (s *Service) SetClock(now func() time.Time) { s.now = now }

func newRawToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func humanDuration(d time.Duration) string {
	switch {
	case d%(24*time.Hour) == 0 && d >= 24*time.Hour:
		return plural(int(d/(24*time.Hour)), "day")
	case d%time.Hour == 0 && d >= time.Hour:
		return plural(int(d/time.Hour), "hour")
	default:
		return plural(int(d/time.Minute), "minute")
	}
}

func plural(n int, unit string) string {
	if n == 1 {
		return "1 " + unit
	}
	return fmt.Sprintf("%d %ss", n, unit)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

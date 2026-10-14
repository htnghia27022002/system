package dependency

import (
	"log/slog"

	"be/internal/mail"
	"be/internal/repository"
	accountsvc "be/internal/services/account"
)

// NewMailer builds the MAIL_DRIVER transport. Production sends asynchronously so
// request latency (and timing) never depends on the mail server.
func NewMailer(infra *Infra) (mail.Mailer, error) {
	m, err := mail.New(infra.Config.Mail, slog.Default())
	if err != nil {
		return nil, err
	}
	return mail.Async(m, slog.Default()), nil
}

// NewAccountService wires password reset, email verification and invitations,
// and hooks verification into sign-up / login.
func NewAccountService(infra *Infra, auth *AuthServices, mailer mail.Mailer) *accountsvc.Service {
	svc := accountsvc.NewService(
		newUserRepository(infra.DB),
		newRoleRepository(infra.DB),
		repository.NewUserTokenRepository(infra.DB),
		newAuthRepository(infra.DB),
		repository.NewLoginAttemptRepository(infra.DB),
		infra.DB,
		mailer,
		accountsvc.DefaultConfig(infra.Config.AppName, infra.Config.AppURL),
	)
	auth.Auth.SetEmailVerification(svc, infra.Config.Auth.RequireEmailVerification)
	return svc
}

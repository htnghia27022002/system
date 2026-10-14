// Package mail sends transactional email. Transports: "log" (dev — writes the
// message to the structured log) and "smtp". Services depend on Mailer only.
package mail

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"be/internal/config"
)

// Message is one outgoing email. HTML is optional; Text is always sent.
type Message struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

// Mailer delivers a message.
type Mailer interface {
	Send(ctx context.Context, msg Message) error
}

// New builds the transport selected by MAIL_DRIVER.
func New(cfg config.MailConfig, logger *slog.Logger) (Mailer, error) {
	switch cfg.Driver {
	case "", "log":
		return LogMailer{Logger: logger}, nil
	case "smtp":
		if cfg.SMTPHost == "" {
			return nil, fmt.Errorf("mail: MAIL_DRIVER=smtp requires SMTP_HOST")
		}
		return NewSMTPMailer(cfg), nil
	default:
		return nil, fmt.Errorf("mail: unknown MAIL_DRIVER %q (want log or smtp)", cfg.Driver)
	}
}

// LogMailer writes messages to the log instead of sending them. Use it in dev
// and tests: links in account emails show up in the BE log.
type LogMailer struct {
	Logger *slog.Logger
}

func (m LogMailer) Send(ctx context.Context, msg Message) error {
	logger := m.Logger
	if logger == nil {
		logger = slog.Default()
	}
	logger.InfoContext(ctx, "mail (log driver)", "to", msg.To, "subject", msg.Subject, "text", msg.Text)
	return nil
}

// asyncMailer hands each message to a goroutine so request latency (and timing)
// does not depend on the mail server. Failures are logged.
type asyncMailer struct {
	next    Mailer
	logger  *slog.Logger
	timeout time.Duration
}

// Async wraps m so Send returns immediately.
func Async(m Mailer, logger *slog.Logger) Mailer {
	if logger == nil {
		logger = slog.Default()
	}
	return asyncMailer{next: m, logger: logger, timeout: 30 * time.Second}
}

func (a asyncMailer) Send(_ context.Context, msg Message) error {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), a.timeout)
		defer cancel()
		if err := a.next.Send(ctx, msg); err != nil {
			a.logger.Error("mail send failed", "to", msg.To, "subject", msg.Subject, "err", err)
		}
	}()
	return nil
}

// sanitizeHeader strips CR/LF so user-controlled values cannot inject headers.
func sanitizeHeader(v string) string {
	return strings.NewReplacer("\r", " ", "\n", " ").Replace(v)
}

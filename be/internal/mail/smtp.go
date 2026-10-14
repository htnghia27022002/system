package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"time"

	"be/internal/config"
)

// SMTPMailer sends through an SMTP relay (Mailpit in dev, SES/Postmark/... in prod).
type SMTPMailer struct {
	cfg config.MailConfig
}

func NewSMTPMailer(cfg config.MailConfig) *SMTPMailer {
	return &SMTPMailer{cfg: cfg}
}

func (m *SMTPMailer) Send(ctx context.Context, msg Message) error {
	from, err := mail.ParseAddress(m.cfg.From)
	if err != nil {
		return fmt.Errorf("mail: invalid MAIL_FROM: %w", err)
	}
	to, err := mail.ParseAddress(msg.To)
	if err != nil {
		return fmt.Errorf("mail: invalid recipient: %w", err)
	}
	body, err := buildMIME(from, to, msg)
	if err != nil {
		return err
	}

	addr := net.JoinHostPort(m.cfg.SMTPHost, strconv.Itoa(m.cfg.SMTPPort))
	dialer := &net.Dialer{Timeout: 10 * time.Second}
	tlsConfig := &tls.Config{ServerName: m.cfg.SMTPHost, MinVersion: tls.VersionTLS12}

	var conn net.Conn
	if m.cfg.SMTPTLS == "tls" {
		conn, err = tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	} else {
		conn, err = dialer.DialContext(ctx, "tcp", addr)
	}
	if err != nil {
		return fmt.Errorf("mail: dial %s: %w", addr, err)
	}
	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}

	client, err := smtp.NewClient(conn, m.cfg.SMTPHost)
	if err != nil {
		conn.Close()
		return fmt.Errorf("mail: smtp handshake: %w", err)
	}
	defer client.Close()

	if m.cfg.SMTPTLS == "starttls" {
		if ok, _ := client.Extension("STARTTLS"); !ok {
			return fmt.Errorf("mail: server does not offer STARTTLS (set SMTP_TLS=none for a local relay)")
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("mail: starttls: %w", err)
		}
	}
	if m.cfg.SMTPUsername != "" {
		auth := smtp.PlainAuth("", m.cfg.SMTPUsername, m.cfg.SMTPPassword, m.cfg.SMTPHost)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("mail: auth: %w", err)
		}
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("mail: MAIL FROM: %w", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("mail: RCPT TO: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mail: DATA: %w", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("mail: write body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail: end DATA: %w", err)
	}
	return client.Quit()
}

// buildMIME renders a multipart/alternative message (text + optional HTML).
func buildMIME(from, to *mail.Address, msg Message) ([]byte, error) {
	var buf bytes.Buffer
	header := func(k, v string) { fmt.Fprintf(&buf, "%s: %s\r\n", k, v) }

	header("From", from.String())
	header("To", to.String())
	header("Subject", mime.QEncoding.Encode("utf-8", sanitizeHeader(msg.Subject)))
	header("Date", time.Now().Format(time.RFC1123Z))
	header("Message-ID", fmt.Sprintf("<%s@%s>", randomID(), domainOf(from.Address)))
	header("MIME-Version", "1.0")

	part := func(contentType, content string) error {
		header("Content-Type", contentType+"; charset=utf-8")
		header("Content-Transfer-Encoding", "quoted-printable")
		buf.WriteString("\r\n")
		qp := quotedprintable.NewWriter(&buf)
		if _, err := qp.Write([]byte(content)); err != nil {
			return err
		}
		if err := qp.Close(); err != nil {
			return err
		}
		buf.WriteString("\r\n")
		return nil
	}

	if msg.HTML == "" {
		if err := part("text/plain", msg.Text); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	boundary := "b_" + randomID()
	header("Content-Type", fmt.Sprintf("multipart/alternative; boundary=%q", boundary))
	buf.WriteString("\r\n")
	fmt.Fprintf(&buf, "--%s\r\n", boundary)
	if err := part("text/plain", msg.Text); err != nil {
		return nil, err
	}
	fmt.Fprintf(&buf, "--%s\r\n", boundary)
	if err := part("text/html", msg.HTML); err != nil {
		return nil, err
	}
	fmt.Fprintf(&buf, "--%s--\r\n", boundary)
	return buf.Bytes(), nil
}

func randomID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func domainOf(address string) string {
	for i := len(address) - 1; i >= 0; i-- {
		if address[i] == '@' {
			return address[i+1:]
		}
	}
	return "localhost"
}

package mail_test

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/textproto"
	"strconv"
	"strings"
	"testing"
	"time"

	"be/internal/config"
	"be/internal/mail"
)

func TestNewRejectsBadConfig(t *testing.T) {
	if _, err := mail.New(config.MailConfig{Driver: "smtp"}, nil); err == nil {
		t.Fatal("smtp without host must fail fast")
	}
	if _, err := mail.New(config.MailConfig{Driver: "carrier-pigeon"}, nil); err == nil {
		t.Fatal("unknown driver must fail fast")
	}
	if m, err := mail.New(config.MailConfig{}, nil); err != nil || m == nil {
		t.Fatalf("default driver is log: %v", err)
	}
}

func TestTemplatesEscapeHTMLAndKeepLink(t *testing.T) {
	msg, err := mail.InviteEmail("to@example.com", mail.AccountEmail{
		AppName: "Acme", Name: "<script>x</script>", Link: "https://app.example.com/accept-invite?token=abc", ExpiresIn: "7 days",
	})
	if err != nil {
		t.Fatal(err)
	}
	if msg.Subject != "You are invited to Acme" {
		t.Fatalf("subject %q", msg.Subject)
	}
	if strings.Contains(msg.HTML, "<script>") || !strings.Contains(msg.HTML, "&lt;script&gt;") {
		t.Fatal("names must be HTML-escaped")
	}
	if !strings.Contains(msg.Text, "https://app.example.com/accept-invite?token=abc") || !strings.Contains(msg.HTML, "token=abc") {
		t.Fatal("link missing")
	}
}

// fakeSMTP accepts one plaintext session and returns what it received.
func fakeSMTP(t *testing.T) (port int, received chan string) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { ln.Close() })
	received = make(chan string, 1)
	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()
		tp := textproto.NewConn(conn)
		var transcript strings.Builder
		_ = tp.PrintfLine("220 fake ESMTP")
		for {
			line, err := tp.ReadLine()
			if err != nil {
				received <- transcript.String()
				return
			}
			transcript.WriteString(line + "\n")
			switch cmd := strings.ToUpper(strings.SplitN(line, " ", 2)[0]); cmd {
			case "EHLO", "HELO":
				_ = tp.PrintfLine("250 fake")
			case "DATA":
				_ = tp.PrintfLine("354 go ahead")
				body, _ := io.ReadAll(tp.DotReader())
				transcript.Write(body)
				_ = tp.PrintfLine("250 queued")
			case "QUIT":
				_ = tp.PrintfLine("221 bye")
				received <- transcript.String()
				return
			default:
				_ = tp.PrintfLine("250 ok")
			}
		}
	}()
	_, portStr, _ := net.SplitHostPort(ln.Addr().String())
	port, _ = strconv.Atoi(portStr)
	return port, received
}

func TestSMTPMailerSendsMultipartMessage(t *testing.T) {
	port, received := fakeSMTP(t)
	m := mail.NewSMTPMailer(config.MailConfig{
		From: "Acme <no-reply@acme.test>", SMTPHost: "127.0.0.1", SMTPPort: port, SMTPTLS: "none",
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := m.Send(ctx, mail.Message{
		To: "user@example.com", Subject: "Hello\r\nBcc: evil@example.com", Text: "plain body", HTML: "<p>html body</p>",
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}
	got := <-received
	for _, want := range []string{"MAIL FROM:<no-reply@acme.test>", "RCPT TO:<user@example.com>", "multipart/alternative", "plain body", "html body"} {
		if !strings.Contains(got, want) {
			t.Fatalf("transcript missing %q:\n%s", want, got)
		}
	}
	r := textproto.NewReader(bufio.NewReader(strings.NewReader(got[strings.Index(got, "From:"):])))
	headers, _ := r.ReadMIMEHeader()
	if headers.Get("Bcc") != "" {
		t.Fatal("CR/LF in the subject must not inject headers")
	}
}

func TestSMTPMailerRequiresStartTLSWhenConfigured(t *testing.T) {
	port, _ := fakeSMTP(t)
	m := mail.NewSMTPMailer(config.MailConfig{From: "a@acme.test", SMTPHost: "127.0.0.1", SMTPPort: port, SMTPTLS: "starttls"})
	err := m.Send(context.Background(), mail.Message{To: "user@example.com", Subject: "x", Text: "x"})
	if err == nil || !strings.Contains(err.Error(), "STARTTLS") {
		t.Fatalf("expected a STARTTLS error, got %v", err)
	}
}

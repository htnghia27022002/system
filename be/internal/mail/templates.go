package mail

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"
)

// AccountEmail is the data for account lifecycle emails (reset, verify, invite).
type AccountEmail struct {
	AppName string
	Name    string
	Link    string
	// ExpiresIn is a human label, e.g. "1 hour".
	ExpiresIn string
	// InvitedBy is set for invitations.
	InvitedBy string
}

type template struct {
	subject string
	text    *texttemplate.Template
	html    *htmltemplate.Template
}

func newTemplate(subject, text, html string) template {
	return template{
		subject: subject,
		text:    texttemplate.Must(texttemplate.New("text").Parse(text)),
		html:    htmltemplate.Must(htmltemplate.New("html").Parse(layoutHTML(html))),
	}
}

func (t template) render(to string, data AccountEmail) (Message, error) {
	var text, html bytes.Buffer
	if err := t.text.Execute(&text, data); err != nil {
		return Message{}, err
	}
	if err := t.html.Execute(&html, data); err != nil {
		return Message{}, err
	}
	return Message{To: to, Subject: fmt.Sprintf(t.subject, data.AppName), Text: text.String(), HTML: html.String()}, nil
}

func layoutHTML(content string) string {
	return `<!doctype html><html><body style="margin:0;padding:24px;background:#f6f7f9;font-family:system-ui,-apple-system,Segoe UI,Roboto,sans-serif;color:#111827">
<div style="max-width:480px;margin:0 auto;background:#ffffff;border:1px solid #e5e7eb;border-radius:12px;padding:24px">
<p style="margin:0 0 16px;font-weight:600">{{.AppName}}</p>` + content + `
</div></body></html>`
}

const buttonHTML = `<p style="margin:24px 0"><a href="{{.Link}}" style="display:inline-block;background:#3d72e0;color:#ffffff;text-decoration:none;padding:10px 16px;border-radius:8px">%s</a></p>
<p style="font-size:12px;color:#6b7280">Or open this link: {{.Link}}</p>`

var (
	passwordResetTemplate = newTemplate(
		"Reset your %s password",
		`Hi {{.Name}},

Someone asked to reset the password for your {{.AppName}} account.
Open this link to choose a new password (valid for {{.ExpiresIn}}):

{{.Link}}

If you did not ask for this, ignore this email; your password stays the same.
`,
		`<p>Hi {{.Name}},</p><p>Someone asked to reset the password for your account. The link is valid for {{.ExpiresIn}}.</p>`+
			fmt.Sprintf(buttonHTML, "Reset password")+
			`<p style="font-size:12px;color:#6b7280">If you did not ask for this, ignore this email; your password stays the same.</p>`,
	)
	verifyEmailTemplate = newTemplate(
		"Verify your %s email",
		`Hi {{.Name}},

Confirm this email address for your {{.AppName}} account (link valid for {{.ExpiresIn}}):

{{.Link}}
`,
		`<p>Hi {{.Name}},</p><p>Confirm this email address for your account. The link is valid for {{.ExpiresIn}}.</p>`+
			fmt.Sprintf(buttonHTML, "Verify email"),
	)
	inviteTemplate = newTemplate(
		"You are invited to %s",
		`Hi {{.Name}},

{{if .InvitedBy}}{{.InvitedBy}} invited you{{else}}You are invited{{end}} to join {{.AppName}}.
Open this link to set your password (valid for {{.ExpiresIn}}):

{{.Link}}
`,
		`<p>Hi {{.Name}},</p><p>{{if .InvitedBy}}{{.InvitedBy}} invited you{{else}}You are invited{{end}} to join {{.AppName}}. The link is valid for {{.ExpiresIn}}.</p>`+
			fmt.Sprintf(buttonHTML, "Accept invitation"),
	)
)

// PasswordResetEmail renders the "reset your password" email.
func PasswordResetEmail(to string, data AccountEmail) (Message, error) {
	return passwordResetTemplate.render(to, data)
}

// VerifyEmail renders the "verify your email" email.
func VerifyEmail(to string, data AccountEmail) (Message, error) {
	return verifyEmailTemplate.render(to, data)
}

// InviteEmail renders the invitation email.
func InviteEmail(to string, data AccountEmail) (Message, error) {
	return inviteTemplate.render(to, data)
}

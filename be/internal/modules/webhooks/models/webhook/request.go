package webhook

import "time"

const (
	CaptureStatusOK        = "ok"
	CaptureStatusOversized = "oversized"
	CaptureStatusError     = "error"
)

// Request is one inbound HTTP call captured into an inbox.
type Request struct {
	ID            string         `json:"id" db:"id"`
	InboxID       string         `json:"inboxId" db:"inbox_id"`
	Method        string         `json:"method" db:"method"`
	URL           string         `json:"url" db:"url"`
	ClientIP      string         `json:"clientIp" db:"client_ip"`
	Headers       map[string]any `json:"headers" db:"headers"`
	Query         map[string]any `json:"query" db:"query"`
	Form          map[string]any `json:"form" db:"form"`
	Body          []byte         `json:"-" db:"body"`
	ContentType   string         `json:"contentType" db:"content_type"`
	BodyTruncated bool           `json:"bodyTruncated" db:"body_truncated"`
	CaptureStatus string         `json:"captureStatus" db:"capture_status"`
	IsRead        bool           `json:"isRead" db:"is_read"`
	SoftDeletedAt *time.Time     `json:"softDeletedAt,omitempty" db:"soft_deleted_at"`
	CreatedAt     time.Time      `json:"createdAt" db:"created_at"`
}

func (Request) TableName() string {
	return "webhook_requests"
}

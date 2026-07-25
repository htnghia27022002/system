package webhook

import "time"

const (
	CaptureStatusOK        = "ok"
	CaptureStatusOversized = "oversized"
	CaptureStatusError     = "error"
)

// Request is one inbound HTTP call captured into an inbox.
type Request struct {
	ID            string         `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	InboxID       string         `json:"inboxId" gorm:"column:inbox_id;type:uuid;not null;index"`
	Method        string         `json:"method" gorm:"type:varchar(16);not null"`
	URL           string         `json:"url" gorm:"type:text;not null"`
	ClientIP      string         `json:"clientIp" gorm:"column:client_ip;type:varchar(64);not null;default:''"`
	Headers       map[string]any `json:"headers" gorm:"type:jsonb;serializer:json;not null"`
	Query         map[string]any `json:"query" gorm:"type:jsonb;serializer:json;not null"`
	Form          map[string]any `json:"form" gorm:"type:jsonb;serializer:json;not null"`
	Body          []byte         `json:"-" gorm:"type:bytea"`
	ContentType   string         `json:"contentType" gorm:"column:content_type;type:varchar(255);not null;default:''"`
	BodyTruncated bool           `json:"bodyTruncated" gorm:"column:body_truncated;not null;default:false"`
	CaptureStatus string         `json:"captureStatus" gorm:"column:capture_status;type:varchar(32);not null;default:'ok'"`
	IsRead        bool           `json:"isRead" gorm:"column:is_read;not null;default:false"`
	SoftDeletedAt *time.Time     `json:"softDeletedAt,omitempty" gorm:"column:soft_deleted_at"`
	CreatedAt     time.Time      `json:"createdAt" gorm:"autoCreateTime"`
}

func (Request) TableName() string {
	return "webhook_requests"
}

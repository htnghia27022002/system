package interfaces

import (
	"context"

	webhookmodel "be/internal/models/webhook"
)

// WebhookInboxRepository persists per-account webhook inboxes.
type WebhookInboxRepository interface {
	GetByUserID(ctx context.Context, userID string) (*webhookmodel.Inbox, error)
	GetByPublicUUID(ctx context.Context, publicUUID string) (*webhookmodel.Inbox, error)
	Create(ctx context.Context, inbox *webhookmodel.Inbox) error
	Update(ctx context.Context, inbox *webhookmodel.Inbox) error
	UpdateCounters(ctx context.Context, inboxID string, lifetimeReceived, activeCount int) error
}

// WebhookRequestListFilter filters active (non-soft-deleted) requests.
type WebhookRequestListFilter struct {
	InboxID string
	Method  string
	Q       string
	// ReadFilter: "" or "all" = no filter; "read"; "unread".
	ReadFilter string
	Page       int
	Limit      int
}

// WebhookRequestRepository persists captured inbound requests.
type WebhookRequestRepository interface {
	Create(ctx context.Context, req *webhookmodel.Request) error
	GetByIDAndInbox(ctx context.Context, id, inboxID string) (*webhookmodel.Request, error)
	ListActive(ctx context.Context, filter WebhookRequestListFilter) ([]webhookmodel.Request, int64, error)
	SoftDelete(ctx context.Context, id, inboxID string) (bool, error)
	SoftDeleteAllActive(ctx context.Context, inboxID string) (int64, error)
	SetRead(ctx context.Context, id, inboxID string, isRead bool) (bool, error)
	CountByInbox(ctx context.Context, inboxID string) (int64, error)
	DeleteOldest(ctx context.Context, inboxID string, limit int) ([]webhookmodel.Request, error)
}

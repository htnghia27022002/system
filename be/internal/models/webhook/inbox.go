package webhook

import "time"

// Inbox is the per-account webhook capture endpoint identity.
type Inbox struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"userId" db:"user_id"`
	PublicUUID       string    `json:"publicUuid" db:"public_uuid"`
	LifetimeReceived int       `json:"lifetimeReceived" db:"lifetime_received"`
	ActiveCount      int       `json:"activeCount" db:"active_count"`
	CreatedAt        time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt" db:"updated_at"`
}

func (Inbox) TableName() string {
	return "webhook_inboxes"
}

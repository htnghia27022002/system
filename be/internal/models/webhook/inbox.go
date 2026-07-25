package webhook

import "time"

// Inbox is the per-account webhook capture endpoint identity.
type Inbox struct {
	ID               string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID           string    `json:"userId" gorm:"column:user_id;type:uuid;not null;uniqueIndex"`
	PublicUUID       string    `json:"publicUuid" gorm:"column:public_uuid;type:uuid;not null;uniqueIndex"`
	LifetimeReceived  int       `json:"lifetimeReceived" gorm:"column:lifetime_received;not null;default:0"`
	ActiveCount      int       `json:"activeCount" gorm:"column:active_count;not null;default:0"`
	CreatedAt        time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}

func (Inbox) TableName() string {
	return "webhook_inboxes"
}

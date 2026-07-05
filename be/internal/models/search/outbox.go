package searchmodel

import "time"

const (
	OutboxStatusPending    = "pending"
	OutboxStatusProcessing = "processing"
	OutboxStatusCompleted  = "completed"
	OutboxStatusFailed     = "failed"

	OutboxOpUpsert = "upsert"
	OutboxOpDelete = "delete"
)

type OutboxEntry struct {
	ID           string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EntityType   string     `json:"entityType" gorm:"column:entity_type;type:varchar(50);not null"`
	EntityID     string     `json:"entityId" gorm:"column:entity_id;type:uuid;not null"`
	Operation    string     `json:"operation" gorm:"type:varchar(20);not null"`
	Status       string     `json:"status" gorm:"type:varchar(20);not null;default:pending"`
	AttemptCount int        `json:"attemptCount" gorm:"column:attempt_count;not null;default:0"`
	LastError    *string    `json:"lastError,omitempty" gorm:"column:last_error;type:text"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
	ProcessedAt  *time.Time `json:"processedAt,omitempty" gorm:"column:processed_at"`
}

func (OutboxEntry) TableName() string {
	return "search_outbox"
}

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
	ID           string     `json:"id" db:"id"`
	EntityType   string     `json:"entityType" db:"entity_type"`
	EntityID     string     `json:"entityId" db:"entity_id"`
	Operation    string     `json:"operation" db:"operation"`
	Status       string     `json:"status" db:"status"`
	AttemptCount int        `json:"attemptCount" db:"attempt_count"`
	LastError    *string    `json:"lastError,omitempty" db:"last_error"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
	ProcessedAt  *time.Time `json:"processedAt,omitempty" db:"processed_at"`
}

func (OutboxEntry) TableName() string {
	return "search_outbox"
}

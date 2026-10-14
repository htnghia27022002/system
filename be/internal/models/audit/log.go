package audit

import "time"

// Log is one row of audit_logs.
type Log struct {
	ID           string         `json:"id" db:"id"`
	ActorID      *string        `json:"actorId,omitempty" db:"actor_id"`
	ActorEmail   string         `json:"actorEmail" db:"actor_email"`
	Action       string         `json:"action" db:"action"`
	ResourceType string         `json:"resourceType" db:"resource_type"`
	ResourceID   string         `json:"resourceId" db:"resource_id"`
	Metadata     map[string]any `json:"metadata" db:"metadata"`
	IPAddress    string         `json:"ipAddress" db:"ip_address"`
	UserAgent    string         `json:"userAgent" db:"user_agent"`
	RequestID    string         `json:"requestId" db:"request_id"`
	CreatedAt    time.Time      `json:"createdAt" db:"created_at"`
}

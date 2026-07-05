package search

import "time"

const (
	EntityUser       = "user"
	EntityRole       = "role"
	EntityPermission = "permission"

	IndexName = "system-search"
)

type Document struct {
	EntityType     string            `json:"entityType"`
	EntityID       string            `json:"entityId"`
	Title          string            `json:"title"`
	SearchableText string            `json:"searchableText"`
	PermissionKeys []string          `json:"permissionKeys"`
	Metadata       map[string]string `json:"metadata,omitempty"`
	UpdatedAt      time.Time         `json:"updatedAt"`
	SemanticSummary string           `json:"semanticSummary,omitempty"`
}

func DocumentID(entityType, entityID string) string {
	return entityType + ":" + entityID
}

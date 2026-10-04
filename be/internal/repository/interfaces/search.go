package interfaces

import (
	"context"
	"time"
)

// Searchable entity types for admin global search.
const (
	SearchEntityUser       = "user"
	SearchEntityRole       = "role"
	SearchEntityPermission = "permission"
)

// SearchHit is one row returned by the Postgres-backed global search.
type SearchHit struct {
	EntityType string
	EntityID   string
	Title      string
	Snippet    string
	Metadata   map[string]string
	UpdatedAt  time.Time
}

type SearchRepository interface {
	// Search matches term (case-insensitive substring) across the given entity types,
	// newest first. entityTypes must be non-empty and already permission-filtered.
	Search(ctx context.Context, term string, entityTypes []string, offset, limit int) ([]SearchHit, int64, error)
}

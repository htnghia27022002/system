package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"be/internal/repository/interfaces"
	"be/pkg/postgres"
)

// SearchRepository runs admin global search directly against Postgres (ILIKE over
// users, roles, and permissions). It replaces the former Elasticsearch index.
type SearchRepository struct {
	db *postgres.Postgres
}

var _ interfaces.SearchRepository = (*SearchRepository)(nil)

func NewSearchRepository(db *postgres.Postgres) *SearchRepository {
	return &SearchRepository{db: db}
}

// Each branch yields: entity_type, entity_id, title, snippet, meta_a, meta_b, updated_at.
// $1 is the ILIKE pattern. Seeded rows may have NULL timestamps (Insert skips zero
// times), so updated_at falls back to created_at, then epoch, like list queries do.
var searchBranches = map[string]string{
	interfaces.SearchEntityUser: `
SELECT 'user', id::text, full_name, email, email, status, COALESCE(updated_at, created_at, 'epoch'::timestamptz)
FROM users
WHERE deleted_at IS NULL AND status = 'active'
  AND (full_name ILIKE $1 OR email ILIKE $1)`,
	interfaces.SearchEntityRole: `
SELECT 'role', id::text, name, COALESCE(description, ''), COALESCE(slug, ''), '', COALESCE(updated_at, created_at, 'epoch'::timestamptz)
FROM roles
WHERE name ILIKE $1 OR COALESCE(slug, '') ILIKE $1 OR COALESCE(description, '') ILIKE $1`,
	interfaces.SearchEntityPermission: `
SELECT 'permission', id::text, name, COALESCE(description, ''), key, "group", COALESCE(updated_at, created_at, 'epoch'::timestamptz)
FROM permissions
WHERE key ILIKE $1 OR name ILIKE $1 OR "group" ILIKE $1 OR COALESCE(description, '') ILIKE $1`,
}

func (r *SearchRepository) Search(
	ctx context.Context,
	term string,
	entityTypes []string,
	offset, limit int,
) ([]interfaces.SearchHit, int64, error) {
	branches := make([]string, 0, len(entityTypes))
	for _, t := range entityTypes {
		if sql, ok := searchBranches[t]; ok {
			branches = append(branches, sql)
		}
	}
	if len(branches) == 0 {
		return []interfaces.SearchHit{}, 0, nil
	}

	union := strings.Join(branches, "\nUNION ALL\n")
	pattern := "%" + escapeLike(term) + "%"
	db := r.db.Querier(ctx)

	var total int64
	if err := db.QueryRow(ctx, "SELECT COUNT(*) FROM ("+union+") AS hits", pattern).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("search count: %w", err)
	}
	if total == 0 {
		return []interfaces.SearchHit{}, 0, nil
	}

	rows, err := db.Query(ctx,
		"SELECT * FROM ("+union+") AS hits ORDER BY 7 DESC, 2 ASC OFFSET $2 LIMIT $3",
		pattern, offset, limit,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("search query: %w", err)
	}
	defer rows.Close()

	hits := make([]interfaces.SearchHit, 0, limit)
	for rows.Next() {
		var (
			hit          interfaces.SearchHit
			metaA, metaB string
			updatedAt    time.Time
		)
		if err := rows.Scan(&hit.EntityType, &hit.EntityID, &hit.Title, &hit.Snippet, &metaA, &metaB, &updatedAt); err != nil {
			return nil, 0, fmt.Errorf("search scan: %w", err)
		}
		hit.UpdatedAt = updatedAt
		hit.Metadata = searchMetadata(hit.EntityType, metaA, metaB)
		hits = append(hits, hit)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("search rows: %w", err)
	}
	return hits, total, nil
}

func searchMetadata(entityType, a, b string) map[string]string {
	switch entityType {
	case interfaces.SearchEntityUser:
		return map[string]string{"email": a, "status": b}
	case interfaces.SearchEntityRole:
		return map[string]string{"slug": a}
	case interfaces.SearchEntityPermission:
		return map[string]string{"key": a, "group": b}
	default:
		return nil
	}
}

// escapeLike escapes ILIKE wildcards so user input matches literally.
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

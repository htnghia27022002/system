package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"

	searchmodel "be/internal/models/search"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type SearchOutboxRepository struct {
	*repo.Repository[searchmodel.OutboxEntry]
}

var _ interfaces.SearchOutboxRepository = (*SearchOutboxRepository)(nil)

func NewSearchOutboxRepository(db *postgres.Postgres) *SearchOutboxRepository {
	return &SearchOutboxRepository{
		Repository: repo.New[searchmodel.OutboxEntry](db, repo.Opts{
			Table: "search_outbox",
			PK:    "id",
		}),
	}
}

func (r *SearchOutboxRepository) Enqueue(ctx context.Context, entityType, entityID, operation string) (string, error) {
	entry := searchmodel.OutboxEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Operation:  operation,
		Status:     searchmodel.OutboxStatusPending,
	}
	if err := r.Insert(ctx, &entry); err != nil {
		return "", err
	}
	return entry.ID, nil
}

func (r *SearchOutboxRepository) GetByID(ctx context.Context, id string) (*searchmodel.OutboxEntry, error) {
	return r.FindByID(ctx, id)
}

func (r *SearchOutboxRepository) ClaimByID(ctx context.Context, id string) (*searchmodel.OutboxEntry, error) {
	var entry *searchmodel.OutboxEntry
	err := r.DB().WithTx(ctx, func(ctx context.Context) error {
		sb := r.Select().
			Where(squirrel.Eq{postgres.QuoteIdent("id"): id}).
			Where(squirrel.Eq{postgres.QuoteIdent("status"): []string{
				searchmodel.OutboxStatusPending,
				searchmodel.OutboxStatusFailed,
			}}).
			Where(squirrel.Expr(postgres.QuoteIdent("attempt_count")+" < ?", 10)).
			Limit(1).
			Suffix("FOR UPDATE")
		found, err := r.QueryOne(ctx, sb)
		if err != nil || found == nil {
			entry = found
			return err
		}
		found.Status = searchmodel.OutboxStatusProcessing
		found.AttemptCount++
		if err := r.Update(ctx, found); err != nil {
			return err
		}
		entry = found
		return nil
	})
	return entry, err
}

func (r *SearchOutboxRepository) ClaimPending(ctx context.Context, limit int) ([]searchmodel.OutboxEntry, error) {
	var entries []searchmodel.OutboxEntry
	err := r.DB().WithTx(ctx, func(ctx context.Context) error {
		sb := r.Select().
			Where(squirrel.Eq{postgres.QuoteIdent("status"): []string{
				searchmodel.OutboxStatusPending,
				searchmodel.OutboxStatusFailed,
			}}).
			Where(squirrel.Expr(postgres.QuoteIdent("attempt_count")+" < ?", 10)).
			OrderBy("created_at ASC").
			Limit(uint64(limit)).
			Suffix("FOR UPDATE")
		found, err := r.Query(ctx, sb)
		if err != nil {
			return err
		}
		for i := range found {
			found[i].Status = searchmodel.OutboxStatusProcessing
			found[i].AttemptCount++
			if err := r.Update(ctx, &found[i]); err != nil {
				return err
			}
		}
		entries = found
		return nil
	})
	return entries, err
}

func (r *SearchOutboxRepository) MarkCompleted(ctx context.Context, id string) error {
	now := time.Now()
	_, err := r.UpdateMap(ctx, id, map[string]any{
		"status":       searchmodel.OutboxStatusCompleted,
		"processed_at": now,
		"last_error":   nil,
	})
	return err
}

func (r *SearchOutboxRepository) MarkFailed(ctx context.Context, id string, errMsg string) error {
	_, err := r.UpdateMap(ctx, id, map[string]any{
		"status":     searchmodel.OutboxStatusFailed,
		"last_error": errMsg,
	})
	return err
}

func (r *SearchOutboxRepository) ResetFailed(ctx context.Context, id string) error {
	_, err := r.UpdateWhere(ctx,
		query.New(1, 1).
			WhereEqual("id", id).
			WhereEqual("status", searchmodel.OutboxStatusFailed),
		map[string]any{
			"status":     searchmodel.OutboxStatusPending,
			"last_error": nil,
		},
	)
	return err
}

func (r *SearchOutboxRepository) Stats(ctx context.Context) (*interfaces.OutboxStats, error) {
	stats := &interfaces.OutboxStats{}
	pending, err := r.Count(ctx, query.New(1, 1).WhereEqual("status", searchmodel.OutboxStatusPending))
	if err != nil {
		return nil, err
	}
	stats.PendingCount = pending

	failed, err := r.Count(ctx, query.New(1, 1).WhereEqual("status", searchmodel.OutboxStatusFailed))
	if err != nil {
		return nil, err
	}
	stats.FailedCount = failed

	oldest, err := r.FindOne(ctx, query.New(1, 1).
		WhereEqual("status", searchmodel.OutboxStatusPending).
		OrderBy("created_at ASC"))
	if err != nil {
		return nil, err
	}
	if oldest != nil {
		stats.OldestPendingAgeSeconds = int64(time.Since(oldest.CreatedAt).Seconds())
	}
	return stats, nil
}

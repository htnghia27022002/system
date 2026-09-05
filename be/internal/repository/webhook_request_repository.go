package repository

import (
	"context"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"

	webhookmodel "be/internal/models/webhook"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type WebhookRequestRepository struct {
	*repo.Repository[webhookmodel.Request]
}

var _ interfaces.WebhookRequestRepository = (*WebhookRequestRepository)(nil)

func NewWebhookRequestRepository(db *postgres.Postgres) *WebhookRequestRepository {
	return &WebhookRequestRepository{
		Repository: repo.New[webhookmodel.Request](db, repo.Opts{
			Table: "webhook_requests",
			PK:    "id",
		}),
	}
}

func (r *WebhookRequestRepository) Create(ctx context.Context, req *webhookmodel.Request) error {
	if req.Headers == nil {
		req.Headers = map[string]any{}
	}
	if req.Query == nil {
		req.Query = map[string]any{}
	}
	if req.Form == nil {
		req.Form = map[string]any{}
	}
	return r.Insert(ctx, req)
}

func (r *WebhookRequestRepository) GetByIDAndInbox(ctx context.Context, id, inboxID string) (*webhookmodel.Request, error) {
	return r.FindOne(ctx, query.New(1, 1).
		WhereEqual("id", id).
		WhereEqual("inbox_id", inboxID))
}

func (r *WebhookRequestRepository) ListActive(ctx context.Context, filter interfaces.WebhookRequestListFilter) ([]webhookmodel.Request, int64, error) {
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	page, limit := query.NormalizePage(filter.Page, limit)

	q := query.New(page, limit).
		WhereEqual("inbox_id", filter.InboxID).
		Where("soft_deleted_at", query.OpIsNull, nil).
		OrderBy("created_at DESC")

	if method := strings.TrimSpace(filter.Method); method != "" {
		q.WhereRaw("UPPER(method) = ?", strings.ToUpper(method))
	}
	switch strings.ToLower(strings.TrimSpace(filter.ReadFilter)) {
	case "read":
		q.WhereRaw("is_read = TRUE")
	case "unread":
		q.WhereRaw("is_read = FALSE")
	}
	if term := strings.TrimSpace(filter.Q); term != "" {
		like := "%" + term + "%"
		q.WhereRaw("(url ILIKE ? OR encode(COALESCE(body, '\\x'::bytea), 'escape') ILIKE ?)", like, like)
	}

	return r.Paginate(ctx, q)
}

func (r *WebhookRequestRepository) SoftDelete(ctx context.Context, id, inboxID string) (bool, error) {
	now := time.Now().UTC()
	n, err := r.UpdateWhere(ctx,
		query.New(1, 1).
			WhereEqual("id", id).
			WhereEqual("inbox_id", inboxID).
			Where("soft_deleted_at", query.OpIsNull, nil),
		map[string]any{"soft_deleted_at": now},
	)
	return n > 0, err
}

func (r *WebhookRequestRepository) SoftDeleteAllActive(ctx context.Context, inboxID string) (int64, error) {
	now := time.Now().UTC()
	return r.UpdateWhere(ctx,
		query.New(1, 1).
			WhereEqual("inbox_id", inboxID).
			Where("soft_deleted_at", query.OpIsNull, nil),
		map[string]any{"soft_deleted_at": now},
	)
}

func (r *WebhookRequestRepository) SetRead(ctx context.Context, id, inboxID string, isRead bool) (bool, error) {
	n, err := r.UpdateWhere(ctx,
		query.New(1, 1).
			WhereEqual("id", id).
			WhereEqual("inbox_id", inboxID).
			Where("soft_deleted_at", query.OpIsNull, nil),
		map[string]any{"is_read": isRead},
	)
	return n > 0, err
}

func (r *WebhookRequestRepository) CountByInbox(ctx context.Context, inboxID string) (int64, error) {
	return r.Count(ctx, query.New(1, 1).WhereEqual("inbox_id", inboxID))
}

func (r *WebhookRequestRepository) DeleteOldest(ctx context.Context, inboxID string, limit int) ([]webhookmodel.Request, error) {
	if limit <= 0 {
		return nil, nil
	}
	oldest, err := r.Query(ctx, r.Select().
		Where(squirrel.Eq{postgres.QuoteIdent("inbox_id"): inboxID}).
		OrderBy("created_at ASC").
		Limit(uint64(limit)))
	if err != nil {
		return nil, err
	}
	if len(oldest) == 0 {
		return nil, nil
	}
	ids := make([]string, 0, len(oldest))
	for _, row := range oldest {
		ids = append(ids, row.ID)
	}
	if _, err := r.ExecBuilder(ctx, r.DB().Builder.
		Delete(postgres.QuoteIdent("webhook_requests")).
		Where(squirrel.Eq{postgres.QuoteIdent("id"): ids})); err != nil {
		return nil, err
	}
	return oldest, nil
}

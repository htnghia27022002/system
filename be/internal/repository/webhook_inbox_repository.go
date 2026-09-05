package repository

import (
	"context"

	webhookmodel "be/internal/models/webhook"
	"be/internal/repository/interfaces"
	"be/pkg/postgres"
	"be/pkg/query"
	"be/pkg/repo"
)

type WebhookInboxRepository struct {
	*repo.Repository[webhookmodel.Inbox]
}

var _ interfaces.WebhookInboxRepository = (*WebhookInboxRepository)(nil)

func NewWebhookInboxRepository(db *postgres.Postgres) *WebhookInboxRepository {
	return &WebhookInboxRepository{
		Repository: repo.New[webhookmodel.Inbox](db, repo.Opts{
			Table: "webhook_inboxes",
			PK:    "id",
		}),
	}
}

func (r *WebhookInboxRepository) GetByUserID(ctx context.Context, userID string) (*webhookmodel.Inbox, error) {
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("user_id", userID))
}

func (r *WebhookInboxRepository) GetByPublicUUID(ctx context.Context, publicUUID string) (*webhookmodel.Inbox, error) {
	return r.FindOne(ctx, query.New(1, 1).WhereEqual("public_uuid", publicUUID))
}

func (r *WebhookInboxRepository) Create(ctx context.Context, inbox *webhookmodel.Inbox) error {
	return r.Insert(ctx, inbox)
}

func (r *WebhookInboxRepository) Update(ctx context.Context, inbox *webhookmodel.Inbox) error {
	return r.Repository.Update(ctx, inbox)
}

func (r *WebhookInboxRepository) UpdateCounters(ctx context.Context, inboxID string, lifetimeReceived, activeCount int) error {
	_, err := r.UpdateMap(ctx, inboxID, map[string]any{
		"lifetime_received": lifetimeReceived,
		"active_count":      activeCount,
	})
	return err
}

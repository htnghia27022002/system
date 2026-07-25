package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	webhookmodel "be/internal/models/webhook"
	"be/internal/repository/interfaces"
)

type WebhookInboxRepository struct {
	db *gorm.DB
}

var _ interfaces.WebhookInboxRepository = (*WebhookInboxRepository)(nil)

func NewWebhookInboxRepository(db *gorm.DB) *WebhookInboxRepository {
	return &WebhookInboxRepository{db: db}
}

func (r *WebhookInboxRepository) GetByUserID(ctx context.Context, userID string) (*webhookmodel.Inbox, error) {
	var inbox webhookmodel.Inbox
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&inbox).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *WebhookInboxRepository) GetByPublicUUID(ctx context.Context, publicUUID string) (*webhookmodel.Inbox, error) {
	var inbox webhookmodel.Inbox
	if err := r.db.WithContext(ctx).Where("public_uuid = ?", publicUUID).First(&inbox).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inbox, nil
}

func (r *WebhookInboxRepository) Create(ctx context.Context, inbox *webhookmodel.Inbox) error {
	return r.db.WithContext(ctx).Create(inbox).Error
}

func (r *WebhookInboxRepository) Update(ctx context.Context, inbox *webhookmodel.Inbox) error {
	return r.db.WithContext(ctx).Save(inbox).Error
}

func (r *WebhookInboxRepository) UpdateCounters(ctx context.Context, inboxID string, lifetimeReceived, activeCount int) error {
	return r.db.WithContext(ctx).Model(&webhookmodel.Inbox{}).
		Where("id = ?", inboxID).
		Updates(map[string]any{
			"lifetime_received": lifetimeReceived,
			"active_count":      activeCount,
		}).Error
}

package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	webhookmodel "be/internal/models/webhook"
	"be/internal/repository/interfaces"
)

type WebhookRequestRepository struct {
	db *gorm.DB
}

var _ interfaces.WebhookRequestRepository = (*WebhookRequestRepository)(nil)

func NewWebhookRequestRepository(db *gorm.DB) *WebhookRequestRepository {
	return &WebhookRequestRepository{db: db}
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
	return r.db.WithContext(ctx).Create(req).Error
}

func (r *WebhookRequestRepository) GetByIDAndInbox(ctx context.Context, id, inboxID string) (*webhookmodel.Request, error) {
	var req webhookmodel.Request
	if err := r.db.WithContext(ctx).
		Where("id = ? AND inbox_id = ?", id, inboxID).
		First(&req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &req, nil
}

func (r *WebhookRequestRepository) ListActive(ctx context.Context, filter interfaces.WebhookRequestListFilter) ([]webhookmodel.Request, int64, error) {
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	q := r.db.WithContext(ctx).Model(&webhookmodel.Request{}).
		Where("inbox_id = ? AND soft_deleted_at IS NULL", filter.InboxID)

	if method := strings.TrimSpace(filter.Method); method != "" {
		q = q.Where("UPPER(method) = ?", strings.ToUpper(method))
	}
	switch strings.ToLower(strings.TrimSpace(filter.ReadFilter)) {
	case "read":
		q = q.Where("is_read = TRUE")
	case "unread":
		q = q.Where("is_read = FALSE")
	}
	if term := strings.TrimSpace(filter.Q); term != "" {
		like := "%" + term + "%"
		// encode(escape) is binary-safe; also match URL/path for free-text search.
		q = q.Where("(url ILIKE ? OR encode(COALESCE(body, '\\x'::bytea), 'escape') ILIKE ?)", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []webhookmodel.Request
	offset := (page - 1) * limit
	if err := q.Order("created_at DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *WebhookRequestRepository) SoftDelete(ctx context.Context, id, inboxID string) (bool, error) {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&webhookmodel.Request{}).
		Where("id = ? AND inbox_id = ? AND soft_deleted_at IS NULL", id, inboxID).
		Update("soft_deleted_at", now)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *WebhookRequestRepository) SoftDeleteAllActive(ctx context.Context, inboxID string) (int64, error) {
	now := time.Now().UTC()
	result := r.db.WithContext(ctx).Model(&webhookmodel.Request{}).
		Where("inbox_id = ? AND soft_deleted_at IS NULL", inboxID).
		Update("soft_deleted_at", now)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *WebhookRequestRepository) SetRead(ctx context.Context, id, inboxID string, isRead bool) (bool, error) {
	result := r.db.WithContext(ctx).Model(&webhookmodel.Request{}).
		Where("id = ? AND inbox_id = ? AND soft_deleted_at IS NULL", id, inboxID).
		Update("is_read", isRead)
	if result.Error != nil {
		return false, result.Error
	}
	return result.RowsAffected > 0, nil
}

func (r *WebhookRequestRepository) CountByInbox(ctx context.Context, inboxID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&webhookmodel.Request{}).
		Where("inbox_id = ?", inboxID).
		Count(&count).Error
	return count, err
}

func (r *WebhookRequestRepository) DeleteOldest(ctx context.Context, inboxID string, limit int) ([]webhookmodel.Request, error) {
	if limit <= 0 {
		return nil, nil
	}
	var oldest []webhookmodel.Request
	if err := r.db.WithContext(ctx).
		Where("inbox_id = ?", inboxID).
		Order("created_at ASC").
		Limit(limit).
		Find(&oldest).Error; err != nil {
		return nil, err
	}
	if len(oldest) == 0 {
		return nil, nil
	}
	ids := make([]string, 0, len(oldest))
	for _, row := range oldest {
		ids = append(ids, row.ID)
	}
	if err := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&webhookmodel.Request{}).Error; err != nil {
		return nil, err
	}
	return oldest, nil
}

package repository

import (
	"context"
	"errors"
	"time"

	searchmodel "be/internal/models/search"
	"be/internal/repository/interfaces"

	"gorm.io/gorm"
)

type SearchOutboxRepository struct {
	db *gorm.DB
}

var _ interfaces.SearchOutboxRepository = (*SearchOutboxRepository)(nil)

func NewSearchOutboxRepository(db *gorm.DB) *SearchOutboxRepository {
	return &SearchOutboxRepository{db: db}
}

func (r *SearchOutboxRepository) Enqueue(ctx context.Context, entityType, entityID, operation string) (string, error) {
	entry := searchmodel.OutboxEntry{
		EntityType: entityType,
		EntityID:   entityID,
		Operation:  operation,
		Status:     searchmodel.OutboxStatusPending,
	}
	if err := r.db.WithContext(ctx).Create(&entry).Error; err != nil {
		return "", err
	}
	return entry.ID, nil
}

func (r *SearchOutboxRepository) GetByID(ctx context.Context, id string) (*searchmodel.OutboxEntry, error) {
	var entry searchmodel.OutboxEntry
	if err := r.db.WithContext(ctx).First(&entry, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

func (r *SearchOutboxRepository) ClaimByID(ctx context.Context, id string) (*searchmodel.OutboxEntry, error) {
	var entry searchmodel.OutboxEntry
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("id = ?", id).
			Where("status IN ?", []string{searchmodel.OutboxStatusPending, searchmodel.OutboxStatusFailed}).
			Where("attempt_count < ?", 10).
			First(&entry).Error; err != nil {
			return err
		}
		entry.Status = searchmodel.OutboxStatusProcessing
		entry.AttemptCount++
		return tx.Save(&entry).Error
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entry, nil
}

func (r *SearchOutboxRepository) ClaimPending(ctx context.Context, limit int) ([]searchmodel.OutboxEntry, error) {
	var entries []searchmodel.OutboxEntry
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Where("status IN ?", []string{searchmodel.OutboxStatusPending, searchmodel.OutboxStatusFailed}).
			Where("attempt_count < ?", 10).
			Order("created_at ASC").
			Limit(limit).
			Find(&entries).Error; err != nil {
			return err
		}
		for i := range entries {
			entries[i].Status = searchmodel.OutboxStatusProcessing
			entries[i].AttemptCount++
			if err := tx.Save(&entries[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return entries, err
}

func (r *SearchOutboxRepository) MarkCompleted(ctx context.Context, id string) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&searchmodel.OutboxEntry{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":       searchmodel.OutboxStatusCompleted,
			"processed_at": now,
			"last_error":   nil,
		}).Error
}

func (r *SearchOutboxRepository) MarkFailed(ctx context.Context, id string, errMsg string) error {
	return r.db.WithContext(ctx).Model(&searchmodel.OutboxEntry{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":     searchmodel.OutboxStatusFailed,
			"last_error": errMsg,
		}).Error
}

func (r *SearchOutboxRepository) ResetFailed(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&searchmodel.OutboxEntry{}).
		Where("id = ? AND status = ?", id, searchmodel.OutboxStatusFailed).
		Updates(map[string]any{
			"status":     searchmodel.OutboxStatusPending,
			"last_error": nil,
		}).Error
}

func (r *SearchOutboxRepository) Stats(ctx context.Context) (*interfaces.OutboxStats, error) {
	stats := &interfaces.OutboxStats{}
	if err := r.db.WithContext(ctx).Model(&searchmodel.OutboxEntry{}).
		Where("status = ?", searchmodel.OutboxStatusPending).
		Count(&stats.PendingCount).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&searchmodel.OutboxEntry{}).
		Where("status = ?", searchmodel.OutboxStatusFailed).
		Count(&stats.FailedCount).Error; err != nil {
		return nil, err
	}

	var oldest searchmodel.OutboxEntry
	err := r.db.WithContext(ctx).
		Where("status = ?", searchmodel.OutboxStatusPending).
		Order("created_at ASC").
		First(&oldest).Error
	if err == nil {
		stats.OldestPendingAgeSeconds = int64(time.Since(oldest.CreatedAt).Seconds())
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return stats, nil
}

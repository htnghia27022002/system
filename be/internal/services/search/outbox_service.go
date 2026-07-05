package searchsvc

import (
	"context"

	searchmodel "be/internal/models/search"
	"be/internal/handlers/publisher"
	"be/internal/repository/interfaces"
)

type OutboxService struct {
	repo      interfaces.SearchOutboxRepository
	publisher *publisher.Publisher
}

func NewOutboxService(repo interfaces.SearchOutboxRepository, pub *publisher.Publisher) *OutboxService {
	if pub == nil {
		pub = publisher.NewNoop()
	}
	return &OutboxService{repo: repo, publisher: pub}
}

func (s *OutboxService) EnqueueUpsert(ctx context.Context, entityType, entityID string) error {
	if s == nil || s.repo == nil {
		return nil
	}
	id, err := s.repo.Enqueue(ctx, entityType, entityID, searchmodel.OutboxOpUpsert)
	if err != nil {
		return err
	}
	return s.publisher.PublishSearchOutbox(ctx, id)
}

func (s *OutboxService) EnqueueDelete(ctx context.Context, entityType, entityID string) error {
	if s == nil || s.repo == nil {
		return nil
	}
	id, err := s.repo.Enqueue(ctx, entityType, entityID, searchmodel.OutboxOpDelete)
	if err != nil {
		return err
	}
	return s.publisher.PublishSearchOutbox(ctx, id)
}

func (s *OutboxService) Replay(ctx context.Context, id string) error {
	if s == nil || s.repo == nil {
		return nil
	}
	if err := s.repo.ResetFailed(ctx, id); err != nil {
		return err
	}
	return s.publisher.PublishSearchOutbox(ctx, id)
}

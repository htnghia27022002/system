package searchsvc

import (
	"context"
	"fmt"

	searchmodel "be/internal/models/search"
	searchpkg "be/internal/search"
	"be/internal/repository/interfaces"
)

type IndexProcessor struct {
	outbox  interfaces.SearchOutboxRepository
	builder *DocumentBuilder
	client  *searchpkg.Client
}

func NewIndexProcessor(
	outbox interfaces.SearchOutboxRepository,
	builder *DocumentBuilder,
	client *searchpkg.Client,
) *IndexProcessor {
	return &IndexProcessor{outbox: outbox, builder: builder, client: client}
}

func (p *IndexProcessor) ProcessByID(ctx context.Context, id string) error {
	if p.client == nil || !p.client.Enabled() {
		return nil
	}

	entry, err := p.outbox.ClaimByID(ctx, id)
	if err != nil {
		return err
	}
	if entry == nil {
		return nil
	}

	if err := p.processEntry(ctx, *entry); err != nil {
		_ = p.outbox.MarkFailed(ctx, entry.ID, err.Error())
		return err
	}
	return p.outbox.MarkCompleted(ctx, entry.ID)
}

func (p *IndexProcessor) ProcessBatch(ctx context.Context, limit int) (int, error) {
	if p.client == nil || !p.client.Enabled() {
		return 0, nil
	}

	entries, err := p.outbox.ClaimPending(ctx, limit)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, entry := range entries {
		if err := p.processEntry(ctx, entry); err != nil {
			_ = p.outbox.MarkFailed(ctx, entry.ID, err.Error())
			continue
		}
		if err := p.outbox.MarkCompleted(ctx, entry.ID); err != nil {
			return processed, err
		}
		processed++
	}
	return processed, nil
}

func (p *IndexProcessor) processEntry(ctx context.Context, entry searchmodel.OutboxEntry) error {
	switch entry.Operation {
	case searchmodel.OutboxOpDelete:
		return p.client.Delete(ctx, entry.EntityType, entry.EntityID)
	case searchmodel.OutboxOpUpsert:
		doc, err := p.builder.Build(ctx, entry.EntityType, entry.EntityID)
		if err != nil {
			return err
		}
		if doc == nil {
			return p.client.Delete(ctx, entry.EntityType, entry.EntityID)
		}
		return p.client.Upsert(ctx, *doc)
	default:
		return fmt.Errorf("unknown outbox operation: %s", entry.Operation)
	}
}

func (p *IndexProcessor) Reindex(ctx context.Context, entityType string) (int, error) {
	if p.client == nil || !p.client.Enabled() {
		return 0, fmt.Errorf("search is disabled")
	}
	if err := p.client.EnsureIndex(ctx); err != nil {
		return 0, err
	}

	switch entityType {
	case "", "all":
		count := 0
		for _, t := range []string{searchpkg.EntityUser, searchpkg.EntityRole, searchpkg.EntityPermission} {
			n, err := p.reindexType(ctx, t)
			if err != nil {
				return count, err
			}
			count += n
		}
		return count, nil
	default:
		return p.reindexType(ctx, entityType)
	}
}

func (p *IndexProcessor) reindexType(ctx context.Context, entityType string) (int, error) {
	count := 0
	switch entityType {
	case searchpkg.EntityUser:
		users, err := p.builder.AllUsers(ctx)
		if err != nil {
			return 0, err
		}
		for _, user := range users {
			doc, err := p.builder.Build(ctx, searchpkg.EntityUser, user.ID)
			if err != nil {
				return count, err
			}
			if doc == nil {
				if err := p.client.Delete(ctx, searchpkg.EntityUser, user.ID); err != nil {
					return count, err
				}
				continue
			}
			if err := p.client.Upsert(ctx, *doc); err != nil {
				return count, err
			}
			count++
		}
	case searchpkg.EntityRole:
		roles, err := p.builder.AllRoles(ctx)
		if err != nil {
			return 0, err
		}
		for _, role := range roles {
			doc, err := p.builder.Build(ctx, searchpkg.EntityRole, role.ID)
			if err != nil {
				return count, err
			}
			if doc == nil {
				continue
			}
			if err := p.client.Upsert(ctx, *doc); err != nil {
				return count, err
			}
			count++
		}
	case searchpkg.EntityPermission:
		permissions, err := p.builder.AllPermissions(ctx)
		if err != nil {
			return 0, err
		}
		for _, item := range permissions {
			doc, err := p.builder.Build(ctx, searchpkg.EntityPermission, item.ID)
			if err != nil {
				return count, err
			}
			if doc == nil {
				continue
			}
			if err := p.client.Upsert(ctx, *doc); err != nil {
				return count, err
			}
			count++
		}
	default:
		return 0, fmt.Errorf("unsupported entity type: %s", entityType)
	}
	return count, nil
}

func (p *IndexProcessor) OutboxStats(ctx context.Context) (*interfaces.OutboxStats, error) {
	return p.outbox.Stats(ctx)
}

func (p *IndexProcessor) ReplayFailed(ctx context.Context, id string) error {
	return p.outbox.ResetFailed(ctx, id)
}
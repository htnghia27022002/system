package dependency

import (
	"be/internal/repository"
	searchsvc "be/internal/services/search"
)

// SearchStack groups search indexing, query, and outbox services.
type SearchStack struct {
	Service   *searchsvc.Service
	Processor *searchsvc.IndexProcessor
	Outbox    *searchsvc.OutboxService
}

func NewSearchStack(infra *Infra) *SearchStack {
	outboxRepo := repository.NewSearchOutboxRepository(infra.DB)
	userRepo := repository.NewUserRepository(infra.DB)
	roleRepo := repository.NewRoleRepository(infra.DB)
	permissionRepo := repository.NewPermissionRepository(infra.DB)

	documentBuilder := searchsvc.NewDocumentBuilder(userRepo, roleRepo, permissionRepo)
	indexProcessor := searchsvc.NewIndexProcessor(outboxRepo, documentBuilder, infra.SearchClient)
	outboxService := searchsvc.NewOutboxService(outboxRepo, infra.Publisher)
	searchService := searchsvc.NewService(infra.SearchClient, userRepo, documentBuilder)

	return &SearchStack{
		Service:   searchService,
		Processor: indexProcessor,
		Outbox:    outboxService,
	}
}

package dependency

import (
	"be/internal/repository"
	searchsvc "be/internal/services/search"
)

func NewSearchService(infra *Infra) *searchsvc.Service {
	return searchsvc.NewService(repository.NewSearchRepository(infra.DB))
}

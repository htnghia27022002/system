package subscribers

import (
	"fmt"

	"be/internal/queue"
	searchsvc "be/internal/services/search"
)

type Registry struct {
	byName map[string]queue.Handler
}

func NewRegistry(processor *searchsvc.IndexProcessor) *Registry {
	items := []queue.Handler{
		NewProcessSearchOutbox(processor),
	}

	byName := make(map[string]queue.Handler, len(items))
	for _, item := range items {
		byName[item.Name()] = item
	}

	return &Registry{byName: byName}
}

func (r *Registry) Get(name string) (queue.Handler, error) {
	if r == nil {
		return nil, fmt.Errorf("subscribers: registry is nil")
	}

	item, ok := r.byName[name]
	if !ok {
		return nil, fmt.Errorf("subscribers: handler %q not registered", name)
	}
	return item, nil
}

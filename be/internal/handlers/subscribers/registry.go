package subscribers

import (
	"fmt"

	"be/internal/queue"
	ingestsvc "be/internal/services/maps/ingest"
)

type Registry struct {
	byName map[string]queue.Handler
}

func NewRegistry(ingest *ingestsvc.Service) *Registry {
	items := []queue.Handler{
		NewProcessMapsIngest(ingest),
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

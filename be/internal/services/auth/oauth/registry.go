package oauth

import (
	"fmt"

	apperrors "be/internal/common/errors"
	"be/internal/config"
)

// Registry maps provider IDs to their adapter implementations.
type Registry struct {
	providers map[string]Provider
}

// NewRegistry builds a provider registry from config and allowed provider IDs.
// Register new adapters here when adding a provider.
func NewRegistry(cfg config.Config) *Registry {
	all := []Provider{
		NewGoogleProvider(cfg),
	}

	byID := make(map[string]Provider, len(all))
	for _, p := range all {
		byID[p.ID()] = p
	}
	return &Registry{providers: byID}
}

func (r *Registry) Get(id string) (Provider, error) {
	p, ok := r.providers[id]
	if !ok {
		return nil, fmt.Errorf("%w: unsupported oauth provider", apperrors.ErrBadRequest)
	}
	return p, nil
}

// All returns every registered provider adapter (configured or not).
func (r *Registry) All() []Provider {
	out := make([]Provider, 0, len(r.providers))
	for _, p := range r.providers {
		out = append(out, p)
	}
	return out
}

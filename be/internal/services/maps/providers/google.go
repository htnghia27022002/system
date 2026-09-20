package providers

import (
	"context"

	"be/internal/config"
	"be/internal/services/maps"
)

func init() {
	maps.Register("google", func(config.Config, maps.HTTPDoer) maps.Provider {
		return NewGoogleProvider()
	})
}

// GoogleProvider is a P1 stub. Replace Search with the Maps Geocoding API later.
type GoogleProvider struct{}

func NewGoogleProvider() *GoogleProvider { return &GoogleProvider{} }

func (p *GoogleProvider) ID() string { return "google" }

func (p *GoogleProvider) Search(_ context.Context, _, _ string) ([]maps.Hit, error) {
	return nil, nil
}

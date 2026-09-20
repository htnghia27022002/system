package maps

import (
	"context"
	"net/http"
)

// Hit is one vendor-agnostic geocode result (lat/lng only).
type Hit struct {
	ExternalID string
	Title      string
	Subtitle   string
	Lat        float64
	Lng        float64
}

// Provider looks up map addresses. One adapter per vendor file in providers/.
type Provider interface {
	ID() string
	Search(ctx context.Context, query, countryCode string) ([]Hit, error)
}

// HTTPDoer fetches Nominatim (tests inject a fake).
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

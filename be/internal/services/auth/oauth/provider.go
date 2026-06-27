package oauth

import (
	"context"
	"net/http"

	"golang.org/x/oauth2"
)

// Profile holds normalized identity data from an OAuth provider.
type Profile struct {
	ID    string
	Email string
	Name  string
}

// Provider is the adapter contract for a single OAuth provider.
// Implement one struct per provider file (e.g. google_provider.go).
type Provider interface {
	// ID returns the provider slug used in routes (e.g. "google").
	ID() string
	// Configured reports whether credentials are present for this provider.
	Configured() bool
	// OAuthConfig builds the oauth2 config for the given redirect URI.
	OAuthConfig(redirectURI string) (*oauth2.Config, error)
	// FetchProfile retrieves user identity using the authorized HTTP client.
	FetchProfile(ctx context.Context, client *http.Client) (*Profile, error)
}

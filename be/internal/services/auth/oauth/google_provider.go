package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	apperrors "be/internal/common/errors"
	"be/internal/config"
)

const ProviderGoogle = "google"

// GoogleProvider implements OAuth for Google accounts.
type GoogleProvider struct {
	cfg config.Config
}

func NewGoogleProvider(cfg config.Config) *GoogleProvider {
	return &GoogleProvider{cfg: cfg}
}

func (p *GoogleProvider) ID() string { return ProviderGoogle }

func (p *GoogleProvider) Configured() bool {
	return p.cfg.OAuthGoogleClientID != "" && p.cfg.OAuthGoogleClientSecret != ""
}

func (p *GoogleProvider) OAuthConfig(redirectURI string) (*oauth2.Config, error) {
	if !p.Configured() {
		return nil, fmt.Errorf("%w: google oauth is not configured", apperrors.ErrBadRequest)
	}
	return &oauth2.Config{
		ClientID:     p.cfg.OAuthGoogleClientID,
		ClientSecret: p.cfg.OAuthGoogleClientSecret,
		RedirectURL:  redirectURI,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}, nil
}

func (p *GoogleProvider) FetchProfile(ctx context.Context, client *http.Client) (*Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("google profile request failed: %s", strings.TrimSpace(string(body)))
	}

	var payload struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	if payload.ID == "" || payload.Email == "" {
		return nil, fmt.Errorf("%w: incomplete oauth profile", apperrors.ErrBadRequest)
	}
	return &Profile{
		ID:    payload.ID,
		Email: payload.Email,
		Name:  payload.Name,
	}, nil
}

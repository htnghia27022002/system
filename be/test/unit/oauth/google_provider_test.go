package oauth_test

import (
	"testing"

	"be/internal/config"
	oauthprovider "be/internal/services/auth/oauth"
)

func TestGoogleProvider_NotConfigured(t *testing.T) {
	p := oauthprovider.NewGoogleProvider(config.Config{})
	if p.Configured() {
		t.Fatal("expected google provider to be unconfigured without credentials")
	}
	_, err := p.OAuthConfig("http://localhost/callback")
	if err == nil {
		t.Fatal("expected error when google oauth is not configured")
	}
}

func TestRegistry_GetGoogle(t *testing.T) {
	reg := oauthprovider.NewRegistry(config.Config{})
	p, err := reg.Get(oauthprovider.ProviderGoogle)
	if err != nil {
		t.Fatalf("get google: %v", err)
	}
	if p.ID() != oauthprovider.ProviderGoogle {
		t.Fatalf("got id %q", p.ID())
	}
}

func TestRegistry_GetUnsupported(t *testing.T) {
	reg := oauthprovider.NewRegistry(config.Config{})
	_, err := reg.Get("github")
	if err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestValidateRedirectURI(t *testing.T) {
	if err := oauthprovider.ValidateRedirectURI("http://localhost/auth/callback"); err != nil {
		t.Fatalf("valid uri: %v", err)
	}
	if err := oauthprovider.ValidateRedirectURI("not-a-url"); err == nil {
		t.Fatal("expected error for invalid uri")
	}
}

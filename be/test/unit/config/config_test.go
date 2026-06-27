package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"be/internal/config"
)

func TestLoadFromYAMLWithEnvSecretOverride(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(configPath, []byte(`
server:
  port: "9090"
database:
  host: yaml-host
jwt:
  accessTtl: 30m
cors:
  origins:
    - http://yaml.local
oauth:
  redirectUrl: http://yaml.local/callback
  allowedProviders:
    - google
`), 0o644); err != nil {
		t.Fatal(err)
	}

	clearDeployEnv(t)
	t.Setenv("CONFIG_FILE", configPath)
	t.Setenv("DB_PASSWORD", "secret-db-pass")
	t.Setenv("JWT_SECRET", "secret-jwt")
	t.Setenv("OAUTH_GOOGLE_CLIENT_ID", "google-id")
	t.Setenv("OAUTH_GOOGLE_CLIENT_SECRET", "google-secret")
	t.Setenv("DB_HOST", "env-host")

	cfg := config.Load()

	if cfg.Port != "9090" {
		t.Fatalf("port from yaml: got %q", cfg.Port)
	}
	if cfg.DBHost != "env-host" {
		t.Fatalf("DB_HOST env should override yaml: got %q", cfg.DBHost)
	}
	if cfg.DBPass != "secret-db-pass" {
		t.Fatalf("DB password from env: got %q", cfg.DBPass)
	}
	if cfg.JWTSecret != "secret-jwt" {
		t.Fatalf("JWT secret from env: got %q", cfg.JWTSecret)
	}
	if cfg.JWTAccessTTL != 30*time.Minute {
		t.Fatalf("access ttl: got %v", cfg.JWTAccessTTL)
	}
	if len(cfg.CORSOrigins) != 1 || cfg.CORSOrigins[0] != "http://yaml.local" {
		t.Fatalf("cors from yaml: %v", cfg.CORSOrigins)
	}
	if cfg.OAuthGoogleClientID != "google-id" {
		t.Fatalf("oauth client id from env: got %q", cfg.OAuthGoogleClientID)
	}
	if cfg.OAuthRedirectURL != "http://yaml.local/callback" {
		t.Fatalf("oauth redirect from yaml: got %q", cfg.OAuthRedirectURL)
	}
	if len(cfg.OAuthAllowedProviders) != 1 || cfg.OAuthAllowedProviders[0] != "google" {
		t.Fatalf("oauth providers from yaml: %v", cfg.OAuthAllowedProviders)
	}
}

func TestLoadMissingYAMLUsesDefaults(t *testing.T) {
	clearDeployEnv(t)
	t.Setenv("CONFIG_FILE", filepath.Join(t.TempDir(), "missing.yaml"))
	t.Setenv("DB_PASSWORD", "")
	t.Setenv("JWT_SECRET", "")

	cfg := config.Load()
	if cfg.Port != "8080" {
		t.Fatalf("default port: got %q", cfg.Port)
	}
	if cfg.DBHost != "localhost" {
		t.Fatalf("default db host: got %q", cfg.DBHost)
	}
}

// clearDeployEnv removes Docker-injected vars so unit tests assert yaml/default behavior.
func clearDeployEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{
		"PORT", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_SSLMODE",
		"JWT_ACCESS_TTL", "JWT_REFRESH_TTL", "CORS_ORIGINS", "REDIS_URL",
		"OAUTH_REDIRECT_URL", "OAUTH_ALLOWED_PROVIDERS",
	} {
		t.Setenv(key, "")
	}
}

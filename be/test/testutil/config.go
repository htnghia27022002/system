package testutil

import (
	"testing"
	"time"

	"be/internal/config"
)

// UnitConfig returns minimal config for unit tests (JWT, OAuth stubs).
func UnitConfig() config.Config {
	return config.Config{
		JWTSecret:    "test-secret",
		JWTAccessTTL: time.Hour,
	}
}

// SkipIfShort skips integration/e2e tests when -short is passed.
func SkipIfShort(t *testing.T) {
	t.Helper()
	if testing.Short() {
		t.Skip("skipped in short mode")
	}
}

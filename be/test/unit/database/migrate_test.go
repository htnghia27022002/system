package database_test

import (
	"testing"

	"be/internal/config"
	"be/internal/database"
)

func TestMigrationURL(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		DBHost:    "postgres",
		DBPort:    "5432",
		DBName:    "myapp_pg",
		DBUser:    "postgres",
		DBPass:    "secret",
		DBSSLMode: "disable",
	}

	got := database.MigrationURL(cfg)
	want := "postgres://postgres:secret@postgres:5432/myapp_pg?sslmode=disable"
	if got != want {
		t.Fatalf("MigrationURL() = %q, want %q", got, want)
	}
}

func TestParseForceVersion(t *testing.T) {
	t.Parallel()

	v, err := database.ParseForceVersion("3")
	if err != nil {
		t.Fatalf("ParseForceVersion: %v", err)
	}
	if v != 3 {
		t.Fatalf("got %d want 3", v)
	}

	if _, err := database.ParseForceVersion("bad"); err == nil {
		t.Fatal("expected error for invalid version")
	}
}

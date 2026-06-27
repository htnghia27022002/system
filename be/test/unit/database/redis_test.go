package database_test

import (
	"testing"

	"be/internal/config"
	"be/internal/database"
)

func TestConnectRedisEmptyURL(t *testing.T) {
	t.Parallel()

	client, err := database.ConnectRedis(config.Config{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if client != nil {
		t.Fatal("expected nil client when redis url is empty")
	}
}

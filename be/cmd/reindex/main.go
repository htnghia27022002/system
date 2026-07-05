package main

// Usage: go run ./cmd/reindex -entity=all
import (
	"context"
	"flag"
	"log"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
)

func main() {
	entityType := flag.String("entity", "all", "entity type to reindex: all, user, role, permission")
	flag.Parse()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	container := app.NewContainer(cfg, db)
	count, err := container.SearchProcessor.Reindex(context.Background(), *entityType)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("reindex complete: entity=%s indexed=%d", *entityType, count)
}

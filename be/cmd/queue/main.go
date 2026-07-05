package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/internal/handlers/subscribers"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	container := app.NewContainer(cfg, db)
	defer container.Close()

	ctx := context.Background()
	if container.SearchClient != nil && container.SearchClient.Enabled() {
		if err := container.SearchService.EnsureIndex(ctx); err != nil {
			log.Printf("search index ensure failed: %v", err)
		}
	}

	if count, err := container.SearchProcessor.ProcessBatch(ctx, 100); err != nil {
		log.Printf("queue worker initial drain failed: %v", err)
	} else if count > 0 {
		log.Printf("queue worker drained %d pending outbox entries", count)
	}

	registry := subscribers.NewRegistry(container.SearchProcessor)

	if container.QueueClient != nil && container.QueueClient.Enabled() {
		if err := container.QueueClient.EnsureInfrastructure(ctx); err != nil {
			log.Fatal(err)
		}
		if err := container.QueueClient.StartConsumers(ctx, registry); err != nil {
			log.Fatal(err)
		}
	} else {
		go runPollingWorker(container)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
}

func runPollingWorker(container *app.Container) {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	log.Printf("queue worker polling outbox (NATS disabled)")
	for range ticker.C {
		count, err := container.SearchProcessor.ProcessBatch(context.Background(), 25)
		if err != nil {
			log.Printf("queue worker poll failed: %v", err)
			continue
		}
		if count > 0 {
			log.Printf("queue worker processed %d outbox entries", count)
		}
	}
}

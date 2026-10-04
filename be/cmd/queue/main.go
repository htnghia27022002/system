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
	defer db.Close()

	container := app.NewContainer(cfg, db)
	defer container.Close()

	ctx := context.Background()
	registry := subscribers.NewRegistry(container.IngestService)

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
	log.Printf("queue worker polling maps ingest (NATS disabled)")
	if container.IngestService == nil {
		return
	}
	for range ticker.C {
		ingested, err := container.IngestService.ProcessQueued(context.Background(), 10)
		if err != nil {
			log.Printf("queue worker maps ingest poll failed: %v", err)
			continue
		}
		if ingested > 0 {
			log.Printf("queue worker processed %d maps ingest runs", ingested)
		}
	}
}

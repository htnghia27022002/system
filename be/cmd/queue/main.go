// Command queue is the background worker. It consumes the JetStream streams that
// enabled modules declare (module.StreamProvider) or, when NATS is disabled, calls
// each module.Poller on a ticker. With no module declaring either, it idles.
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"be/internal/app"
	"be/internal/config"
	"be/internal/database"
	"be/internal/logging"
	"be/internal/module"
)

const pollInterval = 2 * time.Second

func main() {
	cfg := config.Load()
	logging.Setup(cfg.LogFormat, cfg.LogLevel)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	container, err := app.NewContainer(cfg, db)
	if err != nil {
		log.Fatal(err)
	}
	defer container.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	streams := container.Modules.Streams()
	pollers := container.Modules.Pollers()

	switch {
	case container.QueueClient != nil && container.QueueClient.Enabled() && len(streams) > 0:
		if err := container.QueueClient.EnsureInfrastructure(ctx, streams); err != nil {
			log.Fatal(err)
		}
		if err := container.QueueClient.StartConsumers(ctx, streams); err != nil {
			log.Fatal(err)
		}
		log.Printf("queue worker consuming %d stream(s)", len(streams))
	case len(pollers) > 0:
		log.Printf("queue worker polling %d module(s) (NATS disabled)", len(pollers))
		go runPollers(ctx, pollers)
	default:
		log.Printf("queue worker idle: no enabled module declares streams or pollers (modules: %v)",
			container.Modules.Names())
	}

	<-ctx.Done()
}

func runPollers(ctx context.Context, pollers []module.Poller) {
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, p := range pollers {
				processed, err := p.Poll(ctx)
				if err != nil {
					log.Printf("queue worker poll failed: %v", err)
					continue
				}
				if processed > 0 {
					log.Printf("queue worker processed %d job(s)", processed)
				}
			}
		}
	}
}

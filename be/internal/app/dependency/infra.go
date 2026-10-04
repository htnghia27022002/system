package dependency

import (
	"log"

	jwtmanager "be/common/jwt"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/queue"
	"be/pkg/postgres"
)

// Infra holds process-wide infrastructure shared by service resolvers.
type Infra struct {
	Config      config.Config
	DB          *postgres.Postgres
	Queue       queue.Config
	QueueClient *queue.Client
	JWT         *jwtmanager.Manager
	Publisher   *publisher.Publisher
}

func NewInfra(cfg config.Config, db *postgres.Postgres) *Infra {
	queueCfg := queue.LoadConfig()
	queueClient, err := queue.New(queueCfg)
	if err != nil {
		log.Printf("queue client init failed: %v", err)
		queueClient = nil
	}

	return &Infra{
		Config:      cfg,
		DB:          db,
		Queue:       queueCfg,
		QueueClient: queueClient,
		JWT:         jwtmanager.NewManager(cfg),
		Publisher:   publisher.New(queueClient),
	}
}

func (i *Infra) Close() {
	if i.Publisher != nil {
		i.Publisher.Close()
	}
}

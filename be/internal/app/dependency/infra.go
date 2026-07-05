package dependency

import (
	"log"

	jwtmanager "be/internal/common/jwt"
	"be/internal/config"
	"be/internal/handlers/publisher"
	"be/internal/queue"
	searchpkg "be/internal/search"

	"gorm.io/gorm"
)

// Infra holds process-wide infrastructure shared by service resolvers.
type Infra struct {
	Config       config.Config
	DB           *gorm.DB
	Queue        queue.Config
	QueueClient  *queue.Client
	JWT          *jwtmanager.Manager
	Publisher    *publisher.Publisher
	SearchClient *searchpkg.Client
}

func NewInfra(cfg config.Config, db *gorm.DB) *Infra {
	queueCfg := queue.LoadConfig()
	queueClient, err := queue.New(queueCfg)
	if err != nil {
		log.Printf("queue client init failed: %v", err)
		queueClient = nil
	}

	searchClient, err := searchpkg.NewClient(cfg.Elasticsearch.URL, cfg.Elasticsearch.Enabled)
	if err != nil {
		log.Printf("search client init failed: %v", err)
		searchClient, _ = searchpkg.NewClient(cfg.Elasticsearch.URL, false)
	}

	return &Infra{
		Config:       cfg,
		DB:           db,
		Queue:        queueCfg,
		QueueClient:  queueClient,
		JWT:          jwtmanager.NewManager(cfg),
		Publisher:    publisher.New(queueClient),
		SearchClient: searchClient,
	}
}

func (i *Infra) Close() {
	if i.Publisher != nil {
		i.Publisher.Close()
	}
}

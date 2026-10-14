package main

import (
	"log"

	"be/internal/config"
	"be/internal/database"
	"be/internal/logging"
	"be/public"
)

func main() {
	cfg := config.Load()
	logging.Setup(cfg.LogFormat, cfg.LogLevel)

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	redisURL := cfg.RedisURL
	if cfg.Cache.Enabled && cfg.Cache.Driver == "redis" && cfg.Cache.RedisURL != "" {
		redisURL = cfg.Cache.RedisURL
	}

	redisClient, err := database.ConnectRedisURL(redisURL)
	if err != nil {
		log.Fatal(err)
	}
	if redisClient != nil {
		defer func() { _ = redisClient.Close() }()
	}

	if err := public.Run(cfg, db, redisClient); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"

	"be/internal/config"
	"be/internal/database"
	"be/public"
)

func main() {
	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := public.Run(cfg, db); err != nil {
		log.Fatal(err)
	}
}

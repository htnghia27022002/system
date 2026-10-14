package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"be/internal/config"
	"be/internal/database"
	"be/internal/module"
	"be/internal/modules"
)

func main() {
	class := flag.String("class", "DatabaseSeeder", "Seeder class name (e.g. PermissionSeeder, DatabaseSeeder)")
	flag.Parse()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mods, err := modules.Build(&module.Kernel{Config: cfg, DB: db})
	if err != nil {
		log.Fatal(err)
	}

	if err := database.SeedClass(context.Background(), db, *class, mods.SeedContributions()...); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Seeder %q completed.\n", *class)
}

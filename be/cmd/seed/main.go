package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"be/internal/config"
	"be/internal/database"
)

func main() {
	class := flag.String("class", "DatabaseSeeder", "Seeder class name (e.g. PermissionSeeder, DatabaseSeeder)")
	flag.Parse()

	cfg := config.Load()
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := database.SeedClass(context.Background(), db, *class); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Seeder %q completed.\n", *class)
}

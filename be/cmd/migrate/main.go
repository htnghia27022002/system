package main

import (
	"fmt"
	"log"
	"os"

	"be/internal/config"
	"be/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.Load()
	cmd := os.Args[1]

	var err error
	switch cmd {
	case "up":
		err = database.RunMigrations(cfg)
	case "down":
		err = database.MigrateDown(cfg)
	case "drop":
		err = database.MigrateDrop(cfg)
	case "version":
		var version uint
		var dirty bool
		version, dirty, err = database.MigrateVersion(cfg)
		if err == nil {
			fmt.Printf("version: %d dirty: %v\n", version, dirty)
		}
	case "force":
		if len(os.Args) < 3 {
			log.Fatal("usage: go run ./cmd/migrate force <version>")
		}
		var version int
		version, err = database.ParseForceVersion(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		err = database.MigrateForce(cfg, version)
	case "steps":
		if len(os.Args) < 3 {
			log.Fatal("usage: go run ./cmd/migrate steps <n>")
		}
		var steps int
		steps, err = database.ParseForceVersion(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		err = database.MigrateSteps(cfg, steps)
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	fmt.Println(`Usage: go run ./cmd/migrate <command>

Commands:
  up        Apply all pending migrations
  down      Roll back the last migration
  drop      Drop all objects (runs down migrations)
  version   Print current migration version
  force N   Set version without running SQL (recovery)
  steps N   Run N migrations (negative = down)`)
}

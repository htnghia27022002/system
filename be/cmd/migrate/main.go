package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"be/internal/config"
	"be/internal/database"
	"be/internal/module"
	"be/internal/modules"
)

func main() {
	setName := flag.String("module", database.CoreMigrationSet,
		"migration set for down|force|steps: core or an enabled module name")
	flag.Usage = printUsage
	flag.Parse()
	if flag.NArg() < 1 {
		printUsage()
		os.Exit(1)
	}

	cfg := config.Load()
	// Modules only need Config to expose their migrations (constructors do no I/O).
	mods, err := modules.Build(&module.Kernel{Config: cfg})
	if err != nil {
		log.Fatal(err)
	}
	moduleSets := mods.MigrationSets()
	sets := append([]database.MigrationSet{database.CoreMigrations()}, moduleSets...)

	cmd := flag.Arg(0)
	switch cmd {
	case "up":
		err = database.RunMigrations(cfg, moduleSets...)
	case "version":
		for _, set := range sets {
			version, dirty, verr := database.MigrateVersion(cfg, set)
			if verr != nil {
				log.Fatal(verr)
			}
			fmt.Printf("%-10s version: %d dirty: %v\n", set.Name, version, dirty)
		}
	case "down":
		err = database.MigrateDown(cfg, pick(sets, *setName))
	case "drop":
		// Drop removes every table in the database (all sets), then their version tables.
		err = database.MigrateDrop(cfg, database.CoreMigrations())
	case "force":
		if flag.NArg() < 2 {
			log.Fatal("usage: go run ./cmd/migrate [-module name] force <version>")
		}
		var version int
		if version, err = database.ParseForceVersion(flag.Arg(1)); err == nil {
			err = database.MigrateForce(cfg, pick(sets, *setName), version)
		}
	case "steps":
		if flag.NArg() < 2 {
			log.Fatal("usage: go run ./cmd/migrate [-module name] steps <n>")
		}
		var steps int
		if steps, err = database.ParseForceVersion(flag.Arg(1)); err == nil {
			err = database.MigrateSteps(cfg, pick(sets, *setName), steps)
		}
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func pick(sets []database.MigrationSet, name string) database.MigrationSet {
	for _, set := range sets {
		if set.Name == name {
			return set
		}
	}
	log.Fatalf("unknown migration set %q (core or an enabled module)", name)
	return database.MigrationSet{}
}

func printUsage() {
	fmt.Println(`Usage: go run ./cmd/migrate [-module name] <command>

Each module has its own version table (schema_migrations_<module>); core uses schema_migrations.

Commands:
  up        Apply pending migrations: core, then every enabled module
  version   Print the version of core and every enabled module
  down      Roll back the last migration of -module (default core)
  drop      Drop all tables in the database
  force N   Set -module's version without running SQL (recovery)
  steps N   Run N migrations of -module (negative = down)`)
}

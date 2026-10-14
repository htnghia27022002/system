// Command openapi writes the OpenAPI document (core + registry default modules)
// to stdout or -o FILE. It needs no database: modules are built from Config only.
//
//	go run ./cmd/openapi -o api/openapi.json
package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"be/internal/config"
	"be/internal/module"
	"be/internal/modules"
	"be/public/routes"
)

func main() {
	out := flag.String("o", "", "output file (default stdout)")
	flag.Parse()

	cfg := config.Load()
	// Document the registry defaults, independent of the local MODULES value.
	cfg.ModulesSet = false
	mods, err := modules.Build(&module.Kernel{Config: cfg})
	if err != nil {
		log.Fatal(err)
	}
	data, err := json.MarshalIndent(routes.OpenAPI(mods), "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	data = append(data, '\n')
	if *out == "" {
		_, _ = os.Stdout.Write(data)
		return
	}
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		log.Fatal(err)
	}
}

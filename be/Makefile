.PHONY: run build test test-unit test-integration test-e2e test-all lint tidy migrate-up migrate-down migrate-version migrate-create seed seed-class

run:
	go run .

# Build main API binary (be/main.go) → app (migrations embedded in binary)
build:
	CGO_ENABLED=0 go build -o app .

test: test-unit

test-unit:
	go test ./test/unit/...

test-integration:
	go test -tags=integration ./test/integration/...

test-e2e:
	go test -tags=e2e ./test/e2e/...

test-all: test-unit test-integration test-e2e

lint:
	golangci-lint run ./...

tidy:
	go mod tidy

# Migration sets: core (migrations/) and one per module (internal/modules/<name>/migrations).
# Pass module=<name> to target a module set; default core.
module ?= core
migrations_dir = $(if $(filter core,$(module)),migrations,internal/modules/$(module)/migrations)

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate -module $(module) down

migrate-version:
	go run ./cmd/migrate version

# Usage: make migrate-create name=add_users_index [module=webhooks]
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=your_migration_name [module=name]"; exit 1; fi
	@test -d $(migrations_dir) || { echo "No migrations folder $(migrations_dir)"; exit 1; }
	@next=$$(ls $(migrations_dir)/*.up.sql 2>/dev/null | sed 's/.*\///;s/_.*//' | sort -n | tail -1); \
	if [ -z "$$next" ]; then next=0; fi; \
	next=$$(expr $$next + 1); \
	padded=$$(printf "%06d" $$next); \
	touch $(migrations_dir)/$${padded}_$(name).up.sql $(migrations_dir)/$${padded}_$(name).down.sql; \
	echo "Created $(migrations_dir)/$${padded}_$(name).up.sql and .down.sql"

seed:
	go run ./cmd/seed

# Usage: make seed-class class=PermissionSeeder
seed-class:
	@if [ -z "$(class)" ]; then echo "Usage: make seed-class class=PermissionSeeder"; exit 1; fi
	go run ./cmd/seed --class=$(class)

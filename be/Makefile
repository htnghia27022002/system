.PHONY: run test test-unit test-integration test-e2e test-all lint tidy migrate-up migrate-down migrate-version migrate-create seed seed-class

run:
	go run .

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

migrate-up:
	go run ./cmd/migrate up

migrate-down:
	go run ./cmd/migrate down

migrate-version:
	go run ./cmd/migrate version

# Usage: make migrate-create name=add_users_index
migrate-create:
	@if [ -z "$(name)" ]; then echo "Usage: make migrate-create name=your_migration_name"; exit 1; fi
	@next=$$(ls migrations/*.up.sql 2>/dev/null | sed 's/.*\///;s/_.*//' | sort -n | tail -1); \
	if [ -z "$$next" ]; then next=0; fi; \
	next=$$((10#$$next + 1)); \
	padded=$$(printf "%06d" $$next); \
	touch migrations/$${padded}_$(name).up.sql migrations/$${padded}_$(name).down.sql; \
	echo "Created migrations/$${padded}_$(name).up.sql and .down.sql"

seed:
	go run ./cmd/seed

# Usage: make seed-class class=PermissionSeeder
seed-class:
	@if [ -z "$(class)" ]; then echo "Usage: make seed-class class=PermissionSeeder"; exit 1; fi
	go run ./cmd/seed --class=$(class)

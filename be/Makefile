.PHONY: run test test-unit test-integration test-e2e test-all lint tidy migrate-up

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
	@echo "Apply SQL migrations manually with psql, e.g.:"
	@echo "  psql $$DATABASE_URL -f migrations/001_init_schema.sql"
	@echo "  psql $$DATABASE_URL -f migrations/002_align_rbac_fields.sql"
	@echo "  psql $$DATABASE_URL -f migrations/003_seed_rbac.sql"

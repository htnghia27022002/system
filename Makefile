.PHONY: up up-d down logs prod psql redis env test test-be test-fe lint

up:
	docker compose up --build

up-d:
	docker compose up --build -d

down:
	docker compose down

logs:
	docker compose logs -f

prod:
	docker compose --profile prod up --build

psql:
	docker compose exec postgres psql -U $${POSTGRES_USER:-postgres} -d $${POSTGRES_DB:-myapp_pg}

redis:
	docker compose exec redis redis-cli ping

env:
	@test -f .env || cp .env.docker.example .env
	@echo "Created .env from .env.docker.example"

test: test-be test-fe

test-be:
	docker compose exec -T be go test ./test/unit/...

test-be-integration:
	docker compose exec -T be go test -tags=integration ./test/integration/...

test-be-e2e:
	docker compose exec -T be go test -tags=e2e ./test/e2e/...

test-be-all:
	docker compose exec -T be sh -c "go test ./test/unit/... && go test -tags=integration ./test/integration/... && go test -tags=e2e ./test/e2e/..."

test-fe:
	pnpm test

lint:
	pnpm lint

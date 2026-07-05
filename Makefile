.PHONY: up up-d down logs prod psql redis env test test-be test-fe lint

# be/.env for DB_* substitution (postgres) + be/queue containers. fe/.env for prod FE build args.
COMPOSE := docker compose --env-file be/.env
COMPOSE_PROD := docker compose --env-file be/.env --env-file fe/.env

up:
	$(COMPOSE) up --build

up-d:
	$(COMPOSE) up --build -d

down:
	docker compose down

logs:
	$(COMPOSE) logs -f

prod:
	$(COMPOSE_PROD) --profile prod up --build

psql:
	$(COMPOSE) exec postgres psql -U $${DB_USER:-postgres} -d $${DB_NAME:-myapp_pg}

redis:
	$(COMPOSE) exec redis redis-cli ping

env:
	@test -f .env || cp .env.example .env
	@test -f be/.env || cp be/.env.example be/.env
	@test -f fe/.env || cp fe/.env.example fe/.env
	@echo "Created missing env files (.env, be/.env, fe/.env)"

test: test-be test-fe

test-be:
	$(COMPOSE) exec -T be go test ./test/unit/...

test-be-integration:
	$(COMPOSE) exec -T be go test -tags=integration ./test/integration/...

test-be-e2e:
	$(COMPOSE) exec -T be go test -tags=e2e ./test/e2e/...

test-be-all:
	$(COMPOSE) exec -T be sh -c "go test ./test/unit/... && go test -tags=integration ./test/integration/... && go test -tags=e2e ./test/e2e/..."

test-fe:
	pnpm test

lint:
	pnpm lint

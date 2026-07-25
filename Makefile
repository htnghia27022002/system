.PHONY: up up-d down logs build prod psql redis env test test-be test-fe lint

# be/.env for DB_* substitution (postgres) + be/queue containers.
# Profiles: dev = full stack; prod = nginx + go-app + go-queue + nats + redis.
# Prod pulls GO_APP_IMAGE from Docker Hub (set in root .env) — no source build on server.
COMPOSE := docker compose --env-file be/.env --profile dev
COMPOSE_PROD := docker compose --env-file be/.env --profile prod
GO_APP_IMAGE ?= $(shell sed -n 's/^GO_APP_IMAGE=//p' .env 2>/dev/null)

up:
	$(COMPOSE) up --build

up-d:
	$(COMPOSE) up --build -d

down:
	docker compose --profile dev --profile prod down

logs:
	$(COMPOSE) logs -f

# Build Go prod image + push to Docker Hub (multi-stage; final layer = binaries only).
build:
	@test -n "$(GO_APP_IMAGE)" || { echo "Set GO_APP_IMAGE in .env first (e.g. youruser/go-app:1.0.0)"; exit 1; }
	docker build -f docker/be/Dockerfile.prod -t $(GO_APP_IMAGE) .
	docker push $(GO_APP_IMAGE)

# Prod server: pull Hub image + start (no --build, no Go source tree required).
prod:
	@test -n "$(GO_APP_IMAGE)" || { echo "Set GO_APP_IMAGE in .env first (e.g. youruser/go-app:1.0.0)"; exit 1; }
	$(COMPOSE_PROD) pull go-app
	$(COMPOSE_PROD) up -d

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

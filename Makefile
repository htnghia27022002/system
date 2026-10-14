.PHONY: bootstrap module openapi up up-d down logs build prod psql redis env test test-be test-fe lint

# be/.env for DB_* substitution (postgres) + be/queue containers.
# Profiles: dev = postgres + be + fe + nginx; prod = nginx + go-app (pulls GO_APP_IMAGE).
# Optional add-ons: cache (redis), queue (nats + worker), prod-queue (nats + go-queue), mail (mailpit).
#   make up-d PROFILES="dev queue"      make prod PROD_PROFILES="prod prod-queue"
PROFILES ?= dev
PROD_PROFILES ?= prod
ALL_PROFILES := --profile dev --profile prod --profile cache --profile queue --profile prod-queue --profile mail
COMPOSE := docker compose --env-file be/.env $(addprefix --profile ,$(PROFILES))
COMPOSE_PROD := docker compose --env-file be/.env $(addprefix --profile ,$(PROD_PROFILES))
GO_APP_IMAGE ?= $(shell sed -n 's/^GO_APP_IMAGE=//p' .env 2>/dev/null)

up:
	$(COMPOSE) up --build

up-d:
	$(COMPOSE) up --build -d

down:
	docker compose --env-file be/.env $(ALL_PROFILES) down

logs:
	$(COMPOSE) logs -f

# Build Go prod image + push to Docker Hub (multi-stage; final layer = binaries only).
build:
	@test -n "$(GO_APP_IMAGE)" || { echo "Set GO_APP_IMAGE in .env first (e.g. youruser/go-app:1.0.0)"; exit 1; }
	docker build -f be/Dockerfile -t $(GO_APP_IMAGE) be
	docker push $(GO_APP_IMAGE)

# Prod server: pull Hub image + start (no --build, no Go source tree required).
prod:
	@test -n "$(GO_APP_IMAGE)" || { echo "Set GO_APP_IMAGE in .env first (e.g. youruser/go-app:1.0.0)"; exit 1; }
	$(COMPOSE_PROD) pull go-app
	$(COMPOSE_PROD) up -d

psql:
	$(COMPOSE) exec postgres psql -U $${DB_USER:-postgres} -d $${DB_NAME:-myapp_pg}

redis:
	docker compose --env-file be/.env --profile cache exec redis redis-cli ping

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
	pnpm -C fe test:run

lint:
	pnpm -C fe lint

# New project from this base (run once after cloning):
#   make bootstrap name="Acme Portal" slug=acme domain=acme.local [image=acme/acme-api]
bootstrap:
	@test -n "$(name)" -a -n "$(slug)" || { echo 'usage: make bootstrap name="Acme Portal" slug=acme [domain=acme.local] [image=org/api]'; exit 1; }
	node scripts/bootstrap.mjs --name "$(name)" --slug "$(slug)" $(if $(domain),--domain "$(domain)") $(if $(image),--image "$(image)")

# Scaffold a CRUD module on both sides: make module name=invoices [label="Invoices"] [singular=invoice]
module:
	@test -n "$(name)" || { echo 'usage: make module name=invoices [label="Invoices"] [singular=invoice]'; exit 1; }
	node scripts/new-module.mjs "$(name)" $(if $(label),--label "$(label)") $(if $(singular),--singular "$(singular)")

# API contract: BE route docs -> be/api/openapi.json -> fe/src/types/api.gen.ts
# (BE step runs in the be container; use `cd be && make openapi` with Go on the host).
openapi:
	$(COMPOSE) exec -T be go run ./cmd/openapi -o api/openapi.json
	pnpm -C fe gen:api

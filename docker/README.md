# Docker Environment

Run the stack with Docker Compose profiles: **dev** (full stack + hot reload) or **prod** (API edge only).

## Quick start

```bash
cp .env.example .env

# Add local hostname (once)
echo "127.0.0.1 system.local" | sudo tee -a /etc/hosts

make up
```

Open **http://system.local** or **http://localhost** (nginx listens on host port 80 by default).

- Frontend: `http://system.local`
- API: `http://system.local/api`
- Seed admin: `admin@example.com` / `admin1234`

If port 80 is already in use, set `NGINX_HTTP_PORT=8080` in `.env` and use `http://system.local:8080`.

## Nginx config layout

```text
docker/nginx/
├── nginx.dev.conf          # dev entry (includes config/)
├── nginx.prod.conf         # prod entry (API-only)
└── config/
    ├── upstreams.dev.conf  # fe + be upstreams (dev)
    ├── upstreams.prod.conf # go-app upstream
    ├── be.conf             # /api/ → backend
    ├── webhooks-capture.conf
    ├── fe.dev.conf         # / + Next.js HMR (dev)
    └── fe.prod.conf        # retained for optional FE-behind-nginx setups
```

`nginx.dev.conf` / `nginx.prod.conf` mount as `/etc/nginx/conf.d/default.conf` and `include` files from `/etc/nginx/config/`.

## Services

### Dev profile (`make up` / `make up-d`)

| Service | Container | Internal port | Notes |
|---------|-----------|---------------|-------|
| nginx | `nginx` | 80 → host `NGINX_HTTP_PORT` | Reverse proxy, `server_name system.local` |
| frontend | `fe` | 3000 | Next.js dev server |
| backend | `be` | 8080 | Go + Gin API (publishes to NATS, no consumers) |
| queue | `queue` | — | `go run ./cmd/queue` — JetStream consumers |
| postgres | `postgres` | 5432 | Data: `docker/data/postgres/` |
| redis | `redis` | 6379 | Data: `docker/data/redis/` |
| elasticsearch | `elasticsearch` | 9200 | Search index; data: `docker/data/elasticsearch/` |
| nats | `nats` | 4222 (client), 8222 (monitor) | JetStream message broker |

### Prod profile (`make prod`)

API edge only — FE, Postgres, and Elasticsearch are **not** in this stack (configure hosts in `be/.env`).

Go services **pull a pre-built image from Docker Hub** (`GO_APP_IMAGE` in root `.env`). The prod host does not copy Go sources or run `docker build`.

| Service | Container | Notes |
|---------|-----------|-------|
| nginx | `nginx-prod` | Proxies `/api/` + webhook capture rewrite |
| backend | `go-app` | Pulls `GO_APP_IMAGE`, runs `./server` |
| queue | `go-queue` | Same image, runs `./queue` |
| redis | `redis` | Shared with prod profile |
| nats | `nats` | Shared with prod profile |

Point `DB_HOST`, `ELASTICSEARCH_URL`, etc. in `be/.env` at external services before `make prod`.

### Publish Go image (CI / laptop)

```bash
# In root .env:
# GO_APP_IMAGE=youruser/go-app:1.0.0

make build   # docker build + docker push
```

Final image contains `./server`, `./queue`, and `config.yaml` only — no Go source tree.

## Commands

| Command | Description |
|---------|-------------|
| `make up` | Start **dev** stack (foreground, rebuild) |
| `make up-d` | Start **dev** stack in background |
| `make down` | Stop and remove containers (dev + prod profiles) |
| `make logs` | Follow dev service logs |
| `make build` | Build + push `GO_APP_IMAGE` to Docker Hub |
| `make prod` | Pull `GO_APP_IMAGE` + start prod profile (no source build) |
| `make psql` | Open psql shell (dev) |
| `make redis` | Ping Redis |

## Profiles

**Dev** — hot reload with volume mounts:

```bash
docker compose --env-file be/.env --profile dev up --build
# or: make up-d
```

**Prod** — pull Hub image + redis + nats + nginx (no FE, no Go build on server):

```bash
# .env must set GO_APP_IMAGE=youruser/go-app:1.0.0
make prod
```

## Routing (nginx)

| Path | Config file | Upstream |
|------|-------------|----------|
| `/api/*` | `config/be.conf` | Go backend |
| `/tools/webhooks/{uuid}` | `config/webhooks-capture.conf` | Go capture API |
| `/_next/webpack-hmr` | `config/fe.dev.conf` | Next.js (dev WebSocket) |
| `/*` | `config/fe.dev.conf` (dev) / inline (prod) | FE (dev) or API status JSON (prod) |

## Environment

Three files — **no app config in root `.env`**.

| File | Purpose |
|------|---------|
| `.env` | `NGINX_HTTP_PORT`, `GO_APP_IMAGE` (compose auto-loads) |
| `be/.env` | Full BE + queue: `DB_*`, `REDIS_URL`, `ELASTICSEARCH_*`, `NATS_*`, JWT, OAuth, CORS |
| `fe/.env` | Full FE: `NEXT_PUBLIC_*` (dev profile) |

Setup: `make env` or copy `.env.example`, `be/.env.example`, `fe/.env.example`.

Compose uses `env_file: ./be/.env` on `be` / `queue` / `go-app` / `go-queue`, and `env_file: ./fe/.env` on `fe`. Postgres init reads `${DB_*}` from `be/.env` via `make` → `docker compose --env-file be/.env`.

## Persistent data

Database and cache data are stored under [`data/`](data/) as bind mounts:

```text
docker/data/
├── postgres/       # PostgreSQL files
├── redis/          # Redis AOF/RDB (if enabled)
├── elasticsearch/  # Elasticsearch index data
└── nats/           # NATS JetStream data
```

`make down` stops containers only — data in `docker/data/` remains on disk.

To reset from scratch:

```bash
make down
rm -rf docker/data/postgres docker/data/redis
make up-d
```

If you previously used the old named volume `postgres_data`, copy data manually or re-seed after reset.

## WSL2 notes

- `WATCHPACK_POLLING=true` is set for Next.js file watching in Docker on WSL.
- If `make up` fails on port 80, set `NGINX_HTTP_PORT=8080` in `.env`.

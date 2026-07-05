# Docker Environment

Run the full stack (nginx, Go BE, Next.js FE, PostgreSQL, Redis) with one command.

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
├── nginx.prod.conf         # prod entry
└── config/
    ├── upstreams.dev.conf  # fe + be upstreams (dev)
    ├── upstreams.prod.conf # fe-prod + be-prod upstreams
    ├── be.conf             # /api/ → backend
    ├── fe.dev.conf         # / + Next.js HMR (dev)
    └── fe.prod.conf        # / + static cache (prod)
```

`nginx.dev.conf` / `nginx.prod.conf` mount as `/etc/nginx/conf.d/default.conf` and `include` files from `/etc/nginx/config/`.

## Services

| Service | Dev container | Internal port | Notes |
|---------|---------------|---------------|-------|
| nginx | `nginx` | 80 → host `NGINX_HTTP_PORT` | Reverse proxy, `server_name system.local` |
| frontend | `fe` | 3000 | Next.js dev server |
| backend | `be` | 8080 | Go + Gin API (publishes to NATS, no consumers) |
| queue | `queue` | — | `go run ./cmd/queue` — JetStream consumers |
| postgres | `postgres` | 5432 | Data: `docker/data/postgres/` |
| redis | `redis` | 6379 | Data: `docker/data/redis/` |
| elasticsearch | `elasticsearch` | 9200 | Search index; data: `docker/data/elasticsearch/` |
| nats | `nats` | 4222 (client), 8222 (monitor) | JetStream message broker |

## Commands

| Command | Description |
|---------|-------------|
| `make up` | Start dev stack (foreground, rebuild) |
| `make up-d` | Start dev stack in background |
| `make down` | Stop and remove containers |
| `make logs` | Follow all service logs |
| `make prod` | Start production profile |
| `make psql` | Open psql shell |
| `make redis` | Ping Redis |

## Profiles

**Dev (default)** — hot reload with volume mounts:

```bash
docker compose up --build
```

**Production-like** — compiled Go binary + Next.js build:

```bash
docker compose --profile prod up --build
```

## Routing (nginx)

| Path | Config file | Upstream |
|------|-------------|----------|
| `/api/*` | `config/be.conf` | Go backend |
| `/_next/webpack-hmr` | `config/fe.dev.conf` | Next.js (dev WebSocket) |
| `/*` | `config/fe.dev.conf` | Next.js frontend |

## Environment

Three files — **no app config in root `.env`**.

| File | Purpose |
|------|---------|
| `.env` | `NGINX_HTTP_PORT` only (compose auto-loads) |
| `be/.env` | Full BE + queue: `DB_*`, `REDIS_URL`, `ELASTICSEARCH_*`, `NATS_*`, JWT, OAuth, CORS |
| `fe/.env` | Full FE: `NEXT_PUBLIC_*` |

Setup: `make env` or copy `.env.example`, `be/.env.example`, `fe/.env.example`.

Compose uses `env_file: ./be/.env` on `be` / `queue` / `be-prod`, and `env_file: ./fe/.env` on `fe` / `fe-prod`. Postgres init reads `${DB_*}` from `be/.env` via `make` → `docker compose --env-file be/.env`.

## Persistent data

Database and cache data are stored under [`data/`](data/) as bind mounts:

```text
docker/data/
├── postgres/       # PostgreSQL files
├── redis/          # Redis AOF/RDB (if enabled)
└── elasticsearch/  # Elasticsearch index data
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

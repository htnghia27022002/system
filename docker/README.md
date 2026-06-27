# Docker Environment

Run the full stack (nginx, Go BE, Next.js FE, PostgreSQL, Redis) with one command.

## Quick start

```bash
cp .env.docker.example .env

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
| backend | `be` | 8080 | Go + Gin API |
| postgres | `postgres` | 5432 | Data: `docker/data/postgres/` |
| redis | `redis` | 6379 | Data: `docker/data/redis/` |

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

Copy [`.env.docker.example`](../.env.docker.example) to `.env` at the repo root.

| Variable | Purpose |
|----------|---------|
| `APP_HOST` | Local hostname (`system.local`) |
| `NGINX_HTTP_PORT` | Host port for nginx (default `80`) |
| `NEXT_PUBLIC_API_BASE_URL` | FE API base (`http://system.local/api`) |
| `NEXT_PUBLIC_USE_MOCK_API` | `false` = real backend |
| `CORS_ORIGINS` | Browser origins allowed by BE (include host **with port**, e.g. `http://system.local:8080`) |

For non-Docker local dev, use `be/.env` and `fe/.env` separately.

## Persistent data

Database and cache data are stored under [`data/`](data/) as bind mounts:

```text
docker/data/
├── postgres/   # PostgreSQL files
└── redis/      # Redis AOF/RDB (if enabled)
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

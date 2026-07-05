# BE implement: Authentication

**Feature:** `docs/features/001-auth/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md)

## Analysis summary

Core auth (login, register, refresh, logout, me, OAuth service) is implemented. Gaps: FE needs to know if Google is configured; logout must receive refresh token from FE.

## API (existing + new)

| Method | Path | Notes |
|--------|------|-------|
| POST | `/api/auth/login` | Done |
| POST | `/api/auth/register` | Done |
| POST | `/api/auth/refresh` | Done — rotates refresh token |
| POST | `/api/auth/logout` | Body: `{ refreshToken }` |
| GET | `/api/auth/me` | Bearer JWT |
| GET | `/api/auth/oauth/:provider/start?redirect_uri=` | Redirect to Google |
| POST | `/api/auth/oauth/:provider/callback` | Body: `{ code, redirectUri }` |
| GET | `/api/auth/oauth/providers` | **New** — `{ providers: ["google"] }` |

## Files changed

| Path | Change |
|------|--------|
| `be/internal/services/auth/oauth_service.go` | `ListConfiguredProviders()` |
| `be/public/handlers/auth_handler.go` | `OAuthProviders` handler |
| `be/public/routes/auth.go` | Register providers route |

## OAuth env

| Run mode | File | Variables |
|----------|------|-----------|
| **Docker / monorepo compose** | `be/.env` | `OAUTH_*`, `JWT_*`, `CORS_*`, `DB_*`, `REDIS_URL`, … |
| **Local BE** | `be/.env` | Same file; use `DB_HOST=localhost`, `OAUTH_REDIRECT_URL=http://localhost:3000/auth/callback` |

After changing **`be/.env`**, recreate BE: `docker compose up -d be` (or `make up-d`).

Google Console **Authorized redirect URI** must match the active `OAUTH_REDIRECT_URL` (e.g. `http://system.local:8080/auth/callback` for Docker on port 8080).

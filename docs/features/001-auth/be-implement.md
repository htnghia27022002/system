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

## OAuth env (Docker)

| Variable | Purpose |
|----------|---------|
| `OAUTH_GOOGLE_CLIENT_ID` | Google client ID |
| `OAUTH_GOOGLE_CLIENT_SECRET` | Google secret |
| `OAUTH_ALLOWED_PROVIDERS` | `google` |
| `OAUTH_REDIRECT_URL` | Must match FE `/auth/callback` URL registered in Google Console |

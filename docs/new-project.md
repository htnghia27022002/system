# New project checklist

Starting a product from this base. Work top to bottom; each step is small.

## 1. Clone and name it

```bash
git clone <this repo> acme && cd acme
make bootstrap name="Acme Portal" slug=acme domain=acme.local   # optional: image=acme/acme-api
echo "127.0.0.1 acme.local" | sudo tee -a /etc/hosts
```

`scripts/bootstrap.mjs` renames the display name, database (`acme_db` / `acme_test`),
dev hostname, logo wordmark, OpenAPI title and prod image, and creates `.env`,
`be/.env`, `fe/.env` with a random `JWT_SECRET`. Re-running is harmless. It does not
touch `docs/features/` (history of the base).

Then replace the artwork: `fe/src/assets/favicon/*`, `fe/public/icons/*`
(`pnpm -C fe generate:pwa-icons`), and the landing copy in `fe/src/features/landing/content.ts`.

## 2. Decide the optional pieces

| Piece | Keep? | How to drop / enable |
|-------|-------|----------------------|
| Webhooks tool (reference module) | Usually no | [docs/modules.md → Remove a module](modules.md#remove-a-module-eg-in-a-cloned-project), or `MODULES=` in `be/.env` |
| Redis cache | Only if you need it | compose profile `cache`, `REDIS_URL` |
| NATS queue + worker | Only for async work | compose profile `queue`, `NATS_ENABLED=true` |
| Self sign-up | Product decision | Admin → Settings → "Allow self sign-up" (invite-only when off) |
| Email verification gate | Product decision | `AUTH_REQUIRE_EMAIL_VERIFICATION=true` |
| Google OAuth | If needed | `OAUTH_GOOGLE_CLIENT_ID/SECRET`, `OAUTH_REDIRECT_URL` |

## 3. Wire the environment (`be/.env`, `fe/.env`)

- **Mail:** dev uses `MAIL_DRIVER=log` (links appear in the BE log) or Mailpit
  (`make up-d PROFILES="dev mail"`). Production: `MAIL_DRIVER=smtp`, `SMTP_*`,
  `MAIL_FROM` on a domain with SPF/DKIM.
- **Links in emails:** `APP_URL` = the public FE origin.
- **Uploads:** `STORAGE_DRIVER=local` for one server; `s3` (+ `S3_*`) for several
  replicas or containers without persistent disks.
- **Secrets:** `JWT_SECRET`, `DB_PASSWORD`, OAuth/SMTP/S3 credentials come from the
  environment only — never commit `.env`.
- **Proxy & cookies:** `TRUSTED_PROXIES` = your load balancer; FE on another site →
  `AUTH_COOKIE_SAMESITE=none` (HTTPS) and the FE origin in `CORS_ORIGINS`.

## 4. Build features as modules

```bash
make module name=invoices label="Invoices"     # BE + FE CRUD scaffold, registered
make openapi                                   # spec + FE types after changing DTOs
```

Business features never go into the core. Specs, plans and contracts live in
`docs/features/NNN-name/` (see `docs/workflow/overview.md`). Every new route needs
an `apidoc.Op` (CI fails otherwise); every admin page needs `{resource}:view` /
`{resource}:modify` keys and a `PermissionGuard`.

## 5. Before the first deploy

- [ ] `make test` green locally; CI green (BE, FE, scaffold, Docker images).
- [ ] Seed users: change or remove `admin@example.com` / `demo@example.com`
      (`be/internal/database/seeders/user_seeder.go`) — they are for development.
- [ ] `AUTO_MIGRATE=false` with several replicas; run `go run ./cmd/migrate up` and
      `go run ./cmd/seed` as a release step.
- [ ] Health checks: orchestrator probes `/api/ready`; logs in JSON (`LOG_FORMAT=json`).
- [ ] Set `AUTH_*` limits for your traffic (rate limit, lockout).
- [ ] Review Admin → Settings (sign-up, support email, audit retention).
- [ ] Back up Postgres (and the upload volume or bucket).

## 6. Keep in sync with the base (optional)

Add the base as a remote (`git remote add base <url>`) and merge `base/main` when you
want core fixes. Conflicts stay small when product code lives in modules.

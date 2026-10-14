# Next.js + React + TypeScript

**Next.js 15 App Router** + React 19 + TypeScript + Tailwind CSS v4 — the core admin/app
shell of the base. Optional features plug in as modules (`src/config/modules.ts`, guide:
[`../docs/modules.md`](../docs/modules.md)). This folder is a standalone pnpm project
(own `pnpm-lock.yaml`, `Dockerfile`, `.env`).

## Structure

```
src/
  app/                  ← Route wiring only (metadata, layout guards, thin page.tsx)
  components/
    common/             ← Shared presentational components (AuthLayout, Breadcrumbs, ...)
    providers/          ← React context providers
    ui/                 ← shadcn/ui components
  config/               ← env.ts, i18n.ts, modules.ts (module registry)
  features/             ← All domain logic + page components
    landing/            ← Public `/` (LandingPage)
    access-control/     ← UsersTable, RolesTable, permission hooks, admin nav
    admin-dashboard/
    user-profile/
    tools/              ← Module: Tools hub + Webhooks (module.ts descriptor)
    auth/
  hooks/                ← Shared custom hooks
  layouts/              ← AdminLayout
  lib/                  ← utils.ts (cn), permissions.ts (definePermissions)
  locales/              ← en/, vi/
  services/             ← api-client, auth-token-service, mock/
  store/                ← Zustand stores
  styles/               ← index.css (Tailwind + CSS variables)
  types/                ← Cross-feature TypeScript types (modules.ts = AppModule)
  utils/                ← Helpers
```

### Key design principle

`src/app/` files are **thin**: they handle routing, export `metadata` for SEO, and wire layout guards. All actual UI and business logic lives in `src/features/`.

```
app/admin/page.tsx               ← 12 lines: metadata + render feature component
  └─ features/access-control/   ← all the real code
```

No intermediate `src/screens/` or `src/pages/` layer — those are redundant.

## App Router

```
src/app/
  layout.tsx                ← root layout (AppProviders, AuthHydrator, Toaster)
  not-found.tsx             ← 404
  sitemap.ts / robots.ts    ← SEO
  (public)/
    layout.tsx              ← passthrough
    page.tsx                ← / → features/landing/LandingPage
  (auth)/
    layout.tsx              ← GuestGuard
    login/page.tsx          ← /login → features/auth/LoginPage
    register/page.tsx       ← /register → features/auth/RegisterPage
  admin/
    layout.tsx              ← ProtectedGuard + AdminGuard + AdminLayout
    page.tsx                ← /admin → features/admin-dashboard + PermissionGuard
    users/page.tsx          ← /admin/users → features/access-control/UsersTable
    roles/page.tsx          ← /admin/roles → features/access-control/RolesTable
    profile/page.tsx        ← /admin/profile → features/user-profile
    tools/webhooks/page.tsx ← /admin/tools/webhooks → features/tools (module)
```

## Getting Started

**Docker (full stack, from repo root):**

```bash
cp .env.example .env && make up
# → http://localhost:8080
```

**Local FE only:**

```bash
cp .env.example .env
pnpm install
pnpm dev
```

**Production image (this folder only):** `docker build -t my-fe --build-arg NEXT_PUBLIC_API_BASE_URL=https://api.example.com/api .`

### Demo accounts

| Role | Email | Password |
|------|-------|----------|
| Admin | `admin@example.com` | `admin1234` |
| User | `demo@example.com` | `password123` |

Works with mock API or real BE when `NEXT_PUBLIC_USE_MOCK_API=false`.

## Env Variables

| Run mode | Env file |
|----------|----------|
| **Docker stack** | `fe/.env` via `env_file` on `fe` container |
| **Local host** | `fe/.env` |

Key vars (via `src/config/env.ts`):

| Variable | Docker (`fe/.env`) | Local split dev (`fe/.env`) |
|----------|-------------------|----------------------------|
| `NEXT_PUBLIC_API_BASE_URL` | `http://system.local:8080/api` | `http://localhost:8080/api` |
| `NEXT_PUBLIC_SITE_URL` | `http://system.local:8080` | `http://localhost:3000` |
| `NEXT_PUBLIC_USE_MOCK_API` | `false` | `false` (or `true` for mock-only) |

## Related docs

- Architecture rules: [`AGENTS.md`](AGENTS.md)
- Monorepo + Docker: [`../README.md`](../README.md), [`../docker/README.md`](../docker/README.md)

## Scripts

```bash
pnpm dev       # next dev --turbopack
pnpm build     # next build
pnpm start     # next start
pnpm test      # vitest watch
pnpm test:run  # vitest run (CI)
pnpm lint      # eslint
```

## Stack

- Next.js 15 · React 19 · TypeScript ~6.0
- Tailwind CSS v4 · shadcn/ui (Radix)
- TanStack Query v5 · Zustand v5
- React Hook Form · Zod v4
- i18next (en / vi) · next-themes
- Axios with JWT refresh interceptor (`withCredentials` + `X-Requested-With`; the BE keeps the refresh token in an httpOnly cookie, so only the access token is in localStorage)
- Vitest + Testing Library

## API types & errors

- `pnpm gen:api` regenerates `src/types/api.gen.ts` from `../be/api/openapi.json` (or `OPENAPI_SOURCE=<url or file>` when the FE lives alone). Import request/response types from `@/types/api.gen` instead of hand-writing them; CI fails when the file is stale.
- Failed calls: branch on `apiErrorCode(error)` (`@/lib/api-error`), and show `translateApiError(error, t, fallback)` — it maps BE codes to `common:errors.<code>` in EN/VI.

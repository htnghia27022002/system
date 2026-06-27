# Next.js + React + TypeScript

**Next.js 15 App Router** + React 19 + TypeScript + Tailwind CSS v4.

## Structure

```
src/
  app/                  ← Route wiring only (metadata, layout guards, thin page.tsx)
  components/
    common/             ← Shared presentational components (AuthLayout, Breadcrumbs, ...)
    providers/          ← React context providers
    ui/                 ← shadcn/ui components
  config/               ← env.ts, i18n.ts
  features/             ← All domain logic + page components
    auth/               ← LoginPage, RegisterPage, guards, forms, hooks
    home/               ← HomePage
    admin-dashboard/    ← AdminDashboardOverview, charts
    access-control/     ← UsersTable, RolesTable, permission hooks
  hooks/                ← Shared custom hooks
  layouts/              ← MainLayout, AdminLayout (route layout shells)
  lib/                  ← utils.ts (cn)
  locales/              ← en/, vi/
  services/             ← api-client, auth-token-service, mock/
  store/                ← Zustand stores
  styles/               ← index.css (Tailwind + CSS variables)
  types/                ← Cross-feature TypeScript types
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
    layout.tsx              ← MainLayout
    page.tsx                ← / → features/home/HomePage
  (auth)/
    layout.tsx              ← GuestGuard
    login/page.tsx          ← /login → features/auth/LoginPage
    register/page.tsx       ← /register → features/auth/RegisterPage
  admin/
    layout.tsx              ← ProtectedGuard + AdminGuard + AdminLayout
    page.tsx                ← /admin → features/admin-dashboard + PermissionGuard
    users/page.tsx          ← /admin/users → features/access-control/UsersTable
    roles/page.tsx          ← /admin/roles → features/access-control/RolesTable
```

## Getting Started

**Docker (full stack, from repo root):**

```bash
cp .env.docker.example .env && make up
# → http://localhost:8080
```

**Local FE only:**

```bash
cp .env.example .env
pnpm install
pnpm dev
```

### Demo accounts

| Role | Email | Password |
|------|-------|----------|
| Admin | `admin@example.com` | `admin1234` |
| User | `demo@example.com` | `password123` |

Works with mock API or real BE when `NEXT_PUBLIC_USE_MOCK_API=false`.

## Env Variables

See `fe/.env.example`. Key vars (via `src/config/env.ts`):

| Variable | Docker | Local + real BE |
|----------|--------|-----------------|
| `NEXT_PUBLIC_API_BASE_URL` | `http://localhost:8080/api` | `http://localhost:8080/api` |
| `NEXT_PUBLIC_USE_MOCK_API` | `false` | `false` |
| `NEXT_PUBLIC_APP_NAME` | `System App` | `System App` |

For mock-only local dev: `NEXT_PUBLIC_USE_MOCK_API=true`, `NEXT_PUBLIC_API_BASE_URL=http://localhost:3000/api`.

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
- Axios with JWT refresh interceptor
- Vitest + Testing Library

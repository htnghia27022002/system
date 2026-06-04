# Next.js + React + TypeScript Source Structure

This project uses **Next.js 15 App Router** with React 19, TypeScript, and Tailwind CSS v4.

## Project Root

```txt
fe/
  public/
  src/
    app/          ← Next.js App Router (routing, layouts, pages)
    assets/       ← images, fonts, static media
    components/   ← reusable UI components (shared presentation)
    config/       ← environment and runtime configuration
    features/     ← feature/domain modules (auth, admin-dashboard, access-control)
    hooks/        ← shared custom hooks
    layouts/      ← reusable layout shells used inside app/ layouts
    lib/          ← utility helpers (cn, etc.)
    locales/      ← i18n JSON files (en/, vi/)
    screens/      ← page-level React components imported by app/ pages
    services/     ← base API clients and external integrations
    store/        ← global state (Zustand)
    styles/       ← global styles and CSS variables
    types/        ← cross-feature TypeScript types
    utils/        ← helper functions and constants
  .env.example
  next.config.ts
  postcss.config.mjs
  tsconfig.json
  vitest.config.ts
```

## App Router Structure

```txt
src/app/
  layout.tsx              ← root Server Component layout (wraps AppProviders)
  not-found.tsx           ← 404 page
  sitemap.ts              ← generates /sitemap.xml
  robots.ts               ← generates /robots.txt
  (public)/
    layout.tsx            ← MainLayout wrapper
    page.tsx              ← / (home)
  (auth)/
    layout.tsx            ← GuestGuard (redirects authenticated users)
    login/page.tsx        ← /login
    register/page.tsx     ← /register
  admin/
    layout.tsx            ← ProtectedGuard + AdminGuard + AdminLayout
    page.tsx              ← /admin (dashboard overview)
    users/page.tsx        ← /admin/users
    roles/page.tsx        ← /admin/roles
```

## `src` Structure

```txt
src/
  app/                  ← Next.js App Router (DO NOT put page components here)
  assets/               ← images, fonts, static media
  components/
    common/             ← shared presentational components
    providers/          ← React context providers ('use client')
    ui/                 ← shadcn/ui components ('use client' where needed)
  config/               ← env.ts (Zod-validated NEXT_PUBLIC_* vars)
  features/
    auth/               ← login, register, auth guards, JWT token mgmt
    admin-dashboard/    ← admin overview page components
    access-control/     ← users, roles, permissions CRUD
  hooks/                ← shared custom hooks ('use client')
  layouts/              ← MainLayout, AuthLayout, AdminLayout (used in app/ layouts)
  lib/                  ← utils.ts (cn helper)
  locales/              ← en/common.json, en/admin.json, vi/...
  screens/              ← page-level components (imported by app/ pages)
  services/             ← api-client, auth-token-service, health-api, mock/
  store/                ← auth-store, global-loading-store (Zustand, 'use client')
  styles/               ← index.css (Tailwind + CSS variables)
  types/                ← auth.ts, navigation.ts
  utils/                ← mock-jwt.ts, mock-delay.ts
```

> **Note**: `src/screens/` is where page-level React components live (previously `src/pages/`).
> `src/pages/` is intentionally avoided because Next.js treats it as the Pages Router directory.

## Folder Responsibility Rules

- `src/app/` handles routing. Page files export `metadata` and render screen components.
- `src/screens/` contains the actual page UI components (`'use client'`).
- `src/layouts/` contains layout shells (header, sidebar, auth wrapper). Used inside `src/app/**/layout.tsx`.
- `src/features/` owns business logic for each domain. Each feature has components, hooks, services, types.
- `src/components/` contains truly reusable UI pieces, not feature-specific logic.
- `src/services/` hosts shared API plumbing (HTTP client, JWT interceptors).
- `src/store/` holds global Zustand stores.

## Client vs. Server Components

- Files in `src/app/**/layout.tsx` and `src/app/**/page.tsx` are **Server Components** by default.
- Add `'use client'` to any file using React hooks (`useState`, `useEffect`, `useQuery`, TanStack Query, Zustand, etc.).
- Providers, stores, hooks, feature components, and layout shells are Client Components.
- Public pages in `src/app/` should export `metadata` for SEO.

## Getting Started

```bash
cd fe
cp .env.example .env.local
pnpm install
pnpm dev
```

Auth routes: `/login`, `/register`
Admin back office: `/admin`

### Demo accounts (mock API)

When `NEXT_PUBLIC_USE_MOCK_API=true` (default in `.env.example`):

| Role | Email | Password | After login |
|------|-------|----------|-------------|
| Admin | `admin@example.com` | `admin1234` | `/admin` |
| User | `demo@example.com` | `password123` | `/` |

Non-admin users are redirected away from `/admin`. Mock responses use `NEXT_PUBLIC_MOCK_API_DELAY_MS` (default `1200`).

## Stack

- **Next.js 15** App Router + React 19 + TypeScript ~6.0
- **Tailwind CSS v4** + shadcn/ui (Radix)
- **TanStack Query v5**, Zustand v5
- **i18next** (`en` / `vi`)
- **Axios** API client with JWT refresh interceptor
- **Sonner** toasts, **Recharts**, **react-day-picker** Calendar
- **Vitest** + **Testing Library** (config: `vitest.config.ts`)
- **pnpm** v10

## Scripts

```bash
pnpm dev       # next dev --turbopack (hot reload)
pnpm build     # next build (production)
pnpm start     # next start (serve production build)
pnpm test      # vitest (watch mode)
pnpm test:run  # vitest run (CI)
pnpm lint      # eslint .
```

## Path Alias

- Use the `@/*` alias to reference `src/*` (configured in `tsconfig.json`).
- Prefer `@/features/...` style imports over long relative chains (`../../../`).

## Environment Variables

All client-accessible env vars must be `NEXT_PUBLIC_*`:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:3000/api
NEXT_PUBLIC_APP_NAME=System App
NEXT_PUBLIC_USE_MOCK_API=true
NEXT_PUBLIC_MOCK_API_DELAY_MS=1200
NEXT_PUBLIC_SITE_URL=https://yourdomain.com
```

Accessed via `src/config/env.ts` (Zod-validated).

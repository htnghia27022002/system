# FE Agent Rules (Next.js + React + TypeScript)

This file defines strict implementation rules for AI agents working in `fe/`.
Follow these rules before creating, editing, moving, or deleting files.

## 1) Source of Truth

- Use `README.md` in this folder as the architecture baseline.
- Use `src/styles/index.css` (`:root` / `.dark` CSS variables) as the visual theme baseline for all UI work.
- `DESIGN.md` is deprecated legacy reference only — do not apply its Apple tokens to new work.

## 2) Stack

**Next.js 15 App Router** + React 19 + TypeScript.

- **Routing**: `src/app/` (App Router, file-based). No `react-router-dom`.
- **Links / navigation**: `Link` from `next/link`; `useRouter`, `usePathname`, `useParams` from `next/navigation`.
- **Build**: `pnpm dev` (`next dev --turbopack`) | `pnpm build` (`next build`)
- **Env vars**: `NEXT_PUBLIC_*` → `process.env.*`. Centralized in `src/config/env.ts`.
- **CSS**: Tailwind CSS v4 via PostCSS (`postcss.config.mjs`).
- **Tests**: Vitest + Testing Library, config in `vitest.config.ts`.

## 3) `src` Structure

```
src/
  app/            ← Next.js App Router (routing, metadata, layout wiring only)
  components/
    common/       ← shared presentational components (AuthLayout lives here)
    providers/    ← React context providers ('use client')
    ui/           ← shadcn/ui components
  config/         ← env.ts (Zod-validated NEXT_PUBLIC_* vars), i18n.ts
  features/       ← all domain logic + page-level components
    auth/         ← LoginPage, RegisterPage, guards, forms, hooks
    home/         ← HomePage
    admin-dashboard/
    access-control/
  hooks/          ← shared custom hooks ('use client')
  layouts/        ← MainLayout, AdminLayout (true route layout shells only)
  lib/            ← utils.ts (cn helper)
  locales/        ← en/, vi/ i18n JSON
  services/       ← api-client, auth-token-service, mock/
  store/          ← Zustand stores ('use client')
  styles/         ← index.css (Tailwind + CSS variables)
  types/          ← cross-feature TypeScript types
  utils/          ← helper functions
```

**Key rules:**
- `src/app/` files are thin: route wiring, `metadata` exports, layout guards. No business logic.
- `src/features/` owns all domain logic **and** page-level components.
- `src/layouts/` only contains `MainLayout` and `AdminLayout` — shells used inside `src/app/**/layout.tsx`.
- `src/components/common/` contains truly shared UI including `AuthLayout`.
- **Do NOT create `src/pages/` or `src/screens/`** — `src/pages/` conflicts with Next.js Pages Router; `src/screens/` was removed as redundant.

## 4) App Router Structure

```
src/app/
  layout.tsx                ← root Server Component + AppProviders
  not-found.tsx             ← 404
  sitemap.ts / robots.ts    ← SEO
  (public)/
    layout.tsx              ← wraps with MainLayout
    page.tsx                ← imports HomePage from @/features/home
  (auth)/
    layout.tsx              ← GuestGuard
    login/page.tsx          ← imports LoginPage from @/features/auth
    register/page.tsx       ← imports RegisterPage from @/features/auth
  admin/
    layout.tsx              ← ProtectedGuard + AdminGuard + AdminLayout
    page.tsx                ← imports AdminDashboardOverview + PermissionGuard
    users/page.tsx          ← imports UsersTable + PermissionGuard
    roles/page.tsx          ← imports RolesTable + PermissionGuard
```

## 5) Client vs. Server Components

- `src/app/**/layout.tsx` and `src/app/**/page.tsx` are **Server Components** by default.
  - They export `metadata` for SEO and import Client Components.
  - They do NOT use React hooks.
- Any file using hooks (`useState`, `useEffect`, `useQuery`, `useRouter`, Zustand, etc.) must start with `'use client'`.
- All feature components, stores, providers, hooks, and layout shells are Client Components.

## 6) Feature-First Rule

- All domain logic and page-level components belong in `src/features/<feature-name>/`.
- Each feature owns: `components/`, `hooks/`, `services/`, `schemas/`, `types.ts`, `index.ts`.
- `src/app/` pages directly import from `src/features/` — no intermediate layer needed.
- Shared UI only goes to `src/components/` when it is truly cross-feature.

## 7) Auth Guards

All guards are Client Components in `src/features/auth/` (or `access-control/`):

| Guard | Where used | Behavior |
|---|---|---|
| `GuestGuard` | `app/(auth)/layout.tsx` | Redirects logged-in users |
| `ProtectedGuard` | `app/admin/layout.tsx` | Redirects unauthenticated to `/login` |
| `AdminGuard` | `app/admin/layout.tsx` | Redirects non-admin to `/` |
| `PermissionGuard` | Inside each admin `page.tsx` | Redirects without permission to `/admin` |

## 8) File Placement Rules

| What | Where |
|---|---|
| Route entry (metadata + render) | `src/app/**/page.tsx` |
| Route layout shell | `src/app/**/layout.tsx` |
| Domain page component | `src/features/<feature>/components/` |
| Layout shell (MainLayout, AdminLayout) | `src/layouts/` |
| Auth card wrapper | `src/components/common/auth-layout.tsx` |
| Shared UI primitives | `src/components/ui/` (shadcn) |
| Shared presentational components | `src/components/common/` |
| Global hooks | `src/hooks/` |
| Zustand stores | `src/store/` |
| API client / interceptors | `src/services/` |
| Env config | `src/config/env.ts` |
| i18n JSON | `src/locales/` |
| Global styles | `src/styles/index.css` |
| Tests | Colocated `*.test.ts` / `*.test.tsx` |

## 9) SEO Rules

- Every public page in `src/app/**/page.tsx` must export `metadata`.
- Admin and auth pages: `robots: { index: false, follow: false }`.
- `src/app/sitemap.ts` lists public URLs.
- `src/app/robots.ts` disallows `/admin/`.

## 10) Import Rules

- Use `@/*` alias for all imports (mapped to `src/*`).
- Import features through `@/features/<feature>` (public `index.ts`), not deep internal files.
- `src/app/` imports from `features/`, `components/`, `layouts/`.
- `src/features/` imports from `components/`, `hooks/`, `services/`, `utils/`, `types/`, `config/`.
- `src/layouts/` imports from `components/`, `hooks/`, `utils/`, `types/`, `config/`.
- Features must NOT import from other features' internals.
- **Do NOT import from `react-router-dom`** — removed. Use `next/link` and `next/navigation`.

## 11) Env Variables

- All client-side env vars must be `NEXT_PUBLIC_*`.
- Access via `src/config/env.ts` — do NOT use `process.env.*` directly in feature code.

## 12) Naming Conventions

- Folders: `kebab-case`
- React component files: `PascalCase.tsx`
- Hook files: `useXxx.ts`
- Utility files: `camelCase.ts`

## 13) Before Finishing Any Task

1. Is the file placed in the correct folder?
2. Is it a Client Component if it uses hooks? (has `'use client'`)
3. Does a new public `app/page.tsx` export `metadata`?
4. Are feature boundaries respected?
5. Is content written in English?
6. Was a shadcn component used before writing custom markup?

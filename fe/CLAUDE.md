# FE Agent Rules (Next.js + React + TypeScript)

This file defines strict implementation rules for Claude when working in `fe/`. Claude Code loads it automatically for any file under `fe/`; project-wide rules live in `.claude/rules/`, skills in `.claude/skills/` (`fe-develop`, `design-taste-frontend`), and subagents in `.claude/agents/` (`fe`).
Follow these rules before creating, editing, moving, or deleting files.

## 0) Self-contained package

`fe/` is a **standalone frontend package**. It can be moved out of this workspace without changing internal import paths.

- Package root: `fe/package.json` + `fe/pnpm-lock.yaml` + `fe/.npmrc` (not a pnpm workspace member) and `fe/Dockerfile` (build context = `fe/`)
- Internal imports: `@/*` → `src/*` only
- API integration: HTTP via `NEXT_PUBLIC_API_BASE_URL` — do **not** import from `be/`
- Run all FE commands from `fe/` (or repo root scripts that delegate to `fe/`)

When relocating the frontend, move the entire `fe/` directory. Update only external wiring (Docker, CI, env) at the new location.

## 1) Source of Truth

- Use `README.md` in this folder as the architecture baseline.
- Use **[`DESIGN.md`](./DESIGN.md)** for all visual / UI design rules (theme tokens, page spacing, mobile, clickable table links, component consistency).
- Theme tokens live in `src/styles/index.css` (`:root` / `.dark`) — see `DESIGN.md`.

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
  config/         ← env.ts (Zod-validated NEXT_PUBLIC_* vars), i18n.ts, modules.ts (module registry)
  features/       ← all domain logic + page-level components
    auth/
    landing/      ← public `/`
    access-control/
    admin-dashboard/
    user-profile/
    tools/        ← optional module (module.ts descriptor; Webhooks tool)
  hooks/          ← shared custom hooks ('use client')
  layouts/        ← AdminLayout (admin chrome only)
  lib/            ← utils.ts (cn helper), permissions.ts (definePermissions)
  locales/        ← en/, vi/ i18n JSON
  services/       ← api-client, auth-token-service, mock/
  store/          ← Zustand stores ('use client')
  styles/         ← index.css (Tailwind + CSS variables)
  types/          ← cross-feature TypeScript types (modules.ts = AppModule contract)
  utils/          ← helper functions
```

**Key rules:**
- `src/app/` files are thin: route wiring, `metadata` exports, layout guards. No business logic.
- `src/features/` owns all domain logic **and** page-level components.
- `src/layouts/` contains `AdminLayout` only — used from `src/app/admin/layout.tsx`. Public marketing/tools chrome lives in `PublicSiteHeader` / landing nav.
- `src/components/common/` contains truly shared UI including `AuthLayout`.
- **Do NOT create `src/pages/` or `src/screens/`** — `src/pages/` conflicts with Next.js Pages Router.

## 4) App Router Structure

```
src/app/
  layout.tsx                ← root Server Component + AppProviders
  not-found.tsx             ← 404
  sitemap.ts / robots.ts    ← SEO
  (public)/
    layout.tsx              ← passthrough (LandingNav / PublicSiteHeader inside features)
    page.tsx                ← imports LandingPage from @/features/landing
  (auth)/
    layout.tsx              ← GuestGuard
    login/page.tsx          ← imports LoginPage from @/features/auth
    register/page.tsx       ← imports RegisterPage from @/features/auth
  admin/
    layout.tsx              ← ProtectedGuard + AdminGuard + AdminLayout
    page.tsx                ← imports AdminDashboardOverview + PermissionGuard
    users/page.tsx          ← imports UsersTable + PermissionGuard
    roles/page.tsx          ← imports RolesTable + PermissionGuard
    profile/page.tsx        ← ProfilePage
    tools/webhooks/page.tsx ← module route (tools)
```

## 4.1) Modules (optional features)

Core never imports a module. An optional/business feature is a module:

1. `src/features/<name>/module.ts` exports an `AppModule` (`src/types/modules.ts`): admin nav entries, breadcrumbs, permission catalog rows, public links, sitemap paths, and `backendModule` (the BE module it needs).
2. `src/features/<name>/permissions.ts` defines its keys with `definePermissions('<resource>')` (`@/lib/permissions`).
3. Its routes live under `src/app/` (file routing cannot be registered dynamically) and use `PermissionGuard`.
4. Append it to `src/config/modules.ts` — the only file that lists modules. Sidebar, breadcrumbs, public header, sitemap and the mock RBAC catalog are composed from it; the sidebar also hides a module whose `backendModule` is missing from `GET /api/app/info`.

`src/config/modules.ts` is the one place allowed to deep-import `features/<name>/module.ts`. Remove a module by deleting its feature folder, its `app/` routes and its line there. Guide: [`../docs/modules.md`](../docs/modules.md).

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

### 7.1) New feature permissions + admin menu (mandatory)

When adding an admin page or sidebar entry:

1. Core resource: extend `PermissionKeys` in `src/features/access-control/permission-keys.ts`. Module resource: `definePermissions` in the module's `permissions.ts`.
2. Wrap the page with `PermissionGuard` using the **view** key.
3. Core page: add the nav item in `src/features/access-control/hooks/use-admin-nav-items.ts`. Module page: add it to the module's `adminNav` in `module.ts`. Either way it carries the **view** permission and is filtered with `hasPermission` — never show the link without the key.
4. Gate mutate actions with `PermissionGate` + **modify** key.
5. Align with `docs/features/.../contracts/permissions.md`.

See `.claude/rules/feature-permissions.md`.

## 8) File Placement Rules

| What | Where |
|---|---|
| Route entry (metadata + render) | `src/app/**/page.tsx` |
| Route layout shell | `src/app/**/layout.tsx` |
| Domain page component | `src/features/<feature>/components/` |
| Layout shell (AdminLayout) | `src/layouts/` |
| Auth card wrapper | `src/components/common/auth-layout.tsx` |
| Shared UI primitives | `src/components/ui/` (shadcn) |
| Shared presentational components | `src/components/common/` |
| Global hooks | `src/hooks/` |
| Zustand stores | `src/store/` |
| API client / interceptors | `src/services/` |
| Env config | `src/config/env.ts` |
| i18n JSON | `src/locales/` (`en/` English; `vi/` UTF-8 Vietnamese **with diacritics** — see `.claude/rules/locale-vi-utf8.md`) |
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
- Exception: `store` and some services may import `@/features/auth/services/auth-api` to avoid circular barrels.
- Admin chrome (`components/common` sidebar/nav, `layouts`) may import `@/features/access-control` and `@/features/auth` **barrels** (RBAC/session kernel). One-way only.
- `src/app/` imports from `features/`, `components/`, `layouts/`.
- `src/features/` imports from `components/`, `hooks/`, `services/`, `utils/`, `types/`, `config/`.
- Features must NOT import from other features' internals.
- **Do NOT import from `react-router-dom`** — removed. Use `next/link` and `next/navigation`.

## 11) Env Variables

- All client-side env vars must be `NEXT_PUBLIC_*`.
- Access via `src/config/env.ts` — do NOT use `process.env.*` directly in feature code.

## 12) Naming Conventions

Match shadcn/ui and [bulletproof-react](https://github.com/alan2207/bulletproof-react):

- Folders: `kebab-case`
- React components: `kebab-case.tsx` (`users-table.tsx`, `login-page.tsx`)
- Hooks: `use-users.ts`
- Utilities: `kebab-case.ts`

Do not introduce `PascalCase.tsx` page files. Existing `LandingPage.tsx` / `ToolsHubPage.tsx` are legacy; new files follow kebab-case.

## 13) Before Finishing Any Task

1. Is the file placed in the correct folder?
2. Is it a Client Component if it uses hooks? (has `'use client'`)
3. Does a new public `app/page.tsx` export `metadata`?
4. Are feature boundaries respected?
5. Is content written in English?
6. Does UI follow [`DESIGN.md`](./DESIGN.md) and `.claude/rules/fe-ux-product.md`? (theme tokens, matching page padding, mobile, clickable cells as links, loading/empty/error, shadcn first)
7. Does a new feature follow `.claude/rules/fe-feature-module.md`?

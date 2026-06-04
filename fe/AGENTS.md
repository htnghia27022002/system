# FE Agent Rules (Next.js + React + TypeScript)

This file defines strict implementation rules for AI agents working in `fe/`.
Follow these rules before creating, editing, moving, or deleting files.

## 1) Source of Truth

- Use `README.md` in this folder as the architecture baseline.
- Use `src/styles/index.css` (`:root` / `.dark` CSS variables) as the visual theme baseline for all UI work.
- `DESIGN.md` is deprecated legacy reference only — do not apply its Apple tokens to new work.
- Do not introduce a different folder strategy unless explicitly requested by the user.
- Prefer consistency with existing patterns over personal preference.

## 2) Stack

This project uses **Next.js 15 App Router** with React 19 and TypeScript.

- **Routing**: `src/app/` (App Router, file-based). No `react-router-dom`.
- **Links / navigation**: `Link` from `next/link`; `useRouter`, `usePathname`, `useParams` from `next/navigation`.
- **Build**: `pnpm dev` → `next dev --turbopack` | `pnpm build` → `next build`
- **Env vars**: `NEXT_PUBLIC_*` (reads from `process.env.*`). Centralized in `src/config/env.ts`.
- **CSS**: Tailwind CSS v4 via PostCSS (`postcss.config.mjs`), no Vite plugin.
- **Tests**: Vitest + Testing Library, config in `vitest.config.ts` (separate from Next.js build).

## 3) Required `src` Structure

Agents must keep code organized under this structure:

- `src/app` — Next.js App Router pages and layouts (file-based routing)
- `src/assets`
- `src/components`
- `src/features`
- `src/hooks`
- `src/layouts` — reusable layout shells (used inside `src/app/` layouts)
- `src/screens` — page-level React components (imported by `src/app/` pages)
- `src/services`
- `src/store`
- `src/styles`
- `src/types`
- `src/utils`
- `src/config`

> **Important**: Do NOT use `src/pages/`. Next.js treats `src/pages/` as the Pages Router
> directory. Page components live in `src/screens/` instead.

Do not create random top-level folders in `src` without user approval.

## 4) App Router File Structure

```
src/app/
  layout.tsx             ← root layout (Server Component), imports AppProviders
  not-found.tsx          ← 404 page
  sitemap.ts             ← SEO sitemap
  robots.ts              ← SEO robots
  (public)/
    layout.tsx           ← wraps with MainLayout
    page.tsx             ← home page
  (auth)/
    layout.tsx           ← GuestGuard (redirects logged-in users)
    login/page.tsx
    register/page.tsx
  admin/
    layout.tsx           ← ProtectedGuard + AdminGuard + AdminLayout
    page.tsx             ← admin overview
    users/page.tsx
    roles/page.tsx
```

## 5) Client vs. Server Components

- **Root layout** (`src/app/layout.tsx`) is a Server Component. It imports client providers via `AppProviders`.
- **All page files** in `src/app/**/page.tsx` are Server Components by default. They can export `metadata` (SEO) and import Client Components.
- **Client Components** must have `'use client'` at the top of the file.
- **Any file using React hooks** (`useState`, `useEffect`, `useQuery`, `useRouter`, etc.) must be a Client Component.
- **Stores, hooks, providers, feature components** — add `'use client'` if they use any hooks.
- **shadcn/ui components** that use hooks already have `'use client'` added.

## 6) Feature-First Rule

- Business/domain logic belongs in `src/features/<feature-name>/`.
- Each feature should own its local `components`, `hooks`, `services`, and `types`.
- Shared UI only goes to `src/components` when it is truly cross-feature.
- Do not place feature-specific logic in shared folders.

## 7) Auth Guards

Route protection is done via Client Component wrappers in layout files:

| Guard | File | Behavior |
|---|---|---|
| `GuestGuard` | `src/features/auth` | Redirects logged-in users away from auth routes |
| `ProtectedGuard` | `src/features/auth` | Redirects unauthenticated users to `/login` |
| `AdminGuard` | `src/features/auth` | Redirects non-admin users to `/` |
| `PermissionGuard` | `src/features/access-control` | Redirects users without a specific permission to `/admin` |

All guards use Zustand auth store + `useRouter` from `next/navigation`.

## 8) File Placement Rules

- Route entry points go in `src/app/` (Next.js convention).
- Page-level React components go in `src/screens/`.
- Layout shells go in `src/layouts/` (reused inside `src/app/**/layout.tsx`).
- Global API plumbing (base client/interceptors/adapters) goes in `src/services/`.
- Global reusable hooks go in `src/hooks/`.
- Global types go in `src/types/`; feature types stay inside feature folders.
- Runtime/environment configuration goes in `src/config/`.
- Generic helpers go in `src/utils/`.
- App-wide styles/themes go in `src/styles/`.
- Tests are colocated next to the file under test (`*.test.ts` / `*.test.tsx`).

## 9) SEO Rules

- Every public page file in `src/app/**/page.tsx` must export a `metadata` object or `generateMetadata` function.
- Admin and auth pages should set `robots: { index: false, follow: false }`.
- `src/app/sitemap.ts` lists all public URLs.
- `src/app/robots.ts` disallows `/admin/` for crawlers.

## 10) Import and Boundary Rules

- Avoid deep imports into another feature's internals.
- Prefer exporting a feature's public surface via `src/features/<feature>/index.ts`.
- `src/app/` pages should compose from `src/screens/`, `src/features/` and shared modules.
- Use the `@/*` path alias (mapped to `src/*`) instead of long relative chains (`../../../`).
- Import features through their public entry (`@/features/<feature>`), not deep internal files.
- **Do NOT import from `react-router-dom`** — it is removed. Use `next/link` and `next/navigation`.

## 11) Env Variables

- All env vars must be `NEXT_PUBLIC_*` for client-side access.
- Access via `src/config/env.ts` — do NOT use `process.env.*` directly in feature code.
- Internal property names in `env` object are kept for backward compatibility (e.g., `env.VITE_API_BASE_URL`).

## 12) Naming Conventions

- Folder names: `kebab-case`
- React component files: `PascalCase.tsx`
- Hook files: `useXxx.ts`
- Utility files: `camelCase.ts`
- Keep names explicit and domain-oriented.

## 13) Editing Behavior for Agents

- Make minimal, targeted changes.
- Do not refactor unrelated modules in the same task.
- Do not move files across layers without a clear architectural reason.
- Do not add new dependencies unless the task requires them.
- Preserve backward compatibility unless breaking changes are requested.

## 14) Quality Guardrails

- Keep components focused (UI rendering first; move heavy logic into hooks/services).
- Avoid API calls directly inside deeply presentational components.
- Co-locate files that change together.
- Avoid excessive folder nesting unless complexity justifies it.

## 15) UI Library and Design Reference

- **Library-first:** Use installed stack components (`src/components/ui` from shadcn/ui, Radix primitives, TanStack Query, next-themes, etc.).
- **Registry before markup:** For any UI primitive, check `src/components/ui/` first. If missing, run `pnpm dlx shadcn@latest add <name> -y` in `fe/`.
- **Theme via CSS variables:** Customize colors, radius, surfaces in `src/styles/index.css`.

## 16) Language Rule

- Any newly created or modified file content must be written in English.
- Do not mix Vietnamese and English in source text unless explicitly requested.

## 17) Before Finishing Any Task

Agents should self-check:

1. Is the file placed in the correct folder?
2. Does it respect feature boundaries?
3. Are naming conventions followed?
4. Did we avoid unrelated refactors?
5. Is the content in English?
6. For UI work: was a matching shadcn component used or added before custom markup?
7. Does any new component using hooks have `'use client'`?
8. Does any new public page export `metadata`?

If any answer is "no", fix it before finalizing.

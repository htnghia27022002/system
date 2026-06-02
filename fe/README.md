# React + TypeScript Source Structure

This version follows the structure proposed in the DEV article and adapts it for a scalable React + TypeScript codebase.

## Project Root

```txt
my-app/
  public/
  src/
  .gitignore
  package.json
  README.md
  tsconfig.json
  vite.config.ts
```

## `src` Structure

```txt
src/
  assets/                # images, fonts, static media
  components/            # reusable UI components (shared presentation)
  features/              # feature/domain modules (auth, dashboard, profile...)
    auth/
      components/
      hooks/
      services/
      types.ts
      index.ts
  hooks/                 # shared custom hooks
  layouts/               # app layouts (main, admin, dashboard)
  pages/                 # route-level pages
  services/              # base API clients and external integrations
  store/                 # global state configuration
  styles/                # global styles and theme files
  types/                 # cross-feature TypeScript types
  utils/                 # helper functions and constants
  config/                # environment and runtime configuration
  App.tsx
  main.tsx
  router.tsx
```

Tests are colocated with the code they cover (for example `Button.tsx` + `Button.test.tsx` in the same folder), not gathered in a separate top-level test tree.

## Folder Responsibility Rules

- `features` owns business logic for each domain and can contain its own local components/hooks/services.
- `components` should contain only truly reusable UI pieces, not feature-specific business logic.
- `pages` compose features and shared components for route screens.
- `services` at root should host shared API plumbing (HTTP client, interceptors, adapters).
- `store` can be centralized or split by feature if the project grows.

## Scaling Rules

- Keep files colocated with the feature when they change together.
- Avoid deep import paths into another feature's internals.
- Prefer exposing feature APIs via `features/<feature>/index.ts`.
- Keep the tree shallow and readable; split a feature only when it becomes too large.

## Suggested Naming

- Folders: `kebab-case`
- Components: `PascalCase.tsx`
- Hooks: `useXxx.ts`
- Utility files: `camelCase.ts`
- Tests: `*.test.ts` / `*.test.tsx`, colocated next to the file under test

## Path Alias

- Use the `@/*` alias to reference `src/*` (configured in `vite.config.ts` and `tsconfig.json`).
- Prefer `@/feature/...` style imports over long relative chains (`../../../`).
- Import features through their public entry (`@/features/<feature>`), not deep internal files.

## Getting Started

```bash
cd fe
cp .env.example .env
pnpm install
pnpm dev
```

Auth routes: `/login`, `/register` (feature module under `src/features/auth`).

Admin back office: `/admin` (overview dashboard, `src/features/admin-dashboard`).

### Demo accounts (mock API)

When `VITE_USE_MOCK_API=true` (default in `.env.example`):

| Role | Email | Password | After login |
|------|-------|----------|-------------|
| Admin | `admin@example.com` | `admin1234` | `/admin` |
| User | `demo@example.com` | `password123` | `/` |

Non-admin users are redirected away from `/admin`. Mock responses use `VITE_MOCK_API_DELAY_MS` (default `1200`) so you can see button loading on auth and nav `Progress` + card `Skeleton` on the dashboard.

If roles behave unexpectedly after upgrading, clear `localStorage` key `mock_auth_users` and sign in again.

Stack initialized in this repo:

- Vite + React + TypeScript
- Tailwind CSS v4 + shadcn/ui (Radix)
- TanStack Query, React Router, Zustand
- i18next (`en` / `vi`)
- Axios API client with JWT refresh interceptor
- Sonner toasts, Recharts, react-day-picker Calendar

## Reference

- [Recommended Folder Structure for React 2025](https://dev.to/pramod_boda/recommended-folder-structure-for-react-2025-48mc)

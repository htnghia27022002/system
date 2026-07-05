# FE Implementation: Elasticsearch Search (Admin Unified Search)

**Feature ID:** `002-elasticsearch-search`  
**Scope:** Navbar typeahead, full search page, permission-aware navigation, table search UX

## Overview

Admin unified search lives in `fe/src/features/admin-search/`. Server-side list filters (Users/Roles tables) reuse `ServerListSearch` in access-control. Both search surfaces skip the global navbar loading bar and show loading locally.

## Routes

| Route | Component | Guard |
|-------|-----------|-------|
| `/admin/search` | `AdminSearchPage` | `SearchAccessGuard` — requires `users:view`, `roles:view`, or `permissions:view` |
| Navbar (all admin pages) | `AdminSearchBar` | Hidden link targets when caller lacks search access |

## Key paths

```text
fe/src/features/admin-search/
├── components/
│   ├── admin-search-bar.tsx      # Navbar combobox (200ms debounce)
│   ├── admin-search-page.tsx     # Full search + filters + pagination
│   ├── search-hit-list.tsx       # Result cards on search page
│   ├── search-access-guard.tsx   # Redirect to /admin if no entity view permission
│   └── search-entity-badge.tsx
├── hooks/
│   ├── use-admin-search.ts       # React Query → GET /admin/search
│   └── use-search-navigation.ts  # Permission-aware hit links
├── lib/
│   └── search-routes.ts          # searchHitHref, canNavigateToSearchHit
└── services/
    └── search-api.ts

fe/src/features/access-control/components/
└── server-list-search.tsx        # Debounced server-side table search

fe/src/hooks/use-nav-loading.ts   # Skips queries/requests with skipNavLoading
```

## Navbar search (`AdminSearchBar`)

- **Debounce:** 200ms via `useDebouncedValue`; Enter flushes immediately; Escape closes panel.
- **Stale debounce guard:** Updates `activeQuery` only when `debouncedQuery === query` to prevent search/clear loops.
- **UX:** Clear (X) button, local spinner, full-width dropdown, “View all results” footer → `/admin/search?q=…`.
- **Loading:** `meta.skipNavLoading` on React Query + `skipNavLoading: true` on axios GET — no top `NavLoadingBar`.

## Search page (`AdminSearchPage`)

- URL-driven: `?q=`, `types=`, `page=`.
- Entity-type filter options respect caller permissions (`useSearchNavigation`).
- Initial load: spinner in results area only; refetch: overlay on results block.
- Same `skipNavLoading` pattern as navbar search.

## Hit navigation (permission-aware)

`search-routes.ts` + `useSearchNavigation`:

| Entity | Target | Required permission |
|--------|--------|---------------------|
| User | `/admin/users?search=…` | `users:view` |
| Role | `/admin/roles?search=…` | `roles:view` |
| Permission | `/admin/roles?permissionKey=…` | `roles:view` |

Hits without list-page access render as non-clickable rows (BE already filters by entity view permission; FE adds defense in depth).

## Table search (Users / Roles lists)

`ServerListSearch` in `users-table.tsx` / `roles-table.tsx`:

- Debounced URL params via `useListQueryParams` (300ms default).
- **Loop fix:** Calls `onSearch` only when debounced value matches draft and differs from URL; callbacks stored in refs.
- **URL dedup:** `setParams` skips `router.push` when normalized query string unchanged.
- **Loading:** `skipNavLoading` on list APIs + `isRefreshing` overlay on `DataTable` (not navbar).

## Global loading skip mechanism

| Layer | Mechanism |
|-------|-----------|
| React Query | `meta: { skipNavLoading: true }` on search + list queries |
| Axios | `skipNavLoading: true` on request config (`fe/src/global.d.ts`) |
| Nav bar | `useNavLoading` predicate excludes skipped queries; interceptor skips tracked GETs |

## Related docs

- [be-implement.md](be-implement.md) — search API, outbox, queue worker
- [spec.md](spec.md) — requirements and RBAC rules

# Feature: Elasticsearch Search (Admin Unified Search)

**ID:** `002-elasticsearch-search`  
**Status:** Implemented (BE + FE)

## Workflow progress

| Phase | File | Status |
|-------|------|--------|
| 1 Analyze | [spec.md](spec.md) | Done |
| 2 Decompose | tasks.md | Pending |
| 3 Design | plan.md, be-implement.md, fe-implement.md | Done |
| 4 Implement | `be/` code | Done |
| 4 Implement | `fe/` code | Done |
| 5 Test | qa-checklist.md | Pending |

## Summary

Permission-aware unified admin search across users, roles, and permissions; reliable PostgreSQL-to-search-index sync via **transactional outbox + NATS JetStream**; bulk reindex for recovery; and graceful degradation when Elasticsearch is unavailable.

## Architecture (implemented)

```text
API (be)     → OutboxService → handlers/publisher → NATS (SubjectSearchOutbox)
Worker (queue) → Ensure stream/consumer → handlers/subscribers → IndexProcessor → Elasticsearch
```

- **Constants:** `be/internal/queue/constants.go` — all stream/subject/consumer names
- **Options:** `be/internal/queue/nats.json` — JetStream retention, ack policy, etc.
- **Worker:** `be/cmd/queue` (Docker service `queue`)
- **DI:** `be/internal/app/dependency/` — per-domain service resolvers; `container.go` orchestrates

Full design:

- BE: **[be-implement.md](be-implement.md)**
- FE: **[fe-implement.md](fe-implement.md)**

## Key paths

| Area | Path |
|------|------|
| Outbox table | `be/migrations/000004_search_outbox.up.sql` |
| DI resolvers | `be/internal/app/dependency/` |
| Queue infra | `be/internal/queue/` |
| Publish | `be/internal/handlers/publisher/` |
| Consume | `be/internal/handlers/subscribers/` |
| Search services | `be/internal/services/search/` |
| Admin API | `be/public/handlers/search_handler.go` |
| FE feature | `fe/src/features/admin-search/` |
| FE table search | `fe/src/features/access-control/components/server-list-search.tsx` |

## Recent updates (2026-07-05)

### Backend

- Refactored DI from monolithic `container.go` into `internal/app/dependency/` (infra, search stack, auth/user/role/permission resolvers, HTTP handlers).

### Frontend

- Navbar search: 200ms debounce, clear button, local loading, permission-aware hit links.
- Search page `/admin/search`: enabled with `SearchAccessGuard`; entity filters respect RBAC.
- Skip global navbar loading for search and list-table GET requests (`skipNavLoading`).
- Fixed debounce/URL sync loops in `ServerListSearch` and `AdminSearchBar`.
- Table search (Users/Roles): loading overlay on data table only.

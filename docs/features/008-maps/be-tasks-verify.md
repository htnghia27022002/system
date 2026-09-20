# BE tasks verify: Maps

**Feature:** `docs/features/008-maps/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@be`  
**Date:** 2026-09-12

## Summary

All `[BE]` tasks for 008-maps are implemented against the latest contracts (shared tables, integer enums, ingest counters `success` / `error` / `itemsSuccess` / `itemsError`). `make test-be` from the repo root passed, including `be/test/unit/maps/`. Restart compose `be` and `queue` before live API checks so migration `000010` and catalog/geo seed apply.

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T001 | Migration number `000010`; no map keys in env | Done | `be/migrations/000010_geo_and_data_sources.{up,down}.sql`; no Maps vars added to root/`.env` or `be/.env` |
| T005 | Seed `maps:view` / `maps:modify` | Done | `catalog.go` UUIDs `…019` / `…020`; `TestDefaultPermissionsUseViewModifyKeys` |
| T006 | Shared geo/place/source schema | Done | Migration creates `countries`, `administrative_divisions`, `locations`, `places`, `categories`, `news`, `data_sources`, `data_ingest_runs`, `data_ingest_run_sources` (no `map_*`) |
| T006a | VN + categories seed | Done | `CategorySeeder` (5 unique keys); `GeoSeeder` + `fixtures/vn_divisions.json` — see Seed scope |
| T007 | Reusable models | Done | `models/geo`, `location`, `place`, `category`, `news`, `datasource` — SMALLINT enums, `created_by` / `updated_by` |
| T008 | Maps DTOs | Done | `internal/dto/maps/` (pin, detail, source, ingest, mapping, geo) |
| T009 | Repositories | Done | Interfaces + pgx/squirrel impls for geo, location, place, category, news, sources, ingest runs |
| T010 | Services | Done | `services/maps/source`, `place`, `ingest` (interfaces only) |
| T011 | Handlers, routes, DI | Done | `maps_handler.go`, `routes/maps.go` under `/admin/maps`; additive `container.go` / `dependency/` |
| T012 | Queue extend | Done | `system.maps.ingest`, `nats.json` stream `maps`, `PublishMapsIngest`, `process_maps_ingest.go`, `cmd/queue` poll drain |
| T025 | Empty places list | Done | `GET /places` returns `{ items: [] }` when none pinnable |
| T026 | Places query + catalogs | Done | Filters `category`, `q`, `adminDivisionId` (path prefix), `countryCode`; categories + geo endpoints |
| T030 | Place detail + news | Done | `GET /places/:id` includes location + news attribution |
| T034 | Sources CRUD | Done | GET/POST/PATCH/DELETE `/sources` with URL/method/mapping validation |
| T035 | Mapping allow-list | Done | `ValidateFieldMapping`; `TestValidateFieldMappingRejectsUnknownTargets` |
| T040 | Start ingest + 409 | Done | `POST /ingest`; `TestConcurrentStartReturnsConflict` |
| T041 | Fetch enabled Sources only | Done | 30s timeout, 5 MiB cap; `TestDisabledSourceIsNotFetched`, `TestUnsavedURLIsNeverFetched` |
| T042 | Merge persist + per-source outcomes | Done | Location → Place → News; pending without coords; isolated source failures |
| T043 | Ingest status GET | Done | Latest + get-by-id include `sources[]` |
| T044 | Ingest unit tests | Done | `test/unit/maps/ingest_test.go` |
| T047 | Location reuse rules | Done | `findOrCreateLocation` (key → country+admin+street → formatted → rounded coords + formatted) |
| T048 | Place identity | Done | Unique `(location_id, place_key)`; coords never merge Places |
| T049 | Chung cư merge tests | Done | `test/unit/maps/merge_test.go` |
| T051 | Route RBAC | Done | `RequireView("maps")` / `RequireModify("maps")` after Auth |
| T052 | Place status PATCH | Done | `active`/`hidden` only; reject `active` without valid coords |
| T054 | RBAC unit tests | Done | `test/unit/maps/rbac_test.go` (403 / view-only / Super admin / 401) |
| T056 | Mapping + unsaved URL tests | Done | `mapping_test.go` + `TestUnsavedURLIsNeverFetched` |
| T058 | This verify report | Done | After `make test-be` |

## Verification commands

```bash
make test-be
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-be` | Pass | All `be/test/unit/...` packages OK, including `be/test/unit/maps` (0.012s) |
| `go build ./cmd/queue/` | Pass | Compiled inside compose `be` container |

GitNexus MCP was unavailable. Kernel edits (`DefaultPermissions`, `RegisterAdminRoutes`, `NewContainer`, queue constants / `nats.json` / `NewRegistry`, publisher) were **additive only**.

## Seed scope (T006a)

| Seeded | Count / keys | Notes |
|--------|----------------|-------|
| Country | VN (`VNM`, Việt Nam) | Idempotent upsert by `code` |
| Provinces | **All 34** 2025 provincial units | Official 2-digit codes; `path` like `VN/48/` |
| Communes / wards | **10-row subset** (not 3,321) | Hà Nội 2, Đà Nẵng 2, HCMC 3, Hải Phòng 1, Huế 1, Cần Thơ 1 — sample codes for filter/path testing |
| Categories | `room_rental`, `restaurant`, `hotel`, `eatery`, `uncategorized` | Unique `key`; `ON CONFLICT DO NOTHING` |

The full official 3,321-commune 2025 list is **not** vendored. `@qa` / later work can replace the fixture without a schema change.

## Acceptance coverage (BE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| Empty pinnable list | T025 list handler + empty items | Pass (unit / contract) |
| Category / q / division / country filters | T026 place list | Pass (implementation) |
| Pin detail + news attribution | T030 | Pass (implementation) |
| Sources CRUD + mapping allow-list | T034–T035, mapping tests | Pass |
| Load data enabled-only, 409, pending coords | T040–T044 | Pass |
| Two Places at one Location; unique URL | T047–T049 | Pass |
| View vs modify vs Super admin vs 401 | T051, T054 | Pass |
| Independent Tests on `system.local` | After `be`/`queue` restart | N/A — `@qa` |

## Gaps / follow-ups

- [ ] Restart compose `be` and `queue` so migration `000010` + `CategorySeeder` / `GeoSeeder` / catalog keys apply before `@fe` or `@qa` live checks
- [ ] Full 3,321-commune Vietnam fixture not seeded (documented subset only)
- [ ] Integration/e2e HTTP tests against Postgres were not required by `[BE]` tasks (`make test-be` is unit only)

## Sign-off (BE)

- [x] All claimed `[BE]` tasks done or explicitly deferred above
- [x] Tests listed above passed
- [x] Matches `plan.md` BE sections (or deviations documented)

**Ready for `@fe`:** yes — `/api/admin/maps/*` contracts are implemented (camelCase JSON). FE can bind places, geo, sources, and ingest polling. Restart `be` first so new routes and seed exist.

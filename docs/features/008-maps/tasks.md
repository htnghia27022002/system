# Tasks: Maps (Admin Street Map, Sources, and Ingest)

**Input**: [spec.md](./spec.md)  
**Plan**: [plan.md](./plan.md)  
**Contracts**: [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)  
**Prerequisites**: `spec.md` (Status: Draft — ready for tasks). Generated from **spec.md** first (project override: tasks before plan).

**Tests**: Spec does not require TDD. Prefer BE unit tests for Location / Place / News merge rules, article URL uniqueness, pending-without-coords, and concurrent ingest reject. Prefer FE unit tests for `MapProvider` registry (OSM vs Google stub) and category filter helpers. QA owns live Independent Tests after `be` restart.

**Organization**: Phases follow user stories (US1–US7, all P1). Task labels use `[BE]` / `[FE]` / `[QA]` per constitution; `[USn]` maps to spec user stories.

**Architecture note (map provider)**: The admin map UI MUST talk only to a **`MapProvider` port** (center/zoom, add/update/remove pins, pin click, fit bounds). P1 ships a live **OpenStreetMap** adapter (default) plus a **Google Maps stub** behind `NEXT_PUBLIC_MAP_PROVIDER`. OSM/Google types MUST NOT leak into pages, hooks, or domain models. Persistence stores **lat/lng only**. See plan.md Research Decision R1.

**GitNexus blast radius (design note — HIGH, additive only)**: Extending `DefaultPermissions`, `PermissionKeys` / `PermissionResource`, `useAdminNavItems`, `RegisterAdminRoutes`, `app.NewContainer` / `internal/app/dependency/*`, `subscribers.NewRegistry`, `queue` constants / `nats.json`, and `fe/src/config/env.ts` is **HIGH**. Keep changes **additive**. Do not rename existing keys, routes, or JSON fields. GitNexus MCP was **not available** in this design session — treat those symbols as HIGH.

**P1 only**: No public guest map, no scheduled ingest, no open-web scrape, no dynamic schema columns, no BE map-tile proxy.

## Format: `[ID] [Prefix] [P?] [Story?] Description`

- **`[BE]` / `[FE]` / `[QA]`**: Role ownership
- **`[P]`**: Parallelizable (different files, no incomplete dependency)
- **`[USn]`**: User story from spec (story phases only)
- Include exact file paths in descriptions

## Path Conventions

- Backend: `be/` — `public/routes` → `public/handlers` → `internal/services` → `internal/repository`
- Frontend: `fe/src/features/maps/` + thin `fe/src/app/admin/maps/page.tsx`
- Queue: existing `be/internal/queue` + `be/cmd/queue` (extend; do not add a second worker stack)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm migration number, env ownership, and feature module placement.

- [x] T001 [BE] Confirm next migration after `be/migrations/000009_user_is_super_admin` is `be/migrations/000010_geo_and_data_sources.{up,down}.sql`; add no Maps keys to root `.env` or `be/.env` except existing NATS/queue vars already used by `cmd/queue`
- [x] T002 [FE] [P] Confirm new feature module `fe/src/features/maps/` and thin route `fe/src/app/admin/maps/page.tsx` (PermissionGuard). Do **not** add a public guest map route
- [x] T003 [FE] [P] Document FE-owned map env in `fe/.env.example` only: `NEXT_PUBLIC_MAP_PROVIDER=osm` (allowed `osm` \| `google`) and optional `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY`. Never put map keys in root `.env` or `be/.env`
- [ ] T004 [QA] [P] Note Independent Test matrix from spec US1–US7 for later `qa-checklist.md` (menu gate, category filter, pin detail, Sources CRUD, Load data, chung cư two-place rule, RBAC + Super admin)

**Checkpoint**: Migration number `000010_geo_and_data_sources`, FE module path, and map env ownership agreed.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: RBAC seed, schema, domain layers, MapProvider port + OSM adapter + Google stub. **BLOCKS all user stories.**

**⚠️ CRITICAL**: No story UI or ingest work depends on unfinished migration, catalog keys, or the MapProvider port.

- [x] T005 [BE] Add `maps:view` / `maps:modify` to `be/internal/database/seeders/catalog.go` `DefaultPermissions()` with stable UUIDs `10000001-0000-4000-8000-000000000019` and `…000020` (see [permissions.md](./contracts/permissions.md)). Admin role receives them via existing `RolePermissionSeeder`
- [x] T006 [BE] Add migration `be/migrations/000010_geo_and_data_sources.up.sql` (+ `.down.sql`) for shared tables `countries`, `administrative_divisions`, `locations`, `places`, `categories`, `news`, `data_sources`, `data_ingest_runs`, `data_ingest_run_sources` per [database.md](./contracts/database.md). Do **not** create `map_*` tables
- [x] T006a [BE] Seed `countries` (at least VN) + Vietnam 2025 2-level `administrative_divisions` (official codes, `path`) via seeder or SQL fixture; seed `categories` by unique `key` (`room_rental`, `restaurant`, `hotel`, `eatery`, `uncategorized`) — idempotent, never a second row per key
- [x] T007 [BE] [P] Add models under reusable packages: `be/internal/models/geo/` (country, division), `location/`, `place/`, `category/`, `news/`, `datasource/` (data source + ingest run) with camelCase JSON tags matching contracts. Persist closed enums as SMALLINT; map to JSON text in DTO/handler. Include `created_by` / `updated_by`. Do not put shared entities only under a `maps/` folder that implies Maps ownership
- [x] T008 [BE] [P] Add DTOs under `be/internal/dto/maps/` for pin list item, place detail, source CRUD, field mapping, ingest run / per-source outcome
- [x] T009 [BE] Add repository interfaces in `be/internal/repository/interfaces/` + pgx/squirrel impls in `be/internal/repository/` for countries/divisions, locations (find-or-create by site identity), places (find by location+place_key, list pinnable + admin path filter), categories (get-or-reuse by unique `key`, never duplicate), news (upsert by original URL), data_sources (CRUD + list enabled), data_ingest_runs (insert, claim in-progress, update, latest)
- [x] T010 [BE] Add `be/internal/services/maps/` packages: `source` (CRUD), `place` (list/detail/status), `ingest` (merge rules + HTTP fetch of saved Sources only). Services depend on repository interfaces only
- [x] T011 [BE] Add handlers `be/public/handlers/maps_handler.go` + routes `be/public/routes/maps.go`; register from `be/public/routes/admin.go` under `/admin/maps` with `RequireView("maps")` / `RequireModify("maps")`; wire DI in `be/internal/app/dependency/` + `be/internal/app/container.go` (**additive**)
- [x] T012 [BE] Extend existing queue (do **not** add a second worker): constants in `be/internal/queue/constants.go`, stream/consumer options in `be/internal/queue/nats.json`, `PublishMapsIngest` on `be/internal/handlers/publisher/publisher.go`, subscriber `be/internal/handlers/subscribers/process_maps_ingest.go` registered in `registry.go`. When NATS is disabled, extend `be/cmd/queue/main.go` polling drain to process queued `data_ingest_runs` the same way search outbox is drained
- [x] T013 [FE] Extend `PermissionKeys` + `PermissionResource` with `maps` in `fe/src/features/access-control/permission-keys.ts`; update `permission-keys.test.ts` and mock catalog if used
- [x] T014 [FE] Add Maps nav item in `fe/src/features/access-control/hooks/use-admin-nav-items.ts`: `href: '/admin/maps'`, `permission: PermissionKeys.maps.view`, filtered by `hasPermission` (same chrome family as Users / Roles / Webhooks). Add `nav.maps` in `fe/src/locales/en/admin.json` and `fe/src/locales/vi/admin.json` (Vietnamese UTF-8 with diacritics)
- [x] T015 [FE] Add MapProvider **port** (vendor-agnostic types only) in `fe/src/features/maps/providers/map-provider.ts`: `LatLng`, `MapPin`, `MapViewport`, `MapProvider` with `mount` / `unmount`, `setCenter`, `setZoom`, `addPin`, `updatePin`, `removePin`, `clearPins`, `onPinClick`, `fitBounds`. No OSM or Google types in this file
- [x] T016 [FE] [P] Implement live **OSM adapter** in `fe/src/features/maps/providers/osm-map-provider.ts` (Leaflet or MapLibre + OSM raster tiles). Confine all OSM/Leaflet/MapLibre imports to this file
- [x] T017 [FE] [P] Add **Google Maps stub** in `fe/src/features/maps/providers/google-map-provider.ts` implementing the same port: if selected and key missing or stub-only, `mount` shows a clear “Google Maps is not configured” state and does not throw uncaught. Document how to replace the stub in file header comments
- [x] T018 [FE] Add `createMapProvider()` registry in `fe/src/features/maps/providers/create-map-provider.ts` switching on `NEXT_PUBLIC_MAP_PROVIDER`. Extend `fe/src/config/env.ts` Zod schema with `MAP_PROVIDER` (`osm` default) and optional `GOOGLE_MAPS_API_KEY` mapped from `NEXT_PUBLIC_*` only
- [x] T019 [FE] [P] Add feature types + API client under `fe/src/features/maps/` (`types.ts`, `services/maps-api.ts`) matching camelCase contracts
- [x] T020 [FE] [P] Add thin page `fe/src/app/admin/maps/page.tsx` wrapping a Maps page component with `PermissionGuard` + `PermissionKeys.maps.view` (`robots: { index: false }`)

**Checkpoint**: Migration applies; catalog keys seed; `/api/admin/maps/*` routes exist (may return empty lists); FE registry can construct OSM provider; Google stub is registered; Maps menu hidden without `maps:view`.

---

## Phase 3: User Story 1 — Operator with maps:view opens Maps from the admin menu (Priority: P1) 🎯 MVP

**Goal**: Signed-in operator with `maps:view` (or Super admin) sees Maps in admin nav, opens `/admin/maps`, and sees a public street map even when no Places are pinned.

**Independent Test**: Sign in as a user with `maps:view`. Confirm Maps appears in admin navigation, `/admin/maps` opens, and a public street map is visible (empty pins are acceptable). Guest is sent to login.

### Implementation

- [x] T021 [FE] [US1] Build `fe/src/features/maps/components/maps-page.tsx` (kebab-case) that mounts `MapProvider` via `createMapProvider()` only — never import OSM/Google types in this file
- [x] T022 [FE] [US1] Add `fe/src/features/maps/components/map-canvas.tsx` that owns the map container element, calls `mount` / `unmount`, and exposes pin/viewport operations through the port
- [x] T023 [FE] [P] [US1] Empty-pins state on the Maps surface explaining that nothing is pinned yet (i18n EN + VI) — map still loads
- [x] T024 [FE] [P] [US1] Guest hitting `/admin/maps` is redirected to login by existing `ProtectedGuard` / admin layout (verify; do not weaken guards)
- [x] T025 [BE] [P] [US1] `GET /api/admin/maps/places` returns `{ items: [] }` when no active pinnable Places (JWT + `maps:view`)

**Checkpoint**: US1 Independent Test passes (FR-001, FR-002, SC-001).

---

## Phase 4: User Story 2 — Filter and search by category shows only matching active Places (Priority: P1)

**Goal**: Operator narrows the map to P1 categories (`roomRental`, `restaurant`, `hotel`, `eatery`). Only matching **active** Places with valid coordinates remain pinned.

**Independent Test**: Seed or ingest at least two categories. Apply a single-category filter and a search. Confirm only matching active pins remain; hidden and pending stay off the map. Two Places at one Location in different categories filter independently.

### Implementation

- [x] T026 [BE] [US2] Finalize `GET /api/admin/maps/places` query: `category`, `q` (name / street / formatted), `adminDivisionId`, `countryCode`; descendant filter via division `path`; return only `status=active` with valid lat/lng per [endpoints.md](./contracts/endpoints.md). Add `GET /api/admin/maps/categories`, `GET /api/admin/maps/geo/countries` and `GET /api/admin/maps/geo/divisions`
- [x] T027 [FE] [US2] Category + admin-division filter + search in `fe/src/features/maps/components/map-filters.tsx`; TanStack Query hook `fe/src/features/maps/hooks/use-map-places.ts` (no `useEffect` server fetch)
- [x] T028 [FE] [US2] Sync query results to the map through MapProvider `clearPins` / `addPin` / `fitBounds` only
- [x] T029 [FE] [P] [US2] Empty-filter state (no matching pins) distinct from “no pins yet” and from ingest error (i18n)

**Checkpoint**: US2 Independent Test passes (FR-004–FR-006, SC-003).

---

## Phase 5: User Story 3 — Click pin shows Place detail and news with attribution (Priority: P1)

**Goal**: Selecting a pin opens Place detail (name, category, status, details, Location identity) and news with **source name** + **original article URL**.

**Independent Test**: Open a Place that has at least one News row. Confirm detail fields, headlines, source name, and a working original URL. Empty news is not an error. View-only user cannot change status or Load data.

### Implementation

- [x] T030 [BE] [US3] `GET /api/admin/maps/places/:id` returns place + location + `news` (`sourceName`, `originalUrl`, title) per [endpoints.md](./contracts/endpoints.md)
- [x] T031 [FE] [US3] Bind MapProvider `onPinClick` in map-canvas / maps-page; open detail panel `fe/src/features/maps/components/place-detail-panel.tsx`
- [x] T032 [FE] [US3] Render news with source name and external original-URL link; empty-news state; hide status-change and Load data behind `PermissionGate` + `PermissionKeys.maps.modify`
- [x] T033 [FE] [P] [US3] Hook `fe/src/features/maps/hooks/use-place-detail.ts` (TanStack Query)

**Checkpoint**: US3 Independent Test passes (FR-007, FR-008, SC-006).

---

## Phase 6: User Story 4 — Operator with maps:modify creates and edits a Source (Priority: P1)

**Goal**: `maps:modify` (or Super admin) CRUD Sources: name, enabled, HTTP method, URL, headers, query params, body, field mapping onto the **fixed** schema. View-only users see a read-only list.

**Independent Test**: Create a Source with all required fields, disable it, edit mapping, delete a throwaway Source. `maps:view`-only user can see the list but cannot change it.

### Implementation

- [x] T034 [BE] [US4] Implement Sources CRUD: `GET/POST /api/admin/maps/sources`, `GET/PATCH/DELETE /api/admin/maps/sources/:id` with validation (method, URL, mapping targets only fixed fields) in `be/internal/services/maps/` + handler
- [x] T035 [BE] [P] [US4] Field-mapping validator: allow only documented Location / Place / News / details paths; reject unknown stored-property targets (FR-011)
- [x] T036 [FE] [US4] Sources tab/panel on the same `/admin/maps` surface (`fe/src/features/maps/components/sources-panel.tsx` + `source-form-dialog.tsx` / `source-form-fields.tsx`). Do **not** add a second admin menu item
- [x] T037 [FE] [US4] Create/edit/delete + enabled toggle gated by `PermissionGate` (`maps.modify`); view-only list for `maps:view`
- [x] T038 [FE] [P] [US4] Mapping editor UI targets fixed schema fields only (location name/locationKey/countryCode/adminCode/street/postalCode/formatted/lat/lng; place name/category/placeKey/unit/lat/lng/details; news title/originalUrl). No “add column” control
- [x] T039 [FE] [P] [US4] Hooks `use-map-sources.ts` + mutations; i18n for Sources chrome

**Checkpoint**: US4 Independent Test passes (FR-009–FR-012).

---

## Phase 7: User Story 5 — Load data ingests enabled Sources and pins Places with coordinates (Priority: P1)

**Goal**: Load data fetches **only enabled saved Sources**, maps each list item, persists Location / Place / News, pins active Places with valid coordinates, and shows in-progress / success / empty / error states.

**Independent Test**: One enabled Source (structured list + coords) and one disabled Source. Run Load data. Only the enabled Source is fetched. Places with coords appear as pins. Progress then completion. A failing Source errors without wiping other pins. No enabled Sources → clear empty/error. View-only denied.

### Implementation

- [x] T040 [BE] [US5] `POST /api/admin/maps/ingest` (`maps:modify`): if another run is `queued` or `running`, return **409** with a clear in-progress message; otherwise insert `data_ingest_runs` and enqueue via existing publisher (or sync `ProcessRun` when NATS and queue worker are both unavailable — same service method)
- [x] T041 [BE] [US5] Ingest processor in `be/internal/services/maps/`: HTTP fetch **only** `enabled=true` persisted Sources (method, URL, headers, params, body); timeout and size limit; treat body as a JSON list (root array or mapped `listPath`); no open-web scrape
- [x] T042 [BE] [US5] Per-item mapping + persist using merge rules (create Location if needed, Place by place key at Location, News unique on `originalUrl`; missing/invalid coords → Place `pending` not pinned). Per-source outcomes on `data_ingest_run_sources`; one Source failure does not roll back other Sources
- [x] T043 [BE] [US5] `GET /api/admin/maps/ingest` (latest) and `GET /api/admin/maps/ingest/:id` for FE polling (`maps:view`)
- [x] T044 [BE] [P] [US5] Unit tests under `be/test/unit/maps/` for: disabled Source skipped; non-list payload fails that Source; duplicate `originalUrl` does not insert a second News row; duplicate category `key` does not insert a second Category; no-coords → pending; 409 on concurrent start
- [x] T045 [FE] [US5] Load data action (`PermissionGate` modify) + ingest status chrome `fe/src/features/maps/components/ingest-status.tsx` (in-progress vs empty vs error vs success). Poll latest run while `queued`/`running`
- [x] T046 [FE] [US5] After ingest completes, invalidate place queries so pins refresh; do not treat in-progress as a finished empty map

**Checkpoint**: US5 Independent Test passes (FR-013–FR-015, FR-019, FR-021, FR-023, SC-002, SC-005, SC-010).

---

## Phase 8: User Story 6 — One Location can host many Places; two shops in one chung cư stay distinct (Priority: P1)

**Goal**: Two businesses at the same building share a Location but remain two Places and two pins. Merge a News row onto an existing Place only when a stable place key matches **at that Location**.

**Independent Test**: Ingest two shops with the same building address/coordinates and different place keys. Confirm two Places, two pins, correct detail per pin. Re-ingest one article URL → attaches to the existing Place, no duplicate Place.

### Implementation

- [x] T047 [BE] [US6] Location reuse: match existing Location by mapped `locationKey`, else country + admin unit + street, else normalized `formatted` (optional same-site coordinate identity). **Never** use coordinates as a Place merge key
- [x] T048 [BE] [US6] Place identity: unique `(location_id, place_key)` where `place_key` is mapped external id or normalized `name+unit`. Different keys at the same Location → two Places. Missing place key at a known Location → **new** Place. Unknown site → create Location, then Place, then News
- [x] T049 [BE] [P] [US6] Unit tests: two shops same coords → two Places; same place key at Location → one Place + extra News if URL new; same article URL → no second News row (refresh mapped fields in place)
- [x] T050 [FE] [P] [US6] Overlapping pins allowed: MapProvider must keep distinct pin ids at the same lat/lng; click still opens the correct Place (no clustering in P1)

**Checkpoint**: US6 Independent Test passes (FR-016–FR-018, SC-004, SC-005).

---

## Phase 9: User Story 7 — Permissions: view sees the map; modify changes Sources and ingest (Priority: P1)

**Goal**: Standard view/modify split. Missing `maps:view` hides menu and blocks view APIs. Missing `maps:modify` blocks Source writes, Load data, and Place status. Super admin bypass after sign-in only.

**Independent Test**: (a) no Maps keys, (b) view only, (c) view+modify, (d) Super admin with empty role permissions. Menu, page, view APIs, and modify actions match the spec matrix.

### Implementation

- [x] T051 [BE] [US7] Confirm every Maps route uses `RequireView("maps")` or `RequireModify("maps")` after `Auth` (Super admin bypass already in `RequirePermission`). Unauthenticated → 401
- [x] T052 [BE] [US7] `PATCH /api/admin/maps/places/:id` for operator status `active` \| `hidden` only (`maps:modify`); reject setting `active` when coordinates are still invalid (stay `pending`)
- [x] T053 [FE] [US7] Sidebar already filtered (T014). Double-check no ungated Maps link in `use-admin-nav-items.ts` or other nav. Page `PermissionGuard` uses view. Mutate actions use `PermissionGate` modify
- [x] T054 [BE] [P] [US7] Unit or handler tests: 403 without view; 403 view-only on POST ingest / POST source / PATCH place status; Super admin allowed without catalog keys
- [x] T055 [FE] [P] [US7] Place status control in detail panel for modify users only; i18n for status labels

**Checkpoint**: US7 Independent Test passes (FR-003, FR-012, FR-020, FR-022, SC-007, SC-008).

---

## Phase 10: Polish & Cross-Cutting

**Purpose**: Tests, verify reports, QA gate. No application code from architect.

- [x] T056 [BE] [P] Additional unit tests under `be/test/unit/maps/` for mapping allow-list and “URL not saved as Source is never fetched”
- [x] T057 [FE] [P] Vitest for `createMapProvider` (osm vs google stub / missing key) and category filter helpers under `fe/src/features/maps/`
- [x] T058 [BE] Write `docs/features/008-maps/be-tasks-verify.md` after `make test-be` from repo root (restart compose `be` and `queue` before live API checks)
- [x] T059 [FE] Write `docs/features/008-maps/fe-tasks-verify.md` after `make test-fe`; browser-verify `/admin/maps` (menu, empty map, filter, pin detail, Sources, Load data) before claiming FE done
- [ ] T060 [QA] Manual Independent Tests US1–US7 + SC-001–SC-010 on `http://system.local:8080` after verify docs — see later `qa-checklist.md`

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)** → **Foundational (Phase 2)** → story phases
- **US1** needs Phase 2 catalog + MapProvider + empty places API + gated page
- **US2** needs US1 map canvas + places list query
- **US3** needs US2 pins + place detail API
- **US4** needs Phase 2 sources routes; can parallel US2/US3 once foundation exists
- **US5** needs US4 Sources + ingest service + queue/sync processor
- **US6** is ingest merge-rule hardening (depends on US5 processor)
- **US7** can start after routes exist; full matrix after US4–US5 actions exist
- **QA** after BE/FE verify

### User Story Dependencies

| Story | Depends on | Notes |
|-------|------------|--------|
| US1 | Phase 2 | MVP: gated street map + empty state |
| US2 | US1 + GET places | Category filter/search |
| US3 | US2 pins + GET place :id | Pin detail + news |
| US4 | Phase 2 sources API | Sources tab on same page |
| US5 | US4 + ingest + queue | Load data |
| US6 | US5 processor | Chung cư two-place rule |
| US7 | Routes + FE gates | RBAC matrix; Super admin reuse |

### Parallel Opportunities

```text
After Phase 2:
  @be: places list/detail + sources CRUD + ingest processor
  @fe: MapProvider OSM + Maps page (US1) in parallel with Sources form (US4)
  @fe: filters (US2) and pin detail (US3) after GET places exists
  Queue subscriber can parallel HTTP handler once IngestService.ProcessRun exists
```

### Suggested MVP

1. Phase 1–2 (RBAC, schema, MapProvider port + OSM, Google stub, empty APIs)
2. US1 gated map
3. US4 Sources CRUD (so ingest has something to run)
4. US5 Load data → pins appear
5. US2 filter + US3 detail + US6 merge tests + US7 matrix

---

## Implementation Strategy

### Recommended start order for agents

1. **`@be` first**: T001, T005–T012, T025–T026, T030, T034–T035, T040–T044, T047–T049, T051–T052, T054, T056 → `be-tasks-verify.md`
2. **`@fe` after catalog + contracts stable**: T002–T003, T013–T020 (foundation UI), then US1–US5 UI (T021–T023, T027–T029, T031–T033, T036–T039, T045–T046, T050, T053, T055, T057) → `fe-tasks-verify.md`
3. **`@qa`**: T060 after verify docs

### Notes

- `[P]` = different files, no incomplete dependency
- JSON camelCase everywhere (`placeKey`, `originalUrl`, `locationKey`, `roomRental`)
- Map vendor types stay inside `providers/osm-map-provider.ts` and `providers/google-map-provider.ts`
- Do not implement P2 guest map, scheduled ingest, or HTML scrape
- Dev BE is `go run .` — restart `be` and `queue` after BE edits before Independent Tests

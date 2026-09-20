# Implementation Plan: Maps (Admin Street Map, Sources, and Ingest)

**Feature**: `008-maps` | **Date**: 2026-09-12 | **Spec**: [spec.md](./spec.md) | **Tasks**: [tasks.md](./tasks.md)

**Contracts** (authoritative): [database](./contracts/database.md) · [endpoints](./contracts/endpoints.md) · [permissions](./contracts/permissions.md)

**Input**: Feature specification from `/docs/features/008-maps/spec.md`

## Summary

Add an admin **Maps** module at `/admin/maps`. Operators with `maps:view` see a public street map, filter P1 categories, and open Place pins with news attribution. Operators with `maps:modify` manage HTTP **Sources** on the same surface and run **Load data** to ingest structured lists into a fixed **Location → Place → News** model.

**Backend** owns persistence, merge rules, Sources CRUD, and ingest (authorized fetch of **saved enabled Sources only**). **Frontend** owns the admin map UI and a **MapProvider** port: P1 ships a live **OpenStreetMap** adapter (default) and a **Google Maps stub**. Switching providers MUST NOT rewrite Places, ingest, or admin chrome.

P2 (guest map, scheduled ingest, HTML scrape, dynamic columns, routing/booking) is deferred.

## Technical Context

**Language/Version**: Go 1.22 (BE); TypeScript strict on Next.js 15 App Router + React 19 (FE)

**Primary Dependencies**: Gin, pgx, squirrel, golang-migrate, PostgreSQL, existing NATS JetStream queue (BE); TanStack Query, shadcn/ui, react-i18next, Leaflet or MapLibre **inside the OSM adapter only** (FE)

**Storage**: PostgreSQL shared tables `countries`, `administrative_divisions`, `locations`, `places`, `categories`, `news`, `data_sources`, `data_ingest_runs`, `data_ingest_run_sources`. Chain: Location → Place → Category → News. Coordinates are **lat/lng only**. No `map_*` prefix. No `listings` table.

**Testing**: `make test-be` / `make test-fe`; BE unit tests for merge rules, URL uniqueness, pending-without-coords, concurrent ingest; FE Vitest for MapProvider registry; QA Independent Tests

**Target Platform**: Docker stack (`make up-d`) and standalone BE/FE deploys; modern browsers; platform admin users (not guests in P1)

**Project Type**: Full-stack monorepo feature (`be/` + `fe/` independent packages). HTTP only between packages

**Performance Goals**: Warm Maps open under 1 minute (SC-001). Typical Load data of hundreds of items completes while the operator waits, with in-progress UI. Pin filter/search interactive on the admin page

**Constraints**: Package independence; camelCase JSON; English docs; `maps:view` / `maps:modify` only; no BE tile proxy; no OSM/Google types in feature pages, hooks, or domain models; map keys only in `fe/.env`; HIGH RBAC/nav/queue symbols stay additive

**Scale/Scope**: One admin page; P1 categories room rental / restaurant / hotel / eatery; operator-configured Sources (not unbounded crawl)

## Constitution Check

*GATE: Must pass before design lock. Re-checked after design below.*

| Principle | Status | Notes |
|-----------|--------|--------|
| I. Package independence | **Pass** | BE owns ingest + DB; FE owns map UI + providers; HTTP only. Map env is FE-owned (`NEXT_PUBLIC_MAP_PROVIDER`, optional Google key). No `NEXT_PUBLIC_*` in `be/.env` or root `.env` |
| II. Role-owned artifacts | **Pass** | Architect owns `tasks.md` / `plan.md` / `contracts/*`. No app code in this phase |
| III. Spec before code | **Pass** | `spec.md` → `tasks.md` → `plan.md` + contracts (project override) |
| IV. API contract alignment | **Pass** | camelCase `/api/admin/maps/*`; `maps:view` / `maps:modify`; Super admin bypass reused |
| V. English documentation | **Pass** | Feature docs English; UI EN/VI via i18n (`nav.maps` VI = Bản đồ) |

**Post-design re-check**: Still pass. MapProvider is an FE port (same independence story as OAuth adapters on BE). Queue reuse extends `cmd/queue`; it is not a third app package. No constitution violations requiring Complexity Tracking exceptions.

## Complexity Tracking

GitNexus MCP **was not available**. Treat as HIGH and keep **additive**: `DefaultPermissions`, `PermissionKeys` / `PermissionResource`, `useAdminNavItems`, `RegisterAdminRoutes`, `app.NewContainer` / dependency resolvers, `subscribers.NewRegistry`, `queue` constants / `nats.json`, `fe/src/config/env.ts`.

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| — | — | Hard-coding Leaflet/OSM in `maps-page.tsx` would force a rewrite to add Google and leaks vendor types into the product UI |

---

## Research Decisions (Phase 0)

### R1 — Map integration is a provider port (hard constraint)

- **Decision**:
  1. Product UI (Maps page, pins, filters, pin detail, Sources, Load data) talks only to a **`MapProvider` port**: `mount` / `unmount`, `setCenter`, `setZoom`, `addPin`, `updatePin`, `removePin`, `clearPins`, `onPinClick`, `fitBounds`. Shared types: `LatLng`, `MapPin`, `MapViewport` — **lat/lng only**.
  2. P1 implements a live **OpenStreetMap adapter** (default). Confine Leaflet/MapLibre and OSM tile URLs to `fe/src/features/maps/providers/osm-map-provider.ts`.
  3. P1 implements a **Google Maps stub** at `google-map-provider.ts` that satisfies the same port. If `NEXT_PUBLIC_MAP_PROVIDER=google` and the key is missing or the adapter is still a stub, `mount` shows a clear “Google Maps is not configured” state (no uncaught throw).
  4. Registry `createMapProvider()` switches on FE env. Adding Google later = implement the stub + set env. **Do not** rewrite Places, ingest, Sources, or admin chrome.
  5. Persistence and APIs store/return **lat/lng only**. If OSM tiles disappear, operators set the Google provider + key; pins and filters still work.
  6. **No BE map-tile proxy.** Tiles load in the browser via the client adapter.
  7. **Address catalog is a dedicated BE API** for `countries` and `administrative_divisions`: `GET /api/address/countries` and `GET /api/address/divisions`. Maps pin search stays `GET /api/admin/maps/search` (saved Places/Locations + geocode Provider in `be/internal/services/maps/`). Place and Location services sit beside Maps at `be/internal/services/place/` and `be/internal/services/location/`. Env for geocode lives in `be/.env`: `MAP_GEOCODE_PROVIDER`, optional `MAP_NOMINATIM_URL` / `MAP_NOMINATIM_USER_AGENT`. The browser must not call Nominatim.
- **Rationale**: Spec stays basemap-agnostic (“public street map”). The user locked a strategy/port so Google is a config switch, not a redesign.
- **Alternatives considered**:
  - **Hard-wire react-leaflet in the page** — fastest P1; rejected because Google would require rewriting the feature module.
  - **BE tile proxy** — hides OSM from the browser; rejected (ops cost, couples BE to a vendor, breaks FE-only deploy of the map UI).
  - **Full Google adapter in P1** — allowed if small, but prefer stub + documented registration so `@fe` can finish Google without a new plan.

**How to add Google later (implementer recipe)**

1. Replace the stub body in `fe/src/features/maps/providers/google-map-provider.ts` with a real `MapProvider` (load Maps JavaScript API using `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY`; map `LatLng` / `MapPin` to Google objects **only inside this file**).
2. Keep registration in `create-map-provider.ts` (`case 'google'`).
3. Set in **`fe/.env` / `fe/.env.example` only**:
   - `NEXT_PUBLIC_MAP_PROVIDER=google`
   - `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=<key>`
4. Extend `fe/src/config/env.ts` only if new optional FE vars are needed. Never add map keys to root `.env` or `be/.env`.
5. No BE migration, no ingest change, no contract change.

### R2 — Admin module placement and RBAC

- **Decision**: `/admin/maps` (not Tools catalog). Keys `maps:view` / `maps:modify`. Seed UUIDs `…019` / `…020`. Sidebar via `useAdminNavItems` + `hasPermission`. Super admin bypass reused. Sources live on the same page (tab/panel), not a second menu item.
- **Rationale**: Spec Q1/Q4; matches Users / Roles / Webhooks chrome family.

### R3 — Location / Place / News merge (do not contradict spec)

- **Decision**: Shared tables `locations`, `places`, `categories`, `news`. Chain **Location → Place → Category → News**. **Never** merge Places on coordinates or Location alone. Place unique `(location_id, place_key)`. `place_key` = mapped external id, else `lower(trim(name)) + '|' + lower(trim(unit))`. Category unique on `key` (reuse row; never a second danh mục). News unique on `original_url` (refresh in place; do not insert a second row). News.category_id = Place.category_id at ingest. Location reuse: `locationKey`, else country + admin unit + street, else normalized `formatted`, else rounded coords **plus** same `formatted`. Missing Location identity → skip item. Missing coords → Place `pending`. New Place with valid coords → `active`. Hidden is not auto-cleared on re-ingest.
- **Rationale**: Spec Q3 / FR-016–FR-019 / chung cư clarification.

### R9 — Integer enums and audit columns (user confirm 2026-09-12)

- **Decision**: Closed enums persist as `SMALLINT` (1, 2, 3, …). Application maps to JSON text (`pending`/`active`/`hidden`, HTTP methods, ingest statuses). See [database.md](./contracts/database.md) Integer enums table. Every feature table has `created_by` / `updated_by` (UUID, NULL on seed) plus existing timestamps.
- **Rationale**: User locked int storage + audit actors. VARCHAR enum labels are not stored.
- **Alternatives rejected**: PostgreSQL native ENUM types (harder to port / migrate). VARCHAR status (rejected by product).

### R8 — Shared table names and address tree (user confirm 2026-09-12)

- **Decision**:
  1. Do **not** use `map_locations` / `map_sources` / other `map_*` entity tables. Names: `countries`, `administrative_divisions`, `locations`, `places`, `categories`, `news`, `data_sources`, `data_ingest_runs`, `data_ingest_run_sources`. Do **not** name the article table `listings`.
  2. Address is a tree + Location fields (country, leaf division, street, postal, formatted, path on the division). `level` is ordinal per country; `type` is semantic. Vietnam P1 seed is 2 official levels (2025 reform).
  3. Maps P1 APIs stay under `/api/admin/maps/*`; they are the first writer. Other features may FK the same tables later.
- **Rationale**: User locked reuse + Option A address. Hard-coding `map_` would force a rename when Profile or another module needs a Location.
- **Alternatives rejected**: Keep `map_*` prefixes (not reusable). Fixed `province`/`district`/`ward` columns (breaks VN 2-level and other countries).

### R4 — Ingest uses existing queue; sync fallback when NATS is off

- **Decision**:
  1. `POST /api/admin/maps/ingest` inserts `data_ingest_runs` (`queued`), then publishes `system.maps.ingest` with `{ "id": "<runId>" }` via existing `handlers/publisher`.
  2. New subscriber `process_maps_ingest` calls `IngestService.ProcessRun` (same as search outbox → `ProcessByID`).
  3. Extend `be/internal/queue/constants.go` + `nats.json` stream `maps` / consumer `maps_ingest`. Register in `subscribers.NewRegistry`.
  4. When NATS is disabled, extend `cmd/queue` **polling** to drain queued map ingest runs (same process as search outbox). Do **not** add a second worker binary.
  5. If the API process cannot publish and no worker is running (local BE-only), the handler MAY call `ProcessRun` **synchronously** and return `200` with the finished run. FE still understands `202` + poll `GET /api/admin/maps/ingest`.
  6. Concurrent start: if any run is `queued` or `running`, return **409**.
  7. Fetch **only** `enabled` persisted Sources. Timeout ~30s/Source; body cap ~5 MiB. HTML / non-list → that Source fails; other Sources continue.
- **Rationale**: Spec Load data can be slow; repo already has NATS + `cmd/queue`. Webhooks capture is inbound HTTP, not a model for outbound fetch. A second worker stack would violate BE layout.
- **Alternatives considered**: HTTP-only sync always — simpler but blocks the API on multi-source loads and loses in-progress polling when the request is still open. Dedicated ingest microservice — rejected.

### R5 — Field mapping is configuration, not schema

- **Decision**: `field_mapping` JSONB with `listPath`, `location.*`, `place.*`, `news.*`, `details` allow-list (`description`, `phone`, `hours`, `website`, `priceRange`), optional `categoryMap`. Reject unknown persisted property names. No per-source columns.
- **Rationale**: FR-011; P2 owns dynamic storage.

### R6 — FE module and env

- **Decision**: New `fe/src/features/maps/` (kebab-case). Thin `fe/src/app/admin/maps/page.tsx` + `PermissionGuard`. Env in `fe/.env.example`:
  - `NEXT_PUBLIC_MAP_PROVIDER=osm` (or `google`)
  - `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY=` (optional)
  Zod in `fe/src/config/env.ts` (`MAP_PROVIDER` default `osm`). Feature code reads `env`, not `process.env`.
- **Rationale**: App independence + environment rules.

### R7 — GitNexus / blast radius

GitNexus MCP was unavailable. Implementers **must** run impact before editing shared symbols:

| Area | Likely symbols / files | Risk note |
|------|------------------------|-----------|
| Permission catalog | `DefaultPermissions` | Additive rows `…019` / `…020` |
| Admin routes | `RegisterAdminRoutes` | New `/admin/maps` group |
| DI | `NewContainer`, dependency resolvers | New maps + ingest wiring |
| Queue | constants, `nats.json`, `NewRegistry`, `Publisher`, `cmd/queue` | Additive stream + poll drain |
| FE RBAC / nav | `PermissionKeys`, `useAdminNavItems` | Additive `maps` + one nav item |
| FE env | `env.ts` | Additive optional map vars |

Do not rename existing permission keys or nav hrefs.

---

## Data Model (Phase 1 design)

> Authoritative schema: [contracts/database.md](./contracts/database.md). Summary below.

| Entity | Table | Role |
|--------|-------|------|
| Country | `countries` | ISO catalog; P1 seeds VN |
| Admin division | `administrative_divisions` | Tree: parent, level, type, official code, `path` |
| Location | `locations` | Shared site: structured address + lat/lng + location_key |
| Place | `places` | Shared POI / pin: place_key, `category_id`, status, lat/lng, details |
| Category | `categories` | Shared danh mục; unique `key`; P1 seed 5 rows |
| News | `news` | Article: `place_id` + `category_id`; unique `original_url` |
| Data Source | `data_sources` | HTTP method/URL/headers/params/body + field_mapping |
| Ingest run | `data_ingest_runs` | Job status + counters |
| Ingest source | `data_ingest_run_sources` | Per–Data Source outcome |

```text
Country (1) ---- (*) AdministrativeDivision (self parent)
Country / Division ---- (*) Location (1) ---- (*) Place (1) ---- (*) News
Category (1) ---- (*) Place
Category (1) ---- (*) News
DataSource (1) ---- (*) News          (SET NULL on source delete)
User (1) ---- (*) DataIngestRun (1) ---- (*) DataIngestRunSource
```

### State transitions

```text
Load data --> IngestRun(queued) --> running --> completed | failed
Enabled Source --> fetch list --> map items
Item + site unknown --> skip (items_error)
Item + new site --> Location then Place then News
Item + known Location + new place_key --> new Place
Item + known (Location, place_key) --> attach News if originalUrl new
Same originalUrl --> update News in place (no duplicate)
Place coords invalid --> pending (not pinned)
Place coords valid (new) --> active (pinned)
Operator --> hidden | active (active requires valid coords)
```

---

## Contracts

> **Authoritative sources:** [contracts/database.md](./contracts/database.md), [contracts/endpoints.md](./contracts/endpoints.md), [contracts/permissions.md](./contracts/permissions.md).  
> Summary below is a convenience mirror; prefer the contract files when they differ.

Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. All JSON **camelCase**.  
**Permissions:** `maps:view` / `maps:modify` — see [permissions.md](./contracts/permissions.md).

| Method | Path | Auth |
|--------|------|------|
| GET | `/api/admin/maps/categories` | JWT + `maps:view` |
| GET | `/api/admin/maps/places` | JWT + `maps:view` |
| GET | `/api/admin/maps/places/:id` | JWT + `maps:view` |
| PATCH | `/api/admin/maps/places/:id` | JWT + `maps:modify` |
| GET | `/api/address/countries` | JWT |
| GET | `/api/address/divisions` | JWT |
| GET | `/api/admin/maps/sources` | JWT + `maps:view` |
| GET | `/api/admin/maps/sources/:id` | JWT + `maps:view` |
| POST | `/api/admin/maps/sources` | JWT + `maps:modify` |
| PATCH | `/api/admin/maps/sources/:id` | JWT + `maps:modify` |
| DELETE | `/api/admin/maps/sources/:id` | JWT + `maps:modify` |
| POST | `/api/admin/maps/ingest` | JWT + `maps:modify` |
| GET | `/api/admin/maps/ingest` | JWT + `maps:view` |
| GET | `/api/admin/maps/ingest/:id` | JWT + `maps:view` |

UI: `/admin/maps` (`PermissionGuard` view). Guest → login.

---

## Project Structure

### Documentation (this feature)

```text
docs/features/008-maps/
├── spec.md
├── tasks.md
├── plan.md                 ← this file
├── contracts/
│   ├── database.md
│   ├── endpoints.md
│   └── permissions.md
├── checklists/requirements.md
├── be-tasks-verify.md      ← @be (not this phase)
├── fe-tasks-verify.md      ← @fe (not this phase)
└── qa-checklist.md         ← @qa (not this phase)
```

### Backend (new domain + additive kernel)

```text
be/migrations/000010_geo_and_data_sources.{up,down}.sql
be/internal/models/geo/
be/internal/models/location/
be/internal/models/place/
be/internal/models/category/
be/internal/models/news/
be/internal/models/datasource/
be/internal/dto/maps/
be/internal/repository/interfaces/… (geo, location, place, category, news, datasource)
be/internal/repository/…
be/internal/services/location/         # Location CRUD
be/internal/services/place/            # Place CRUD / pin list
be/internal/services/maps/             # Search + geocode port/providers, Sources, ingest
be/public/handlers/maps_handler.go
be/public/routes/maps.go
be/public/routes/admin.go              # Register maps group (additive)
be/internal/app/container.go
be/internal/app/dependency/
be/internal/database/seeders/catalog.go
be/internal/queue/constants.go
be/internal/queue/nats.json
be/internal/handlers/publisher/publisher.go
be/internal/handlers/subscribers/process_maps_ingest.go
be/internal/handlers/subscribers/registry.go
be/cmd/queue/main.go                   # poll drain for queued ingest runs
be/test/unit/maps/
```

### Frontend (new feature module + additive kernel)

```text
fe/.env.example                        # NEXT_PUBLIC_MAP_PROVIDER, NEXT_PUBLIC_GOOGLE_MAPS_API_KEY
fe/src/config/env.ts
fe/src/app/admin/maps/page.tsx
fe/src/features/access-control/permission-keys.ts
fe/src/features/access-control/hooks/use-admin-nav-items.ts
fe/src/locales/{en,vi}/admin.json
fe/src/features/maps/
  index.ts
  types.ts
  services/maps-api.ts
  hooks/use-map-places.ts
  hooks/use-place-detail.ts
  hooks/use-map-sources.ts
  hooks/use-map-ingest.ts
  providers/map-provider.ts            # port + LatLng / MapPin
  providers/osm-map-provider.ts        # live OSM
  providers/google-map-provider.ts     # stub
  providers/create-map-provider.ts
  components/maps-page.tsx
  components/map-canvas.tsx
  components/map-filters.tsx
  components/place-detail-panel.tsx
  components/sources-panel.tsx
  components/source-form-dialog.tsx
  components/source-form-fields.tsx
  components/ingest-status.tsx
```

**Structure Decision**: Full-stack admin feature. BE owns ingest + persistence. FE owns map UI + provider adapters. Existing NATS/`cmd/queue` is extended. No FSD. No Kratos split. No BE knowledge of OSM or Google.

---

## Env & deploy notes

| Var | Owner file | Purpose |
|-----|------------|---------|
| `NEXT_PUBLIC_MAP_PROVIDER` | `fe/.env` | `osm` (default) or `google` |
| `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY` | `fe/.env` | Required only when provider is `google` and a real adapter is used |
| `NEXT_PUBLIC_API_BASE_URL` | `fe/.env` | Existing Maps API calls |
| `MAP_GEOCODE_PROVIDER` | `be/.env` | `osm` (default) or `google` (stub). Admin search only |
| `MAP_NOMINATIM_URL` | `be/.env` | Nominatim base URL when provider is `osm` |
| `MAP_NOMINATIM_USER_AGENT` | `be/.env` | Required identifying User-Agent for Nominatim |
| `NATS_ENABLED` / `NATS_URL` | `be/.env` | Existing queue; ingest publisher/subscriber |
| `NGINX_HTTP_PORT` | root `.env` | Existing; **no** map keys here |

**Standalone BE**: copy `be/` + `be/.env`. Ingest and APIs work without FE. Run `cmd/queue` (or rely on sync `ProcessRun` fallback). No map SDK on the server.

**Standalone FE**: copy `fe/` + `fe/.env`. Point API at BE. Choose map provider locally. If OSM tiles are blocked, set Google provider + key after the adapter is implemented.

---

## Quickstart validation (for implementers / QA)

### Prerequisites

- `make env` / configured `be/.env`, `fe/.env`, root `.env`
- `make up-d` (restart `be` and `queue` after BE edits)
- User with `maps:view` / `maps:modify` (seeded admin is Super admin and bypasses)

### Smoke scenarios

1. Guest `/admin/maps` → login.
2. User without `maps:view` (and not Super admin) → no Maps menu; GET places → 403.
3. `maps:view` → menu + street map + empty pins state; Sources read-only; Load data hidden/denied.
4. `maps:modify` → create Source (structured JSON list) → Load data → in-progress then pins.
5. Filter one P1 category → other categories hidden; two Places at one building stay two pins.
6. Pin click → source name + original URL.
7. Re-run Load data → no duplicate News row for the same URL.
8. Missing coords → Place `pending`, not pinned.
9. Set `NEXT_PUBLIC_MAP_PROVIDER=google` without a real adapter/key → configured-not-ready UI; Places API unchanged.

### Commands

```bash
make test-be
make test-fe
```

---

## Implementation order

1. **`@be`**: catalog + migration `000010_geo_and_data_sources` + VN division seed → models/DTOs/repos/services → admin routes → ingest + queue subscriber → merge-rule tests → `be-tasks-verify.md`
2. **`@fe`**: PermissionKeys + sidebar + env + MapProvider port + OSM + Google stub → Maps page → filters/detail/Sources/ingest UI → `fe-tasks-verify.md`
3. **`@qa`**: checklist + Independent Tests on `system.local`

P1 done when US1–US7 acceptance criteria are met. Do not implement P2.

## Risks

- HIGH: additive RBAC/nav/queue edits must not rename existing keys or break Super admin bypass.
- HIGH: Place merge-on-coords would violate the chung cư rule — covered by unique `(location_id, place_key)` and unit tests.
- MEDIUM: OSM tile availability is an ops concern; provider switch is the mitigation (R1).
- MEDIUM: SSRF-style fetches of operator-saved URLs are intentional (admin Sources). Do not fetch unsaved URLs.
- Dev BE is `go run .` — **restart `be` and `queue` after BE edits** before Independent Tests.

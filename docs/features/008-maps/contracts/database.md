# Database contract: Maps (shared geo, places, data sources)

**Feature**: `008-maps`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative persistence contract. Implementers follow this file; `plan.md` links here.  
> Coordinates are vendor-agnostic **lat/lng**. No OpenStreetMap, Google, tile, or map-provider columns.  
> **Table names are shared platform tables**, not Maps-only. Do **not** prefix them with `map_`. Maps is the first consumer; later features may foreign-key the same rows.

## Scope

- **New tables**: `countries`, `administrative_divisions`, `locations`, `places`, `categories`, `news`, `data_sources`, `data_ingest_runs`, `data_ingest_run_sources`
- **Altered tables**: none
- **Renamed vs first draft**: `map_locations` → `locations`; `map_places` → `places`; `map_listings` → `news` (not `listings` — that name collides with classified ads); `map_sources` → `data_sources`; `map_ingest_*` → `data_ingest_*`
- **N/A**: not applicable — feature owns PostgreSQL persistence

## Shared vs Maps-owned

| Table | Reuse |
|-------|--------|
| `countries`, `administrative_divisions` | Geo catalog for any feature (profile, shipping, Maps filters) |
| `locations` | Any site (building, campus, warehouse) — not Maps-only |
| `places` | Any point-of-interest at a Location |
| `categories` | Shared danh mục catalog. Unique `key` — ingest never creates a second row for the same category |
| `news` | Sourced article about a Place, under that Place’s Category. Unique `original_url` |
| `data_sources`, `data_ingest_runs`, `data_ingest_run_sources` | Generic HTTP ingest configuration and run history |

P1 Maps APIs remain under `/api/admin/maps/*` and `maps:view` / `maps:modify`. Other features must not assume they own these tables; they attach new FKs later.

## Integer enums (DB) → text (API)

Closed enums are stored as **`SMALLINT`**. JSON and UI map to camelCase / HTTP text in the application layer. Do **not** persist enum labels as VARCHAR.

| Column | 1 | 2 | 3 | 4 | 5 |
|--------|---|---|---|---|---|
| `places.status` | `pending` | `active` | `hidden` | — | — |
| `data_sources.http_method` | `GET` | `POST` | `PUT` | `PATCH` | — |
| `data_ingest_runs.status` | `queued` | `running` | `completed` | `failed` | — |
| `data_ingest_run_sources.status` | `pending` | `running` | `completed` | `failed` | `skipped` |

CHECK constraints use the integer set only. Open vocabularies stay text (`administrative_divisions.type`, ISO country codes, `categories.key`).

## Audit columns

Every table in this contract includes:

| Column | Type | Null | Notes |
|--------|------|------|-------|
| `created_by` | UUID | YES | FK → `users(id)` ON DELETE SET NULL. NULL on catalog seed |
| `updated_by` | UUID | YES | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` on insert |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` on insert/update |

Rules: ingest sets both to the user who started Load data; Source / Place-status writes set `updated_by` (and `created_by` on insert) to the operator. Seed rows for `countries`, `administrative_divisions`, and `categories` MAY leave both NULL.

## Tables

### `countries`

ISO country catalog.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `code` | VARCHAR(2) | NO | | ISO 3166-1 alpha-2 (`VN`, `US`). UNIQUE |
| `code3` | VARCHAR(3) | NO | | ISO 3166-1 alpha-3 (`VNM`, `USA`). UNIQUE |
| `name` | TEXT | NO | | English display (`Vietnam`) |
| `name_local` | TEXT | NO | `''` | Local script (`Việt Nam`) |
| `phone_code` | TEXT | NO | `''` | Optional (`+84`). Unused by Maps P1 |
| `currency` | VARCHAR(3) | NO | `''` | Optional ISO 4217. Unused by Maps P1 |
| `is_active` | BOOLEAN | NO | `TRUE` | |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `countries_code_key` on `(code)`
- UNIQUE: `countries_code3_key` on `(code3)`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

P1 seeds at least **VN**. Other countries may be added without a column migration.

### `administrative_divisions`

Variable-depth admin tree **per country**. `level` is an ordinal inside that country (1 = first subdivision below country). `type` is the semantic label. **Do not** treat level 2 as “district” globally — Vietnam (from 1 Jul 2025) is 2 levels (province → commune/ward/special zone); other countries may be 2–5.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `country_id` | UUID | NO | | FK → `countries(id)` ON DELETE RESTRICT |
| `parent_id` | UUID | YES | | FK → `administrative_divisions(id)` ON DELETE RESTRICT. NULL = top subdivision |
| `level` | SMALLINT | NO | | Ordinal in this country. Not a global meaning |
| `code` | TEXT | NO | | Official admin code (VN province 2 digits, commune 5 digits) |
| `name` | TEXT | NO | | Local / official short name |
| `name_en` | TEXT | NO | `''` | |
| `full_name` | TEXT | NO | `''` | Display (`Thành phố Đà Nẵng`) |
| `type` | TEXT | NO | | Neutral token: `province`, `municipality`, `commune`, `ward`, `special_zone`, `state`, `county`, `city`, `district` (legacy / other countries). Not VN-only slugs as the only system |
| `path` | TEXT | NO | | Materialized path `{country_code}/{code}/…/` e.g. `VN/48/00123/`. Used for descendant filter |
| `lat` | DOUBLE PRECISION | YES | | Optional centroid (fit-bounds / roll-up) |
| `lng` | DOUBLE PRECISION | YES | | Optional centroid |
| `is_active` | BOOLEAN | NO | `TRUE` | `FALSE` for abolished units (e.g. former VN districts) |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `administrative_divisions_country_code_key` on `(country_id, code)`
- FK: `country_id` → `countries(id)` ON DELETE RESTRICT
- FK: `parent_id` → `administrative_divisions(id)` ON DELETE RESTRICT
- CHECK: `level >= 1`
- INDEX: `idx_admin_div_country_parent` on `(country_id, parent_id)`
- INDEX: `idx_admin_div_country_level` on `(country_id, level)`
- INDEX: `idx_admin_div_path` on `(path)`
- INDEX: `idx_admin_div_path_prefix` on `(path text_pattern_ops)` for `LIKE 'VN/48/%'`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

If a future country reuses the same `code` under different parents, change uniqueness to `(country_id, parent_id, code)` in a later migration. VN official codes are unique per country.

**Filter / calculate**

- Filter Places in a unit: Location.admin_division_id’s `path` equals that unit’s path **or** is a descendant (`path LIKE unit.path || '%'`).
- Roll-up counts: `GROUP BY` ancestor via `path`.
- Fit map to a unit: use division `lat`/`lng` when present; else derive from child Place coords in P1 if needed.

### `locations`

One geographic **site** (chung cư, campus, depot). Reusable beyond Maps. Many Places belong to one Location. Structured address lives here — not on Place (Place keeps `unit`).

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `name` | TEXT | NO | | Site display name |
| `location_key` | TEXT | YES | | Optional stable site identity from mapping. UNIQUE when not null |
| `country_id` | UUID | YES | | FK → `countries(id)` ON DELETE RESTRICT |
| `admin_division_id` | UUID | YES | | FK → `administrative_divisions(id)` ON DELETE RESTRICT. Prefer **leaf** (commune/ward). NULL if ingest cannot resolve |
| `street` | TEXT | NO | `''` | Thoroughfare / house / building line |
| `postal_code` | TEXT | NO | `''` | |
| `formatted` | TEXT | NO | `''` | Single-line display (was draft `address`) |
| `lat` | DOUBLE PRECISION | YES | | WGS84; null if unknown |
| `lng` | DOUBLE PRECISION | YES | | WGS84; null if unknown |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `locations_location_key_key` on `(location_key)` WHERE `location_key IS NOT NULL`
- FK: `country_id` → `countries(id)` ON DELETE RESTRICT
- FK: `admin_division_id` → `administrative_divisions(id)` ON DELETE RESTRICT
- INDEX: `idx_locations_formatted_norm` on `(lower(btrim(formatted)))` WHERE `formatted <> ''`
- INDEX: `idx_locations_street_norm` on `(lower(btrim(street)))` WHERE `street <> ''`
- INDEX: `idx_locations_coords` on `(lat, lng)` WHERE `lat IS NOT NULL AND lng IS NOT NULL`
- INDEX: `idx_locations_country_id` on `(country_id)`
- INDEX: `idx_locations_admin_division_id` on `(admin_division_id)`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

**Location reuse (application, not a unique constraint on coords)**

1. If mapping supplies `locationKey` and a row exists with that key → reuse.
2. Else if `country_id` + `admin_division_id` + normalized `street` all set and a row matches → reuse.
3. Else if mapped `formatted` is non-empty and a row exists with the same normalized `formatted` → reuse.
4. Else if both mapped lat/lng are valid and a row exists with the same rounded coords (6 decimal places) **and** the same normalized `formatted` (including both empty) → reuse.
5. Else create a new Location.

Coordinates alone MUST NOT collapse two named sites with different streets/formatted lines. Coordinates MUST NEVER be a Place merge key.

Unresolved admin (unknown code/name) → Location may still be created with `admin_division_id` NULL; pin filter by admin excludes it.

### `categories`

Shared danh mục catalog. Chain: **Location → Place → Category → News**.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `key` | TEXT | NO | | Stable token: `room_rental`, `restaurant`, `hotel`, `eatery`, `uncategorized`. UNIQUE |
| `name` | TEXT | NO | | English display |
| `name_local` | TEXT | NO | `''` | Vietnamese or other local label |
| `sort_order` | INTEGER | NO | `0` | Filter UI order |
| `is_active` | BOOLEAN | NO | `TRUE` | Inactive keys stay on existing Places; omitted from P1 filter chips |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `categories_key_key` on `(key)` — **same danh mục never inserts a second row**
- INDEX: `idx_categories_active_sort` on `(is_active, sort_order, name)`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

**P1 seed** (idempotent upsert by `key`): `room_rental`, `restaurant`, `hotel`, `eatery`, `uncategorized`. JSON exposes camelCase keys (`roomRental`). Ingest `categoryMap` resolves to `categories.key` then `categories.id`. Unknown mapped values → `uncategorized` (reuse that row).

### `places`

One business or unit at a Location. Map pins are per Place. Reusable POI table. Each Place has **one** Category.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `location_id` | UUID | NO | | FK → `locations(id)` ON DELETE RESTRICT |
| `category_id` | UUID | NO | | FK → `categories(id)` ON DELETE RESTRICT |
| `place_key` | TEXT | NO | | Stable key **at this Location**: mapped external id, else normalized `name` + optional `unit` |
| `name` | TEXT | NO | | Display name |
| `unit` | TEXT | NO | `''` | Unit / shop number when mapped |
| `status` | SMALLINT | NO | `1` | Int enum: `1` pending, `2` active, `3` hidden. JSON text after map |
| `lat` | DOUBLE PRECISION | YES | | Pin latitude; may copy Location coords when the item has none |
| `lng` | DOUBLE PRECISION | YES | | Pin longitude |
| `details` | JSONB | NO | `'{}'::jsonb` | Mapped extras from the **fixed** allow-list only |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `places_location_place_key_key` on `(location_id, place_key)`
- FK: `location_id` → `locations(id)` ON DELETE RESTRICT
- FK: `category_id` → `categories(id)` ON DELETE RESTRICT
- CHECK: `status IN (1, 2, 3)`
- INDEX: `idx_places_pinnable` on `(category_id, status)` WHERE `status = 2 AND lat IS NOT NULL AND lng IS NOT NULL`
- INDEX: `idx_places_location_id` on `(location_id)`
- INDEX: `idx_places_category_id` on `(category_id)`
- INDEX: `idx_places_name` on `(lower(name))`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

**Status rules**

- Ingest with **invalid or missing** lat/lng (after inheriting Location coords if needed) → status `1` (`pending`). Not pinned.
- Ingest with **valid** coords on a **new** Place → status `2` (`active`).
- Subsequent ingest that adds valid coords to a pending Place → may become `2`.
- Operator hidden (`3`) is preserved on re-ingest (do not auto-unhide).
- Operator may set `2` or `3` via API (JSON `active` / `hidden`); `1` is system-owned. Setting `active` without valid coords is rejected.

**Valid coordinates**: both `lat` and `lng` non-null; `lat` in `[-90, 90]`; `lng` in `[-180, 180]`.

### `news`

One sourced article about a Place, under that Place’s Category. Do **not** name this `listings`.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `place_id` | UUID | NO | | FK → `places(id)` ON DELETE CASCADE |
| `category_id` | UUID | NO | | FK → `categories(id)` ON DELETE RESTRICT. P1: same as `places.category_id` at ingest |
| `source_id` | UUID | YES | | FK → `data_sources(id)` ON DELETE SET NULL |
| `source_name` | TEXT | NO | | Denormalized Data Source name at ingest time |
| `original_url` | TEXT | NO | | Canonical article URL; **globally unique** |
| `title` | TEXT | NO | `''` | |
| `details` | JSONB | NO | `'{}'::jsonb` | Optional mapped extras from the news allow-list |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- UNIQUE: `news_original_url_key` on `(original_url)` — **same URL never inserts a second News row**
- FK: `place_id` → `places(id)` ON DELETE CASCADE
- FK: `category_id` → `categories(id)` ON DELETE RESTRICT
- FK: `source_id` → `data_sources(id)` ON DELETE SET NULL
- INDEX: `idx_news_place_created` on `(place_id, created_at DESC)`
- INDEX: `idx_news_category_id` on `(category_id)`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

**Merge**: same `original_url` (trimmed) → update mapped fields in place (`title`, `details`, `source_*`); never insert a second row. Do not re-parent the News row to a different Place or Category on refresh.

### `data_sources`

Operator-managed HTTP ingest endpoint (not Maps-specific). Only saved rows may be fetched.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `name` | TEXT | NO | | Display name |
| `enabled` | BOOLEAN | NO | `TRUE` | Load data skips `FALSE` |
| `http_method` | SMALLINT | NO | `1` | Int enum: `1` GET, `2` POST, `3` PUT, `4` PATCH. JSON text after map |
| `url` | TEXT | NO | | Absolute http(s) URL |
| `headers` | JSONB | NO | `'{}'::jsonb` | Object of header name → string |
| `query_params` | JSONB | NO | `'{}'::jsonb` | Object of query name → string |
| `body` | TEXT | NO | `''` | Request body for POST/PUT/PATCH; ignored for GET |
| `field_mapping` | JSONB | NO | | See Field mapping below |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- CHECK: `http_method IN (1, 2, 3, 4)`
- CHECK: `url ~* '^https?://'`
- INDEX: `idx_data_sources_enabled` on `(enabled)`
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL

P1 does **not** vault header secrets. Treat Data Sources as admin-only configuration. P1 Maps UI is the first editor.

### `data_ingest_runs`

One operator-triggered Load data execution.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `status` | SMALLINT | NO | `1` | Int enum: `1` queued, `2` running, `3` completed, `4` failed |
| `triggered_by` | UUID | NO | | Who started the run; also copied to `created_by` on insert |
| `error_message` | TEXT | NO | `''` | |
| `sources_total` | INTEGER | NO | `0` | Enabled Data Sources counted at start |
| `success` | INTEGER | NO | `0` | Sources that completed |
| `error` | INTEGER | NO | `0` | Sources that failed |
| `items_success` | INTEGER | NO | `0` | Mapped items that persisted |
| `items_error` | INTEGER | NO | `0` | Items skipped or mapping-failed |
| `started_at` | TIMESTAMPTZ | YES | | |
| `finished_at` | TIMESTAMPTZ | YES | | |
| `created_by` | UUID | YES | | Same as `triggered_by` on insert |
| `updated_by` | UUID | YES | | Worker or operator |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- CHECK: `status IN (1, 2, 3, 4)`
- FK: `triggered_by` → `users(id)` ON DELETE RESTRICT
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL
- INDEX: `idx_data_ingest_runs_created` on `(created_at DESC)`
- Partial INDEX: `idx_data_ingest_runs_active` on `(status)` WHERE `status IN (1, 2)`

A run is `failed` when **every** enabled Data Source failed, or there were no enabled Data Sources. Mixed outcomes → `completed`.

### `data_ingest_run_sources`

Per–Data Source outcome for one ingest run.

| Column | Type | Null | Default | Notes |
|--------|------|------|---------|-------|
| `id` | UUID | NO | `gen_random_uuid()` | PK |
| `run_id` | UUID | NO | | FK → `data_ingest_runs(id)` ON DELETE CASCADE |
| `source_id` | UUID | YES | | FK → `data_sources(id)` ON DELETE SET NULL |
| `source_name` | TEXT | NO | | Snapshot |
| `status` | SMALLINT | NO | `1` | Int enum: `1` pending, `2` running, `3` completed, `4` failed, `5` skipped |
| `error_message` | TEXT | NO | `''` | |
| `items_success` | INTEGER | NO | `0` | |
| `items_error` | INTEGER | NO | `0` | |
| `created_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `updated_by` | UUID | YES | | FK → `users(id)` ON DELETE SET NULL |
| `created_at` | TIMESTAMPTZ | NO | `NOW()` | |
| `updated_at` | TIMESTAMPTZ | NO | `NOW()` | |

**Constraints / indexes**

- PK: `id`
- CHECK: `status IN (1, 2, 3, 4, 5)`
- FK: `run_id` → `data_ingest_runs(id)` ON DELETE CASCADE
- FK: `source_id` → `data_sources(id)` ON DELETE SET NULL
- FK: `created_by` → `users(id)` ON DELETE SET NULL
- FK: `updated_by` → `users(id)` ON DELETE SET NULL
- INDEX: `idx_data_ingest_run_sources_run` on `(run_id)`

Disabled Data Sources are **not** inserted. `skipped` is reserved for a documented non-fetch of an enabled source (rare in P1).

## Field mapping (`data_sources.field_mapping`)

JSON object. Mapping copies fields from each list item onto the **fixed** schema. It MUST NOT create new tables or columns.

```json
{
  "listPath": "data.items",
  "location": {
    "name": "building.name",
    "locationKey": "building.id",
    "countryCode": "building.country",
    "adminCode": "building.wardCode",
    "street": "building.street",
    "postalCode": "building.postal",
    "formatted": "building.address",
    "lat": "building.lat",
    "lng": "building.lng"
  },
  "place": {
    "name": "shop.name",
    "category": "shop.type",
    "placeKey": "shop.id",
    "unit": "shop.unit",
    "lat": "shop.lat",
    "lng": "shop.lng"
  },
  "news": {
    "title": "article.title",
    "originalUrl": "article.url"
  },
  "details": {
    "description": "shop.desc",
    "phone": "shop.phone",
    "hours": "shop.hours",
    "website": "shop.website",
    "priceRange": "shop.price"
  },
  "categoryMap": {
    "nha-hang": "restaurant",
    "khach-san": "hotel"
  }
}
```

| Key | Required | Meaning |
|-----|----------|---------|
| `listPath` | no | Dot path from JSON root to the array. Empty / omitted = root is the array |
| `location.countryCode` | no | ISO 3166-1 alpha-2. Resolve `countries.code` |
| `location.adminCode` | no | Official division `code` (leaf preferred). Resolve `administrative_divisions` in that country |
| `location.street` | no | Line / house / building |
| `location.postalCode` | no | |
| `location.formatted` | no | Display line. Legacy mapping key `address` MAY be accepted as an alias for `formatted` |
| `location.name` / `locationKey` / lat / lng | recommended | Site identity |
| `place.name` | yes for a successful item | |
| `place.placeKey` | no | External id. If absent, `place_key` = `lower(trim(name)) + '\|' + lower(trim(unit))` |
| `news.originalUrl` | yes for a News row | Unique article URL |
| `details.*` | no | Only keys in the allow-list below |
| `categoryMap` | no | Source category string → stored category |

**Admin resolve**: `countryCode` + `adminCode` → `admin_division_id`. If only a name is mapped (no code), P1 may leave `admin_division_id` NULL (do not guess). Legacy district codes that are `is_active = false` do not match unless a later mapping table is added (P2).

**Details allow-list** (stored on `places.details`): `description`, `phone`, `hours`, `website`, `priceRange`. Unknown detail keys are ignored.

**Category**: apply `categoryMap` when present; else accept already-canonical `categories.key` values. Resolve to an existing `categories` row (upsert-by-key is **forbidden** for unknown keys — use `uncategorized`). Never insert a duplicate `categories` row for a key that already exists.

**Partial item**: article URL but no Place identity → create a new Place under the Location if the site can be determined; if Location cannot be determined, skip the item (`items_error++`).

## Relationships

```text
countries (1) ---- (*) administrative_divisions
administrative_divisions (1) ---- (*) administrative_divisions   (parent_id)
countries (1) ---- (*) locations
administrative_divisions (1) ---- (*) locations
locations (1) ---- (*) places
categories (1) ---- (*) places
categories (1) ---- (*) news
places (1) ---- (*) news
data_sources (1) ---- (*) news          (optional; SET NULL on source delete)
users (1) ---- (*) data_ingest_runs
data_ingest_runs (1) ---- (*) data_ingest_run_sources
data_sources (1) ---- (*) data_ingest_run_sources  (optional; SET NULL)
```

## Migration

- Up/down: `be/migrations/000010_geo_and_data_sources.{up,down}.sql`
- Next number after `000009_user_is_super_admin`
- Notes:
  - Enable `pgcrypto` / `gen_random_uuid()` if not already.
  - Down drops in FK-safe order: `data_ingest_run_sources`, `data_ingest_runs`, `news`, `places`, `data_sources`, `locations`, `administrative_divisions`, `countries`, `categories`.
  - Seed P1 `categories` rows by unique `key` in the same migration or a seeder (idempotent). Never insert a second row for an existing `key`.
  - No seed Places or Data Sources required for P1 (empty map is a valid US1 state).
  - Seed **VN** country + official 2-level divisions (Decision 19/2025/QĐ-TTg: 34 provinces, 3,321 communes/wards) via a **seeder or SQL fixture**, not hand-typed rows in the migration. Inactive legacy districts are out of P1 seed.

## Retention / soft-delete

- Data Sources: hard delete. News rows keep `source_name`; `source_id` becomes NULL.
- Places / Locations: no operator delete in P1 (`hidden` hides pins).
- Ingest runs: retain; no purge in P1.
- Admin divisions: deactivate (`is_active`), do not delete rows that Locations reference.
- No map-vendor metadata to retain.

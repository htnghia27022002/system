# Endpoints contract: Maps

**Feature**: `008-maps`  
**Owner**: `@technical-architect` (phase 3)  
**Status**: Ready for implement

> Authoritative HTTP API contract for this feature. Implementers follow this file; `plan.md` links here.  
> Base: `NEXT_PUBLIC_API_BASE_URL` → `/api`. JSON fields are **camelCase**.  
> Persistence and APIs are **map-vendor-agnostic** (lat/lng only). No OSM or Google types in request/response.  
> Sources routes read/write the shared **`data_sources`** table. Places/Locations use shared **`places`** / **`locations`**.

## Scope

- **New routes**: JWT admin group under `/api/admin/maps/*`
- **Changed routes**: none
- **N/A**: not applicable — feature adds HTTP APIs

## Conventions

- Prefix: `/api`
- Auth: JWT via `middleware.Auth`, plus RBAC `RequireView("maps")` / `RequireModify("maps")` (see [permissions.md](./permissions.md))
- Super admin bypass from `007-super-admin` applies after sign-in; it does **not** replace authentication
- List defaults: `page=1`, `limit=50` (max 200) where paginated
- Errors: `{ "error": "message" }` via existing `response.HandleError`
- Coordinates: numbers (`lat`, `lng`). Valid when both present, `lat` ∈ [-90, 90], `lng` ∈ [-180, 180]
- Category JSON: `category` is the catalog `key` in camelCase (`roomRental` | `restaurant` | `hotel` | `eatery` | `uncategorized`). Also return `categoryId` (UUID). Table is shared `categories`
- Place status JSON: `pending` | `active` | `hidden` (DB `places.status` SMALLINT 1/2/3 — map in the app)
- Ingest run status JSON: `queued` | `running` | `completed` | `failed` (DB 1/2/3/4)
- Per-source ingest status JSON: `pending` | `running` | `completed` | `failed` | `skipped` (DB 1–5)
- HTTP method JSON: `GET` | `POST` | `PUT` | `PATCH` (DB `data_sources.http_method` 1–4)
- Audit JSON (when present): `createdBy`, `updatedBy` (user UUID or `null`)

## Shared shapes

**LocationSummary**

```json
{
  "id": "uuid",
  "name": "Chung cu Example",
  "locationKey": "bldg-1",
  "countryCode": "VN",
  "adminDivisionId": "uuid",
  "adminPath": "VN/79/26734/",
  "street": "12 Nguyen Hue",
  "postalCode": "",
  "formatted": "12 Nguyen Hue, Phuong Ben Nghe, TP Ho Chi Minh",
  "lat": 10.7769,
  "lng": 106.7009
}
```

`locationKey`, `countryCode`, `adminDivisionId`, `adminPath`, `lat`, and `lng` may be `null` or empty as documented. `formatted` is the display line (replaces the former single `address` string).

**AdminDivision**

```json
{
  "id": "uuid",
  "countryCode": "VN",
  "parentId": "uuid",
  "level": 1,
  "code": "48",
  "name": "Da Nang",
  "nameEn": "Da Nang",
  "fullName": "Thanh pho Da Nang",
  "type": "municipality",
  "path": "VN/48/",
  "lat": 16.0544,
  "lng": 108.2022
}
```

`parentId` is `null` for top-level divisions.

**PlacePin** (map list item)

```json
{
  "id": "uuid",
  "name": "Pho Shop",
  "categoryId": "uuid",
  "category": "eatery",
  "status": "active",
  "lat": 10.7769,
  "lng": 106.7009,
  "locationId": "uuid",
  "locationName": "Chung cu Example"
}
```

**NewsItem**

```json
{
  "id": "uuid",
  "title": "New pho shop opens",
  "originalUrl": "https://news.example/article-1",
  "categoryId": "uuid",
  "category": "eatery",
  "sourceName": "City Feed",
  "sourceId": "uuid",
  "createdAt": "2026-09-12T00:00:00Z",
  "updatedAt": "2026-09-12T00:00:00Z"
}
```

`sourceId` may be `null` if the Source was deleted.

**SourceRecord**

```json
{
  "id": "uuid",
  "name": "City Feed",
  "enabled": true,
  "httpMethod": "GET",
  "url": "https://feeds.example/places.json",
  "headers": { "Authorization": "Bearer …" },
  "queryParams": { "city": "hcm" },
  "body": "",
  "fieldMapping": {},
  "createdBy": "uuid",
  "updatedBy": "uuid",
  "createdAt": "2026-09-12T00:00:00Z",
  "updatedAt": "2026-09-12T00:00:00Z"
}
```

P1 returns stored headers to `maps:view` (admin-only surface; no vault). Do not log header values.

**IngestRun**

```json
{
  "id": "uuid",
  "status": "running",
  "errorMessage": "",
  "sourcesTotal": 2,
  "success": 1,
  "error": 0,
  "itemsSuccess": 40,
  "itemsError": 2,
  "sources": [
    {
      "id": "uuid",
      "sourceId": "uuid",
      "sourceName": "City Feed",
      "status": "completed",
      "errorMessage": "",
      "itemsSuccess": 40,
      "itemsError": 2
    }
  ],
  "startedAt": "2026-09-12T00:00:00Z",
  "finishedAt": null,
  "createdAt": "2026-09-12T00:00:00Z"
}
```

`sources` is omitted or empty on list-of-runs if a compact list is used; **latest** and **get by id** MUST include per-source outcomes.

---

## Endpoints

### `GET /api/admin/maps/places`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | Query: `category` (`roomRental` \| `restaurant` \| `hotel` \| `eatery`, optional), `q` (optional free-text on place name, location `street` / `formatted`), `adminDivisionId` (optional UUID — include Places whose Location division `path` equals or is a descendant of that unit), `countryCode` (optional ISO alpha-2), `page` (default 1), `limit` (default 50, max 200) |
| **Success** | `200` — pinnable Places only: `status=active` **and** valid lat/lng. Non-P1 / `uncategorized` excluded when `category` is set; when `category` is omitted, still exclude `pending` / `hidden` and invalid coords. `uncategorized` active pins may appear only when no category filter is applied |
| **Errors** | `401`; `403`; `400` invalid `category` |

**Success example**

```json
{
  "items": [
    {
      "id": "uuid",
      "name": "Pho Shop",
      "categoryId": "uuid",
      "category": "eatery",
      "status": "active",
      "lat": 10.7769,
      "lng": 106.7009,
      "locationId": "uuid",
      "locationName": "Chung cu Example"
    }
  ],
  "page": 1,
  "limit": 50,
  "total": 1,
  "hasMore": false
}
```

Empty `items` is success (empty map / empty filter), not an error.

---

### `GET /api/admin/maps/places/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | Path `:id` = place UUID |
| **Success** | `200` — Place detail + Location + news (any status; used for pin click and ingest follow-up) |
| **Errors** | `401`; `403`; `404` |

**Success example**

```json
{
  "id": "uuid",
  "name": "Pho Shop",
  "unit": "A-12",
  "placeKey": "shop-88",
  "categoryId": "uuid",
  "category": "eatery",
  "status": "active",
  "lat": 10.7769,
  "lng": 106.7009,
  "details": {
    "description": "Pho and drinks",
    "phone": "+84…"
  },
  "location": {
    "id": "uuid",
    "name": "Chung cu Example",
    "locationKey": "bldg-1",
    "countryCode": "VN",
    "adminDivisionId": "uuid",
    "adminPath": "VN/79/26734/",
    "street": "12 Nguyen Hue",
    "postalCode": "",
    "formatted": "12 Nguyen Hue, Phuong Ben Nghe, TP Ho Chi Minh",
    "lat": 10.7769,
    "lng": 106.7009
  },
  "news": [
    {
      "id": "uuid",
      "title": "New pho shop opens",
      "originalUrl": "https://news.example/article-1",
      "categoryId": "uuid",
      "category": "eatery",
      "sourceName": "City Feed",
      "sourceId": "uuid",
      "createdAt": "2026-09-12T00:00:00Z",
      "updatedAt": "2026-09-12T00:00:00Z"
    }
  ],
  "createdAt": "2026-09-12T00:00:00Z",
  "updatedAt": "2026-09-12T00:00:00Z"
}
```

Empty `news` is `[]`, not an error.

---

### `POST /api/admin/maps/places`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | `{ locationId, name, category, placeKey?, unit?, status?, lat?, lng?, details? }` |
| **Success** | `201` — Place detail |
| **Errors** | `401`; `403`; `400`; `409` duplicate `(locationId, placeKey)` |

### `DELETE /api/admin/maps/places/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Success** | `204` — News rows cascade |
| **Errors** | `401`; `403`; `404` |

`GET /places?manage=true` lists Places of any status (optional `locationId`, `status`, `q`) for operator CRUD. Omit `manage` to keep the pinnable map list.

### Location CRUD

| Method | Path | Auth | Notes |
|--------|------|------|-------|
| GET | `/api/admin/maps/locations` | view | Paginated; `q`, `countryCode`, `adminDivisionId` |
| GET | `/api/admin/maps/locations/:id` | view | |
| POST | `/api/admin/maps/locations` | modify | `{ name, locationKey?, countryCode?, adminDivisionId?, street?, postalCode?, formatted?, lat?, lng? }` |
| PATCH | `/api/admin/maps/locations/:id` | modify | Partial body |
| DELETE | `/api/admin/maps/locations/:id` | modify | `409` if the Location still has Places |

### `PATCH /api/admin/maps/places/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | Partial Place fields: `status` (`pending` \| `active` \| `hidden`), `name`, `category`, `locationId`, `placeKey`, `unit`, `lat`, `lng`, `details`. Setting `active` still requires valid coordinates |
| **Success** | `200` — updated Place detail (same shape as GET) |
| **Errors** | `401`; `403`; `404`; `400` if `status` is `pending` or unknown; `400` if setting `active` without valid coordinates |

---

### `GET /api/admin/maps/sources`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | Query: `page`, `limit` (defaults as above) |
| **Success** | `200` — all Sources, newest-updated first |
| **Errors** | `401`; `403` |

**Success example**

```json
{
  "items": [
    {
      "id": "uuid",
      "name": "City Feed",
      "enabled": true,
      "httpMethod": "GET",
      "url": "https://feeds.example/places.json",
      "headers": { "Authorization": "Bearer …" },
      "queryParams": {},
      "body": "",
      "fieldMapping": {},
      "createdAt": "2026-09-12T00:00:00Z",
      "updatedAt": "2026-09-12T00:00:00Z"
    }
  ],
  "page": 1,
  "limit": 50,
  "total": 1,
  "hasMore": false
}
```

---

### `GET /api/admin/maps/sources/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | Path `:id` |
| **Success** | `200` — `SourceRecord` |
| **Errors** | `401`; `403`; `404` |

---

### `POST /api/admin/maps/sources/probe`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | `httpMethod`, `url`, optional `headers`, `queryParams`, `body` (same shapes as create; not persisted) |
| **Success** | `200` — `{ "status": 200, "body": <parsed JSON> }` |
| **Errors** | `401`; `403`; `400` (invalid URL, non-JSON, HTML, HTTP error, timeout, body > 5 MiB) |

Used only to preview a structured list while configuring a Source. Load data still fetches **saved enabled** Sources only.

---

### `POST /api/admin/maps/sources`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | Body fields below |
| **Success** | `201` — `SourceRecord` |
| **Errors** | `401`; `403`; `400` validation |

**Request body**

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| `name` | string | yes | Non-empty, max 200 |
| `enabled` | boolean | no | Default `true` |
| `httpMethod` | string | yes | `GET` \| `POST` \| `PUT` \| `PATCH` |
| `url` | string | yes | Absolute `http://` or `https://` |
| `headers` | object | no | String map; default `{}` |
| `queryParams` | object | no | String map; default `{}` |
| `body` | string | no | Default `""` |
| `fieldMapping` | object | yes | Must only target fixed schema paths (see [database.md](./database.md)) |

Unknown `fieldMapping` target keys → `400`. Mapping does not invent stored properties.

---

### `PATCH /api/admin/maps/sources/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | Partial body (same fields as POST, all optional) |
| **Success** | `200` — updated `SourceRecord` |
| **Errors** | `401`; `403`; `404`; `400` |

Later Load data MUST use the saved values (not a scratch unsaved URL).

---

### `DELETE /api/admin/maps/sources/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | Path `:id` |
| **Success** | `204` |
| **Errors** | `401`; `403`; `404` |

Hard delete. Existing News rows keep `sourceName`; `sourceId` becomes null. Places are not deleted.

---

### `POST /api/admin/maps/ingest`

| | |
|--|--|
| **Auth** | JWT + `maps:modify` |
| **Request** | Empty body. Fetches **only** persisted Sources with `enabled=true`. Never fetches a URL that is not a saved Source |
| **Success** | `202` — `IngestRun` (`status` `queued` or `running`) when the queue worker will process; **or** `200` — completed/failed `IngestRun` when processed synchronously (NATS disabled and queue worker not used) |
| **Errors** | `401`; `403`; `409` if another run is `queued` or `running` |

**409 example**

```json
{ "error": "ingest already in progress" }
```

No enabled Sources: still create a run, do **not** fetch the open web, finish as `failed` with a clear `errorMessage` (e.g. `no enabled sources`). Prefer `200`/`202` with that run over a bare `400` so the UI can show the empty/error ingest state.

Ingest MUST:

1. Treat each Source response as a JSON array (root or `fieldMapping.listPath`).
2. Apply mapping per item (Location → Place → News).
3. Skip disabled Sources entirely (no HTTP call).
4. Isolate Source failures: one failed Source does not delete Places from earlier successful Sources or earlier items.
5. Enforce HTTP timeout (recommend 30s/Source) and response size cap (recommend 5 MiB). Non-JSON / HTML / non-list → that Source `failed`.

---

### `GET /api/admin/maps/ingest`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | None (latest run) |
| **Success** | `200` — latest `IngestRun` including `sources[]`, or `{ "run": null }` when none exist |
| **Errors** | `401`; `403` |

**Success when none**

```json
{ "run": null }
```

**Success when present**

```json
{ "run": { "id": "uuid", "status": "completed", "sources": [] } }
```

(Full `IngestRun` fields as above.)

FE polls this while Load data is in progress.

---

### `GET /api/admin/maps/ingest/:id`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | Path `:id` = run UUID |
| **Success** | `200` — `IngestRun` with `sources[]` |
| **Errors** | `401`; `403`; `404` |

---

### `GET /api/admin/maps/categories`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | None |
| **Success** | `200` — `{ "items": [ { "id", "key", "name", "nameLocal", "sortOrder" } ] }` active categories, `sortOrder` then name. `key` is camelCase in JSON (`roomRental`) |
| **Errors** | `401`; `403` |

Table is shared `categories`. P1 seed rows only. Same `key` is never returned twice.

---

### `GET /api/admin/maps/search`

| | |
|--|--|
| **Auth** | JWT + `maps:view` |
| **Request** | Query: `q` (required, min 2 characters after trim), `countryCode` (optional ISO alpha-2 — biases geocode and local lists) |
| **Success** | `200` — `{ "items": [ SearchHit, … ] }` |
| **Errors** | `401`; `403` |

**SearchHit**

| Field | Type | Notes |
|-------|------|-------|
| `id` | string | Stable hit id (`place:…`, `location:…`, `geocode:{provider}:…`) |
| `kind` | string | `place` \| `location` \| `geocode` (no vendor names) |
| `title` | string | Display name |
| `subtitle` | string | Address or location line |
| `lat` | number | WGS84 |
| `lng` | number | WGS84 |
| `placeId` | string | Present when `kind=place` |
| `locationId` | string | Present when `kind=place` or `kind=location` |

Combines saved Places (pinnable / active with coords), saved Locations with coords, and the configured **geocode Provider** (`be/internal/services/maps/`). A short `q` returns `{ "items": [] }`. Geocode failures do not hide saved hits. The browser must not call Nominatim.

### `GET /api/address/countries`

| | |
|--|--|
| **Auth** | JWT (any signed-in user) |
| **Request** | None |
| **Success** | `200` — `{ "items": [ { "id", "code", "code3", "name", "nameLocal" } ] }` active countries only |
| **Errors** | `401` |

P1 always includes `VN`. Table is `countries` (shared). Not scoped to `maps:view` so other screens can reuse the catalog.

### `GET /api/address/divisions`

| | |
|--|--|
| **Auth** | JWT (any signed-in user) |
| **Request** | Query: `countryCode` (required, ISO alpha-2), `parentId` (optional UUID — children of that node; omit for top-level), `q` (optional name search) |
| **Success** | `200` — `{ "items": [ AdminDivision, … ] }` active divisions only, ordered by `name` |
| **Errors** | `401`; `400` missing/unknown `countryCode` |

Used for cascaded country → division dropdowns. Table is `administrative_divisions` (shared).

---

## Product / edge paths

| Public path | Proxies to | Notes |
|-------------|------------|-------|
| None | — | No guest map and no BE tile proxy in P1. Basemap tiles are loaded by the **FE MapProvider** (OSM adapter talks to public OSM tiles in the browser) |

## UI routes (FE)

| Route | Auth | Behavior |
|-------|------|----------|
| `/admin/maps` | Signed in + `maps:view` or Super admin (`PermissionGuard`) | Street map, filters, pin detail, Sources tab/panel, Load data, ingest status |
| Guest `/admin/maps` | — | Existing admin `ProtectedGuard` → login |

No second admin menu item for Sources.

## Validation / errors

| Case | Status |
|------|--------|
| Missing/expired JWT | `401` |
| Signed in, not Super admin, no `maps:view` on GET | `403` |
| Signed in, not Super admin, no `maps:modify` on write / ingest / status | `403` |
| Unknown place / source / run id | `404` |
| Invalid body, mapping, category, or `active` without coords | `400` |
| Concurrent Load data | `409` |
| Source HTTP failure during ingest | run/source `failed` in body, not necessarily HTTP 5xx on POST if the job finished |

## Out of scope (no routes)

- Public `/maps` guest API
- Tile proxy `/api/admin/maps/tiles`
- Scheduled ingest
- Place directory distinct from pin detail (P1)
- Dynamic schema / per-source columns

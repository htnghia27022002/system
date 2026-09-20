# Feature Specification: Maps (Admin Street Map, Sources, and Ingest)

**Feature ID**: `008-maps`

**Feature Branch**: `008-maps`

**Created**: 2026-09-12

**Status**: Draft — ready for tasks

**Input**: Platform operators need an admin Maps surface on a public street map. They configure HTTP Sources, load structured list data into a fixed Location / Place / News model, pin Places by category, and open pin detail with source name and a link to the original article. No open-web scrape. No public guest map in P1.

## Clarifications

### Session 2026-09-12 (user request — locked P1)

- Q1: Who sees Maps? → A: Maps is a **platform admin menu** item in the same chrome as Users / Roles / Webhooks. Signed-in users with **`maps:view`** see the menu, open the map, filter/search, and open pin detail. Users without `maps:view` must not see the menu and must not access map or source **view** APIs. The Super admin bypass from `007-super-admin` still applies (signed-in Super admin passes Maps gates without catalog keys).
- Q2: How is data loaded? → A: Operators manage a **Sources** table (create, read, update, delete). Each source stores at least: **name**, **enabled** flag, **HTTP method**, **URL**, **headers**, **query parameters**, **body**, and **field mapping**. Mapping tells the product which fields from each source response item save into the product’s **fixed** schema (Location / Place / News / details). Mapping does **not** create arbitrary new database columns per source. **Load data** runs ingest **only** against configured sources. No open-web scrape. No crawling sites that are not saved as Sources. Ingest is an authorized HTTP fetch of operator-configured endpoints that return a **structured list**. The product treats the response as a list of items and applies mapping per item.
- Q3: Pin model (Option A, then chung cư clarification)? → A: The user first chose a single-pin-per-site idea, then clarified that one geographic site (for example an apartment building / chung cư) **may** contain many restaurants, shops, room rentals, and similar units. **Locked model (recommended Option B; analysis proceeds without contradicting that clarification):**
  - **Location** = one geographic site (address + coordinates). Example: a chung cư.
  - **Place** = one business or unit at that Location (eatery, shop, hotel, room rental). Many Places per Location.
  - **News** = one sourced article or post about a Place. Many News rows per Place. Table name is `news`, not `listings`.
  - Place has **details** and **status** (`pending` / `active` / `hidden`).
  - The product **MUST NOT** merge two different Places only because they share coordinates or the same Location.
  - Map pins are **per Place** so category filters work. Overlapping pins at one building are allowed.
  - The same **article URL** must not create a duplicate News row.
  - Attach a new News row to an existing Place only when mapping supplies a stable **place key** (external id and/or name+unit) that already exists **at that Location**. If no match, create a new Place under the Location (create the Location first if the site is new).
  - A Place without valid coordinates stays `pending` and is **not** pinned.
- Q4: RBAC default (standard view/modify; user did not override)? → A:
  - **`maps:view`** — menu, map, filters, pin detail, read-only sources list.
  - **`maps:modify`** — create/edit/delete sources, trigger Load data / ingest, change Place status.
- Q5: Address model? → A: **Option A confirmed.** Shared **Country** + **Administrative division** tree (`parent_id`, ordinal `level`, semantic `type`, materialized `path`). Location stores structured address: country, leaf admin unit, street, postal code, formatted line, lat/lng. Do **not** hard-code level 2 as district (Vietnam is 2 official levels from 1 Jul 2025). Filter/calculate by ancestor `path`. Place keeps unit only.
- Q6: Table names? → A: **Shared platform names**, not `map_*`. `locations`, `places`, `categories`, `news` (not `listings`), `data_sources`, `data_ingest_runs`, `data_ingest_run_sources`, plus `countries` and `administrative_divisions`. Maps is the first consumer; other features may reuse the same tables.
- Q7: Entity chain and uniqueness? → A: **Location → Place → Category → News.** Category is its own table with unique `key` (same danh mục never creates a second category row). News belongs to a Place and that Place’s Category. The same **article `original_url`** never creates a second News row (refresh in place).

## Overview

Operators who already use the admin shell open **Maps** to see Places on a public street map, filter by category, and inspect sourced news. They do not scrape the open web. They save Sources (HTTP endpoints + mapping) and press **Load data** to ingest structured lists into a fixed three-level model:

**Location** (site) → **Place** (business or unit) → **Category** (danh mục) → **News** (sourced article).

| Surface | Path | Auth | Behavior |
|---------|------|------|----------|
| Admin Maps | `/admin/maps` | Signed in + `maps:view` **or Super admin** | Street map, category filter/search, pins, pin detail, read-only sources list, ingest status |
| Sources management | Same Maps surface (panel or tab) | Signed in + `maps:modify` **or Super admin** to change; `maps:view` may read | Create/edit/delete Sources; store method, URL, headers, params, body, mapping, enabled |
| Load data | Action on Maps | Signed in + `maps:modify` **or Super admin** | Ingest enabled Sources only; persist Location / Place / News; update pins and ingest status |
| Place status | Pin detail or Place record | Signed in + `maps:modify` **or Super admin** | Set `active` / `hidden` (and see `pending` when coordinates are missing) |

**Depends on**: `001-auth` (signed-in sessions), existing admin chrome and RBAC (`maps:view` / `maps:modify` follow the standard catalog pattern), `007-super-admin` (bypass after sign-in). **Does not** replace Users, Roles, or Webhooks.

### Product roadmap / phased delivery

| Phase | Scope | Status |
|-------|--------|--------|
| **Phase 1 (P1)** | Admin menu Maps gated by `maps:view`; map view + category and admin-division filter/search + pins + detail with source name and original-article link; Data Sources CRUD (name, enabled, method, URL, headers, params, body, mapping) on shared `data_sources`; Load data / ingest from enabled Data Sources; shared Location / Place / News persistence, Country + admin-division catalog, structured address on Location; empty / error / in-progress ingest states | **In scope** |
| **Phase 2 (P2)** | Public guest map (no login); scheduled / automatic ingest; user-submitted places; turn-by-turn routing; booking / payments; reviews; creating dynamic storage columns from mapping; open-web HTML scraping of arbitrary sites | **Deferred** |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Operator with maps:view opens Maps from the admin menu (Priority: P1)

A signed-in operator who has `maps:view` (or Super admin) sees **Maps** in the admin navigation, opens it, and sees a public street map. They can use the surface even when no Places are pinned yet.

**Why this priority**: Without a gated map surface, no other Maps value is reachable.

**Independent Test**: Sign in as a user with `maps:view`. Confirm Maps appears in admin navigation, `/admin/maps` opens, and a public street map is visible (empty pins are acceptable).

**Acceptance Scenarios**:

1. **Given** a signed-in user with `maps:view`, **When** they view admin navigation, **Then** a Maps item is visible in the same chrome family as Users / Roles / Webhooks.
2. **Given** that user, **When** they open Maps, **Then** they see a public street map on `/admin/maps`.
3. **Given** no active Places with valid coordinates, **When** they open Maps, **Then** the map still loads and an empty state explains that nothing is pinned yet.
4. **Given** a guest (not signed in), **When** they try to open `/admin/maps`, **Then** they are sent to login and see no map data.

---

### User Story 2 - Filter and search by category shows only matching active Places (Priority: P1)

The operator narrows the map to P1 categories: **room rental**, **restaurant**, **hotel**, **eatery**. Only matching **active** Places with valid coordinates remain pinned.

**Why this priority**: Category discovery is the reason pins are per Place, not per building.

**Independent Test**: Seed or ingest at least two categories. Apply a single-category filter and a search that implies a category. Confirm only matching active pins remain; hidden and pending Places stay off the map.

**Acceptance Scenarios**:

1. **Given** active Places in more than one P1 category, **When** the operator filters to one category, **Then** only pins for that category remain.
2. **Given** a category search or filter with no matching active Places, **When** results update, **Then** the map shows no pins and a clear empty-filter state (not a blank failure).
3. **Given** Places that are `hidden` or `pending`, **When** any category filter or search runs, **Then** those Places are not pinned.
4. **Given** several Places at the same Location in different categories, **When** the operator filters to one category, **Then** only Places in that category at that site remain pinned; the other units at the building are not shown.
5. **Given** Places under two different provinces (or equivalent level-1 units), **When** the operator filters to one administrative division, **Then** only Places whose Location is in that unit or a descendant remain pinned.

---

### User Story 3 - Click pin shows Place detail and news with attribution (Priority: P1)

The operator selects a pin and sees that Place’s details, status, and news. Each news item shows the **source name** and a **link to the original article**.

**Why this priority**: Attribution and inspectability are the product’s “news on a map” value.

**Independent Test**: Open a Place that has at least one News row. Confirm detail shows Place fields, headlines or equivalents, source name, and a working original URL.

**Acceptance Scenarios**:

1. **Given** an active Place with valid coordinates, **When** the operator selects its pin, **Then** they see Place details (at least name, category, status, and mapped details) and the Location identity (address and/or site name).
2. **Given** that Place has one or more News rows, **When** detail opens, **Then** each News item shows the source name and a link to the original article URL.
3. **Given** a Place with no News yet, **When** detail opens, **Then** the operator sees an empty-news state, not an error.
4. **Given** a `maps:view` operator who cannot modify, **When** they open detail, **Then** they can read detail and news but cannot change Place status or trigger Load data.

---

### User Story 4 - Operator with maps:modify creates and edits a Source (Priority: P1)

An operator with `maps:modify` (or Super admin) adds or updates a Source: name, enabled flag, HTTP method, URL, headers, query parameters, body, and field mapping onto the fixed schema.

**Why this priority**: Ingest is useless without operator-owned Sources.

**Independent Test**: Create a Source with all required fields, disable it, edit mapping, then delete a throwaway Source. Confirm a `maps:view`-only user can see the list but cannot change it.

**Acceptance Scenarios**:

1. **Given** an operator with `maps:modify`, **When** they create a Source, **Then** they can save name, enabled flag, HTTP method, URL, headers, query parameters, body, and field mapping.
2. **Given** an existing Source, **When** they edit and save, **Then** the stored configuration reflects the new values and later Load data uses the updated settings.
3. **Given** a Source the operator no longer wants, **When** they delete it, **Then** it no longer appears in the Sources list and is not included in the next Load data.
4. **Given** a user with `maps:view` only, **When** they open Sources, **Then** they see a read-only list and cannot create, edit, or delete.
5. **Given** mapping configuration, **When** the operator saves it, **Then** mapping only targets the product’s fixed Location / Place / News / details fields and does not invent new stored properties outside that schema.

---

### User Story 5 - Load data ingests enabled Sources and pins Places with coordinates (Priority: P1)

The operator starts **Load data**. The product fetches each **enabled** Source, maps each list item, persists Location / Place / News using the merge rules, and pins Places that are `active` with valid coordinates. The UI shows in-progress, success, empty, and error states.

**Why this priority**: This is the only P1 way Places appear on the map.

**Independent Test**: Configure one enabled Source that returns a structured list with mappable coordinates and one disabled Source. Run Load data. Confirm only the enabled Source is fetched, Places with coordinates appear as pins, ingest shows progress then completion, and a failing Source shows an error without wiping successful pins.

**Acceptance Scenarios**:

1. **Given** at least one enabled Source, **When** the operator with `maps:modify` starts Load data, **Then** ingest runs only against configured enabled Sources (not against arbitrary websites).
2. **Given** ingest is running, **When** the operator views Maps, **Then** they see an in-progress state and cannot mistake it for a finished empty map.
3. **Given** a Source response that is a list of items, **When** ingest maps each item, **Then** Location / Place / News records persist according to the merge rules, and Places with valid coordinates and `active` status appear as pins.
4. **Given** a Place whose mapped coordinates are missing or invalid, **When** ingest finishes, **Then** that Place is `pending` and is not pinned.
5. **Given** the same article URL is ingested again, **When** Load data completes, **Then** a second News row is not created for that URL (the existing News row is updated).
6. **Given** two items map to the same category key, **When** ingest completes, **Then** only one `categories` row exists for that key.
6. **Given** a Source that fails (unreachable, non-list payload, or mapping failure for all items), **When** ingest finishes, **Then** the operator sees an error for that Source; already persisted Places from other Sources remain.
7. **Given** no enabled Sources, **When** the operator starts Load data, **Then** the product does not fetch the open web and shows a clear empty/error outcome (nothing to ingest).
8. **Given** a user with `maps:view` only, **When** they try to start Load data, **Then** they are denied.

---

### User Story 6 - One Location can host many Places; two shops in one chung cư stay distinct (Priority: P1)

Two different businesses in the same apartment building share a Location (the building) but remain two Places and two pins. They merge into one Place only when a stable place key matches at that Location.

**Why this priority**: The chung cư clarification is the core domain rule; collapsing shops into one pin would break category filters.

**Independent Test**: Ingest (or seed) two shops with the same building address/coordinates and different place keys or names. Confirm two Places, two pins, and that selecting each pin shows the correct shop. Re-ingest one shop’s article URL and confirm it attaches to the existing Place.

**Acceptance Scenarios**:

1. **Given** two sourced items at the same site (same Location) with different place keys or different name+unit values, **When** ingest completes, **Then** the product stores two Places and does not merge them because they share coordinates or the same Location.
2. **Given** those two Places are `active` with valid coordinates, **When** the operator views the map, **Then** both pins are present even if they overlap on the building.
3. **Given** a new News item whose mapping supplies a place key that already exists at that Location, **When** ingest runs, **Then** the News row attaches to that existing Place and a duplicate Place is not created.
4. **Given** a new News item at a known Location with **no** matching place key, **When** ingest runs, **Then** a new Place is created under that Location.
5. **Given** a sourced item for a site that does not exist yet, **When** ingest runs, **Then** the product creates the Location first, then the Place, then the News row.

---

### User Story 7 - Permissions: view sees the map; modify changes Sources and ingest (Priority: P1)

Access follows the standard view/modify split. Missing `maps:view` hides the menu and blocks map/source view APIs. Missing `maps:modify` blocks Source writes, Load data, and Place status changes. Super admin still bypasses after sign-in.

**Why this priority**: Maps is a platform admin module; leaking it to every signed-in user or leaving it ungated would break the RBAC model used by Users / Roles / Webhooks.

**Independent Test**: Sign in as (a) no Maps keys, (b) `maps:view` only, (c) `maps:view` + `maps:modify`, (d) Super admin with empty role permissions. Confirm menu, page, view APIs, and modify actions match the matrix below.

**Acceptance Scenarios**:

1. **Given** a signed-in user who is not Super admin and lacks `maps:view`, **When** they use the admin shell, **Then** Maps is absent from navigation and map/source view APIs deny them.
2. **Given** a signed-in user with `maps:view` and not `maps:modify`, **When** they open Maps, **Then** they can view the map, filter/search, open pin detail, and read the Sources list, but cannot create/edit/delete Sources, start Load data, or change Place status.
3. **Given** a signed-in user with `maps:modify` (and view, or Super admin), **When** they change a Source, start Load data, or change Place status, **Then** those actions succeed (subject to validation).
4. **Given** a signed-in Super admin with no Maps catalog keys, **When** they open Maps and call protected Maps actions, **Then** they are allowed (bypass from `007-super-admin`).
5. **Given** a Super admin who is not signed in, **When** they call a protected Maps API, **Then** they are rejected as unauthorized.

---

### Edge Cases

- Disabled Sources are stored but skipped by Load data.
- A Source URL that is not saved in Sources is never fetched, even if an operator types it only in a scratch field and does not save it.
- Response is not a list (object, HTML page, empty body) → that Source fails ingest; other Sources may still succeed.
- Partial mapping: item has an article URL but no Place identity → create a new Place under the Location if the site can be determined; if Location cannot be determined, skip the item and count it as a mapping failure.
- Same article URL on a later Load data → no duplicate News row; existing News may refresh mapped fields in place.
- Two items share coordinates but represent different Places → two Places; never merge on coordinates alone.
- Two items share a Location and the same place key → one Place; additional News rows attach if their article URLs are new.
- Place later receives valid coordinates on a subsequent ingest → it may leave `pending` and become pinnable once `active`.
- Operator sets a Place to `hidden` → it disappears from the map but remains persisted and readable in detail if the product still allows opening it from a non-map list; P1 does not require a separate Place directory if pin detail plus ingest results are enough.
- Overlapping pins at one building → allowed; selecting a pin still opens the correct Place.
- Concurrent Load data while another ingest is in progress → second start is rejected or queued with a clear in-progress message (no two silent overlapping runs).
- Source with secrets in headers → stored as Source configuration; `maps:view` may see that a header exists; product should not encourage leaking secrets in screenshots, but P1 does not add a vault product.
- Failed ingest does not delete Places created by earlier successful runs.
- Category values outside the P1 set (room rental, restaurant, hotel, eatery) are not filter options in P1; unmapped category → Place is stored but not shown in P1 category filters until mapped into a P1 category (or treated as uncategorized and excluded from those filters).
- Guest / expired session: Maps UI and APIs stay unauthorized.
- Account without Maps keys after Super admin flag is cleared: access falls back to catalog permissions immediately on the next authorized request (same rule as `007-super-admin`).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The admin navigation MUST include a Maps item in the same chrome family as Users / Roles / Webhooks, visible only to signed-in users with `maps:view` or Super admin.
- **FR-002**: `/admin/maps` MUST show a public street map to authorized viewers and MUST reject guests.
- **FR-003**: Users without `maps:view` who are not Super admin MUST NOT see the Maps menu and MUST be denied map and source view APIs.
- **FR-004**: Users with `maps:view` MUST be able to filter and search Places by P1 categories: room rental, restaurant, hotel, eatery, and MUST be able to filter by an administrative division (selected unit includes descendant Locations).
- **FR-005**: The map MUST pin **active** Places that have valid coordinates, one pin per Place. Overlapping pins at one Location MUST be allowed.
- **FR-006**: `pending` and `hidden` Places MUST NOT appear as pins.
- **FR-007**: Selecting a pin MUST open Place detail including name, category, status, details, Location identity, and News.
- **FR-008**: Each News item in detail MUST show the Source name and a link to the original article URL.
- **FR-009**: Operators with `maps:modify` MUST be able to create, edit, and delete Sources.
- **FR-010**: Each Source MUST store at least name, enabled flag, HTTP method, URL, headers, query parameters, body, and field mapping.
- **FR-011**: Field mapping MUST map source item fields onto the fixed Location / Place / News / details schema and MUST NOT create arbitrary new stored properties per source.
- **FR-012**: Users with `maps:view` and without `maps:modify` MUST see Sources as read-only and MUST NOT create, edit, or delete Sources.
- **FR-013**: Load data MUST be available only to `maps:modify` or Super admin and MUST ingest **only** enabled, saved Sources (no open-web scrape; no fetch of URLs that are not saved as Sources).
- **FR-014**: Ingest MUST treat each Source response as a list of items and apply mapping per item.
- **FR-015**: Ingest MUST persist Location, Place, News, details, and Place status using the merge rules in Clarifications Q3.
- **FR-016**: The product MUST NOT merge two Places solely because they share coordinates or the same Location.
- **FR-017**: A News row MUST attach to an existing Place only when mapping supplies a stable place key (external id and/or name+unit) that already exists at that Location; otherwise ingest MUST create a new Place (and create the Location first if the site is new).
- **FR-018**: The same article URL MUST NOT create a duplicate News row (update the existing News row in place).
- **FR-026**: Category MUST be a shared catalog table with a unique key. Ingest MUST reuse the existing Category row when the mapped danh mục matches; it MUST NOT insert a second Category for the same key.
- **FR-019**: A Place without valid coordinates MUST stay `pending` and MUST NOT be pinned.
- **FR-020**: Operators with `maps:modify` MUST be able to change Place status (`active` / `hidden` as operator actions; `pending` remains the incomplete-coordinates state).
- **FR-021**: Maps MUST present empty, in-progress, and error ingest states so operators can tell “no pins yet,” “ingest running,” and “a Source failed.”
- **FR-022**: Super admin bypass from `007-super-admin` MUST apply to Maps menu, pages, and APIs after sign-in; it MUST NOT replace authentication.
- **FR-023**: Load data MUST skip disabled Sources and MUST NOT crawl or scrape sites that are not saved as Sources.
- **FR-024**: Location MUST persist a structured address (country, administrative division when resolved, street, optional postal code, formatted display line) on the shared Location record — not a Maps-only table and not a single undifferentiated address string as the only stored form.
- **FR-025**: Persistence table names MUST be shared (`locations`, `places`, `categories`, `news`, `data_sources`, ingest-run tables, `countries`, `administrative_divisions`) so later features can reuse the same rows. Maps MUST NOT own `map_*` copies of these entities. Do not use a `listings` table.

### Key Entities

- **Country**: ISO country row (`code`, `code3`, names). Shared catalog.
- **Administrative division**: One node in a per-country tree (`parent`, ordinal `level`, `type`, official `code`, materialized `path`). Depth varies by country.
- **Data Source**: Operator-managed HTTP ingest endpoint (shared `data_sources`). Attributes include name, enabled flag, HTTP method, URL, headers, query parameters, body, field mapping, and timestamps. Only saved Data Sources may be fetched.
- **Field mapping**: Configuration that copies fields from one source list item onto the fixed Location / Place / News / details attributes. It does not extend the schema.
- **Location**: One geographic site (structured address + coordinates), for example a chung cư. Shared `locations` table. Many Places belong to one Location.
- **Place**: One business or unit at a Location. Shared `places` table. Attributes include name, **one Category**, details, status (`pending` / `active` / `hidden`), coordinates, and a stable place key. Map pins are per Place.
- **Category**: Shared danh mục (`categories`). Unique `key`. P1 keys: room rental, restaurant, hotel, eatery, plus `uncategorized`.
- **News**: One sourced article about a Place under that Place’s Category. Shared `news` table (not `listings`). Unique `original_url`. Attributes include title, timestamps, Data Source, `place_id`, and `category_id`.
- **Ingest run**: One operator-triggered Load data execution on `data_ingest_runs`. Attributes include progress state, per-Source outcomes, and enough summary to explain empty vs error.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An operator with `maps:view` can open Maps from admin navigation and see the street map in under 1 minute on a warm session (already signed in).
- **SC-002**: After a successful Load data of a Source whose items include valid coordinates and a P1 category, at least 95% of those mappable items appear as the correct category pins on refresh.
- **SC-003**: Filtering to one P1 category hides pins of other categories in 100% of Independent Test cases, including multiple Places at the same Location. Filtering to one admin division hides pins outside that unit and its descendants in 100% of Independent Test cases.
- **SC-004**: In 100% of trials, two shops at the same building with different place keys remain two Places and two pins; they are never merged only because they share a Location or coordinates.
- **SC-005**: Re-running Load data with the same article URL does not create a second News row in 100% of trials.
- **SC-006**: Opening a pin shows source name and a link to the original article for every News row that was ingested with those fields, in 100% of Independent Test cases.
- **SC-007**: A user without `maps:view` (and not Super admin) never sees the Maps menu and cannot read map/source view APIs in 100% of Independent Test attempts.
- **SC-008**: A user with `maps:view` only cannot create a Source, start Load data, or change Place status in 100% of Independent Test attempts.
- **SC-009**: At least 90% of evaluators in an informal walkthrough can start from the admin menu, apply one category filter, open a pin, and reach the original article link without assistance.
- **SC-010**: When ingest is running or a Source fails, 100% of evaluators can distinguish in-progress and error from a successful empty map.

## Assumptions

- The basemap is a **public street map**. The user named **OpenStreetMap** as the intended basemap; planning may use that constraint. Functional requirements stay basemap-agnostic (“public street map”) so the product behavior does not depend on a specific map library.
- Maps is an **admin module**, not a Tools-catalog product and not a guest marketing page in P1. Product path is `/admin/maps`, consistent with Users (`/admin/users`).
- Sources management lives on the Maps surface (panel or tab). P1 does not add a second admin menu item.
- Load data runs against **enabled** Sources only. Disabled Sources remain configured for later use.
- “Authorized HTTP fetch” means the product calls the URL, method, headers, query parameters, and body the operator saved. Authentication to third-party endpoints is whatever those headers/body supply. The product does not add a separate OAuth-broker for Sources in P1.
- Structured list means the Source returns a list of similar items (directly or in a conventional list wrapper). HTML pages and arbitrary websites are out of scope.
- Location reuse: ingest reuses an existing Location when the mapped site identity matches (`locationKey`, country + admin unit + street, formatted line, or rounded coordinates **plus** the same formatted line). A new site creates a new Location. This does **not** justify merging Places.
- Persistence uses shared table names (`locations`, `places`, `categories`, `news`, `data_sources`, …). Maps P1 is the first editor; it does not create `map_*` duplicates.
- Closed status / HTTP-method enums are stored as integers and mapped to text in the API. Rows record `created_by` / `updated_by` when an operator or ingest actor is known.
- New Places with valid coordinates become `active` and pinnable. Missing/invalid coordinates stay `pending`. Operators hide a Place with `hidden` and can return it to `active` when it is complete.
- P1 category filter set is exactly: room rental, restaurant, hotel, eatery. Free-text search may also match Place name and address; it still only returns `active` pinnable Places.
- Existing authentication, admin layout, permission catalog pattern, and Super admin bypass are reused. Administrator seed receives new catalog keys the same way other admin modules do (planning detail).
- Vietnamese UI copy, if any, lives only in the frontend Vietnamese locale files. This spec stays English.
- Independent deployability: ingest, persistence, and Maps APIs are owned by the backend service; the frontend owns the admin map UI and calls the backend over HTTP. Neither package hard-codes the other’s internals.
- Performance expectation: a single Load data of a typical enabled Source (hundreds of items, not unbounded dumps) completes in a time an operator can wait on the page, with in-progress feedback. Exact throughput is a planning concern.
- No requirement in P1 for clustering pins, heat maps, offline maps, or turn-by-turn directions.

## Out of Scope

- Public guest map without login (P2).
- Scheduled or automatic ingest (P2).
- User-submitted Places (P2).
- Turn-by-turn routing, booking, payments, and reviews (P2).
- Creating dynamic storage columns or per-source schemas from mapping (P2).
- Open-web HTML scraping or crawling sites that are not saved as Sources (P2 and not desired as a silent P1 fallback).
- Merging distinct Places because they share a building or coordinates.
- New permission verbs beyond `maps:view` and `maps:modify`.
- Changing Users, Roles, Webhooks, or Super admin rules except consuming the existing bypass.
- Pixel-locked map chrome or a specific map-library mandate in requirements (basemap preference is an assumption).

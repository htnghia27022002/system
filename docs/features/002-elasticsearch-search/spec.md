# Feature Specification: Elasticsearch Search (Admin Unified Search)

**Feature ID**: `002-elasticsearch-search`

**Created**: 2026-07-05

**Status**: Draft

**Input**: Admin unified search across users, roles, and permissions with permission-aware results, reliable PostgreSQL-to-search-index sync via outbox pattern, and graceful degradation when the search service is unavailable.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Admin searches across users, roles, and permissions (Priority: P1)

An authenticated admin opens the admin area and uses a unified search control to find users, roles, and permissions by typing a query (e.g., email fragment, role name, permission key). Results appear in a single list grouped or labeled by entity type, with enough context to identify the correct record and navigate to its detail or list view.

**Why this priority**: Unified search is the primary user-facing value of this feature and replaces fragmented, entity-specific lookup.

**Independent Test**: Sign in as an admin with `search.view` and entity view permissions, enter a known query on the admin search UI, and verify relevant hits from at least two entity types appear with correct titles and links.

**Acceptance Scenarios**:

1. **Given** an admin with `search.view`, `users.view`, and `roles.view`, **When** they search for a term matching both a user email and a role name, **Then** both hits appear in the response with entity type, title, and identifier.
2. **Given** an admin with `search.view`, **When** they filter results to `users` only via entity-type filter, **Then** only user hits are returned.
3. **Given** a search query with no matches, **When** the admin submits the search, **Then** an empty result set is shown with a clear “no results” message (not an error).
4. **Given** paginated results exceeding one page, **When** the admin requests the next page, **Then** additional hits are returned without duplicating prior page items.

---

### User Story 2 - Search index stays in sync after create, update, and delete (Priority: P1)

When an admin (or system process) creates, updates, or deletes a user, role, or permission in PostgreSQL, the corresponding search document is eventually created, updated, or removed in the search index. PostgreSQL remains the source of truth; the index is a derived read model.

**Why this priority**: Stale or missing index data undermines trust in search and is a blocker for production use.

**Independent Test**: Create or update a user via existing admin APIs, wait for sync, then search for the changed record and confirm the index reflects the new values.

**Acceptance Scenarios**:

1. **Given** a new user is created via the admin users API, **When** the outbox worker processes the pending entry, **Then** searching for that user’s email returns a hit with up-to-date title and searchable text.
2. **Given** a role name is updated, **When** sync completes, **Then** search results show the new name and no longer match only the old name.
3. **Given** a permission or user is soft-deleted (per existing domain rules), **When** sync completes, **Then** the entity no longer appears in search results.
4. **Given** the same outbox entry is processed more than once (at-least-once delivery), **When** the worker retries, **Then** the index state remains correct (idempotent upsert or delete).

---

### User Story 3 - Initial bulk reindex and recovery (Priority: P2)

When the search index is empty, corrupted, or drifted from PostgreSQL, an operator can trigger a full or scoped reindex so all searchable entities are rebuilt from the database without manual per-record fixes.

**Why this priority**: Required for first deploy, disaster recovery, and fixing sync failures at scale.

**Independent Test**: Clear or simulate an empty index, run reindex for all v1 entity types, then verify representative records from users, roles, and permissions are searchable.

**Acceptance Scenarios**:

1. **Given** an empty search index and populated PostgreSQL data, **When** a full reindex is triggered, **Then** all active users, roles, and permissions appear in search after completion.
2. **Given** a reindex scoped to one entity type (e.g., `users`), **When** it completes, **Then** only that entity type’s documents are rebuilt; other types remain unchanged.
3. **Given** reindex runs while normal admin traffic continues, **When** mutations occur during reindex, **Then** outbox processing continues and the final index state converges with PostgreSQL (no permanent duplicates or orphans from overlap).
4. **Given** reindex fails partway, **When** the operator retries, **Then** the operation can be safely restarted without corrupting existing index data.

---

### User Story 4 - Permission-aware search results (Priority: P2)

Search never exposes entities the caller is not allowed to view. A caller with `search.view` but without `users.view` does not see user hits; similarly for roles and permissions. Result filtering applies server-side before the response is returned.

**Why this priority**: Prevents information leakage through a new cross-entity surface area.

**Independent Test**: Authenticate as a user with `search.view` and `roles.view` only; search for a term that matches both a role and a user; verify only role hits are returned.

**Acceptance Scenarios**:

1. **Given** a caller with `search.view` and `roles.view` but not `users.view`, **When** they search for a user’s email, **Then** no user hits appear (even if the index contains matching users).
2. **Given** a caller with no entity view permissions but with `search.view`, **When** they search, **Then** they receive an empty hit list (HTTP 200), not an error disclosing protected data.
3. **Given** a caller without `search.view`, **When** they call the search endpoint, **Then** access is denied per existing RBAC middleware behavior (403).
4. **Given** indexed documents carry required permission keys for visibility, **When** results are assembled, **Then** every returned hit corresponds to a permission the caller holds.

---

### User Story 5 - Graceful degradation when search is unavailable (Priority: P3)

When the search service is disabled in configuration or temporarily unreachable, the admin experience degrades gracefully: the search endpoint returns a controlled response and, where a legacy database search path exists (e.g., user list filter), that path may be used as a limited fallback instead of a hard failure for the whole admin UI.

**Why this priority**: Keeps admin workflows usable during outages or local dev without search infrastructure.

**Independent Test**: Disable or stop the search service, invoke search from the admin UI or API, and verify a documented fallback or empty/degraded response without 500 errors for expected cases.

**Acceptance Scenarios**:

1. **Given** search is disabled via configuration, **When** an admin with `search.view` searches, **Then** the API responds with a clear degraded indicator and, for user-only queries where applicable, may return results from existing database search (email and full name) within the same permission rules.
2. **Given** the search service is down, **When** the search endpoint is called, **Then** the caller receives a service-unavailable style response or controlled fallback—not a stack trace or silent empty success that implies no data exists when the index was never queried.
3. **Given** fallback is active for users only, **When** the admin searches for a role name, **Then** role hits are omitted (not falsely reported as “no results” without explanation) unless full search is restored.
4. **Given** search is restored after an outage, **When** outbox backlog is processed, **Then** index and PostgreSQL reconcile without manual per-record intervention for standard mutations.

---

### Edge Cases

- **Index lag**: A record was just updated; search may briefly return stale title until outbox entry is processed—UI should tolerate slight delay; acceptance tests use a bounded sync window (see Success Criteria).
- **Outbox failures**: Transient worker errors retry; permanently failing entries remain in a failed state for operator visibility and manual replay.
- **Duplicate processing**: Same mutation enqueued twice or worker retry—index must remain idempotent (upsert by entity key, delete by entity key).
- **Search service down**: Degraded mode or explicit unavailable response; no permission leaks via error messages.
- **Empty query**: Submitting search with blank or whitespace-only query returns validation error or empty result per product rule (default: empty result, no full index scan).
- **No results**: Distinct from error—clear empty state in UI.
- **Caller lacks all entity permissions**: Empty hit list with successful response when `search.view` is present.
- **Reindex during live traffic**: Concurrent mutations via outbox; final state matches PostgreSQL after backlog clears.
- **Soft-deleted users**: Must not appear in index or search results after sync.
- **Sensitive fields**: Password hashes, tokens, and secrets never indexed even if present in source rows.

## Requirements *(mandatory)*

### Functional Requirements

#### Index document model

- **FR-001**: System MUST maintain a unified logical **SearchDocument** for each searchable entity instance (v1: `user`, `role`, `permission`) keyed by `entityType` + `entityId`.
- **FR-002**: Each SearchDocument MUST include: `entityType`, `entityId`, `title` (primary display label), `searchableText` (concatenated normalized text used for matching), `permissionKeys` (list of RBAC keys required to view this hit—typically one view key per entity type), `metadata` (non-sensitive structured attributes for display, e.g., status, group, email domain), and `updatedAt` (source record timestamp).
- **FR-003**: SearchDocument schema MUST reserve optional **semantic fields** (e.g., embedding vector placeholder, summary text slot) for future AI/RAG use; v1 MUST NOT populate or query these fields.
- **FR-004**: System MUST NOT index sensitive fields: password hashes, refresh tokens, OAuth secrets, API keys, or raw credentials.
- **FR-005**: Soft-deleted or inactive users (per existing domain rules) MUST be excluded from the index after sync.
- **FR-006**: Schema MUST allow adding new `entityType` values in future features without breaking existing documents (extensible entity type enumeration).

#### Outbox sync (PostgreSQL → search index)

- **FR-007**: System MUST record a **SearchOutboxEntry** in PostgreSQL when a searchable entity is created, updated, or deleted (same transaction as the business mutation where feasible).
- **FR-008**: Each outbox entry MUST capture: entity type, entity id, operation (`upsert` or `delete`), payload version or timestamp, status (`pending`, `processing`, `completed`, `failed`), attempt count, last error (if any), and created/updated timestamps.
- **FR-009**: A background worker MUST process pending outbox entries with **at-least-once** semantics and retry transient failures with backoff.
- **FR-010**: Index writes MUST be **idempotent**: upsert overwrites by entity key; delete removes by entity key and succeeds if already absent.
- **FR-011**: PostgreSQL MUST remain the **source of truth**; the search index is eventually consistent and MUST be rebuildable from PostgreSQL via reindex.
- **FR-012**: System MUST support **full and scoped reindex** (all v1 types or single entity type) invokable by an authorized operator or admin-only maintenance path.

#### Search API

- **FR-013**: System MUST expose an authenticated **admin unified search** endpoint under `/api/admin` requiring JWT and RBAC.
- **FR-014**: Endpoint MUST accept query parameters: `q` (search text), `types` (optional comma-separated entity types filter), `page` (1-based), `pageSize` (bounded max, default reasonable for admin UI).
- **FR-015**: Response MUST use **camelCase** JSON with: `hits` (array of SearchHit), `pagination` (`page`, `pageSize`, `total`, `totalPages`), and optional `degraded` flag when fallback or partial service is active.
- **FR-016**: Each **SearchHit** MUST include: `entityType`, `entityId`, `title`, `snippet` or highlighted excerpt (when available), `metadata` (safe display fields), and `updatedAt`.
- **FR-017**: Search MUST apply text matching across `title` and `searchableText` (and metadata fields explicitly marked searchable in v1 mapping requirements during design phase).
- **FR-018**: When search service is healthy, endpoint MUST NOT use SQL fallback except as explicitly defined in FR-022.

#### RBAC and security

- **FR-019**: System MUST add permission **`search.view`** (“Use admin unified search”) following the existing `{resource}.{action}` pattern in the permission catalog and seeders.
- **FR-020**: Calling the unified search endpoint MUST require **`search.view`** in addition to standard admin authentication.
- **FR-021**: Each hit MUST be filtered server-side: caller MUST hold the document’s required view permission (`users.view`, `roles.view`, or `permissions.view`) before the hit is included in the response.
- **FR-022**: When search is disabled or unavailable, system MAY fall back to existing database search for **users only** (email and full name, matching current list behavior) for callers with `users.view`; roles and permissions MUST NOT be fabricated or leaked via fallback.
- **FR-023**: Reindex and outbox replay operations MUST require elevated permission (default: **`search.modify`** or existing admin role assignment policy—see Assumptions) and MUST be audit-logged when audit infrastructure exists.
- **FR-024**: Error responses MUST NOT reveal existence of entities the caller cannot view (no “permission denied on user X” in search context).

#### Observability and operations

- **FR-025**: System MUST expose metrics or admin-visible indicators for outbox backlog size, oldest pending entry age, and failed entry count.
- **FR-026**: Reindex operations MUST report progress or completion status (started, completed, failed) suitable for operator confirmation.
- **FR-027**: Failed outbox entries MUST be inspectable and replayable without re-mutating source data.

#### Frontend (admin UI)

- **FR-028**: Admin UI MUST provide a unified search entry point (global search bar or dedicated admin search page) visible only to users with `search.view`.
- **FR-029**: UI MUST display hits with entity type badge, title, snippet, and navigation to the existing admin detail or list route for that entity.
- **FR-030**: UI MUST handle loading, empty, error, and degraded states without breaking the admin shell.
- **FR-031**: UI copy MUST be **English** in source files for v1.

### Out of Scope (v1)

- Vector embeddings, semantic/vector search, and similarity ranking.
- AI chat, RAG, or natural-language Q&A over the index.
- Indexing entities beyond **users**, **roles**, and **permissions**.
- Public (non-admin) search endpoints.
- Cross-tenant or multi-tenant search isolation (single-tenant admin assumed).
- Search analytics, click tracking, and query autosuggest (may follow in later features).

### Key Entities

- **SearchDocument**: Logical index record representing one searchable entity; keyed by `entityType` + `entityId`; carries display fields, searchable text, permission keys for ACL filtering, optional reserved semantic slots, and `updatedAt`.
- **SearchOutboxEntry**: Durable sync queue row in PostgreSQL linking a source mutation to an index operation (`upsert`/`delete`) with processing status and retry metadata.
- **SearchResult**: API-level paginated container holding hits, pagination metadata, and optional degradation flag.
- **SearchHit**: One permission-filtered row in a SearchResult with enough data to render and navigate in the admin UI.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Admin users with appropriate permissions can find a target record via unified search in **under 3 seconds** end-to-end (UI submit to rendered results) for typical queries on a dev-scale dataset (&lt; 10k total indexed documents).
- **SC-002**: **95%** of search API requests return within **1 second** when the search service is healthy and the outbox backlog is below the operational threshold (see Assumptions).
- **SC-003**: After a standard create/update/delete, the search index reflects the change within **30 seconds** under normal load (outbox processed, no manual reindex).
- **SC-004**: **Zero** permission-leak scenarios in QA test matrix: users without entity view permission never receive hits for that entity type in API or UI.
- **SC-005**: Full reindex of a dev-scale dataset completes successfully and produces searchable records for all active users, roles, and permissions without manual intervention.
- **SC-006**: When search is disabled, admin users receive a predictable degraded experience (explicit indicator or documented fallback) with **no unhandled 500 errors** on the search path in manual QA.

## Assumptions

- **Single-tenant admin search** for v1; all indexed data belongs to one application tenant.
- **Admin-only**: Unified search is not exposed to unauthenticated or non-admin member routes in v1.
- **English UI** copy in admin search surfaces; i18n may follow later.
- **Docker dev stack** will include a search engine service container alongside postgres, redis, be, fe, and nginx (exact product choice deferred to design phase).
- Existing RBAC middleware (`RequireView` / permission checks on admin routes) remains the authority model; search adds `search.view` plus per-entity view keys for result filtering.
- New permission **`search.modify`** is assumed for reindex/replay operations (parallel to `users.modify` pattern); if not seeded, reindex is restricted to the Admin role only until catalog is updated.
- Default pagination: `pageSize` default **20**, maximum **100**.
- Empty or whitespace-only `q` returns an empty hit list without querying the full index.
- Operational threshold for “healthy” sync: outbox backlog **&lt; 100** pending entries and oldest pending **&lt; 60 seconds** under dev load.
- Audit logging for reindex uses existing or planned admin audit mechanism; absence of audit in baseline does not block search but FR-023 must be satisfied when audit exists.
- Current user list SQL search (email, full name) remains the reference behavior for user fallback.

## Current Baseline (context for design phase)

| Area | Status |
|------|--------|
| Admin CRUD APIs for users, roles, permissions | Implemented with JWT + RBAC |
| RBAC permissions | `dashboard.view`, `users.view/modify`, `roles.view/modify`, `permissions.view` |
| User list text filter | SQL `LIKE` on email and full_name only |
| Search engine in Docker stack | **Not present** |
| Unified admin search UI | **Not implemented** |
| Outbox / search index sync | **Not implemented** |

*Detailed gap analysis, index mapping, and service wiring belong in phase 3 (`plan.md`, `be-implement.md`, `fe-implement.md`) — not in this spec.*

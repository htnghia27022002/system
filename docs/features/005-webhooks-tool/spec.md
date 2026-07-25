# Feature Specification: Tools → Webhooks (Per-Account Public URL + Request Inbox)

**Feature ID**: `005-webhooks-tool`

**Feature Branch**: `005-webhooks-tool`

**Created**: 2026-07-25

**Status**: Draft — ready for tasks

**Input**: Implement the real Webhooks tool under Tools: each account gets one public capture URL; inbound HTTP requests are recorded in an inbox; owner UI inspired by webhook.site (list + detail); soft-delete from the list; lifetime vs active totals; replace the Phase-1 shell from `003-landing-page`.

## Clarifications

### Session 2026-07-25 (user request — encoded defaults)

- Q: Scope? → A: **Real Webhooks product** for Tools: one public URL per account, capture inbound requests, owner inbox UI. Replaces the placeholder shell at `/tools/webhooks` from `003-landing-page`.
- Q: How many URLs per account? → A: **Exactly one** inbox / public URL per account. Auto-created on **first open** of the owner tool page (no manual “create inbox” step in P1).
- Q: Owner UI route vs public capture URL? → A: **Split-by-segment IA** (same path family):
  - **Owner inbox UI**: `/tools/webhooks` — signed-in account owner only.
  - **Public capture URL**: `/tools/webhooks/{uuid}` — unauthenticated; any supported HTTP method is captured into that account’s inbox.
  - Guests hitting `/tools/webhooks` (no UUID) are redirected to login.
  - Hitting `/tools/webhooks/{uuid}` never shows another user’s inbox; it only accepts/captures the request and returns a simple capture acknowledgment (not the owner UI).
  - **Tradeoff**: Same path family is easy to explain and copy (“your webhook URL looks like the tool path + id”). Routing must distinguish “no UUID → owner UI” from “UUID present → public capture.” Alternative (`/hooks/{uuid}` or API-only ingest) keeps FE and capture fully separate but is deferred unless planning finds the shared path impractical.
- Q: Soft delete? → A: Owner can **remove a request from the visible inbox** via **soft delete** (hidden from list/detail for normal use). Soft-deleted rows are **not** hard-wiped for accounting; they still count toward **lifetime received**.
- Q: Totals model? → A: Show **both**:
  - **Lifetime received**: all requests ever captured for the inbox (includes soft-deleted; excludes only hard-purged by retention).
  - **Active (visible) count**: requests not soft-deleted (what the inbox list shows before search filters).
  - UI may present as `showing / active` with lifetime nearby (e.g. “12 active · 40 lifetime”) — exact chrome is UX, not pixel-locked to webhook.site.
- Q: Who can manage the inbox? → A: Only the **owning signed-in account**. Public callers may only hit the capture URL. No shared/team inbox in P1. No admin override required for P1.
- Q: Regenerate URL? → A: **P1** — owner can **regenerate** the UUID. Old public URL immediately stops capturing; new URL is issued. Existing captured requests remain in the inbox (tied to the account/inbox, not invalidated by regenerate).
- Q: Copy URL? → A: **P1** — clear affordance to **copy** the current public capture URL.
- Q: Captured HTTP methods? → A: **GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS** (and treat other methods as capture-eligible if received, when practical).
- Q: What is stored per request? → A: Method, full URL (as received), timestamp, client IP (best effort), headers, query string, form fields when present, and raw body/content (subject to size limit).
- Q: Body size limit? → A: **1 MiB** maximum captured body. Oversized bodies are rejected or truncated with a clear capture status (prefer **reject with error response** and still record a truncated/error marker in the inbox if feasible; at minimum do not store unbounded bodies).
- Q: Sort / search / pagination? → A: Inbox list is **newest-first**. P1 includes **search/filter** (at least by method and free-text over URL/path/snippet) and **pagination** or equivalent “load more” so large inboxes remain usable.
- Q: Retention? → A: **P1 default**: keep at most **200** stored requests per inbox (including soft-deleted). When over cap, **hard-purge oldest** records first. Soft-delete is for viewing hygiene; retention hard-purge affects lifetime counting only for purged rows. **P2**: configurable limit, export, longer retention.
- Q: Real-time updates? → A: **P1** = refresh-friendly inbox (manual refresh and/or periodic refresh is enough). Live push/WebSocket is **P2**.
- Q: Custom response from capture URL? → A: **P1** = fixed simple success acknowledgment (status + short body). Custom status/headers/body editor is **P2**.
- Q: Guests / marketing? → A: Catalog/Tools hub may still link to `/tools/webhooks`; guests are sent to login before using the inbox. Public capture URLs remain usable without login (senders do not need accounts).
- Q: Phasing? → A: See Phase table below.

## Overview

Signed-in users open **Tools → Webhooks** and get a personal public URL that third parties can call. Every inbound request to that URL is stored and browsable in a webhook.site-inspired inbox (list + detail). Owners can copy or regenerate the URL, soft-delete noise from the list, and see both **active** and **lifetime** counts.

| Surface | Path | Auth | Behavior |
|---------|------|------|----------|
| Owner inbox UI | `/tools/webhooks` | Required (account owner) | Inbox list + detail; copy/regenerate URL; soft-delete |
| Public capture | `/tools/webhooks/{uuid}` | None | Capture request into owning inbox; return acknowledgment |

**Depends on**: `001-auth` (accounts / signed-in sessions). **Extends**: `003-landing-page` Tools IA (replaces Webhooks placeholder shell with the real tool).

### Product roadmap / phased delivery

| Phase | Scope | Status |
|-------|--------|--------|
| **Phase 1 (P1)** | Auto-create one inbox + UUID per account; public capture for listed methods; owner UI (list newest-first, search/filter, pagination, detail panes for URL/metadata/headers/query/form/body); copy URL; regenerate UUID; soft-delete; active + lifetime totals; 1 MiB body limit; retention cap 200 | **In scope** |
| **Phase 2 (P2)** | Live updates; custom capture responses; export; configurable retention; binary/multipart rich viewers; shared/team inboxes; rate-limit dashboards | **Deferred** |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Open Webhooks tool and receive a public URL (Priority: P1)

A signed-in user opens `/tools/webhooks`. If they have no inbox yet, one is created automatically with a unique public capture URL. They can copy that URL.

**Why this priority**: Without a stable public URL, the tool has no purpose.

**Independent Test**: Sign in, open `/tools/webhooks`, confirm a public URL containing a UUID is shown and can be copied; reload and confirm the same URL persists until regenerate.

**Acceptance Scenarios**:

1. **Given** a signed-in user with no prior Webhooks inbox, **When** they open `/tools/webhooks`, **Then** an inbox is created and a unique public capture URL of the form `/tools/webhooks/{uuid}` is displayed.
2. **Given** a signed-in user who already has an inbox, **When** they open `/tools/webhooks`, **Then** the same public URL is shown (idempotent; no duplicate inboxes).
3. **Given** the public URL is displayed, **When** the user activates Copy, **Then** the full public URL is placed on the clipboard (or equivalent copy success feedback).
4. **Given** a guest (not signed in), **When** they open `/tools/webhooks`, **Then** they are redirected to login and cannot see an inbox or URL belonging to any account.

---

### User Story 2 - Capture inbound HTTP requests on the public URL (Priority: P1)

An external client (or the owner testing with curl/browser) sends an HTTP request to the account’s public capture URL. The request appears in the owner’s inbox.

**Why this priority**: Capture is the core value of the tool.

**Independent Test**: Create/open inbox, send GET and POST (with headers, query, form, JSON body) to the public URL; as owner, confirm entries appear with correct method and detail fields.

**Acceptance Scenarios**:

1. **Given** a valid public capture URL, **When** a client sends GET, POST, PUT, PATCH, DELETE, HEAD, or OPTIONS, **Then** the system records a request entry for that inbox and returns a simple success acknowledgment without requiring authentication.
2. **Given** a captured POST with query string, headers, and JSON body under the size limit, **When** the owner opens that request in the inbox, **Then** they can inspect URL, method, timestamp, client IP (best effort), headers, query parameters, and body content.
3. **Given** a capture URL UUID that does not exist (or was invalidated by regenerate), **When** a client sends a request, **Then** capture fails with a clear not-found (or equivalent) outcome and nothing is added to any other inbox.
4. **Given** a request body larger than **1 MiB**, **When** it is submitted to the capture URL, **Then** the system does not store an unbounded body; the client receives an error (or documented truncate policy), and the owner is not left with a silently incomplete “success” body without indication.
5. **Given** a guest or any third party, **When** they only know the public capture URL, **Then** they can send requests to it but cannot list, view, or soft-delete the owner’s inbox.

---

### User Story 3 - Browse inbox list and inspect request detail (Priority: P1)

The owner sees a two-pane (or equivalent responsive) experience: left = request list; right = selected request detail. Layout and cues are inspired by webhook.site (method badge, id/snippet, IP, timestamp) without requiring a pixel-perfect clone.

**Why this priority**: Inspectability is how users debug integrations.

**Independent Test**: Capture several requests of mixed methods; open inbox; select items; confirm list metadata and detail sections; confirm newest-first order.

**Acceptance Scenarios**:

1. **Given** an inbox with multiple active requests, **When** the owner views the list, **Then** requests appear **newest first**, each showing at least method, a short id or path/snippet, approximate client IP, and timestamp.
2. **Given** the owner selects a request, **When** the detail pane loads, **Then** they see the request URL, metadata (method, time, IP), headers, query string, form values when present, and body/content.
3. **Given** an empty active inbox, **When** the owner opens the tool, **Then** they see an empty state that still shows the public URL and how to send a first request.
4. **Given** more active requests than one page, **When** the owner browses the list, **Then** pagination or load-more allows reaching older active items.
5. **Given** the owner enters a search/filter (method and/or free text), **When** results update, **Then** only matching **active** (non-soft-deleted) requests are listed.

---

### User Story 4 - Soft-delete a request from the list (Priority: P1)

The owner removes a request from the visible inbox to reduce clutter. The removal is a soft delete.

**Why this priority**: Explicit product requirement for manageable viewing without destroying accounting.

**Independent Test**: Soft-delete one of several requests; confirm it disappears from the list; confirm lifetime total still includes it; confirm active count decreases.

**Acceptance Scenarios**:

1. **Given** an active request in the list, **When** the owner removes/soft-deletes it, **Then** it no longer appears in the default inbox list or detail selection.
2. **Given** a soft-deleted request, **When** totals are shown, **Then** **lifetime received** still includes that request and **active count** excludes it.
3. **Given** soft-deleted requests, **When** the owner uses normal search on the inbox, **Then** soft-deleted items are not returned (no “trash” browser required in P1).

---

### User Story 5 - See active and lifetime totals (Priority: P1)

The inbox chrome shows how many requests are visible/active versus how many have been received over the inbox lifetime (subject to retention hard-purge).

**Why this priority**: User asked for a clear total model; dual counters avoid ambiguity.

**Independent Test**: Capture N requests, soft-delete K, confirm active = N−K (within retention) and lifetime = N (unless hard-purged).

**Acceptance Scenarios**:

1. **Given** an inbox that has received requests, **When** the owner views the tool, **Then** both **active (visible)** and **lifetime received** counts are available in the UI.
2. **Given** soft-deletes only (no retention purge), **When** comparing counts, **Then** lifetime ≥ active and lifetime does not decrease because of soft-delete.
3. **Given** retention hard-purge of oldest records, **When** those records are removed, **Then** lifetime and active counts no longer include the purged records.

---

### User Story 6 - Regenerate public URL (Priority: P1)

The owner regenerates the UUID so the previous public URL stops working. Captured history remains.

**Why this priority**: Cheap security/hygiene control when a URL was leaked; preferred in P1.

**Independent Test**: Note old URL; regenerate; send to old URL (fails); send to new URL (captures); confirm prior requests still listed.

**Acceptance Scenarios**:

1. **Given** an owner on `/tools/webhooks`, **When** they confirm regenerate, **Then** a new UUID/public URL replaces the old one and the UI shows the new URL.
2. **Given** the old URL after regenerate, **When** a client sends a request, **Then** it is not captured into the inbox.
3. **Given** regenerate completed, **When** the owner views the inbox, **Then** previously captured requests remain available (subject to soft-delete and retention).

---

### User Story 7 - Responsive / usable detail on smaller screens (Priority: P2)

On narrow viewports, list and detail remain usable (stacked navigation or equivalent), without blocking P1 desktop-oriented layout.

**Why this priority**: Nice-to-have polish after core capture/inbox works.

**Independent Test**: Open tool on a narrow viewport; select a request; confirm detail is reachable and readable.

**Acceptance Scenarios**:

1. **Given** a narrow viewport, **When** the owner uses the inbox, **Then** they can move between list and detail without losing access to copy-URL and totals.

---

### Edge Cases

- Concurrent captures while the owner has the inbox open → P1 may require manual/periodic refresh; new items appear after refresh.
- Extremely large header sets or many query keys → system still stores within reasonable limits; UI remains scrollable.
- HEAD/OPTIONS with empty body → still listed with method and metadata.
- Malformed body / binary content → stored as opaque/captured bytes or noted as non-text; P1 may show a safe preview or “binary/not previewable” state.
- Double soft-delete or delete of already-hidden item → idempotent success or clear no-op.
- Retention at exactly 200 → capturing the 201st hard-purges the oldest stored request.
- User opens `/tools/webhooks/{own-uuid}` in a browser → treated as **public capture** (records a GET), not as owner UI.
- Account deletion (if supported by product elsewhere) → inbox and captures are removed or orphaned per global account-deletion rules (out of scope to redefine here; assume cascade with account when that exists).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide exactly one Webhooks inbox per account, auto-created on first authenticated visit to `/tools/webhooks`.
- **FR-002**: System MUST expose a unique public capture URL path `/tools/webhooks/{uuid}` per inbox.
- **FR-003**: Owner UI MUST live at `/tools/webhooks` and MUST require authentication; guests MUST be redirected to login.
- **FR-004**: Public capture MUST accept unauthenticated requests for at least GET, POST, PUT, PATCH, DELETE, HEAD, and OPTIONS.
- **FR-005**: Each captured request MUST persist method, URL, timestamp, client IP (best effort), headers, query parameters, form fields when present, and body content subject to the size limit.
- **FR-006**: Captured request body size MUST be limited to **1 MiB**; oversize handling MUST be explicit (reject and/or mark truncated — no silent unbounded storage).
- **FR-007**: Owner MUST be able to list active requests newest-first with search/filter and pagination (or load-more).
- **FR-008**: Owner MUST be able to open a selected request and view URL, metadata, headers, query, form, and body/content.
- **FR-009**: Owner MUST be able to soft-delete a request so it is hidden from the active inbox without removing it from lifetime accounting.
- **FR-010**: UI MUST present **active (visible)** count and **lifetime received** count with the definitions in Clarifications.
- **FR-011**: Owner MUST be able to copy the current public capture URL.
- **FR-012**: Owner MUST be able to regenerate the inbox UUID; the previous URL MUST stop capturing; history MUST remain.
- **FR-013**: Only the owning account MUST list, view detail, soft-delete, copy management actions, or regenerate; public callers MUST NOT access inbox management.
- **FR-014**: System MUST retain at most **200** stored requests per inbox in P1 (including soft-deleted); oldest MUST be hard-purged when over cap.
- **FR-015**: Hitting `/tools/webhooks/{uuid}` MUST NOT render another user’s private inbox UI.
- **FR-016**: The Tools catalog entry for Webhooks MUST continue to navigate to `/tools/webhooks` (real tool, not “coming soon” placeholder for signed-in users).

### Key Entities

- **Webhook inbox**: Per-account capture endpoint identity; attributes include public UUID, created/updated time, lifetime received count, link to owning account.
- **Captured request**: One inbound HTTP call; attributes include method, URL, timestamps, client IP, headers, query, form data, body/preview, soft-deleted flag, relationship to inbox.
- **Capture acknowledgment**: Fixed simple response returned to public callers on successful (or failed) capture — not a customizable mock server in P1.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A signed-in user can open Webhooks and obtain a copyable public URL in under 1 minute on first use (including auto-create).
- **SC-002**: After sending a test request to the public URL, the owner sees that request in the inbox within 10 seconds on refresh (P1; no live-push requirement).
- **SC-003**: Soft-deleting a request removes it from the visible list while lifetime count remains correct in 100% of manual test cases in the QA matrix.
- **SC-004**: Regenerating the URL causes subsequent calls to the old URL to fail capture in 100% of trials, while prior history remains visible.
- **SC-005**: Guests cannot view any inbox data at `/tools/webhooks`; unauthenticated capture to a valid UUID still succeeds.
- **SC-006**: At least 90% of evaluators in an informal walkthrough can find method, headers, and body for a selected request without assistance.

## Assumptions

- Existing authentication (`001-auth`) identifies the account that owns the inbox.
- Tools hub and `/tools/webhooks` route shell already exist from `003-landing-page`; this feature replaces placeholder content with the real product.
- “Account” means the signed-in user account (one inbox per user account); organizations/teams are out of scope for P1.
- Client IP is best-effort behind proxies; exact proxy header strategy is a planning concern.
- P1 UI targets desktop-first webhook.site-like density; mobile stacked UX may be P2-quality.
- No requirement to replay captured requests outbound in P1.
- No requirement for signed webhooks, HMAC verification, or mutual TLS in P1.
- Independent deployability: capture and storage are owned by the backend service; the frontend owns the owner UI and calls the backend over HTTP; public capture may be served via the edge/routing layer to the backend without requiring the owner UI to be online.

## Out of Scope

- Custom programmable responses / full mock HTTP server (P2).
- WebSocket or other non-HTTP capture.
- Team/shared inboxes and ACL sharing.
- Export/download of inbox archives (P2).
- Billing or paid retention tiers.
- Pixel-perfect clone of webhook.site branding or layout.
- Changing global auth, landing marketing, or unrelated Tools catalog entries beyond Webhooks availability.
)

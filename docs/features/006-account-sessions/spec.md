# Feature Specification: Account Sessions (Active Sign-ins on Account Settings)

**Feature ID**: `006-account-sessions`

**Feature Branch**: `006-account-sessions`

**Created**: 2026-09-04

**Updated**: 2026-09-04

**Status**: Draft — ready for tasks

**Input**: Account sessions on the existing Account settings / profile page: signed-in users see their active sign-in sessions (other devices/browsers), identify which one is this device, revoke other sessions one-by-one, and revoke all other sessions at once. Signing out of this device uses the existing Sign out action (not a Revoke button on the current row).

## Clarifications

### Session 2026-09-04 (user request — locked P1)

- Q: Who can see sessions? → A: **Owner-only**. The signed-in account sees **only their own** sessions. No administrator viewing another user’s sessions. No new admin sidebar item. No extra permission catalog beyond **must be signed in**.
- Q: Where does this live? → A: **Fourth section** on the existing Account settings page at **`/admin/profile`**. No new route. Existing sections remain: avatar, personal information, change password, then **sessions**.
- Q: Which sessions appear? → A: The list shows only **active** sessions (not ended, not expired).
- Q: What does each row show? → A: When the session **started**, when it **expires**, **last activity**, **network address**, and **browser/device description**. In P1 the description MAY be the **raw user-agent string** — parsed friendly names such as “Chrome on Windows” are **not** required.
- Q: How is this device shown? → A: The **current session is highlighted**. The UI does **not** offer **Revoke** on the current row. Ending the current session is **Sign out**, not Revoke.
- Q: What does “revoke others” mean? → A: End **all sessions except the current one** in a single action, in addition to revoking **one other session at a time**.
- Q: How does the product know which session is current? → A: After sign-in and after **session renewal** (staying signed in without signing in again), the product MUST identify the current session **without putting secret session credentials in page addresses or shareable URLs** (including not putting secret refresh credentials in URLs).
- Q: What is recorded on each session? → A: **Network address** and **user-agent** when the session is **created or renewed**. **Device fingerprint IDs are out of P1**. **Last activity** updates when the session is **renewed** (not necessarily on every page view).
- Q: Does Sign out really end the session? → A: **Yes.** Admin chrome **Sign out** MUST actually end the **current session on the server**. Clearing only the device’s local signed-in state is **not** sufficient; the sessions list must stay truthful.
- Q: Phasing? → A: **P1** = own active-session list on `/admin/profile`, current-device highlight, revoke one other session, revoke all other sessions, truthful Sign out. **P2** (deferred) = parsed friendly device names, geography from network address, admin viewing another user’s sessions, a separate “sign out everywhere including this device” button, a dedicated sessions route.
- Q: Guests? → A: **N/A** — sessions require authentication; guests cannot open Account settings content.

## Overview

Authenticated users who can open **Account settings** (`/admin/profile`) see a **Sessions** section listing their **active sign-ins** across devices and browsers. They can tell which row is **this device**, end other sessions one at a time, or end every other session at once. Ending **this device** uses the existing **Sign out** control in admin chrome — not a Revoke action on the current row.

| Surface | Path | Auth | Behavior |
|---------|------|------|----------|
| Account settings — Sessions (4th section) | `/admin/profile` | Required (signed-in account) | List own active sessions; highlight this device; revoke others; revoke all others |
| Admin chrome Sign out | Existing Sign out in the admin user menu | Required | Ends the **current** session so it is no longer active |

**Depends on**: `001-auth` (sign-in, stay-signed-in / session renewal, Sign out) and `004-user-profile` (Account settings page at `/admin/profile`). **Does not** add a new admin sidebar item or a new permission catalog.

### Product roadmap / phased delivery

| Phase | Scope | Status |
|-------|--------|--------|
| **Phase 1 (P1)** | Fourth section on `/admin/profile`; list own **active** sessions; started / expires / last activity / network address / browser-device description (raw user-agent OK); highlight current session; no Revoke on current row; revoke one other session; revoke all other sessions; identify current session after login and renewal without secrets in URLs; record network address + user-agent on create/renew; last activity on renew; Sign out ends the current session so the list stays truthful | **In scope** |
| **Phase 2 (P2)** | Parsed friendly device names; geography from network address; admin viewing another user’s sessions; a separate “sign out everywhere including this device” button; a dedicated sessions route | **Deferred** |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - See active sessions on Account settings (Priority: P1)

A signed-in user opens **Account settings** (`/admin/profile`) and scrolls to a **Sessions** section below avatar, personal information, and change password. They see every **active** sign-in for **their own** account: when it started, when it expires, last activity, network address, and browser/device description.

**Why this priority**: Without a truthful list, users cannot audit where they are signed in.

**Independent Test**: Sign in, open Account settings, confirm a Sessions section is visible with the current sign-in listed; confirm guests cannot see it; confirm another account’s sessions never appear.

**Acceptance Scenarios**:

1. **Given** a signed-in user, **When** they open Account settings at `/admin/profile`, **Then** they see a **Sessions** section as the **fourth** section (after avatar, personal information, and change password) with no navigation to a new route.
2. **Given** the user has one or more active sessions, **When** the Sessions section loads, **Then** only **active** sessions for **this account** are listed (ended and expired sessions are omitted).
3. **Given** an active session row, **When** the user reads it, **Then** they see at least: started time, expiry time, last activity, network address, and browser/device description (raw user-agent is acceptable in P1).
4. **Given** a guest (not signed in), **When** they try to open `/admin/profile`, **Then** they are redirected to login and cannot view any session list.
5. **Given** user A is signed in, **When** they view Sessions, **Then** they never see sessions belonging to user B.

---

### User Story 2 - Identify this device (Priority: P1)

The user can tell which listed session is **this device**. That row is highlighted. The product still identifies the current session after a fresh sign-in and after staying signed in through session renewal, without putting secret session credentials in page addresses.

**Why this priority**: Users must not revoke the device they are using by accident, and the list is useless if “this device” is ambiguous.

**Independent Test**: Sign in on two browsers; on each, open Account settings and confirm only that browser’s row is marked as this device; stay signed in until the session is renewed and confirm the highlight remains correct.

**Acceptance Scenarios**:

1. **Given** a signed-in user with at least one active session, **When** they view the Sessions section, **Then** the session for **this device** is visually highlighted as current.
2. **Given** the current session row, **When** the user looks for a Revoke action on that row, **Then** **Revoke is not offered** (ending this device uses Sign out).
3. **Given** a successful sign-in, **When** the user opens Sessions, **Then** the product identifies the current session correctly without placing secret session credentials in the page address or a shareable URL.
4. **Given** the user stays signed in through **session renewal**, **When** they view Sessions, **Then** the same device is still identified as current without putting secret session credentials in the page address.

---

### User Story 3 - Revoke another session one at a time (Priority: P1)

The user finds a session that is **not** this device and ends it. That other sign-in stops being active. This device stays signed in.

**Why this priority**: Selective logout of a lost or unfamiliar device is the primary security action.

**Independent Test**: Sign in on two devices; from device A, revoke device B’s session; confirm B can no longer continue as signed in; confirm A remains signed in and B no longer appears as active.

**Acceptance Scenarios**:

1. **Given** at least one other active session besides this device, **When** the user chooses Revoke on that other row and confirms, **Then** that session is ended and disappears from the active list, and this device remains signed in.
2. **Given** the user starts Revoke on another session, **When** they cancel confirmation, **Then** no session is ended and the list is unchanged.
3. **Given** another device whose session was just revoked, **When** that device tries to continue as signed in, **Then** it cannot remain signed in and must sign in again.
4. **Given** the target session is already ended or expired, **When** the user attempts Revoke, **Then** they see a clear message and the list refreshes to match reality (no silent no-op that still shows the row as active).

---

### User Story 4 - Revoke all other sessions at once (Priority: P1)

The user ends **every active session except this device** in one action. This device stays signed in. Other devices must sign in again.

**Why this priority**: A single “sign out everywhere else” action is an explicit P1 product decision for stolen-device or shared-computer recovery.

**Independent Test**: Sign in on three devices; from device A, revoke all others; confirm A remains signed in and B and C cannot continue; confirm only A remains in the active list.

**Acceptance Scenarios**:

1. **Given** two or more other active sessions besides this device, **When** the user confirms **revoke all other sessions**, **Then** every session except the current one is ended, and this device remains signed in.
2. **Given** only the current session is active (no other sessions), **When** the user views Sessions, **Then** revoke-all is unavailable or clearly a no-op that does not sign them out of this device.
3. **Given** the user starts revoke-all, **When** they cancel confirmation, **Then** no other session is ended.
4. **Given** other devices whose sessions were ended by revoke-all, **When** those devices try to continue as signed in, **Then** they cannot remain signed in and must sign in again.

---

### User Story 5 - Sign out of this device ends the current session (Priority: P1)

The user ends **this device** with the existing admin chrome **Sign out** action (not Revoke on the current row). Sign out actually ends the current session so it no longer appears as active if the user later inspects sessions from another remaining sign-in.

**Why this priority**: The list is untrustworthy if Sign out only clears local state while the session remains active.

**Independent Test**: Sign in on two devices; Sign out on device A via admin chrome; from device B open Sessions and confirm A is gone; confirm A cannot continue as signed in without signing in again.

**Acceptance Scenarios**:

1. **Given** a signed-in user, **When** they use admin chrome **Sign out**, **Then** the **current session is ended** (not merely forgotten on this device) and they are no longer signed in on this device.
2. **Given** the user has another remaining active session on a second device, **When** they Sign out on the first device and then open Sessions on the second, **Then** the first device’s session is **not** listed as active.
3. **Given** the current session row on `/admin/profile`, **When** the user looks for a way to end this device from that row, **Then** they are not offered Revoke; ending this device remains **Sign out**.
4. **Given** Sign out completed, **When** the same device tries to continue as signed in, **Then** they must sign in again.

---

### User Story 6 - Session rows record network address, user-agent, and last activity (Priority: P1)

When a session is created or renewed, the product records **network address** and **user-agent** on that row. **Last activity** updates when the session is **renewed**. Device fingerprint IDs are not part of P1.

**Why this priority**: Users need enough context to recognize unfamiliar sign-ins; capture timing is an explicit lock.

**Independent Test**: Sign in and confirm the new row has network address and user-agent; stay signed in until renewal and confirm last activity (and recorded address/user-agent) update; confirm no fingerprint identifier is required to use the list.

**Acceptance Scenarios**:

1. **Given** a new sign-in, **When** the session is created, **Then** that session row records network address and user-agent (browser/device description may show the raw user-agent).
2. **Given** an existing session, **When** it is **renewed**, **Then** last activity is updated, and network address and user-agent are recorded from that renewal.
3. **Given** P1, **When** a user views Sessions, **Then** they can use the list without any device fingerprint identifier being collected or shown.
4. **Given** network address or user-agent cannot be determined, **When** the row is shown, **Then** the field shows a clear unavailable/unknown state rather than inventing a value.

---

### User Story 7 - Deferred session enhancements (Priority: P2)

After P1, the product MAY add parsed friendly device names, geography from network address, administrator viewing of another user’s sessions, a separate “sign out everywhere including this device” control, and/or a dedicated sessions route — without changing P1 routes or owner-only visibility unless a later spec says so.

**Why this priority**: Explicitly deferred; not required to deliver P1.

**Independent Test**: P2 is out of P1 acceptance; when scheduled, each enhancement is tested as its own story.

**Acceptance Scenarios**:

1. **Given** P1 capabilities ship, **When** P2 work is later applied, **Then** `/admin/profile` remains a valid place to manage own sessions unless a later spec replaces it.
2. **Given** P1, **When** stakeholders ask for parsed device names, geo, admin-of-another-user, sign-out-everywhere-including-this-device, or a dedicated route, **Then** those are treated as **out of P1 scope** (this story), not as P1 defects.

---

### Edge Cases

- Only the current session exists → list shows one highlighted row; Revoke is not offered on it; revoke-all is unavailable or a harmless no-op that does not sign the user out.
- User revokes the last other session → list then shows only this device.
- Session ends or expires while the user is looking at the list → next load/refresh omits it; a revoke attempt on a stale row explains that it is no longer active.
- Two devices revoke overlapping sessions at the same time → each successful revoke ends its target; the remaining list matches what is still active; no other account’s sessions are affected.
- Sign-in from the same browser profile after Sign out → a **new** active session is created; the previous ended session does not reappear as active.
- Session renewal on this device → current highlight remains on this device; last activity (and recorded network address / user-agent) update for this session.
- Missing network address or user-agent → show unknown/unavailable; do not block listing or revoke.
- Very long raw user-agent → still shown (may wrap or truncate visually) without requiring parsed vendor/OS names in P1.
- Unauthenticated access → existing Account settings rule: redirect to login; no session data exposed.
- User without permission to open Account settings beyond being signed in → **none required**; any authenticated user who can open `/admin/profile` can use Sessions.
- Administrator signed in as themselves → they see **their own** sessions only; they cannot open another user’s session list in P1.
- Ended sessions and expired sessions → never listed in P1.
- User confirms revoke-all while only this device is active → this device stays signed in.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST show a **Sessions** section as the **fourth** section on the existing Account settings page at `/admin/profile` (after avatar, personal information, and change password). The feature MUST NOT introduce a new route in P1.
- **FR-002**: The Sessions section MUST be available to **any signed-in user** who can open Account settings. The feature MUST NOT add a new admin sidebar item and MUST NOT add a new permission catalog beyond requiring a signed-in session.
- **FR-003**: A signed-in user MUST see **only their own** sessions. The system MUST NOT allow viewing another user’s sessions in P1 (including by administrators).
- **FR-004**: The list MUST include only **active** sessions. Ended and expired sessions MUST NOT appear.
- **FR-005**: Each listed session MUST show: time it **started**, time it **expires**, **last activity**, **network address**, and **browser/device description**.
- **FR-006**: In P1 the browser/device description MAY be the **raw user-agent string**. Parsed friendly names (for example “Chrome on Windows”) MUST NOT be required.
- **FR-007**: The system MUST highlight the **current session** (this device) in the list.
- **FR-008**: The UI MUST NOT offer **Revoke** on the current session row. Ending the current session MUST use the existing **Sign out** action.
- **FR-009**: Authenticated users MUST be able to **revoke one other** (non-current) active session after confirmation, ending only that session.
- **FR-010**: Authenticated users MUST be able to **revoke all other** active sessions at once after confirmation, ending every session except the current one.
- **FR-011**: After a successful single revoke or revoke-all, this device MUST remain signed in, and ended sessions MUST no longer appear as active.
- **FR-012**: After sign-in and after **session renewal**, the product MUST identify which listed session is current **without putting secret session credentials in page addresses or shareable URLs**.
- **FR-013**: When a session is **created or renewed**, the system MUST record **network address** and **user-agent** on that session.
- **FR-014**: **Last activity** on a session MUST update when that session is **renewed**. Device fingerprint identifiers MUST NOT be required in P1.
- **FR-015**: Admin chrome **Sign out** MUST **end the current session** so it is no longer active. Clearing only this device’s local signed-in state without ending the session is NOT sufficient.
- **FR-016**: Unauthenticated access to Account settings / session list MUST redirect to login and MUST NOT expose session rows.
- **FR-017**: Revoke and revoke-all MUST ask for confirmation before ending sessions; cancelling MUST leave sessions unchanged.
- **FR-018**: If network address or user-agent is unavailable, the row MUST show an unknown/unavailable state rather than a fabricated value.
- **FR-019**: If a revoke targets a session that is no longer active, the system MUST tell the user clearly and refresh the list to match remaining active sessions.
- **FR-020**: Guests have no sessions capability (N/A).

### Key Entities

- **Account (signed-in user)**: The owner of sessions; sees and manages only their own active sign-ins from Account settings.
- **Active session**: A current sign-in that has not been ended and has not expired. Attributes: started time, expiry time, last activity, network address, user-agent / browser-device description, and whether it is the **current** session for this device.
- **Current session (this device)**: The active session belonging to the device/browser the user is using now. Highlighted; not revocable from the list; ended only by **Sign out**.
- **Other session**: An active session for the same account that is not this device. May be revoked individually or as part of revoke-all.
- **Session renewal**: The product keeping the user signed in without a new sign-in; used to refresh last activity and recorded network address / user-agent, and to keep current-session identity accurate.
- **Sign out**: Existing admin chrome action that MUST end the current session so it is no longer active.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A signed-in user can open Account settings and see the Sessions section (fourth section on `/admin/profile`) in under 5 seconds of intentional navigation on a stable connection.
- **SC-002**: In 100% of tested cases with two or more active devices, the user can identify **this device** on the list without using Revoke on that row.
- **SC-003**: Revoking one other session ends only that session: this device stays signed in, and the revoked device cannot continue as signed in, in 100% of tested cases.
- **SC-004**: Revoke-all-other-sessions ends every active session except the current one, and this device stays signed in, in 100% of tested cases with at least two other sessions.
- **SC-005**: After Sign out on this device, that session is no longer listed as active from any remaining sign-in, and this device cannot continue without signing in again, in 100% of tested cases.
- **SC-006**: Unauthenticated users cannot view any session list (redirect to login).
- **SC-007**: A signed-in user never sees another account’s sessions in 100% of tested cross-account checks.
- **SC-008**: 100% of tested new sign-ins show started time, expiry, last activity, network address (or unknown), and browser/device description (raw user-agent acceptable); last activity changes after a tested session renewal.

## Assumptions

- Authentication (`001-auth`) remains the source of sign-in, stay-signed-in / session renewal, and Sign out entry points; this feature makes session listing and revocation truthful rather than replacing login.
- Account settings (`004-user-profile`) remains at `/admin/profile` with existing avatar, personal information, and change-password sections; Sessions is added as the fourth section.
- Any authenticated user who can open Account settings can use Sessions; **no new permission keys** and **no new sidebar item**.
- Session counts per person are typically small; P1 lists all **active** sessions without requiring pagination.
- List order default: **current session first**, then other sessions by **last activity** (most recent first).
- Revoke and revoke-all require an explicit confirmation step (destructive security action).
- “Network address” is the client address associated with create/renew (commonly an IP address); geo-location from that address is P2.
- Raw user-agent is an acceptable browser/device description in P1; wrapping or visual truncation is allowed if the full string remains available to the user (for example via expansion or title text).
- Last activity does **not** need to update on every page view — only on **session renewal**, per the locked decision.
- If Sign out today sometimes only clears this device’s local signed-in state, that is a **gap this feature must close**, not an acceptable P1 behavior.
- UI copy in source defaults to English; Vietnamese appears via existing i18n when locale is VI.
- P2 items do not block declaring P1 done for handoff to architecture/implementation.

## Out of Scope

- Parsed friendly device names (for example “Chrome on Windows”) — P2.
- Geography or map location derived from network address — P2.
- Administrators viewing or revoking **another user’s** sessions — P2.
- A separate control that signs out **everywhere including this device** in one button (P1 uses revoke-all-others plus existing Sign out) — P2.
- A dedicated sessions route (sessions stay on `/admin/profile`) — P2.
- Device fingerprint identifiers.
- Listing ended or expired sessions (history / audit log of past sign-ins).
- New admin sidebar navigation for sessions.
- Extra permission catalog entries beyond being signed in.
- Changing avatar, personal fields, or password (remain `004-user-profile`).
- Public or guest session management.

## Dependencies

- **001-auth**: Sign-in, session renewal (stay signed in), current-user identity, Sign out. This feature requires Sign out to end the current session and requires a way to identify the current session after login and renewal without secrets in URLs.
- **004-user-profile**: Account settings page at `/admin/profile` (avatar, personal information, change password) where Sessions is the fourth section.
- Existing admin chrome **Sign out** in the user menu: same entry point; behavior MUST end the current session.

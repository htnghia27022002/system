# Feature Specification: User Profile, Access-Control User Edit Alignment, Admin Locale Flags

**Feature ID**: `004-user-profile`

**Feature Branch**: `004-user-profile`

**Created**: 2026-07-25

**Updated**: 2026-07-25 (expanded personal fields)

**Status**: Draft — ready for tasks

**Input**: Post-login profile page (edit personal info, upload avatar, change password); align access-control user create/edit with the same class of personal fields; wire admin sidebar “Account settings” to the profile page; replace admin user-menu language switch text with clickable EN/VI flags consistent with public locale selection. Expanded personal fields: phone, General (bio/about), birthday, address, social link(s).

## Clarifications

### Session 2026-07-25 (user request — encoded)

- Q: Build what first after login? → A: **Self-service profile page** with edit personal info, avatar upload, and change password.
- Q: Access-control relationship? → A: Admin **user create/edit form** must support editing the **same class of personal fields** (name, and avatar where in scope), with UX aligned to profile where it makes sense (admin editing another user vs self-service).
- Q: Navigation entry? → A: Admin sidebar bottom **user menu** item labeled **Account settings** (currently disabled) navigates to the profile page.
- Q: Language control in that menu? → A: Replace text toggle (“Switch to Vietnamese” / switch pattern) with **clickable flags** for **EN** and **VI**; selecting applies language immediately and shows which locale is active. Align with public chrome flag-based locale selection behavior.
- Q: Route? → A: **`/admin/profile`** — profile is reached from the admin shell user menu (today’s entry point). Guests cannot access profile.
- Q: Profile fields (P1)? → A: **Full name** (editable), **email** (displayed, **read-only** on self-service profile — changing email is out of this feature’s auth scope), **avatar** image, **change password** (current + new + confirm).
- Q: Avatar rules? → A: Image upload with **size and type limits**; show current avatar or **initials fallback** when no image.
- Q: Admin password for another user? → A: Admin **set/reset password** remains a **separate optional action** on the access-control user form (leave blank to keep current password on edit; required on create) — not the same as self-service “change password” which requires the current password.
- Q: Admin avatar? → A: Access-control create/edit MAY include **optional avatar** upload/change for the managed user in P1 when practical; same display rules (image or initials).
- Q: Phasing? → A: **P1** = profile + Account settings nav + admin locale flags + access-control personal-field/avatar alignment. **P2** = polish (copy, empty/error states, denser UX refinements) without changing core flows.
- Q: Guests? → A: **N/A** — profile requires authentication.

### Session 2026-07-25 (expanded personal fields)

- Q: Additional profile personal fields? → A: Add optional **Phone**, **General** (short free-text bio/about; UI label **General**), **Birthday**, **Address**, and **Social link(s)** alongside existing name, read-only email, avatar, and change password.
- Q: What is “General”? → A: Optional multiline **bio/about** text; field name and preferred UI label: **General**.
- Q: Address structure? → A: **P1 = single optional text field** (full address as one string). Structured multi-field address (street/city/postal/country) is out of scope for P1.
- Q: Social links model? → A: Optional **list** of entries, each `{ label?: string, url: string }` (label optional, URL required per entry). **Max 5** entries. Empty list is allowed. Extensible later without fixed provider slots.
- Q: Field optionality? → A: All new fields are **optional**. Full name remains required (aligned with existing product rules). Email stays read-only on self-service profile.
- Q: Validation defaults? → A:
  - **Phone**: optional string; if provided, non-empty after trim; max length **50**; lenient format (no mandatory E.164 in P1).
  - **General**: optional multiline; max **1000** characters.
  - **Birthday**: optional date at day precision; **reject future dates**.
  - **Address**: optional single text; max **500** characters.
  - **Social links**: each URL MUST be a valid **http** or **https** URL; optional label max **50** characters; max **5** entries.
- Q: Access-control alignment? → A: Admin user create/edit MUST support the **same class** of personal fields: name, phone, General, birthday, address, social links, and optional avatar — where editing another user makes sense. Email/role/status remain admin-managed; admin password set/reset semantics unchanged.
- Q: Persistence note? → A: **Phone** and **avatar** already exist on the user model in the current system. Birthday, address, General, and social links may require **new persistence** — this is a requirements fact for planning; BA does not specify storage technology.

## Overview

Authenticated users who use the **admin shell** can manage their own account on a dedicated **profile** page and open it from **Account settings** in the sidebar user menu. The same menu switches language via **flag selection** (EN/VI). Administrators managing users in **access control** can edit aligned personal fields (including the expanded profile fields and optional avatar) without conflating that flow with self-service password change.

| Route | Role |
|-------|------|
| **`/admin/profile`** | Self-service profile for the signed-in user (admin shell). Requires login; unauthorized visitors are sent to login. |

### Personal field set (P1)

| Field | Self-service profile | Access-control create/edit | Notes |
|-------|----------------------|----------------------------|-------|
| Full name | Editable (required) | Editable (required) | Existing product rule |
| Email | Read-only | Editable (admin) | Email change out of profile scope |
| Avatar | Upload/replace | Optional upload/replace | Image or initials fallback |
| Phone | Optional | Optional | Lenient string; max 50 |
| General | Optional | Optional | Bio/about multiline; max 1000; UI label **General** |
| Birthday | Optional | Optional | Day precision; not in the future |
| Address | Optional | Optional | Single text string; max 500 |
| Social links | Optional list (0–5) | Optional list (0–5) | `{ label?, url }`; http(s) URLs |
| Change password | Self-service only | N/A | Current + new + confirm |
| Set/reset password | N/A | Create required; edit optional | Blank on edit = keep |
| Role / status | Not on profile | Admin-managed | Unchanged |

**Product roadmap / phased delivery**

| Phase | Scope | Status |
|-------|--------|--------|
| **Phase 1 (P1)** | Profile view/edit (name, read-only email, phone, General, birthday, address, social links, avatar, change password); Account settings → `/admin/profile`; admin user-menu flag locale (EN/VI); access-control user form aligned personal fields + optional avatar | **In scope** |
| **Phase 2 (P2)** | UX polish: refined empty/error/success messaging, upload progress affordances, accessibility polish, optional layout refinements — no new primary field capabilities required | **Deferred polish** |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Open profile from Account settings (Priority: P1)

A signed-in user opens the admin sidebar user menu and chooses **Account settings**. They land on `/admin/profile` and see their current personal information.

**Why this priority**: Entry point is explicitly requested; without it the profile page is undiscoverable from the current chrome.

**Independent Test**: Sign in as an admin-capable user, open the bottom sidebar user menu, click Account settings, confirm navigation to `/admin/profile` with current name, email, avatar/initials, and other personal fields (empty or populated) visible.

**Acceptance Scenarios**:

1. **Given** a signed-in user on an admin page, **When** they open the sidebar user menu and select Account settings, **Then** they navigate to `/admin/profile`.
2. **Given** a guest (not signed in), **When** they try to open `/admin/profile`, **Then** they are redirected to login and cannot view profile data.
3. **Given** a signed-in user on `/admin/profile`, **When** the page loads, **Then** they see their current full name, email, avatar image or initials fallback, and the personal fields (phone, General, birthday, address, social links) in their current state (including empty).

---

### User Story 2 - Edit personal information (Priority: P1)

A signed-in user updates personal profile fields on `/admin/profile` and saves. Editable fields: full name (required), phone, General, birthday, address, and social links. Email remains visible but not editable on the profile page.

**Why this priority**: Core self-service identity and contact update after login; expanded field set is an explicit product decision.

**Independent Test**: On `/admin/profile`, change name and optional phone / General / birthday / address / social links to valid values, save, reload, and confirm persistence; attempt invalid values (future birthday, bad URL, over-long General) and confirm validation blocks save.

**Acceptance Scenarios**:

1. **Given** a user on `/admin/profile`, **When** they change full name to a valid value and save, **Then** the update succeeds and the new name is shown after refresh.
2. **Given** a user on `/admin/profile`, **When** they view email, **Then** email is displayed and cannot be changed on this page.
3. **Given** a user on `/admin/profile`, **When** they enter an optional phone within max length and save, **Then** the phone is stored and shown after refresh; clearing phone and saving removes it.
4. **Given** a user on `/admin/profile`, **When** they enter General text within the max length and save, **Then** it persists; text over the max length is rejected with a clear message.
5. **Given** a user on `/admin/profile`, **When** they set a birthday that is today or in the past and save, **Then** the birthday persists at day precision.
6. **Given** a user on `/admin/profile`, **When** they attempt to save a birthday in the future, **Then** validation rejects the save with a clear message.
7. **Given** a user on `/admin/profile`, **When** they enter an optional address within max length and save, **Then** it persists as a single text value.
8. **Given** a user on `/admin/profile`, **When** they add up to five social links with valid http(s) URLs (optional labels) and save, **Then** the list persists; empty list is allowed.
9. **Given** a social link with a missing or non-http(s) URL, **When** the user attempts to save, **Then** validation rejects that entry with a clear message.
10. **Given** invalid name input (empty or too short per product rules), **When** they attempt to save, **Then** field-level validation prevents save and explains the problem in plain language.
11. **Given** a save failure (network or server error), **When** the user submits, **Then** they see a clear error and prior saved data is not silently lost from the form without feedback.

---

### User Story 3 - Upload or change avatar (Priority: P1)

A signed-in user uploads a new profile image (or replaces the current one). The UI shows the image or initials when no image exists. Invalid files are rejected with clear feedback.

**Why this priority**: Explicit product capability; avatar already exists as a user attribute conceptually.

**Independent Test**: Upload a valid small image on profile; confirm it displays; attempt an oversized or wrong-type file and confirm rejection messaging.

**Acceptance Scenarios**:

1. **Given** a user with no avatar, **When** they view profile, **Then** initials (derived from name) are shown as fallback.
2. **Given** a user on `/admin/profile`, **When** they upload a valid image within size and type limits, **Then** the avatar updates and remains visible after reload.
3. **Given** a file that exceeds size limits or is not an allowed image type, **When** upload is attempted, **Then** the system rejects it and shows a clear message without corrupting the previous avatar.
4. **Given** a user with an existing avatar, **When** they replace it with another valid image, **Then** the new image becomes the displayed avatar.

---

### User Story 4 - Change own password (Priority: P1)

A signed-in user changes their password by providing current password, new password, and confirmation. Rules match existing auth password strength expectations (minimum length and confirmation match).

**Why this priority**: Security-sensitive self-service; explicitly requested.

**Independent Test**: Change password with correct current password; sign out and sign in with the new password; attempt wrong current password and mismatched confirmation.

**Acceptance Scenarios**:

1. **Given** a user on `/admin/profile` change-password section, **When** they submit correct current password plus valid matching new password and confirmation, **Then** the password is updated and they can sign in with the new password afterward.
2. **Given** incorrect current password, **When** they submit, **Then** the change is rejected with a clear message and the password is unchanged.
3. **Given** new password and confirmation that do not match, or a new password below minimum length (aligned with auth: at least 8 characters), **When** they submit, **Then** validation errors are shown and no password change occurs.
4. **Given** a successful password change, **When** the form completes, **Then** the user receives clear success feedback (they are not forced to re-login in P1 unless product later requires it).

---

### User Story 5 - Admin edits user personal info in access control (Priority: P1)

An administrator creates or edits a user in access control and can set/update the same class of personal fields as profile: name, phone, General, birthday, address, social links, and optional avatar. Password on create is required; on edit, optional set/reset (leave blank to keep). Email remains manageable in access control as part of user administration (not read-only there).

**Why this priority**: User asked for access-control form alignment with profile personal fields, including the expanded set.

**Independent Test**: From access-control users UI, edit a user’s name and optional personal fields (phone, General, birthday, address, social links, avatar); optionally set a new password; confirm persistence; confirm leave-blank password keeps the old password; confirm the same validation rules as profile for personal fields.

**Acceptance Scenarios**:

1. **Given** an admin with permission to manage users, **When** they open create or edit user, **Then** they can edit full name, phone, General, birthday, address, social links, and (when in scope) avatar with the same display and validation rules as profile.
2. **Given** admin editing an existing user, **When** they leave password blank and save other personal fields, **Then** the user’s password is unchanged and personal field updates persist.
3. **Given** admin editing an existing user, **When** they provide a new valid password and save, **Then** that user’s password is set to the new value without requiring the user’s current password.
4. **Given** admin creating a user, **When** they submit without a password, **Then** validation requires a password.
5. **Given** admin create/edit, **When** they set email and role/status as today, **Then** those administrative fields remain available (profile self-service does not remove admin email/role/status management).
6. **Given** admin enters invalid personal data (e.g. future birthday or non-http(s) social URL), **When** they save, **Then** validation rejects with clear messages consistent with profile rules.

---

### User Story 6 - Locale via flags in admin user menu (Priority: P1)

In the admin sidebar user menu, language is chosen by clicking **EN** or **VI** flag options (not a single “Switch to …” text toggle). The active locale is visually indicated. Selection applies immediately, consistent with public locale flag selection behavior.

**Why this priority**: Explicit UX request; closes the gap between public chrome and admin menu.

**Independent Test**: Open admin user menu; select the non-active flag; confirm UI language updates immediately and the active flag state is clear; switch back.

**Acceptance Scenarios**:

1. **Given** a signed-in user with English active, **When** they open the user menu, **Then** they see flag-based choices for English and Vietnamese (not “Switch to Vietnamese” as the sole control).
2. **Given** the user menu is open, **When** they select the Vietnamese flag option, **Then** the admin UI language switches to Vietnamese immediately and Vietnamese is shown as active.
3. **Given** Vietnamese is active, **When** they select the English flag option, **Then** the UI switches to English immediately and English is shown as active.
4. **Given** public chrome already uses flag-based locale selection, **When** comparing admin menu locale UX, **Then** behavior is consistent: click to select EN or VI, immediate apply, clear active state.

---

### User Story 7 - Profile UX polish (Priority: P2)

After P1 flows work, refine profile and related forms for clearer success/error states, upload feedback, and accessibility without changing the core field set or routes.

**Why this priority**: Improves quality but is not required to deliver the primary capabilities.

**Independent Test**: Walk profile and access-control personal sections after polish; confirm messages and focus/labels meet accessibility baseline.

**Acceptance Scenarios**:

1. **Given** P1 capabilities ship, **When** P2 polish is applied, **Then** routes and field set remain the same while messaging and interaction feedback improve.
2. **Given** keyboard and screen-reader use, **When** operating profile and locale flag controls, **Then** controls have accessible names and errors are available as text.

---

### Edge Cases

- User with OAuth-only account (no password) attempting change password → clear message that password change is unavailable or guide to set password if product supports it; do not fail silently. **Default for P1**: if the account has no password credential, hide or disable change-password and explain briefly.
- Concurrent edit: admin changes a user’s personal fields while the user is on profile → last successful save wins; user sees updated values on next load.
- Avatar upload interrupted or fails mid-request → previous avatar (or initials) remains; error shown.
- Very long names → validation or truncation rules consistent with existing user name constraints (minimum length aligned with auth registration).
- Inactive user signed out cannot reach profile; admin may still see inactive users in access control per existing rules.
- Selecting the already-active locale flag → no disruptive reload required; remains active.
- Access-control user without permission → existing permission guards still block create/edit.
- Phone with only whitespace → treat as empty (clear) after trim; do not store whitespace-only.
- Birthday equal to today → allowed; tomorrow or later → rejected.
- Social links list with more than 5 entries → rejected; user must remove extras before save.
- Social entry with URL but blank label → allowed; label-only without URL → rejected.
- Clearing optional fields (phone, General, birthday, address, all social links) and saving → stores empty/absent values successfully.
- Duplicate social URLs → allowed in P1 (no uniqueness requirement).

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide an authenticated profile page at `/admin/profile` showing the current user’s full name, email, avatar (image or initials fallback), and personal fields (phone, General, birthday, address, social links) in their current state.
- **FR-002**: Authenticated users MUST be able to update their full name from the profile page (required field).
- **FR-003**: Profile MUST display email as read-only (email change out of scope for this feature).
- **FR-004**: Authenticated users MUST be able to upload or replace a profile avatar image subject to documented size and type limits; invalid uploads MUST be rejected with clear feedback.
- **FR-005**: When no avatar image exists, the UI MUST show an initials fallback derived from the user’s name.
- **FR-006**: Authenticated users who have a password credential MUST be able to change password by providing current password, new password, and confirmation; new password MUST meet the same minimum-length rule as registration (at least 8 characters) and confirmation MUST match.
- **FR-007**: Password change MUST reject incorrect current password without changing the stored password.
- **FR-008**: Accounts without a password credential MUST NOT be offered a broken change-password path (hide/disable with explanation).
- **FR-009**: Admin sidebar user menu **Account settings** MUST navigate to `/admin/profile` (no longer a disabled dead-end).
- **FR-010**: Unauthenticated access to `/admin/profile` MUST redirect to login.
- **FR-011**: Access-control user create/edit MUST allow editing the same personal field class as profile: name, phone, General, birthday, address, social links, and optional avatar, with the same validation and display rules.
- **FR-012**: Access-control create MUST require password; edit MUST allow optional password set/reset (blank = keep existing) without requiring the target user’s current password.
- **FR-013**: Access-control MUST continue to support administrative fields (email, role, status) for managed users; these are not removed because profile keeps email read-only for self-service.
- **FR-014**: Admin user menu MUST replace the text language switch with **flag click-to-select** options for English and Vietnamese.
- **FR-015**: Selecting a locale flag MUST apply that language immediately and visually indicate the active locale.
- **FR-016**: Admin locale flag behavior MUST be consistent with public chrome flag-based locale selection (EN/VI, immediate apply, clear active state).
- **FR-017**: Profile and access-control personal forms MUST show validation and error messages in plain language; successful saves MUST provide clear confirmation.
- **FR-018**: Guests have no profile capability (N/A).
- **FR-019**: Authenticated users and admins (via access control) MUST be able to set, update, or clear optional **phone** (string; max 50 characters; if provided after trim, non-empty; lenient format — no mandatory E.164 in P1).
- **FR-020**: Authenticated users and admins MUST be able to set, update, or clear optional **General** (bio/about multiline text; UI label **General**; max **1000** characters).
- **FR-021**: Authenticated users and admins MUST be able to set, update, or clear optional **birthday** (day-precision date); the system MUST reject birthdays strictly after today.
- **FR-022**: Authenticated users and admins MUST be able to set, update, or clear optional **address** as a single free-text string (max **500** characters). Structured multi-field address is out of scope for P1.
- **FR-023**: Authenticated users and admins MUST be able to manage optional **social links** as a list of zero to **five** entries; each entry MUST have a required **url** that is a valid **http** or **https** URL and MAY have an optional **label** (max **50** characters).
- **FR-024**: Attempting to save personal fields that violate FR-019–FR-023 MUST be blocked with field-level (or list-level) validation messages; no partial silent accept of invalid social entries.

### Key Entities

- **User (profile view)**: Signed-in identity with full name, email, optional avatar, optional phone, optional General (bio/about), optional birthday, optional address, optional social links, password credential presence, role/status (status/role not edited on self-service profile).
- **Avatar**: Optional profile image associated with a user; constrained by allowed types and max size; replaced by initials when absent.
- **Social link**: Optional list item on a user: optional display **label**, required **url** (http/https). Max 5 per user in P1.
- **Password change (self-service)**: Operation requiring current password + new password + confirmation for the signed-in user only.
- **Password set/reset (admin)**: Optional password field on access-control user edit; required on create; does not require the target user’s current password.
- **Locale preference (client)**: Active UI language among supported locales EN and VI, selectable via flags in admin user menu and public chrome.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A signed-in user can open Account settings and reach `/admin/profile` in under 5 seconds of intentional navigation on a stable connection.
- **SC-002**: 100% of tested valid personal-field updates on profile (name plus optional phone, General, birthday, address, social links) persist after a full page reload.
- **SC-003**: Valid avatar uploads within limits display correctly after reload; 100% of tested oversize/wrong-type uploads are rejected with a visible message.
- **SC-004**: Users can complete a successful self-service password change in under 2 minutes when they know the current password; wrong current password never silently succeeds.
- **SC-005**: Admins can update another user’s personal fields (name, phone, General, birthday, address, social links, and avatar if offered) from access control and optionally set a new password without knowing the old one.
- **SC-006**: In the admin user menu, language can be switched EN↔VI via flags with immediate UI update and a visible active indicator; the old sole “Switch to …” text toggle pattern is gone.
- **SC-007**: Unauthenticated users cannot view `/admin/profile` content (redirect to login).
- **SC-008**: 100% of tested invalid personal inputs (future birthday, non-http(s) social URL, over-max General/address/phone, more than 5 social links) are rejected with a visible message and do not persist.

## Assumptions

- Authentication (`001-auth`) remains the source of session and password minimum-length rules (at least 8 characters); this feature does not redesign login/register.
- Changing email, email verification, and forgot-password / reset-via-email flows remain **out of scope**.
- Profile entry is via the **admin shell** user menu; a separate public `/profile` for non-admin members is **out of scope** unless a later feature adds it.
- **Phone** and **avatar** attributes already exist on the user model in the current system; birthday, address, General, and social links may need **new persistence** designed in the plan phase.
- Allowed avatar types default to common web images (e.g. JPEG, PNG, WebP); max size default is **2 MB** unless architecture later documents a different limit — product acceptance is “limits exist and are enforced with clear messages.”
- Access-control permission model is unchanged: only users with existing manage-users capability can edit other users.
- UI copy in source defaults to English; Vietnamese appears via existing i18n when locale is VI. Field label **General** is the preferred English label for the bio/about field.
- Theme controls in the user menu are unchanged by this feature.
- P2 polish does not block declaring P1 done for handoff to architecture/implementation.
- Phone validation remains lenient in P1 (no country-code / E.164 enforcement) unless a later feature tightens it.

## Out of Scope

- Public marketing profile pages or member-facing `/profile` outside admin shell.
- Self-service email change, email verification, or password reset via email link.
- Additional locales beyond EN and VI.
- Editing another user’s profile from the self-service profile page.
- Role/permission changes from the profile page (remain in access control).
- Mandatory re-authentication / global logout of all sessions after password change (optional later hardening).
- Structured multi-field address forms (street, city, postal code, country as separate required fields).
- Fixed social-provider slots only (e.g. mandatory GitHub/LinkedIn/X fields) — P1 uses an extensible labeled URL list instead.
- Enforcing unique social URLs or verifying that URLs resolve / belong to the user.

## Dependencies

- **001-auth**: Session, current-user identity, password rules, login redirect.
- Existing **access-control** users table/forms: extend personal fields/avatar; keep admin password set/reset semantics.
- **003-landing-page** public locale flag pattern: behavioral consistency for admin menu flags (not a hard code dependency).
- Existing user attributes **phone** and **avatar** where already present; new attributes for birthday, address, General, and social links as required by planning.

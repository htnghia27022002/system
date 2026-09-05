# Feature Specification: Super Admin Flag and Full User Edit

**Feature ID**: `007-super-admin`

**Feature Branch**: `007-super-admin`

**Created**: 2026-09-05

**Updated**: 2026-09-05

**Status**: Draft — ready for tasks

**Input**: The Users admin page should work like Account settings: the highest-privilege operator can edit every account field. Add a **Super admin** flag on the user record. Anyone marked Super admin must skip role and permission checks on both the product UI and the API (they still must be signed in).

## Clarifications

### Session 2026-09-05 (user request — locked P1)

- Q: Is Super admin a role in the catalog? → A: **No.** It is a **boolean flag on the user**, independent of role. Do not add a Super admin role or new permission keys.
- Q: Who is the first Super admin? → A: The seeded administrator account (`admin@example.com`) is Super admin after migrate/seed.
- Q: How many Super admins? → A: **Many allowed.** The product **must not** remove, deactivate, or delete the **last** Super admin.
- Q: Who can grant or revoke the flag? → A: **Only a Super admin.** A non–Super admin with `users:modify` can still edit other user fields but cannot change the flag (sending it is rejected).
- Q: What does “skip role and permission checks” mean? → A: After a valid sign-in, Super admin **passes every view/modify gate** on the API and every permission/menu/page gate in the admin UI. Authentication is still required. Owner-only rules that protect the caller from harming themselves (cannot delete/deactivate own account; cannot remove the last Super admin) still apply.
- Q: Does Users edit match Account settings? → A: **Yes, same field class in tabs**: Profile (avatar + personal), Password (set/generate), Account (role, status, Super admin). **Admin viewing another user’s sessions stays deferred** (006 P2).
- Q: Does the Administrator role automatically become Super admin? → A: **No.** Highest role can already edit users through `users:modify`. Super admin is the separate bypass flag.

## Overview

Operators who manage people on **Users** (`/admin/users`) need two things:

1. A **complete edit surface** for another account — the same personal, password, and access fields they already have on **Account settings** for themselves.
2. A **Super admin** mark so a small set of trusted accounts can open every admin surface and call every protected admin API **without** depending on role permission rows.

| Surface | Path | Auth | Behavior |
|---------|------|------|----------|
| Users table | `/admin/users` | Signed in + `users:view` **or Super admin** | List users; show Super admin mark |
| Create / edit user | Dialog on `/admin/users` | Signed in + `users:modify` **or Super admin** | Tabs: Profile, Password, Account; Super admin checkbox only for Super admins |
| All other admin APIs / menus | Existing | Signed in | Super admin skips role/permission gates |

**Depends on**: `001-auth` (JWT, `/auth/me`), `004-user-profile` (Account settings fields), existing Users CRUD. **Does not** add a new sidebar item or permission catalog keys.

### Product roadmap / phased delivery

| Phase | Scope | Status |
|-------|--------|--------|
| **Phase 1 (P1)** | Persist Super admin flag; seed first Super admin; Super admin bypasses RBAC after sign-in (API + UI); Users table shows the flag; only Super admin can set/unset it; protect last Super admin; Users create/edit uses Profile / Password / Account tabs with the same personal and password fields as Account settings | **In scope** |
| **Phase 2 (P2)** | Admin viewing another user’s sessions; parsed device names / geo (still 006 P2) | **Deferred** |

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Super admin opens every admin surface (Priority: P1)

A Super admin signs in with **no** role permissions assigned (or an empty permission list). They can open Dashboard, Users, Roles, and Tools, and the API accepts their admin calls.

**Why this priority**: The flag is useless if the UI or API still requires catalog permissions.

**Independent Test**: Sign in as Super admin (seeded admin). Confirm sidebar shows all admin items and `/admin/users`, `/admin/roles`, `/admin/tools/webhooks` load. Optionally assign that account a role with zero permissions, refresh session, and confirm access remains.

**Acceptance Scenarios**:

1. **Given** a signed-in Super admin, **When** they open any existing admin page that is gated by a view permission, **Then** they are allowed through.
2. **Given** a signed-in Super admin, **When** they call an admin API that requires a modify permission, **Then** the API does not return forbidden for missing permissions.
3. **Given** a Super admin who is not signed in, **When** they call a protected API, **Then** they are rejected as unauthorized (the flag does not replace sign-in).

---

### User Story 2 - Users table shows Super admin (Priority: P1)

An operator who can list users sees which accounts are Super admin.

**Why this priority**: The flag must be visible so operators know who bypasses RBAC.

**Independent Test**: Open `/admin/users` as seeded admin. Confirm the seeded admin row is marked Super admin and a normal member is not.

**Acceptance Scenarios**:

1. **Given** users exist including at least one Super admin, **When** the Users table loads, **Then** Super admin accounts show a Super admin mark and others do not.
2. **Given** a Super admin account, **When** the same list is read from the Users API, **Then** each user payload includes the Super admin flag.

---

### User Story 3 - Highest-privilege operator edits a user like Account settings (Priority: P1)

An operator with user-modify access (or Super admin) opens create/edit user and can change the same personal and password fields as Account settings, plus role, status, and (if they are Super admin) the Super admin flag.

**Why this priority**: The Users page must not be a weaker form of Account settings.

**Independent Test**: As seeded admin, edit another user. Change name, personal fields, password (optional generate), role/status. Save. Reopen and confirm values. Confirm tabs Profile / Password / Account are present. Sessions for that other user are **not** required.

**Acceptance Scenarios**:

1. **Given** an operator who may modify users, **When** they open create or edit user, **Then** they see **Profile**, **Password**, and **Account** tabs.
2. **Given** the Profile tab, **When** they edit, **Then** they can set name, email, avatar (edit), phone, birthday, general, address, and social links.
3. **Given** the Password tab, **When** they set or generate a password, **Then** create requires a password and edit may leave it blank to keep the current password.
4. **Given** the Account tab, **When** they save, **Then** role and status persist.
5. **Given** the actor is Super admin, **When** they open Account, **Then** they can set or clear Super admin on that user (except last Super admin — see US4).
6. **Given** the actor is not Super admin, **When** they open Account, **Then** they do **not** get a working Super admin control, and the API rejects an attempt to change the flag.

---

### User Story 4 - Last Super admin cannot be removed (Priority: P1)

The product always keeps at least one Super admin.

**Why this priority**: Losing the last bypass account would lock operators behind empty permission sets.

**Independent Test**: With only one Super admin, attempt to clear the flag, deactivate, or delete that account. All three fail. Grant Super admin to a second user, then the first can be demoted.

**Acceptance Scenarios**:

1. **Given** exactly one Super admin, **When** anyone tries to clear that flag, deactivate that account, or delete that account, **Then** the action is rejected.
2. **Given** two Super admins, **When** one Super admin clears the flag on the other, **Then** the change succeeds.
3. **Given** a Super admin, **When** they try to delete or deactivate **their own** account, **Then** the existing self-protection rules still apply.

---

### User Story 5 - Ordinary permission checks still apply to everyone else (Priority: P1)

A signed-in member who is **not** Super admin still needs the matching view/modify permission.

**Why this priority**: Bypass must be exclusive to the flag.

**Independent Test**: Sign in as a user with only `dashboard:view`. Confirm Users/Roles/Webhooks stay hidden and those APIs return forbidden.

**Acceptance Scenarios**:

1. **Given** a signed-in user who is not Super admin and lacks `users:view`, **When** they open `/admin/users` or call the Users list API, **Then** they are denied.
2. **Given** a signed-in user who is not Super admin and has `users:view` but not `users:modify`, **When** they try to create or update a user, **Then** they are denied.

## User Story 6 - Super admin grant is not a role assignment (Priority: P1)

Changing a user’s role to Administrator does **not** make them Super admin. Clearing Super admin does **not** remove their role.

**Independent Test**: Promote a member to Administrator without the Super admin flag. Confirm they still need catalog permissions (Administrator seed has them). Confirm they do not bypass if those permissions are later removed. Grant Super admin without changing role; they bypass.

**Acceptance Scenarios**:

1. **Given** a user with the Administrator role and Super admin false, **When** they use the product, **Then** access follows that role’s permissions, not a silent bypass.
2. **Given** a user with a member role and Super admin true, **When** they use the product, **Then** they bypass permission checks.

## Edge cases

- Guest / expired token: still unauthorized.
- JWT issued before the flag existed: treat missing claim as false, then refresh from the user record on each authenticated request when the user store is available.
- Granting Super admin to yourself when you already are Super admin: no-op success.
- Clearing your own Super admin when you are not the last: allowed (you immediately lose bypass on the next request / session refresh).
- Soft-deleted users do not count toward “last Super admin”.
- Search, webhooks capture (public), and login/register stay unchanged.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Persist `is_super_admin` (boolean, default false) on the user record.
- **FR-002**: After migrate/seed, the seeded administrator account is Super admin.
- **FR-003**: Auth payloads (`/auth/me`, login/register/refresh/OAuth user object, access-token claims) include `superAdmin`.
- **FR-004**: Super admin skips every existing view/modify permission middleware check. Sign-in is still required.
- **FR-005**: Super admin skips every existing admin UI permission gate (page guard, nav item, action gate).
- **FR-006**: Only a Super admin may set or clear `superAdmin` on create/update. Others receive forbidden if they send the field.
- **FR-007**: The product rejects clearing, deactivating, or deleting the last Super admin.
- **FR-008**: Users list and get responses include `superAdmin`.
- **FR-009**: Users create/edit UI uses Profile / Password / Account tabs covering the same personal and password fields as Account settings, plus role, status, and Super admin (when allowed).
- **FR-010**: Users table (desktop and mobile card) shows the Super admin mark.
- **FR-011**: Non–Super admin access is unchanged: catalog permissions still apply.
- **FR-012**: Super admin is not a role and does not add catalog keys or a sidebar item.
- **FR-013**: Admin viewing another user’s sessions is **out of P1**.

### Key Entities

- **User**: existing account; gains Super admin flag.
- **Role / Permission**: unchanged catalog; Super admin does not read these for gate decisions.
- **Session (JWT)**: carries Super admin so the UI can hydrate before `/auth/me`; server still re-reads the flag from the user record when possible.

## Success Criteria *(mandatory)*

- **SC-001**: A Super admin with an empty permission list can open all current admin pages and complete a Users edit in one visit without a permission error.
- **SC-002**: A non–Super admin without `users:view` cannot open Users (UI and API) in 100% of Independent Test attempts.
- **SC-003**: Operators can identify Super admin rows on Users without opening the edit dialog.
- **SC-004**: Attempting to remove the last Super admin fails and at least one Super admin remains.
- **SC-005**: Users edit exposes Profile, Password, and Account in the same visit (no extra route).

## Assumptions

- Seeded admin email/id stay the existing demo identifiers.
- Multiple Super admins are acceptable for operations hand-off.
- Re-login or `/auth/me` / access-token refresh is enough for the UI to pick up a newly granted or revoked flag; server middleware re-reads the database flag so API bypass updates without waiting for token expiry.
- Vietnamese UI copy lives only in `fe/src/locales/vi/**`.

## Out of scope (P1)

- New Super admin role or permission keys
- Admin viewing or revoking another user’s sessions
- Audit log of who granted Super admin
- Limiting Super admin to a single reserved account forever
- Changing public webhook capture or guest routes

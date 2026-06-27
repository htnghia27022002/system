# Feature Specification: Authentication (Email, Google OAuth, JWT)

**Feature ID**: `001-auth`

**Created**: 2026-06-27

**Status**: Draft

**Input**: Auth: email register & login, Google OAuth, JWT + refresh, modern friendly login/register UI

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Sign in with email and password (Priority: P1)

A returning user opens the login page, enters a registered email and password, and lands in the app with an active session. Access tokens expire on a schedule; the session continues seamlessly via refresh without forcing re-login while the refresh token is valid.

**Why this priority**: Email/password is the baseline auth path and unblocks all protected features.

**Independent Test**: Register a user (or use seed account), log in at `/login`, verify redirect to the post-login destination and that authenticated API calls succeed.

**Acceptance Scenarios**:

1. **Given** a registered active user, **When** they submit valid email and password on `/login`, **Then** they receive access and refresh tokens and are redirected to the appropriate home or admin route.
2. **Given** a logged-in user whose access token is expired but refresh token is valid, **When** they call a protected API, **Then** the client refreshes the access token automatically and retries the request without showing the login form.
3. **Given** invalid credentials, **When** the user submits the login form, **Then** a clear, non-technical error message is shown and no tokens are stored.
4. **Given** an inactive or disabled account, **When** login is attempted, **Then** access is denied with an appropriate message.

---

### User Story 2 - Register a new account with email (Priority: P1)

A new user creates an account with name, email, and password, then is signed in without a separate login step.

**Why this priority**: Self-service registration is required alongside login for a complete credential-based auth flow.

**Independent Test**: Complete `/register` with a new email, confirm automatic sign-in and access to protected routes according to default role.

**Acceptance Scenarios**:

1. **Given** a visitor on `/register`, **When** they submit valid name, email, password, and matching confirmation, **Then** an account is created, tokens are issued, and they are redirected like a successful login.
2. **Given** an email already registered, **When** registration is submitted, **Then** the user sees a clear error and remains on the register page.
3. **Given** weak or invalid input (email format, password rules, mismatched confirmation), **When** the form is submitted, **Then** field-level validation errors are shown before any API call where possible.

---

### User Story 3 - Sign in with Google (Priority: P2)

A user chooses “Continue with Google” on login or register, completes Google consent, and returns to the app signed in. If the Google email matches an existing account, accounts are linked or the user is signed into that account per linking rules.

**Why this priority**: OAuth reduces friction and matches the product ask; depends on P1 token/session behavior.

**Independent Test**: With Google OAuth configured in the environment, start OAuth from the login page, complete callback, verify session and `/api/auth/me` response.

**Acceptance Scenarios**:

1. **Given** Google OAuth is configured, **When** the user clicks “Continue with Google” on `/login`, **Then** they are redirected to Google and, after success, return to the app with valid JWT tokens.
2. **Given** a Google account whose email is not yet registered, **When** OAuth completes, **Then** a user record is created (or linked) and the user receives tokens with the default member role unless policy says otherwise.
3. **Given** OAuth is cancelled or fails, **When** the user returns to the app, **Then** they see a friendly error on the auth page and no partial session is left.
4. **Given** Google OAuth credentials are missing in the environment, **When** the auth pages load, **Then** the Google button is hidden or disabled with no broken redirect.

---

### User Story 4 - Sign out and session end (Priority: P2)

A signed-in user can log out; refresh tokens are revoked server-side and local tokens are cleared.

**Why this priority**: Required for security and shared-device use.

**Independent Test**: Log in, log out, confirm protected routes redirect to login and refresh no longer works.

**Acceptance Scenarios**:

1. **Given** a signed-in user, **When** they log out, **Then** refresh token is revoked (if provided), local tokens are cleared, and subsequent protected requests require login.
2. **Given** a user logs out on one device, **When** another device uses an old refresh token, **Then** refresh fails and that client must log in again.

---

### User Story 5 - Modern, friendly auth UI (Priority: P1)

Login and register pages are visually clear, accessible, and consistent with the app design system: readable hierarchy, helpful labels, loading states, errors, dark mode support, and obvious paths between login and register.

**Why this priority**: UX is explicitly in scope; poor auth UI blocks adoption even if APIs work.

**Independent Test**: Visual and accessibility review of `/login` and `/register` on mobile and desktop; no mock API in production Docker stack.

**Acceptance Scenarios**:

1. **Given** a visitor on `/login` or `/register`, **When** the page loads, **Then** layout is centered, branded, supports light/dark theme, and primary actions are obvious.
2. **Given** a form submission in progress, **When** the user waits, **Then** the submit button shows a loading state and prevents double submit.
3. **Given** API or network errors, **When** auth fails, **Then** errors appear near the form in plain language (English UI copy in source files).
4. **Given** login page, **When** rendered, **Then** user can navigate to register and vice versa; Google option is visually secondary but discoverable (divider + button pattern).

---

### Edge Cases

- Access token expired and refresh token expired or revoked → redirect to `/login`, clear stored tokens.
- User registers with password then later uses Google with the same email → consistent account identity (link or single sign-in path).
- CORS / wrong API base URL in Docker → login fails with visible error (environment configuration, not silent failure).
- Rate limiting or brute force → generic error message without revealing whether email exists (align with existing BE behavior).
- Password-only account attempting OAuth-only login path → guided to email login or Google as appropriate.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST allow registration with name, email, and password.
- **FR-002**: System MUST allow login with email and password for active accounts.
- **FR-003**: System MUST issue a short-lived JWT access token and a longer-lived refresh token on successful login or registration.
- **FR-004**: System MUST expose a refresh endpoint that rotates refresh tokens and returns a new access token.
- **FR-005**: System MUST revoke refresh tokens on logout when a refresh token is supplied.
- **FR-006**: System MUST expose a current-user endpoint (`/auth/me`) for authenticated clients.
- **FR-007**: System MUST support Google OAuth sign-in (start + callback flow) when provider credentials are configured.
- **FR-008**: Frontend MUST persist tokens securely in browser storage and attach access tokens to API requests.
- **FR-009**: Frontend MUST automatically refresh expired access tokens using the refresh token before failing authenticated requests.
- **FR-010**: Frontend MUST provide `/login` and `/register` routes with form validation and accessible labels.
- **FR-011**: Frontend MUST provide a “Continue with Google” entry point on login (and optionally register) wired to the backend OAuth start URL.
- **FR-012**: Frontend MUST use the real backend API when `NEXT_PUBLIC_USE_MOCK_API=false` (Docker default).
- **FR-013**: Auth pages MUST use the shared auth layout and design tokens (no one-off styling that breaks dark mode).

### Key Entities

- **User**: Identity with email, name, role, status (active/inactive); may have password hash and/or linked OAuth accounts.
- **RefreshToken**: Opaque refresh token (hashed at rest), expiry, revocation timestamp, tied to user.
- **OAuthAccount**: Provider (google), provider user id, linked user, provider email.
- **Session (client)**: Access JWT + refresh token in browser storage; derived user profile in app state.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: A new user can register and reach the post-login destination in under 60 seconds on a stable connection.
- **SC-002**: Login success rate for valid seed credentials is 100% through nginx (`system.local:8080`) with correct CORS origins.
- **SC-003**: Access token refresh succeeds without user interaction for at least one API call after access token expiry (simulated or short TTL in dev).
- **SC-004**: Google OAuth completes end-to-end when `GOOGLE_*` env vars are set (manual QA checklist).
- **SC-005**: Auth pages score no critical accessibility issues: labels linked to inputs, focus visible, errors announced via text.

## Assumptions

- Backend auth API routes under `/api/auth/*` largely exist; this feature completes integration, OAuth UI, polish, and gap fixes rather than greenfield API design.
- Default role for self-registration is **Member**; admin accounts are seeded or provisioned separately.
- Google OAuth requires valid Google Cloud OAuth client IDs and redirect URI matching `OAUTH_REDIRECT_URL` in environment.
- UI copy in source files is **English**; i18n keys in `locales/` may add Vietnamese later without changing this spec’s acceptance language.
- Email verification and password reset are **out of scope** for v1 unless added in a follow-up feature.
- Social providers other than Google are out of scope for v1.

## Current Baseline (context for design phase)

| Area | Status |
|------|--------|
| BE login, register, refresh, logout, me | Implemented |
| BE Google OAuth routes | Scaffolded; requires env credentials |
| FE login/register forms | Implemented (basic) |
| FE token storage + refresh interceptor | Implemented |
| FE Google OAuth button + callback page | **Not implemented** |
| FE auth UI polish (modern friendly) | **Partial** — functional but minimal |

*Detailed gap analysis belongs in phase 3 (`be-implement.md`, `fe-implement.md`) — not in this spec.*

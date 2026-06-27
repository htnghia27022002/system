# Tasks: Authentication (001-auth)

**Input**: [spec.md](spec.md)

## Phase 2 — Work breakdown

### [BE] Backend

- [x] T001 [BE] Add `GET /api/auth/oauth/providers` — list configured OAuth providers
- [x] T002 [BE] Verify OAuth callback links Google account to existing email user

### [FE] Frontend — OAuth & session

- [x] T003 [FE] Add `NEXT_PUBLIC_SITE_URL` to env for OAuth callback URL
- [x] T004 [FE] Add `authApi.oauthCallback`, `getOAuthProviders`, fix `logout` with refresh token
- [x] T005 [FE] Add `/auth/callback` page — exchange code for JWT
- [x] T006 [FE] Add Google sign-in button + divider on login/register
- [x] T007 [FE] Add `useSignOut` hook — revoke refresh token on logout

### [FE] Frontend — UI polish

- [x] T008 [FE] Polish `AuthLayout` with card surface and spacing
- [x] T009 [FE] Add i18n strings for OAuth and form errors
- [x] T010 [FE] Show inline form error banner on login/register failure

### [QA]

- [ ] T011 [QA] Manual: email login, register, logout, refresh via API client
- [ ] T012 [QA] Manual: Google OAuth when credentials configured

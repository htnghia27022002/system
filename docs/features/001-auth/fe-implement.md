# FE implement: Authentication

**Feature:** `docs/features/001-auth/`  
**Based on:** [spec.md](spec.md), [be-implement.md](be-implement.md)

## Analysis summary

Email auth and token refresh work. Missing: Google OAuth UI, callback page, logout with refresh revocation, auth UI polish.

## OAuth flow

```text
Login → "Continue with Google"
  → GET /api/auth/oauth/google/start?redirect_uri={SITE}/auth/callback
  → Google consent
  → redirect to /auth/callback?code=...
  → POST /api/auth/oauth/google/callback { code, redirectUri }
  → signIn → redirect home/admin
```

## Routes

| Route | Component |
|-------|-----------|
| `/login` | `LoginPage` + Google button |
| `/register` | `RegisterPage` + Google button |
| `/auth/callback` | `OAuthCallbackPage` |

## Files changed

| Path | Change |
|------|--------|
| `src/config/env.ts` | `VITE_SITE_URL` |
| `src/features/auth/services/auth-api.ts` | oauth + logout fix |
| `src/features/auth/services/oauth.ts` | start URL, callback URI helpers |
| `src/features/auth/hooks/use-oauth-providers.ts` | fetch providers |
| `src/features/auth/hooks/use-sign-out.ts` | logout API + clear store |
| `src/features/auth/components/google-sign-in-button.tsx` | OAuth CTA |
| `src/features/auth/components/auth-divider.tsx` | "or continue with" |
| `src/features/auth/components/oauth-callback-page.tsx` | callback handler |
| `src/features/auth/components/login-form.tsx` | divider + Google |
| `src/features/auth/components/register-form.tsx` | divider + Google |
| `src/components/common/auth-layout.tsx` | card polish |
| `src/app/(auth)/auth/callback/page.tsx` | route |
| `src/layouts/main-layout.tsx` | useSignOut |
| `src/locales/en/common.json` | OAuth copy |

## Env

| Variable | Example |
|----------|---------|
| `NEXT_PUBLIC_SITE_URL` | `http://system.local:8080` |
| `NEXT_PUBLIC_API_BASE_URL` | `http://system.local:8080/api` |

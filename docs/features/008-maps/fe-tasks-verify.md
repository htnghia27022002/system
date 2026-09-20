# FE tasks verify: Maps

**Feature:** `docs/features/008-maps/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@fe`  
**Date:** 2026-09-12

## Summary

All `[FE]` tasks for 008-maps are implemented. The admin Maps surface lives at `/admin/maps` behind `maps:view`. Product UI talks only to a `MapProvider` port (live OSM adapter + Google stub). Sources, Load data, and Place status are gated by `maps:modify`. `make test-fe` from the repo root passed (60 tests). Guest `/admin/maps` redirects to login. Signed-in menu, map, filter, pin, Sources, and Load data clicks were not completed in this pass (login session did not establish in the browser tool).

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T002 | Feature module + thin `/admin/maps` route; no guest map | Done | `fe/src/features/maps/`; `fe/src/app/admin/maps/page.tsx` |
| T003 | FE-owned map env only | Done | `fe/.env.example` (`NEXT_PUBLIC_MAP_PROVIDER`, `NEXT_PUBLIC_GOOGLE_MAPS_API_KEY`); not added to root `.env` or `be/.env` |
| T013 | `PermissionKeys.maps` + mock catalog | Done | `permission-keys.ts` additive `maps: { view, modify }`; `permission-keys.test.ts`; `access-control.mock.ts` |
| T014 | Nav after Dashboard, before Tools | Done | `use-admin-nav-items.ts` `href: '/admin/maps'` + `hasPermission`; `nav.maps` = Maps / Bản đồ |
| T015 | MapProvider port | Done | `providers/map-provider.ts` — lat/lng types only |
| T016 | Live OSM adapter | Done | `providers/osm-map-provider.ts` (Leaflet + OSM tiles, confined here) |
| T017 | Google stub | Done | `providers/google-map-provider.ts` — “Google Maps is not configured”, no throw |
| T018 | Registry + env Zod | Done | `create-map-provider.ts`; `fe/src/config/env.ts` `MAP_PROVIDER` default `osm` |
| T019 | Types + API client | Done | `types.ts`, `services/maps-api.ts` camelCase (`placeKey`, `originalUrl`, `itemsSuccess`, …) |
| T020 | Thin page + view guard | Done | `PermissionGuard` + `PermissionKeys.maps.view`; `robots: { index: false }` |
| T021 | Maps page via port only | Done | `components/maps-page.tsx` — no OSM/Google types |
| T022 | Map canvas mount/unmount | Done | `components/map-canvas.tsx` |
| T023 | Empty-pins i18n | Done | `maps.emptyPins` EN + VI; map still mounts |
| T024 | Guest redirect | Done | Browser: `/admin/maps` → `/login?from=%2Fadmin%2Fmaps`. `ProtectedGuard` unchanged |
| T027 | Filters + places hook | Done | `map-filters.tsx`, `use-map-places.ts` |
| T028 | Pins via port only | Done | `clearPins` / `addPin` / `fitBounds` |
| T029 | Distinct empty-filter state | Done | `maps.emptyFilter` vs empty pins vs ingest chrome |
| T031 | Pin click → detail | Done | `place-detail-panel.tsx` |
| T032 | News attribution + modify gates | Done | Source name + original URL; status/Load data behind `PermissionGate` |
| T033 | Place detail hook | Done | `use-place-detail.ts` |
| T036 | Sources panel on same page | Done | Tab on `/admin/maps`; no second menu item |
| T037 | Sources CRUD gates | Done | `PermissionGate` modify; view-only list + hint |
| T038 | Fixed mapping fields only | Done | Location / Place / News / details allow-list; no add-column |
| T039 | Sources hooks + i18n | Done | `use-map-sources.ts`; `admin.json` `maps.sources.*` |
| T045 | Load data + ingest status | Done | `ingest-status.tsx`; poll while `queued`/`running` |
| T046 | Invalidate places after ingest | Done | `use-map-ingest.ts`; in-progress is not treated as empty-finished |
| T050 | Overlapping pins | Done | Distinct `pin.id` keys at the same lat/lng; no clustering |
| T053 | No ungated Maps link | Done | Single nav item filtered by `maps.view`; page view guard; mutate gates |
| T055 | Place status control | Done | `active`/`hidden` in detail for modify users |
| T057 | Vitest | Done | `create-map-provider.test.ts`, `category-filters.test.ts` |
| T059 | This verify report | Done | After `make test-fe` + guest browser check |

## Verification commands

```bash
make test-fe
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-fe` | Pass | 18 files / 60 tests, including maps provider + category helpers |
| `pnpm lint` | Fail (pre-existing) | Maps-only issues fixed (`map-canvas` ref, `pins` memo). Remaining errors are landing hero / unused landing imports / `copy-to-clipboard`, not introduced by this feature |

GitNexus MCP was unavailable. Kernel edits (`PermissionKeys`, `useAdminNavItems`, `env.ts`) were **additive only**.

## Acceptance coverage (FE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| Guest `/admin/maps` → login | T024 browser | Pass — `/login?from=%2Fadmin%2Fmaps` |
| Maps menu + street map (empty pins OK) | T014, T021–T023 | Pass (code); signed-in click not completed |
| Category / search / division filters | T027–T029 | Pass (code + unit helpers) |
| Pin detail + news attribution | T031–T033 | Pass (code) |
| Sources CRUD + mapping allow-list | T036–T039 | Pass (code) |
| Load data poll + pin refresh | T045–T046 | Pass (code) |
| Two pins at one lat/lng | T050 | Pass (OSM marker map keyed by id) |
| View vs modify gates | T013, T014, T053, T055 | Pass (code); Super admin reuse via existing `hasPermission` |
| Independent Tests on `system.local` | After sign-in | Partial — see gaps |

## Browser verify

| Step | Result |
|------|--------|
| Open `http://system.local:8080/admin/maps` as guest | Redirected to login (`from=/admin/maps`). Guards not weakened |
| Sign in and click Maps, filters, pins, Sources, Load data | **Not completed.** Browser login did not establish a session after filling the sign-in form. Credentialed API login was not run in this pass |

## Gaps / follow-ups

- [ ] Signed-in walkthrough on `http://system.local:8080/admin/maps` (menu, empty OSM map, filters, pin detail + news, Sources CRUD, Load data polling) — `@qa` or a follow-up FE pass after a working admin session
- [ ] Docker FE `fe_node_modules` volume needed a `pnpm install --config.lockfile=false` so Leaflet resolved; if the container is recreated without that, run install again or rebuild the `fe` image so `leaflet` is present
- [ ] `pnpm lint` still fails on pre-existing landing files (out of scope)

## Sign-off (FE)

- [x] All claimed `[FE]` tasks done or explicitly deferred above
- [x] `make test-fe` passed
- [x] Matches `plan.md` FE sections (MapProvider port, FE-owned env, admin-only `/admin/maps`)

**Ready for `@qa`:** yes for unit/code coverage and guest redirect. Live signed-in Independent Tests still need a logged-in browser session.

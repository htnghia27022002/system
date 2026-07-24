# FE tasks verify: Public Landing Rebuild + Tools IA + Shared Chrome

**Feature:** `docs/features/003-landing-page/`  
**Based on:** [spec.md](spec.md), [tasks.md](tasks.md), [plan.md](plan.md)  
**Agent:** `@fe`  
**Date:** 2026-07-24

## Summary

Greenfield FE rebuild for Phase 1: shared `PublicSiteHeader` + `LocaleSelect` (Select + flags EN/VI), Tools hub with Webhooks nested shell at `/tools/webhooks`, rebuilt landing scrollytelling (tokenomics retired), catalog-driven Tools teaser, and admin **Public home** control. Prior verify claims are obsolete and not reused. `make test-fe` passed (25 tests); `pnpm lint` reported 0 errors (pre-existing warnings only).

## Tasks completed

| Task ID | Description | Status | Evidence |
|---------|-------------|--------|----------|
| T001 | Confirm GSAP; greenfield landing target | Done | `fe/package.json` has `gsap`; landing rebuilt under `features/landing/` |
| T002 | Align tools scaffolding to spec | Done | `catalog.ts`, `ToolsHubPage`, `WebhooksToolPage`, barrel exports |
| T003 | Confirm routes `/`, `/tools`, `/tools/webhooks` | Done | App Router pages under `app/(public)/` |
| T004 | `LocaleSelect` Select + flags | Done | `components/common/locale-select.tsx` |
| T005 | `PublicSiteHeader` shared chrome | Done | `components/common/public-site-header.tsx` + mobile Sheet |
| T006 | Catalog model FR-009 | Done | `features/tools/types.ts` + helpers |
| T007 | Seed Webhooks available + comingSoon | Done | `catalog.ts` Webhooks `href: '/tools/webhooks'` |
| T008 | Webhooks route + tools layout | Done | `app/(public)/tools/webhooks/page.tsx`, `tools/layout.tsx` |
| T009 | LandingLayout only on landing | Done | `LandingPage` wraps `LandingLayout`; tools scroll normally |
| T010 | Tools public barrel | Done | `features/tools/index.ts` |
| T011 | Sitemap `/tools/webhooks` | Done | `app/sitemap.ts` |
| T012 | Rebuild LandingPage + content | Done | New section set; no tokenomics in composition |
| T013 | Retire tokenomics section | Done | Deleted `tokenomics-snap-section.tsx` |
| T014 | GSAP + reduced motion | Done | Section hooks + `useReducedMotion` |
| T015 | ≥3 motion boundaries | Done | Hero, features, tools, CTA scrub/pin paths |
| T016 | LandingPage tests | Done | `LandingPage.test.tsx` |
| T017 | Hero hierarchy | Done | Brand, headline, support, CTAs + hero image |
| T018 | Responsive hero / page wiring | Done | `(public)/page.tsx` → `LandingPage` |
| T019 | Header on landing | Done | `landing-nav.tsx` mounts `PublicSiteHeader` |
| T020 | Header on tools + webhooks | Done | `ToolsHubPage`, `WebhooksToolPage` |
| T021 | Mobile menu destinations | Done | Sheet exposes Home, Tools, auth |
| T022 | Guest vs signed-in chrome | Done | `useAuthStore` + `useSignOut` in header |
| T023 | Logout returns guest chrome | Done | Same shared header on all P1 routes |
| T024 | ThemeToggle in header | Done | `ThemeToggle` always present; `next-themes` |
| T025 | LocaleSelect; remove toggle | Done | No EN↔VI toggle on landing/tools |
| T026 | i18n chrome strings | Done | `locales/en|vi/common.json` (+ admin publicHome) |
| T027 | MainLayout LocaleSelect | Done | `layouts/main-layout.tsx` |
| T028 | Hub listing UI | Done | `ToolsHubPage` + `ToolCard` |
| T029 | Webhooks → nested route | Done | Catalog href + hub Open tool |
| T030 | Coming-soon non-links | Done | Disabled button when not navigable |
| T031 | Webhooks shell page | Done | `webhooks-tool-page.tsx` |
| T032 | Wire webhooks route + metadata | Done | `tools/webhooks/page.tsx` |
| T033 | Webhooks placeholder copy | Done | `tools/content.ts` |
| T034 | Tools teaser from catalog | Done | `tools-snap-section.tsx` → `getToolsTeaser()` |
| T035 | Teaser CTA → `/tools` | Done | `landingContent.tools.primaryCtaHref` |
| T036 | Catalog append docs | Done | Comment block atop `catalog.ts` |
| T037 | Catalog tests | Done | `catalog.test.ts` (Webhooks href, comingSoon) |
| T038 | Barrel-only imports | Done | Landing imports `@/features/tools` |
| T039 | Admin home icon | Done | `app-sidebar.tsx` → `/` with `shell.publicHome` |
| T040 | Dedicated home; admin logo stays `/admin` | Done | Logo → `/admin`; home control → `/` |
| T041 | Remove public toggle UX | Done | Grep clean under landing/tools |
| T042 | Refresh metadata | Done | `/`, `/tools`, `/tools/webhooks` pages |
| T043 | `make test-fe` | Done | 25 passed |
| T044 | Rewrite this verify file | Done | This document |
| T045–T046 | QA manual | Skipped | Owned by `@qa` |
| T047 | US12 deferred | Deferred | Out of Phase 1 |

## Verification commands

```bash
make test-fe
cd fe && pnpm lint
```

| Command | Result | Notes |
|---------|--------|-------|
| `make test-fe` | Pass | 9 files / 25 tests |
| `pnpm lint` (in `fe/`) | Pass | 0 errors; 5 pre-existing warnings unrelated to this feature |
| `pnpm exec tsc --noEmit` (in `fe/`) | Pass | |

Also fixed `fe/src/components/ui/select.tsx` import to `@/lib/utils` so Select resolves under Vitest.

## Acceptance coverage (FE-relevant)

| Spec scenario | Covered by | Result |
|---------------|------------|--------|
| US1 greenfield scrolly + reduced motion | Landing rebuild + GSAP gates | Pass (code) |
| US2 first-viewport hierarchy | Hero section + content | Pass (code) |
| US3 Navigate Home/Tools/auth | `PublicSiteHeader` | Pass (code) |
| US4 Login/Logout chrome | Auth slot in header | Pass (code) |
| US5 Theme persistence | `ThemeToggle` + next-themes | Pass (code) |
| US6 Locale Select + flags | `LocaleSelect` on P1 routes | Pass (code) |
| US7 Hub + open Webhooks | Catalog + hub | Pass (code/tests) |
| US8 Webhooks shell | Route + shell page | Pass (code) |
| US9 Teaser → `/tools` | Teaser CTA + catalog | Pass (tests) |
| US10 Extensible catalog | Append docs + tests | Pass (tests) |
| US11 Admin → public `/` | Sidebar home control | Pass (code) |
| US12 additional tools | Deferred | N/A |

## Gaps / follow-ups

- [x] None blocking for `@qa` Phase 1 manual pass
- [ ] `@qa` should run T045–T046 matrix (guest/signed-in, theme, locale Select, reduced-motion, admin home, direct `/tools/webhooks`)
- [ ] Optional: remove leftover `.landing-section-tokenomics` CSS in `fe/src/styles/index.css` (unused; harmless)
- [ ] Optional: i18n for landing/tools marketing English content modules (P1 used English content modules + chrome i18n)

## Sign-off (FE)

- [x] All claimed `[FE]` Phase 1 tasks done; T047 deferred; QA tasks not claimed
- [x] Tests / lint listed above passed (lint warnings pre-existing)
- [x] Matches `plan.md` FE sections (shared chrome in `components/common`, catalog FE-only, layout scroll split, Webhooks nested route)
